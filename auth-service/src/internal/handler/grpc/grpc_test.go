package grpc

import (
	"testing"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
	"github.com/rs/zerolog"
)

func TestInitGrpcHandler(t *testing.T) {
	log := zerolog.Logger{}
	svc := &service.Service{}

	handler := InitGrpcHandler(log, svc)

	if handler == nil {
		t.Error("expected non-nil handler")
	}
	if handler.svc != svc {
		t.Error("expected service to match")
	}
}
