package config

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/linggaaskaedo/go-kill/common/database"
	"github.com/linggaaskaedo/go-kill/common/logger"
	"github.com/linggaaskaedo/go-kill/common/query"
	"github.com/linggaaskaedo/go-kill/common/server"
)

type Config struct {
	Logger   logger.Config              `yaml:"logger"`
	Server   server.Config              `yaml:"server"`
	Database map[string]database.Config `yaml:"database"`
	Query    query.Config               `yaml:"queries"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
