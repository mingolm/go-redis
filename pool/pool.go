package pool

import (
	"bufio"
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrConnExhaust    = errors.New("connect exhaust")
	ErrMaxPoolSize    = errors.New("connect pool size is max")
	ErrExceedDeadline = errors.New("connect pool exceed deadline")
)

type Pooler interface {
	WithConn(context.Context, func(context.Context, *Conn) error) error
	// Get 获取连接
	Get(context.Context) (*Conn, error)
	Put(context.Context, *Conn) error
	Destroy(context.Context, *Conn) error
}

func NewPool(opt *Options) Pooler {
	if err := opt.init(); err != nil {
		panic("pool config error")
	}

	p := &pool{
		opt:   opt,
		conns: make(chan *Conn, opt.PoolSize),
	}
	p.initConnects()
	return p
}

type pool struct {
	opt   *Options
	conns chan *Conn
}

func (p *pool) initConnects() {
	for i := int64(0); i < p.opt.MinIdleConns; i++ {
		go func() {
			if err := p.addConnect(); err != nil {
				p.opt.Logger.Errorw("pool connect dialer failed",
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
		select {
		case p.conns <- conn:
		default:
		}
	}()

	return fn(ctx, conn)
}

func (p *pool) Get(ctx context.Context) (*Conn, error) {
	var retryTimes = 3
	for {
		select {
		case conn := <-p.conns:
			return conn, nil
		case <-ctx.Done():
			return nil, ErrExceedDeadline
		default:
			if err := p.addConnect(); err != nil {
				retryTimes--
				if retryTimes == 0 {
					if errors.Is(err, ErrMaxPoolSize) {
						return nil, ErrConnExhaust
					}
					return nil, err
				}
				time.Sleep(time.Millisecond * 10)
			}
		}
	}
}

func (p *pool) Put(ctx context.Context, conn *Conn) error {
	select {
	case p.conns <- conn:
		return nil
	default:
		p.opt.Logger.Errorw("pool is full")
		return nil
	}
}

func (p *pool) Destroy(ctx context.Context, conn *Conn) error {
	if err := conn.netConn.Close(); err != nil {
		return err
	}
	conn = nil
	atomic.StoreInt64(&p.opt.PoolSize, 1)
	return nil
}

func (p *pool) addConnect() error {
	size := atomic.LoadInt64(&p.opt.PoolSize)
	if size <= 0 {
		return ErrMaxPoolSize
	}
	conn, err := p.opt.Dialer(context.TODO())
	if err != nil {
		return err
	}
	if err = conn.SetDeadline(time.Now().Add(p.opt.ConnMaxLifetime)); err != nil {
		return err
	}
	select {
	case p.conns <- &Conn{
		netConn:   conn,
		reader:    bufio.NewReader(conn),
		writer:    bufio.NewWriter(conn),
		createdAt: time.Now(),
	}:
		atomic.StoreInt64(&p.opt.PoolSize, -1)
		return nil
	default:
		return nil
	}
}
