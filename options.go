package go_redis

import (
	"context"
	"crypto/tls"
	"go.uber.org/zap"
	"net"
	"runtime"
	"time"
)

type Options struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string

	// Dialer creates new network connection and has priority over
	// Network and Addr options.
	Dialer func(ctx context.Context) (net.Conn, error)

	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string

	// Database to be selected after connecting to the server.
	DB int

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int32
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int32
	// Maximum number of idle connections.
	MaxIdleConns int32
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	ConnMaxLifeTime time.Duration

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config
	// zap logger
	Logger *zap.SugaredLogger
}

func (opt *Options) init() {
	var (
		network         = "tcp"
		addr            = "localhost:6379"
		poolSize        = runtime.GOMAXPROCS(0) * 10
		dialTimeout     = time.Second * 5
		readTimeout     = time.Second * 3
		writeTimeout    = time.Second * 3
		minIdleConnectx = poolSize >> 4
		maxIdleConnectx = poolSize >> 2
		connMaxIdleTime = time.Minute * 30
		connMaxLifeTime = time.Hour
		maxRetries      = 3
		logger          = zap.S()
	)

	if opt.Network == "" {
		opt.Network = network
	}
	if opt.Addr == "" {
		opt.Addr = addr
	}
	if opt.PoolSize == 0 {
		opt.PoolSize = int32(poolSize)
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = dialTimeout
	}
	if opt.ReadTimeout == 0 {
		opt.ReadTimeout = readTimeout
	}
	if opt.WriteTimeout == 0 {
		opt.WriteTimeout = writeTimeout
	}
	if opt.MinIdleConns == 0 {
		opt.MinIdleConns = int32(minIdleConnectx)
	}
	if opt.MaxIdleConns == 0 {
		opt.MaxIdleConns = int32(maxIdleConnectx)
	}
	if opt.ConnMaxIdleTime == 0 {
		opt.ConnMaxIdleTime = connMaxIdleTime
	}
	if opt.ConnMaxLifeTime == 0 {
		opt.ConnMaxLifeTime = connMaxLifeTime
	}
	if opt.MaxRetries == 0 {
		opt.MaxRetries = maxRetries
	}
	if opt.Logger == nil {
		opt.Logger = logger
	}
	if opt.Dialer == nil {
		opt.Dialer = func(ctx context.Context) (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   opt.DialTimeout,
				KeepAlive: 5 * time.Minute,
			}
			if opt.TLSConfig == nil {
				return netDialer.DialContext(ctx, opt.Network, opt.Addr)
			}
			return tls.DialWithDialer(netDialer, opt.Network, opt.Addr, opt.TLSConfig)
		}
	}
}
