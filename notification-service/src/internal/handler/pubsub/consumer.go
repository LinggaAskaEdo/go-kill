package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
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
		// 1. Extract or generate a request ID (correlation ID)
		reqID := extractReqIDFromMessage(msg) // see helper below
		if reqID == "" {
			reqID = xid.New().String()
		}

		// 2. Create a context with the request ID and a logger
		ctx := context.Background()

		// 3. Create a logger with the request ID
		logWithReq := h.log.With().Str("req_id", reqID).Logger()
		ctx = logWithReq.WithContext(ctx)

		// 4. Process the message using the service (which now can use zerolog.Ctx(ctx))
		zerolog.Ctx(ctx).Info().Str("topic", msg.Topic).Bytes("value", msg.Value).Msg("received message")
		// h.processMessage(ctx, msg)
		if err := h.processMessage(ctx, msg); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to process message")
		} else {
			zerolog.Ctx(ctx).Info().Msg("message processed successfully")
		}

		// 5. Mark the message as consumed
		sess.MarkMessage(msg, "")
	}

	return nil
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

	zerolog.Ctx(ctx).Info().Msg(fmt.Sprintf("Processing %s event for order %s", event.EventType, event.Data.OrderID))

	// Check user preferences (error ignored as in original)
	prefs, _ := h.service.Notification.GetUserPreference(ctx, event.Data.UserID)

	// Check rate limit
	if !h.service.Notification.CheckRateLimit(ctx, event.Data.UserID) {
		errStr := fmt.Sprintf("Rate limit exceeded for user %s", event.Data.UserID)
		zerolog.Ctx(ctx).Error().Msg(errStr)
		return x.New(errStr)
	}

	// Process based on event type – each case now only calls the helper
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
