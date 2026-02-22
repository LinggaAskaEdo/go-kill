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
	Host     string `yaml:"host"`     // e.g., "localhost"
	Port     string `yaml:"port"`     // e.g., "27017"
	Database string `yaml:"database"` // database name, e.g., "user_db"
	// You may add more options like auth, TLS, etc.
}

type MongoDBComponent struct {
	log    zerolog.Logger
	cfg    Config
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDBComponent creates a new MongoDB component.
func NewMongoDBComponent(log zerolog.Logger, cfg Config) *MongoDBComponent {
	return &MongoDBComponent{
		log: log,
		cfg: cfg,
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

	m.log.Debug().Str("uri", uri).Str("database", m.cfg.Database).Msg("MongoDB connected")

	// Block until shutdown
	<-ctx.Done()

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
