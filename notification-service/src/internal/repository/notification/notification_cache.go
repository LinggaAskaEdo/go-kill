package notification

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (r *notificationRepository) checkRateLimitCache(ctx context.Context, userID string) bool {
	key := fmt.Sprintf("rate_limit:%s:order_notifications", userID)

	count, err := r.redis0.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Redis error: %v", err)
		return true // Allow on error
	}

	if count == 1 {
		r.redis0.Expire(ctx, key, time.Hour)
	}

	// Max 10 notifications per hour
	return count <= 10
}
