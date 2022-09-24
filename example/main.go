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
	redis := go_redis.NewClient(&go_redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})

	var (
		key = "name"
		ctx = context.Background()
	)

	err := redis.Set(ctx, key, "mingo", time.Second).Err()
	if err != nil {
		panic(err)
	}

	val, err := redis.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	if val != "mingo" {
		panic(errors.New("redis internal error"))
	}

	time.Sleep(time.Second)

	// test key expire
	if err = redis.Get(ctx, key).Err(); err != proto.Nil {
		panic(errors.New("redis internal error"))
	}

	fmt.Println("ok")
}
