package pubsub

import (
	"context"
	"encoding/json"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/service"

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
		h.log.Info().Str("topic", msg.Topic).Bytes("value", msg.Value).Msg("received message")
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

func (h *ConsumerGroupHandler) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var event dto.OrderEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	switch event.EventType {
	case "order.created":
		err := h.service.Analytics.UpdateOrderAnalytics(ctx, event)
		if err != nil {
			return err
		}

		err = h.service.Analytics.UpdateProductAnalytics(ctx, event)
		if err != nil {
			return err
		}
	case "order.cancelled":
		err := h.service.Analytics.UpdateCancellationMetrics(ctx, event)
		if err != nil {
			return err
		}
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
