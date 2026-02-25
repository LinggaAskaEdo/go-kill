package config

import (
	"os"

	"github.com/linggaaskaedo/go-kill/common/component/database"
	"github.com/linggaaskaedo/go-kill/common/component/grpcclient"
	"github.com/linggaaskaedo/go-kill/common/component/grpcserver"
	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/common/component/redis"
	"github.com/linggaaskaedo/go-kill/common/pkg/logger"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Logger     logger.Config                `yaml:"logger"`
	Redis      redis.Config                 `yaml:"redis"`
	Database   map[string]database.Config   `yaml:"database"`
	Query      query.Config                 `yaml:"queries"`
	GRPCClient map[string]grpcclient.Config `yaml:"grpc_client"`
	GRPCServer grpcserver.Config            `yaml:"grpc_server"`
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
