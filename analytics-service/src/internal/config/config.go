package config

import (
	"os"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/common/component/http"
	"github.com/linggaaskaedo/go-kill/common/component/kafkaconsumer"
	"github.com/linggaaskaedo/go-kill/common/component/kafkaproducer"
	"github.com/linggaaskaedo/go-kill/common/component/mongo"
	"github.com/linggaaskaedo/go-kill/common/component/redis"
	"github.com/linggaaskaedo/go-kill/common/component/server"
	"github.com/linggaaskaedo/go-kill/common/pkg/logger"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Logger        logger.Config           `yaml:"logger"`
	Redis         redis.Config            `yaml:"redis"`
	Mongo         map[string]mongo.Config `yaml:"mongo"`
	Http          http.Config             `yaml:"http"`
	Server        server.Config           `yaml:"server"`
	KafkaConsumer kafkaconsumer.Config    `yaml:"kafka_consumer"`
	KafkaProducer kafkaproducer.Config    `yaml:"kafka_producer"`

	Repository repository.Options `yaml:"repository"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Expand environment variables in the YAML content
	expandedData := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expandedData), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
