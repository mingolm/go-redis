package pool

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net"
	"runtime"
	"time"
)

type Options struct {
	Dialer          func(context.Context) (net.Conn, error) // 拨号
	PoolSize        int64                                   // 连接池长度
	MinIdleConns    int64                                   // 最小空闲连接
	MaxIdleConns    int64                                   // 最大空闲连接
	ConnMaxIdleTime time.Duration                           // 连接超时时间
	ConnMaxLifetime time.Duration                           // 连接最大生命时间
	Logger          *zap.SugaredLogger
}

func (opt *Options) init() error {
	if opt.Dialer == nil {
		return errors.New("pool dialer nil")
	}
	if opt.PoolSize == 0 {
		opt.PoolSize = 10 * int64(runtime.GOMAXPROCS(0))
	}
	if opt.MinIdleConns <= 0 {
		opt.MaxIdleConns = opt.PoolSize >> 4
	} else if opt.MinIdleConns > opt.PoolSize {
		opt.MinIdleConns = opt.PoolSize
	}
	if opt.MaxIdleConns <= 0 || opt.MaxIdleConns > opt.PoolSize {
		opt.MaxIdleConns = opt.PoolSize
	}
	if opt.ConnMaxIdleTime == 0 {
		opt.ConnMaxIdleTime = time.Second
	}
	if opt.ConnMaxLifetime == 0 {
		opt.ConnMaxLifetime = time.Hour
	}
	return nil
}
