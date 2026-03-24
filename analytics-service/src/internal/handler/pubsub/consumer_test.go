package pubsub

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/service"
	"github.com/rs/zerolog"
)

func TestNewConsumerGroupHandler(t *testing.T) {
	log := zerolog.Logger{}
	svc := &service.Service{}
	var producer sarama.SyncProducer = nil

	handler := NewConsumerGroupHandler(log, svc, producer)

	if handler == nil {
		t.Error("expected non-nil handler")
	}
	if handler.service != svc {
		t.Error("expected service to match")
	}
	if handler.producer != producer {
		t.Error("expected producer to match")
	}
}

func TestConsumerGroupHandler_Setup(t *testing.T) {
	handler := &ConsumerGroupHandler{}

	err := handler.Setup(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConsumerGroupHandler_Cleanup(t *testing.T) {
	handler := &ConsumerGroupHandler{}

	err := handler.Cleanup(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExtractReqIDFromMessage(t *testing.T) {
	tests := []struct {
		name     string
		headers  []*sarama.RecordHeader
		expected string
	}{
		{
			name: "with req_id header",
			headers: []*sarama.RecordHeader{
				{Key: []byte("req_id"), Value: []byte("test-req-123")},
			},
			expected: "test-req-123",
		},
		{
			name: "with x-request-id header",
			headers: []*sarama.RecordHeader{
				{Key: []byte("x-request-id"), Value: []byte("test-req-456")},
			},
			expected: "test-req-456",
		},
		{
			name:     "no headers",
			headers:  nil,
			expected: "",
		},
		{
			name: "no matching header",
			headers: []*sarama.RecordHeader{
				{Key: []byte("other_key"), Value: []byte("other-value")},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &sarama.ConsumerMessage{
				Headers: tt.headers,
			}
			result := extractReqIDFromMessage(msg)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
