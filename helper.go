package go_redis

import (
	"context"
	"time"
)

func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

func formatMs(ctx context.Context, dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}

func formatSec(ctx context.Context, dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		return 1
	}
	return int64(dur / time.Second)
}
