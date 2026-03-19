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

	filter := bson.M{"date": date}
	update := bson.M{
		"$inc": bson.M{
			"metrics.total_orders":                           1,
			"metrics.total_revenue":                          event.Data.TotalAmount,
			fmt.Sprintf("hourly_breakdown.%d.orders", hour):  1,
			fmt.Sprintf("hourly_breakdown.%d.revenue", hour): event.Data.TotalAmount,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
		"$set": bson.M{
			"updated_at":                  time.Now(),
			"metrics.average_order_value": 0,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to update order analytics")
		return err
	}

	if result.UpsertedCount > 0 {
		if err := r.initializeHourlyBreakdown(ctx, collection, date); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to initialize hourly breakdown")
		}
	}

	return nil
}

func (r *analyticsRepository) initializeHourlyBreakdown(ctx context.Context, collection *mongo.Collection, date time.Time) error {
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
	return err
}

func (r *analyticsRepository) getOrderMongo(ctx context.Context, date time.Time) (*dto.OrderAnalytics, error) {
	var analytics dto.OrderAnalytics

	err := r.mongo0.Collection(r.analyticsOptions.OrderCollection).FindOne(ctx, bson.M{"date": date}).Decode(&analytics)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to fetch analytics")
		return nil, err
	}

	if analytics.Metrics.TotalOrders > 0 {
		analytics.Metrics.AverageOrderVal = analytics.Metrics.TotalRevenue / float64(analytics.Metrics.TotalOrders)
	}

	return &analytics, nil
}

func (r *analyticsRepository) updateProductMongo(ctx context.Context, date time.Time, event dto.OrderEvent) error {
	if len(event.Data.Items) == 0 {
		return nil
	}

	collection := r.mongo0.Collection(r.analyticsOptions.ProductCollection)

	var operations []mongo.WriteModel
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
			"$setOnInsert": bson.M{
				"created_at": time.Now(),
			},
		}

		operations = append(operations, mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true))
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := collection.BulkWrite(ctx, operations, opts)
	if err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to bulk update product analytics")
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
