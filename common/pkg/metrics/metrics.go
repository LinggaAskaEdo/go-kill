package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MessagesReceived = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_messages_received_total",
			Help: "Total number of messages received from Kafka",
		},
		[]string{"topic"},
	)

	MessagesProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_messages_processed_total",
			Help: "Total number of messages successfully processed",
		},
		[]string{"topic", "status"},
	)

	MessageProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "analytics_message_processing_duration_seconds",
			Help:    "Histogram of message processing duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)

	DLQMessagesSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_dlq_messages_total",
			Help: "Total number of messages sent to DLQ",
		},
		[]string{"topic"},
	)

	RetryAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_retry_attempts_total",
			Help: "Total number of retry attempts",
		},
		[]string{"topic"},
	)

	MongoOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_mongo_operations_total",
			Help: "Total number of MongoDB operations",
		},
		[]string{"collection", "operation", "status"},
	)

	MongoOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "analytics_mongo_operation_duration_seconds",
			Help:    "Histogram of MongoDB operation duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"collection", "operation"},
	)

	RedisOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "analytics_redis_operations_total",
			Help: "Total number of Redis operations",
		},
		[]string{"operation", "status"},
	)

	KafkaProducerErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "analytics_kafka_producer_errors_total",
			Help: "Total number of Kafka producer errors",
		},
	)
)
