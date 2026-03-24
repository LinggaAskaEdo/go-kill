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

func TestServe(t *testing.T) {
	log := zerolog.Logger{}
	svc := &service.Service{}

	handler := &Grpc{
		log: log,
		svc: svc,
	}

	methods := handler.Serve()

	if len(methods) != 5 {
		t.Errorf("expected 5 methods, got %d", len(methods))
	}

	expectedMethods := []string{
		"/auth.AuthService/CreateAuthUser",
		"/auth.AuthService/Login",
		"/auth.AuthService/ValidateToken",
		"/auth.AuthService/RefreshToken",
		"/auth.AuthService/Logout",
	}

	for i, expected := range expectedMethods {
		if methods[i] != expected {
			t.Errorf("expected method %s at index %d, got %s", expected, i, methods[i])
		}
	}
}
