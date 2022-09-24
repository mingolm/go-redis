package pool

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrMaxPoolSize    = errors.New("connect pool size is max")
	ErrExceedDeadline = errors.New("connect pool exceed deadline")
)

type Pooler interface {
	WithConn(context.Context, func(context.Context, *Conn) error) error
	Get(context.Context) (*Conn, error)
	Put(context.Context, *Conn) error
}

func NewPool(opt *Options) Pooler {
	p := &pool{
		Options: opt,
		conns:   make(chan *Conn, opt.PoolSize),
	}
	p.PoolSize.Swap(opt.PoolSize)
	p.MinIdleConns.Swap(opt.MinIdleConns)
	p.MaxIdleConns.Swap(opt.MaxIdleConns)
	p.initConnects()
	return p
}

type pool struct {
	*Options
	PoolSize     atomic.Int32 // 连接池长度
	MinIdleConns atomic.Int32 // 最小空闲连接
	MaxIdleConns atomic.Int32 // 最大空闲连接
	conns        chan *Conn
}

func (p *pool) initConnects() {
	for {
		cur := p.MaxIdleConns.Add(-1)
		if cur <= 0 {
			break
		}
		go func() {
			if err := p.addConnect(); err != nil {
				p.Logger.Errorw("pool connect dialer failed",
					"err", err,
				)
			}
		}()
	}
}

func (p *pool) WithConn(ctx context.Context, fn func(context.Context, *Conn) error) error {
	conn, err := p.Get(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err = p.Put(ctx, conn); err != nil {
			p.Logger.Errorw("connect put failed",
				"err", err,
			)
		}
	}()

	return fn(ctx, conn)
}

func (p *pool) Get(ctx context.Context) (*Conn, error) {
	for {
		select {
		case conn := <-p.conns:
			if !p.connHealthCheck(conn) {
				p.Logger.Debugw("connect un-health",
					"created_at", conn.createdAt,
					"used_at", conn.usedAt,
				)
				if err := p.connClose(conn); err != nil {
					p.Logger.Errorw("connect close failed",
						"err", err,
					)
				}
				continue
			}
			return conn, nil
		case <-ctx.Done():
			return nil, ErrExceedDeadline
		default:
			if err := p.addConnect(); err != nil {
				if errors.Is(err, ErrMaxPoolSize) {
					time.Sleep(time.Millisecond * 10)
					continue
				}
				return nil, err
			}
		}
	}
}

func (p *pool) Put(ctx context.Context, conn *Conn) error {
	if !p.connHealthCheck(conn) {
		return p.connClose(conn)
	}

	select {
	case p.conns <- conn:
		return nil
	default:
		p.Logger.Error("connect pool is full")
		return nil
	}
}

func (p *pool) addConnect() error {
	if p.PoolSize.Load() <= 0 {
		return ErrMaxPoolSize
	}

	conn, err := p.Dialer(context.TODO())
	if err != nil {
		return err
	}
	select {
	case p.conns <- NewConnect(conn):
		p.PoolSize.Store(-1)
		return nil
	default:
		p.Logger.Errorw("connect pool add failed",
			"err", errors.New("full"),
		)
		return nil
	}
}

func (p *pool) connHealthCheck(conn *Conn) bool {
	now := time.Now()
	// 最大生命周期
	if p.ConnMaxLifetime > 0 && conn.createdAt.Add(p.ConnMaxLifetime).Before(now) {
		return false
	}

	// 最大空闲时间
	if p.ConnMaxIdleTime > 0 && conn.usedAt.Add(p.ConnMaxIdleTime).Before(now) {
		return false
	}

	if conn.check() != nil {
		return false
	}

	conn.usedAt = time.Now()

	return true
}

func (p *pool) connClose(conn *Conn) error {
	if err := conn.netConn.Close(); err != nil {
		return err
	}

	switch conn.typ {
	case connTypPersistence:
		p.MinIdleConns.Store(-1)
		p.initConnects()
	case connTypBackup:
		p.MaxIdleConns.Store(-1)
	case connTypTmp:
		p.PoolSize.Store(1)
	}

	return nil
}
