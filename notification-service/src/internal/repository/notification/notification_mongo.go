package notification

import (
	"context"
	"fmt"
	"time"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/model/dto"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *notificationRepository) getUserPreferenceMongo(ctx context.Context, userID string) (dto.NotificationPreferences, error) {
	var prefs dto.NotificationPreferences

	err := r.mongo0.Collection(r.opts.NotificationPreferences).FindOne(ctx, bson.M{"user_id": userID}).Decode(&prefs)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to get user preferences, use default")

		prefs = dto.NotificationPreferences{
			UserID:       userID,
			EmailEnabled: true,
			SMSEnabled:   false,
			PushEnabled:  true,
		}
	}

	return prefs, nil
}

func (r *notificationRepository) sendOrderConfirmationMongo(ctx context.Context, event dto.OrderEvent) error {
	var template struct {
		Subject string `bson:"subject"`
		Body    string `bson:"body"`
	}

	if err := r.mongo0.Collection(r.opts.NotificationTemplates).FindOne(
		ctx,
		bson.M{"template_id": "order_confirmation_v1", "type": "email", "active": true},
	).Decode(&template); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Template not found")

		template.Subject = "Order Confirmation - {{order_number}}"
		template.Body = "Your order {{order_number}} has been confirmed. Total: ${{total_amount}}"
	}

	// Replace placeholders
	subject := replaceTemplate(template.Subject, map[string]string{
		"order_number": event.Data.OrderNumber,
	})
	body := replaceTemplate(template.Body, map[string]string{
		"order_number": event.Data.OrderNumber,
		"total_amount": fmt.Sprintf("%.2f", event.Data.TotalAmount),
	})

	// Simulate sending email
	zerolog.Ctx(ctx).Debug().Msg(fmt.Sprintf("Sending email to %s: %s", event.Data.UserEmail, subject))

	// Store notification in MongoDB
	notification := dto.Notification{
		UserID:   event.Data.UserID,
		Type:     "email",
		Category: "order",
		Title:    subject,
		Message:  body,
		Metadata: map[string]interface{}{
			"order_id":     event.Data.OrderID,
			"template_id":  "order_confirmation_v1",
			"sent_to":      event.Data.UserEmail,
			"order_number": event.Data.OrderNumber,
			"total_amount": event.Data.TotalAmount,
		},
		Status:    "sent",
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}

	_, err := r.mongo0.Collection(r.opts.Notifications).InsertOne(ctx, notification)
	if err != nil {
		errStr := "Failed when sendOrderConfirmationMongo"
		zerolog.Ctx(ctx).Error().Err(err).Msg(errStr)
		return x.Wrap(err, errStr)
	}

	return nil
}

func (r *notificationRepository) sendOrderUpdateMongo(ctx context.Context, event dto.OrderEvent) error {
	zerolog.Ctx(ctx).Debug().Msg(fmt.Sprintf("Sending push notification for order %s update", event.Data.OrderNumber))

	notification := dto.Notification{
		UserID:   event.Data.UserID,
		Type:     "push",
		Category: "order",
		Title:    "Order Update",
		Message:  fmt.Sprintf("Your order %s has been updated to status: %s", event.Data.OrderNumber, event.Data.Status),
		Metadata: map[string]interface{}{
			"order_id":     event.Data.OrderID,
			"order_number": event.Data.OrderNumber,
			"new_status":   event.Data.Status,
		},
		Status:    "sent",
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}

	_, err := r.mongo0.Collection(r.opts.Notifications).InsertOne(ctx, notification)
	if err != nil {
		errStr := "Failed when sendOrderUpdateMongo"
		zerolog.Ctx(ctx).Error().Err(err).Msg(errStr)
		return x.Wrap(err, errStr)
	}

	return nil
}

func (r *notificationRepository) sendOrderCancellationMongo(ctx context.Context, event dto.OrderEvent) error {
	zerolog.Ctx(ctx).Debug().Msg(fmt.Sprintf("Sending cancellation email for order %s", event.Data.OrderNumber))

	notification := dto.Notification{
		UserID:   event.Data.UserID,
		Type:     "email",
		Category: "order",
		Title:    "Order Cancelled",
		Message:  fmt.Sprintf("Your order %s has been cancelled", event.Data.OrderNumber),
		Metadata: map[string]interface{}{
			"order_id":     event.Data.OrderID,
			"order_number": event.Data.OrderNumber,
		},
		Status:    "sent",
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}

	_, err := r.mongo0.Collection(r.opts.Notifications).InsertOne(ctx, notification)
	if err != nil {
		errStr := "Failed when sendOrderCancellationMongo"
		zerolog.Ctx(ctx).Error().Err(err).Msg(errStr)
		return x.Wrap(err, errStr)
	}

	return nil
}
