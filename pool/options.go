package pool

import (
	"context"
	"net"
	"time"
)

type Options struct {
	Dialer          func(context.Context) (net.Conn, error) // 拨号
	PoolSize        int                                     // 连接池长度
	MinIdleConns    int                                     // 最小空闲连接
	MaxIdleConns    int                                     // 最大空闲连接
	ConnMaxIdleTime time.Duration                           // 连接超时时间
	ConnMaxLifetime time.Duration                           // 连接最大生命时间
}
