package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"
	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"

	"github.com/rs/zerolog"
)

func (r *analyticsRepository) updateOrderCache(ctx context.Context, date time.Time, analytics *dto.OrderAnalytics) error {
	key := fmt.Sprintf("analytics:daily:%s", date.Format("2006-01-02"))

	cacheData := map[string]interface{}{
		"total_orders":        analytics.Metrics.TotalOrders,
		"total_revenue":       analytics.Metrics.TotalRevenue,
		"average_order_value": analytics.Metrics.AverageOrderVal,
		"cancelled_orders":    analytics.Metrics.CancelledOrders,
	}

	jsonData, _ := json.Marshal(cacheData)
	if err := r.redis0.Set(ctx, key, jsonData, 24*time.Hour).Err(); err != nil {
		return x.WrapWithCode(err, x.CodeCacheSetSimpleKey, "Failed update order cache")
	}

	zerolog.Ctx(ctx).Info().
		Str("date", date.Format("2006-01-02")).
		Int("orders", analytics.Metrics.TotalOrders).
		Float64("revenue", analytics.Metrics.TotalRevenue).
		Msg("Updated Redis cache")

	return nil
}
