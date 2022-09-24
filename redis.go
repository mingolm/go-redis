package go_redis

import (
	"bufio"
	"context"
	"github.com/mingolm/go-redis/pool"
	"github.com/mingolm/go-redis/proto"
	"go.uber.org/zap"
)

func NewClient(opt *Options) *Redis {
	// options 初始化
	opt.init()

	r := &Redis{
		opt: opt,
		connPool: pool.NewPool(&pool.Options{
			Dialer:          opt.Dialer,
			PoolSize:        opt.PoolSize,
			MinIdleConns:    opt.MinIdleConns,
			MaxIdleConns:    opt.MaxIdleConns,
			ConnMaxIdleTime: opt.ConnMaxIdleTime,
			ConnMaxLifetime: opt.ConnMaxLifeTime,
			Logger:          opt.Logger,
		}),
	}
	r.cmdable = r.process

	return r
}

type Redis struct {
	cmdable
	opt      *Options
	connPool pool.Pooler
	logger   *zap.SugaredLogger
}

func (r *Redis) process(ctx context.Context, cmd Cmder) error {
	err := r.connPool.WithConn(ctx, func(ctx context.Context, cn *pool.Conn) error {
		if err := cn.WithWrite(ctx, func(ctx context.Context, wd *bufio.Writer) error {
			return proto.NewWriter(wd).Write(ctx, cmd.Args())
		}); err != nil {
			return err
		}

		if err := cn.WithRead(ctx, func(ctx context.Context, rd *bufio.Reader) error {
			val, err := proto.NewReader(rd).Read()
			if err != nil {
				return err
			}
			return cmd.ReadReply(val)
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
