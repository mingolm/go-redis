package main

import (
	"context"
	"errors"
	"fmt"
	go_redis "github.com/mingolm/go-redis"
	"github.com/mingolm/go-redis/proto"
	"time"
)

func main() {
	cli := go_redis.NewClient(&go_redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})

	var (
		key = "test_key_2"
		ctx = context.Background()
	)

	err := cli.Set(ctx, key, "123", time.Second).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("val: ", cli.Get(ctx, key).String())

	time.Sleep(time.Second)

	if _, err := cli.Get(ctx, key).Result(); err != nil {
		if errors.Is(err, proto.Nil) {
			fmt.Println("ok")
		} else {
			panic(err)
		}
	}
}
