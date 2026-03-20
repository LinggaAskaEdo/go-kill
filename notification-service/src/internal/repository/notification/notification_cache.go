package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

func (r *notificationRepository) checkRateLimitCache(ctx context.Context, userID string) bool {
	key := fmt.Sprintf("rate_limit:%s:order_notifications", userID)

	ttl := r.opts.RateLimit.Window
	if ttl == 0 {
		ttl = time.Hour
	}
	maxPerHour := r.opts.RateLimit.MaxPerHour
	if maxPerHour == 0 {
		maxPerHour = 10
	}

	pipe := r.redis0.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Redis pipeline error")
		return true
	}

	count := incr.Val()
	return count <= int64(maxPerHour)
}
