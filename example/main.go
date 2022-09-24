package main

import (
	"context"
	"fmt"
	go_redis "github.com/mingolm/go-redis"
	"time"
)

func main() {
	cli := go_redis.NewClient(&go_redis.Options{
		Addr:     "",
		Password: "",
	})

	var (
		key = "test_key"
		ctx = context.Background()
	)

	err := cli.Set(ctx, key, "123", time.Minute).Err()
	if err != nil {
		panic(err)
	}

	val, err := cli.Get(ctx, key).String()
	if err != nil {
		panic(err)
	}

	fmt.Println("val: ", val)
}
