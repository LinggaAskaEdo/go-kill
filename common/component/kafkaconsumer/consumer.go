package kafkaconsumer

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
)

type Config struct {
	Brokers                        []string      `yaml:"brokers"`
	GroupID                        string        `yaml:"group_id"`
	Topics                         []string      `yaml:"topics"`
	InitialOffset                  int64         `yaml:"initial_offset"`
	DialTimeout                    time.Duration `yaml:"dial_timeout"`
	ReadTimeout                    time.Duration `yaml:"read_timeout"`
	WriteTimeout                   time.Duration `yaml:"write_timeout"`
	ConsumerGroupSessionTimeout    time.Duration `yaml:"consumer_group_session_timeout"`
	ConsumerGroupHeartbeatInterval time.Duration `yaml:"consumer_group_heartbeat_interval"`
}

type KafkaConsumerComponent struct {
	log     zerolog.Logger
	cfg     Config
	handler sarama.ConsumerGroupHandler
	group   sarama.ConsumerGroup
	ready   chan struct{}
	cancel  context.CancelFunc
}

func NewKafkaConsumerComponent(log zerolog.Logger, cfg Config, handler sarama.ConsumerGroupHandler) *KafkaConsumerComponent {
	return &KafkaConsumerComponent{
		log:     log,
		cfg:     cfg,
		handler: handler,
		ready:   make(chan struct{}),
	}
}

func (k *KafkaConsumerComponent) Start(ctx context.Context) error {
	config := sarama.NewConfig()

	// Use GroupStrategies instead of deprecated Strategy
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}

	if k.cfg.InitialOffset == 0 {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = k.cfg.InitialOffset
	}

	config.Net.DialTimeout = k.cfg.DialTimeout
	config.Net.ReadTimeout = k.cfg.ReadTimeout
	config.Net.WriteTimeout = k.cfg.WriteTimeout
	config.Consumer.Group.Session.Timeout = k.cfg.ConsumerGroupSessionTimeout
	config.Consumer.Group.Heartbeat.Interval = k.cfg.ConsumerGroupHeartbeatInterval

	group, err := sarama.NewConsumerGroup(k.cfg.Brokers, k.cfg.GroupID, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	k.group = group

	consumeCtx, cancel := context.WithCancel(context.Background())
	k.cancel = cancel

	// Run the consumer loop in a separate goroutine
	go k.runConsumerLoop(consumeCtx)

	close(k.ready)
	k.log.Debug().Msg("Kafka consumer started")
	<-ctx.Done()
	k.log.Debug().Msg("Kafka consumer context cancelled – stopping")

	return nil
}

func (k *KafkaConsumerComponent) runConsumerLoop(ctx context.Context) {
	const (
		baseDelay = 100 * time.Millisecond
		maxDelay  = 30 * time.Second
	)
	attempt := 0

	k.log.Debug().Strs("brokers", k.cfg.Brokers).Strs("topics", k.cfg.Topics).Msg("Kafka consumer loop starting")

	for {
		// Check for cancellation before attempting to consume
		select {
		case <-ctx.Done():
			k.log.Debug().Msg("Consumer loop stopped due to context cancellation")
			return
		default:
		}

		err := k.group.Consume(ctx, k.cfg.Topics, k.handler)
		if err == nil {
			// Normal exit (e.g., after rebalance) – reset backoff and continue
			attempt = 0
			continue
		}

		// Log the error
		k.log.Error().Err(err).Msg("Kafka consumer error")

		// Check again for cancellation
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Calculate backoff delay
		delay := exponentialBackoff(baseDelay, maxDelay, attempt)
		k.log.Debug().Dur("delay", delay).Msg("Retrying consumer after error")

		// Wait for the delay or cancellation
		select {
		case <-time.After(delay):
			attempt++
		case <-ctx.Done():
			return
		}
	}
}

func exponentialBackoff(base, max time.Duration, attempt int) time.Duration {
	if attempt < 0 {
		return base
	}
	// Use math.Pow for safe exponentiation (avoids integer shift issues)
	factor := math.Pow(2, float64(attempt))
	delay := float64(base) * factor * (0.5 + rand.Float64())
	if delay > float64(max) {
		delay = float64(max)
	}

	return time.Duration(delay)
}

func (k *KafkaConsumerComponent) Stop(ctx context.Context) error {
	if k.cancel != nil {
		k.cancel()
	}

	if k.group != nil {
		if err := k.group.Close(); err != nil {
			return fmt.Errorf("close consumer group: %w", err)
		}
	}

	k.log.Debug().Msg("Kafka consumer stopped")

	return nil
}

func (k *KafkaConsumerComponent) Ready() <-chan struct{} {
	return k.ready
}
