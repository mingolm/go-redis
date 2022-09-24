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
	args := make([]interface{}, 3, 5)
	args[0] = "SET"
	args[1] = key
	args[2] = val
	if expiration > 0 {
		if usePrecise(expiration) {
			args = append(args, "px", formatMs(ctx, expiration))
		} else {
			args = append(args, "ex", formatSec(ctx, expiration))
		}
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
	String() string
}

type baseCmd struct {
	ctx  context.Context
	args []interface{}
	err  error
}

func (c *baseCmd) Err() error {
	return c.err
}

func (c *baseCmd) Args() []interface{} {
	return c.args
}

type StatusCmd struct {
	*baseCmd
	result string
}

func (cmd *StatusCmd) ReadReply(val interface{}) error {
	cmd.result = val.(string)
	return nil
}

func (cmd *StatusCmd) String() string {
	return cmd.result
}

func (cmd *StatusCmd) Result() (string, error) {
	return cmd.result, cmd.err
}

type StringCmd struct {
	*baseCmd
	result string
}

func (cmd *StringCmd) ReadReply(val interface{}) error {
	cmd.result = val.(string)
	return nil
}

func (cmd *StringCmd) String() string {
	return cmd.result
}

func (cmd *StringCmd) Result() (string, error) {
	return cmd.result, cmd.err
}
