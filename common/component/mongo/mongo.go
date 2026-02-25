package mongo

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type MongoDBComponent struct {
	log    zerolog.Logger
	cfg    Config
	ready  chan struct{}
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDBComponent creates a new MongoDB component.
func NewMongoDBComponent(log zerolog.Logger, cfg Config) *MongoDBComponent {
	return &MongoDBComponent{
		log:   log,
		cfg:   cfg,
		ready: make(chan struct{}),
	}
}

// Start establishes the MongoDB connection, pings the server, and stores the client and database.
// It blocks until the context is cancelled.
func (m *MongoDBComponent) Start(ctx context.Context) error {
	uri := fmt.Sprintf("mongodb://%s:%s", m.cfg.Host, m.cfg.Port)
	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return fmt.Errorf("mongo connect: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("mongo ping: %w", err)
	}

	m.client = client
	m.db = client.Database(m.cfg.Database)

	close(m.ready) // signal readiness
	m.log.Debug().Str("uri", uri).Str("database", m.cfg.Database).Msg("MongoDB connected")
	<-ctx.Done() // Block until shutdown signal
	m.log.Debug().Msg("MongoDB component context cancelled – stopping")

	return nil
}

// Stop disconnects the MongoDB client.
func (m *MongoDBComponent) Stop(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	if err := m.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("mongo disconnect: %w", err)
	}

	m.log.Debug().Msg("MongoDB disconnected")

	return nil
}

// Database returns the MongoDB database handle for use by other components.
// Safe to call only after Start has completed.
func (m *MongoDBComponent) Database() *mongo.Database {
	return m.db
}

// Client returns the underlying mongo.Client.
func (m *MongoDBComponent) Client() *mongo.Client {
	return m.client
}

// Ready returns a channel that is closed when the connection is established.
func (m *MongoDBComponent) Ready() <-chan struct{} {
	return m.ready
}
