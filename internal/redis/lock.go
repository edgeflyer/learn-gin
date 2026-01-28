package redis

import (
	"context"
	"learn-gin/internal/logger"
	"time"
)

func AcquiredLock(ctx context.Context, key string, ttl time.Duration) bool {
	success, err := RDB.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		logger.Log.Error("redis加锁异常")
		return false
	}
	return success
}

func ReleaseLock(ctx context.Context, key string) {
	RDB.Del(ctx, key)
}