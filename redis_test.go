package go_redis

import (
	"context"
	"errors"
	"github.com/mingolm/go-redis/proto"
	"strconv"
	"testing"
	"time"
)

var redis *Redis

func init() {
	redis = NewClient(&Options{
		Addr: "127.0.0.1:6379",
	})
}

func BenchmarkSet(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, err := redis.Set(ctx, strconv.Itoa(i), i, time.Second).Result()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, err := redis.Get(ctx, strconv.Itoa(i)).Result()
		if err != nil {
			if errors.Is(err, proto.Nil) {
				continue
			}
			b.Fatal(err)
		}
	}
}
