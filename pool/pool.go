package pool

import "context"

type Pooler interface {
	WithConn(context.Context, func(context.Context, *Conn) error) error
	// Get 获取连接
	Get(context.Context) (*Conn, error)
	Put(context.Context, *Conn) error
	Destroy(context.Context, *Conn) error
}

func NewPool(opt *Options) Pooler {
	var (
		poolSize = opt.PoolSize
	)

	return &pool{
		conns: make(chan *Conn, poolSize),
	}
}

type pool struct {
	conns chan *Conn
}

func (p *pool) WithConn(context.Context, func(context.Context, *Conn) error) error {
	return nil
}

func (p *pool) Get(context.Context) (*Conn, error) {
	return nil, nil
}

func (p *pool) Put(context.Context, *Conn) error {
	return nil
}

func (p *pool) Destroy(context.Context, *Conn) error {
	return nil
}
