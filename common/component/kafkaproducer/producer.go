package kafkaproducer

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
)

type Config struct {
	Brokers  []string      `yaml:"brokers"`
	RetryMax int           `yaml:"retry_max"`
	Timeout  time.Duration `yaml:"timeout"`
}

type KafkaProducerComponent struct {
	log      zerolog.Logger
	cfg      Config
	producer sarama.SyncProducer
	ready    chan struct{}
}

func NewKafkaProducerComponent(log zerolog.Logger, cfg Config) *KafkaProducerComponent {
	return &KafkaProducerComponent{
		log:   log,
		cfg:   cfg,
		ready: make(chan struct{}),
	}
}

func (k *KafkaProducerComponent) Start(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = k.cfg.RetryMax
	config.Net.DialTimeout = k.cfg.Timeout

	producer, err := sarama.NewSyncProducer(k.cfg.Brokers, config)
	if err != nil {
		k.log.Error().Err(err).Msg("Failed to create Kafka producer")
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	k.producer = producer

	close(k.ready)
	k.log.Debug().Strs("brokers", k.cfg.Brokers).Msg("Kafka producer started")
	<-ctx.Done()
	k.log.Debug().Msg("Kafka producer context cancelled – stopping")

	return nil
}

func (k *KafkaProducerComponent) Stop(ctx context.Context) error {
	if k.producer != nil {
		if err := k.producer.Close(); err != nil {
			return fmt.Errorf("close Kafka producer: %w", err)
		}
	}

	k.log.Debug().Msg("Kafka producer stopped")

	return nil
}

func (k *KafkaProducerComponent) Ready() <-chan struct{} {
	return k.ready
}

func (k *KafkaProducerComponent) Producer() sarama.SyncProducer {
	return k.producer
}

// SendMessage sends a message to the specified topic.
// Safe to call only after the component is ready.
func (k *KafkaProducerComponent) SendMessage(topic string, key, value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	return k.producer.SendMessage(msg)
}
