package ratelimit

import (
	"time"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/queue"
)

func IsGlobalRateLimitExceeded() bool {

	key := "ratelimit:global"
	window := 1 * time.Second

	count, err := queue.RDB.Incr(queue.Ctx, key).Result()

	if err != nil {
		return false
	}

	if count == 1 {
		queue.RDB.Expire(queue.Ctx, key, window)
	}

	return count > 5000
}
