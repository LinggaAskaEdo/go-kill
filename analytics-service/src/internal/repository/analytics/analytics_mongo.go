package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *analyticsRepository) updateOrderMongo(ctx context.Context, date time.Time, event dto.OrderEvent) error {
	hour := event.Timestamp.Hour()
	collection := r.mongo0.Collection(r.analyticsOptions.OrderCollection)

	// Update daily metrics
	filter := bson.M{"date": date}
	update := bson.M{
		"$inc": bson.M{
			"metrics.total_orders":                           1,
			"metrics.total_revenue":                          event.Data.TotalAmount,
			fmt.Sprintf("hourly_breakdown.%d.orders", hour):  1,
			fmt.Sprintf("hourly_breakdown.%d.revenue", hour): event.Data.TotalAmount,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to update order analytics")
		return err
	}

	// If it was an insert, initialize hourly breakdown
	if result.UpsertedCount > 0 {
		r.initializeHourlyBreakdown(ctx, collection, date)
	}

	// Recalculate average order value
	r.recalculateAverageOrderValue(ctx, collection, date)

	return nil
}

func (r *analyticsRepository) initializeHourlyBreakdown(ctx context.Context, collection *mongo.Collection, date time.Time) {
	hourlyData := make([]dto.HourlyMetrics, 24)
	for i := range 24 {
		hourlyData[i] = dto.HourlyMetrics{
			Hour:    i,
			Orders:  0,
			Revenue: 0,
		}
	}

	filter := bson.M{"date": date}
	update := bson.M{
		"$set": bson.M{
			"hourly_breakdown": hourlyData,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to initialize hourly breakdown")
	}
}

func (r *analyticsRepository) recalculateAverageOrderValue(ctx context.Context, collection *mongo.Collection, date time.Time) {
	var analytics dto.OrderAnalytics

	err := collection.FindOne(ctx, bson.M{"date": date}).Decode(&analytics)
	if err != nil {
		return
	}

	if analytics.Metrics.TotalOrders > 0 {
		avgOrderValue := analytics.Metrics.TotalRevenue / float64(analytics.Metrics.TotalOrders)

		_, err := collection.UpdateOne(
			ctx,
			bson.M{"date": date},
			bson.M{
				"$set": bson.M{
					"metrics.average_order_value": avgOrderValue,
				},
			},
		)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to recalculate average order value")
		}
	}
}

func (r *analyticsRepository) getOrderMongo(ctx context.Context, date time.Time) (*dto.OrderAnalytics, error) {
	var analytics dto.OrderAnalytics

	err := r.mongo0.Collection(r.analyticsOptions.OrderCollection).FindOne(ctx, bson.M{"date": date}).Decode(&analytics)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to fetch analytics")
		return nil, err
	}

	return &analytics, nil
}

func (r *analyticsRepository) updateProductMongo(ctx context.Context, date time.Time, event dto.OrderEvent) error {
	collection := r.mongo0.Collection(r.analyticsOptions.ProductCollection)

	for _, item := range event.Data.Items {
		filter := bson.M{
			"product_id": item.ProductID,
			"date":       date,
		}

		update := bson.M{
			"$inc": bson.M{
				"sales_count": item.Quantity,
				"revenue":     item.UnitPrice * float64(item.Quantity),
			},
			"$set": bson.M{
				"product_name": item.ProductName,
				"updated_at":   time.Now(),
			},
		}

		_, err := collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Str("productID", item.ProductID).Msg("Failed to update product analytics")
			return err
		}
	}

	return nil
}

func (r *analyticsRepository) updateCancellationMongo(ctx context.Context, date time.Time) error {
	collection := r.mongo0.Collection(r.analyticsOptions.OrderCollection)

	filter := bson.M{"date": date}
	update := bson.M{
		"$inc": bson.M{
			"metrics.cancelled_orders": 1,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to update cancellation")
		return err
	}

	return nil
}
