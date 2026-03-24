package grpc

import (
	"testing"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"

	"github.com/rs/zerolog"
)

func TestInitGrpcHandler(t *testing.T) {
	log := zerolog.Logger{}
	svc := &service.Service{}

	InitGrpcHandler(log, svc)
}
