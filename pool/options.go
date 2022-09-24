package pool

import (
	"context"
	"go.uber.org/zap"
	"net"
	"time"
)

type Options struct {
	Dialer          func(context.Context) (net.Conn, error) // 拨号
	PoolSize        int32                                   // 连接池长度
	MinIdleConns    int32                                   // 最小空闲连接
	MaxIdleConns    int32                                   // 最大空闲连接
	ConnMaxIdleTime time.Duration                           // 连接超时时间
	ConnMaxLifetime time.Duration                           // 连接最大生命周期
	Logger          *zap.SugaredLogger                      // 日志 debug
}
