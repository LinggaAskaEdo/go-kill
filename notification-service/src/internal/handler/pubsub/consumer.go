package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/notification-service/src/internal/service"

	"github.com/IBM/sarama"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

type ConsumerGroupHandler struct {
	log     zerolog.Logger
	service *service.Service
}

func NewConsumerGroupHandler(log zerolog.Logger, service *service.Service) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		log:     log,
		service: service,
	}
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.processMessageWithContext(sess, msg)
	}

	return nil
}

func (h *ConsumerGroupHandler) processMessageWithContext(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	reqID := extractReqIDFromMessage(msg)
	if reqID == "" {
		reqID = xid.New().String()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logWithReq := h.log.With().Str("req_id", reqID).Logger()
	ctx = logWithReq.WithContext(ctx)

	zerolog.Ctx(ctx).Info().Str("topic", msg.Topic).Msg("received message")

	if err := h.processMessage(ctx, msg); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to process message, will retry")
		return
	}

	zerolog.Ctx(ctx).Info().Msg("message processed successfully")
	sess.MarkMessage(msg, "")
}

func extractReqIDFromMessage(msg *sarama.ConsumerMessage) string {
	for _, header := range msg.Headers {
		if string(header.Key) == "req_id" || string(header.Key) == "x-request-id" {
			return string(header.Value)
		}
	}

	return ""
}

func (h *ConsumerGroupHandler) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event dto.OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed parse json")
		return err
	}

	zerolog.Ctx(ctx).Info().Str("event_type", event.EventType).Str("order_id", event.Data.OrderID).Msg("processing event")

	prefs, err := h.service.Notification.GetUserPreference(ctx, event.Data.UserID)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Str("user_id", event.Data.UserID).Msg("failed to get user preferences, using defaults")
		prefs = dto.NotificationPreferences{
			EmailEnabled: true,
			PushEnabled:  true,
		}
	}

	if !h.service.Notification.CheckRateLimit(ctx, event.Data.UserID) {
		errStr := "rate limit exceeded for user " + event.Data.UserID
		zerolog.Ctx(ctx).Warn().Msg(errStr)
		return errors.New(errStr)
	}

	switch event.EventType {
	case "order.created":
		h.sendNotificationIfEnabled(ctx, prefs.EmailEnabled,
			func() error { return h.service.Notification.SendOrderConfirmation(ctx, event) }, "failed to send order confirmation")
	case "order.updated":
		h.sendNotificationIfEnabled(ctx, prefs.PushEnabled,
			func() error { return h.service.Notification.SendOrderUpdate(ctx, event) }, "failed to send order update")
	case "order.cancelled":
		h.sendNotificationIfEnabled(ctx, prefs.EmailEnabled,
			func() error { return h.service.Notification.SendOrderCancellation(ctx, event) }, "failed to send order cancellation")
	default:
		zerolog.Ctx(ctx).Warn().Str("event_type", event.EventType).Msg("unknown event type, skipping")
	}

	return nil
}

func (h *ConsumerGroupHandler) sendNotificationIfEnabled(ctx context.Context, enabled bool, sendFunc func() error, logMsg string) {
	if !enabled {
		return
	}

	if err := sendFunc(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(logMsg)
	}
}
