package config

import (
	"os"

	"github.com/linggaaskaedo/go-kill/common/database"
	"github.com/linggaaskaedo/go-kill/common/grpcclient"
	"github.com/linggaaskaedo/go-kill/common/grpcserver"
	"github.com/linggaaskaedo/go-kill/common/http"
	"github.com/linggaaskaedo/go-kill/common/logger"
	"github.com/linggaaskaedo/go-kill/common/query"
	"github.com/linggaaskaedo/go-kill/common/redis"
	"github.com/linggaaskaedo/go-kill/common/scheduler"
	"github.com/linggaaskaedo/go-kill/common/server"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Logger     logger.Config                `yaml:"logger"`
	Redis      redis.Config                 `yaml:"redis"`
	Database   map[string]database.Config   `yaml:"database"`
	Query      query.Config                 `yaml:"queries"`
	Scheduler  scheduler.Config             `yaml:"scheduler"`
	GRPCClient map[string]grpcclient.Config `yaml:"grpc_client"`
	GRPCServer grpcserver.Config            `yaml:"grpc_server"`
	Http       http.Config                  `yaml:"http"`
	Server     server.Config                `yaml:"server"`
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
