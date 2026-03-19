package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/service"
	"github.com/linggaaskaedo/go-kill/common/pkg/metrics"

	"github.com/IBM/sarama"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

const (
	MaxRetries      = 3
	RetryBaseDelay  = 1 * time.Second
	DeadLetterTopic = "analytics-service-dlq"
)

type ConsumerGroupHandler struct {
	log      zerolog.Logger
	service  *service.Service
	producer sarama.SyncProducer
}

func NewConsumerGroupHandler(log zerolog.Logger, service *service.Service, producer sarama.SyncProducer) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		log:      log,
		service:  service,
		producer: producer,
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
		startTime := time.Now()
		reqID := extractReqIDFromMessage(msg)
		if reqID == "" {
			reqID = xid.New().String()
		}

		ctx := context.Background()
		logWithReq := h.log.With().Str("req_id", reqID).Logger()
		ctx = logWithReq.WithContext(ctx)

		metrics.MessagesReceived.WithLabelValues(msg.Topic).Inc()
		h.log.Info().Str("topic", msg.Topic).Int32("partition", msg.Partition).Int64("offset", msg.Offset).Msg("received message")

		var lastErr error
		for attempt := 0; attempt <= MaxRetries; attempt++ {
			if attempt > 0 {
				metrics.RetryAttempts.WithLabelValues(msg.Topic).Inc()
				delay := RetryBaseDelay * time.Duration(1<<uint(attempt-1))
				select {
				case <-time.After(delay):
				case <-ctx.Done():
					return ctx.Err()
				}
				zerolog.Ctx(ctx).Info().Int("attempt", attempt).Msg("retrying message processing")
			}

			if err := h.processMessage(ctx, msg); err != nil {
				lastErr = err
				zerolog.Ctx(ctx).Error().Err(err).Int("attempt", attempt+1).Msg("failed to process message")
				continue
			}

			metrics.MessagesProcessed.WithLabelValues(msg.Topic, "success").Inc()
			metrics.MessageProcessingDuration.WithLabelValues(msg.Topic).Observe(time.Since(startTime).Seconds())
			zerolog.Ctx(ctx).Info().Msg("message processed successfully")
			sess.MarkMessage(msg, "")
			break
		}

		if lastErr != nil {
			metrics.MessagesProcessed.WithLabelValues(msg.Topic, "failure").Inc()
			metrics.MessageProcessingDuration.WithLabelValues(msg.Topic).Observe(time.Since(startTime).Seconds())
			h.sendToDLQ(ctx, reqID, msg, lastErr)
			sess.MarkMessage(msg, "")
		}
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

func (h *ConsumerGroupHandler) sendToDLQ(ctx context.Context, reqID string, msg *sarama.ConsumerMessage, originalErr error) {
	dlqMessage := map[string]interface{}{
		"original_topic":     msg.Topic,
		"original_partition": msg.Partition,
		"original_offset":    msg.Offset,
		"timestamp":          time.Now().UTC(),
		"error":              originalErr.Error(),
		"payload":            json.RawMessage(msg.Value),
	}

	dlqData, err := json.Marshal(dlqMessage)
	if err != nil {
		metrics.KafkaProducerErrors.Inc()
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to marshal DLQ message")
		return
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: DeadLetterTopic,
		Key:   sarama.StringEncoder(msg.Key),
		Value: sarama.ByteEncoder(dlqData),
		Headers: []sarama.RecordHeader{
			{Key: []byte("req_id"), Value: []byte(reqID)},
			{Key: []byte("original_error"), Value: []byte(originalErr.Error())},
			{Key: []byte("retry_count"), Value: []byte(fmt.Sprintf("%d", MaxRetries))},
		},
	}

	_, _, err = h.producer.SendMessage(producerMsg)
	if err != nil {
		metrics.KafkaProducerErrors.Inc()
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to send message to DLQ")
		return
	}

	metrics.DLQMessagesSent.WithLabelValues(msg.Topic).Inc()
	zerolog.Ctx(ctx).Warn().
		Err(originalErr).
		Str("topic", DeadLetterTopic).
		Msg("Message sent to DLQ after exhausting retries")
}
