package go_redis

import (
	"context"
	"time"
)

type Cmdable interface {
	Set(ctx context.Context, key string, val any, expiration time.Duration) *StatusCmd
	Get(ctx context.Context, key string) *StringCmd
}

type cmdable func(ctx context.Context, cmd Cmder) error

func (c cmdable) Set(ctx context.Context, key string, val any, expiration time.Duration) *StatusCmd {
	args := []interface{}{"SET", key}
	if expiration > 0 {
		args = append(args, expiration.Microseconds())

	}
	cmd := &StatusCmd{
		baseCmd: &baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
	cmd.err = c(ctx, cmd)
	return cmd
}

func (c cmdable) Get(ctx context.Context, key string) *StringCmd {
	cmd := &StringCmd{
		baseCmd: &baseCmd{
			ctx:  ctx,
			args: []interface{}{"GET", key},
		},
	}
	cmd.err = c(ctx, cmd)
	return cmd
}

type Cmder interface {
	Err() error
	Args() []interface{}
	ReadReply(interface{}) error
	String() (string, error)
}

type baseCmd struct {
	ctx  context.Context
	args []interface{}
	err  error
}

func (c *baseCmd) Err() error {
	return nil
}

func (c *baseCmd) Args() []interface{} {
	return nil
}

type StatusCmd struct {
	*baseCmd
	result string
}

func (cmd *StatusCmd) ReadReply(val interface{}) error {
	cmd.result = val.(string)
	return nil
}

func (cmd *StatusCmd) String() (string, error) {
	return cmd.result, nil
}

type StringCmd struct {
	*baseCmd
	result string
}

func (cmd *StringCmd) ReadReply(val interface{}) error {
	cmd.result = val.(string)
	return nil
}

func (cmd *StringCmd) String() (string, error) {
	return cmd.result, nil
}
