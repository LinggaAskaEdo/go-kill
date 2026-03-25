package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service/auth"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testEmail         = "test@example.com"
	testPassword      = "password123"
	testAuthID        = "auth-123"
	testUserID        = "user-123"
	testAccessToken   = "access-token-123"
	testRefreshToken  = "refresh-token-123"
	testIPAddress     = "192.168.1.1"
	testUserAgent     = "Mozilla/5.0"
	testTokenToRevoke = "token-to-revoke"
)

var _ auth.AuthServiceItf = (*MockAuthService)(nil)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) CreateAuthUser(ctx context.Context, req *dto.CreateAuthUserRequest) (*dto.CreateAuthUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CreateAuthUserResponse), args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, req *dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ValidateTokenResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RefreshTokenResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LogoutResponse), args.Error(1)
}

func setupTestGrpc(mockAuth *MockAuthService) (*Grpc, *service.Service) {
	mockSvc := &service.Service{}
	mockSvc.Auth = mockAuth

	grpcHandler := &Grpc{
		log: zerolog.Logger{},
		svc: mockSvc,
	}

	return grpcHandler, mockSvc
}

func TestCreateAuthUserSuccess(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedResp := &dto.CreateAuthUserResponse{
		Success: true,
		AuthId:  testAuthID,
	}

	mockAuth.On("CreateAuthUser", ctx, mock.MatchedBy(func(req *dto.CreateAuthUserRequest) bool {
		return req.Email == testEmail && req.Password == testPassword
	})).Return(expectedResp, nil)

	req := &authpb.CreateAuthUserRequest{
		Email:    testEmail,
		Password: testPassword,
	}

	resp, err := grpcHandler.CreateAuthUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, testAuthID, resp.AuthId)
	mockAuth.AssertExpectations(t)
}

func TestCreateAuthUserError(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedErr := errors.New("failed to create user")
	mockAuth.On("CreateAuthUser", ctx, mock.AnythingOfType("*dto.CreateAuthUserRequest")).Return(nil, expectedErr)

	req := &authpb.CreateAuthUserRequest{
		Email:    testEmail,
		Password: testPassword,
	}

	resp, err := grpcHandler.CreateAuthUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockAuth.AssertExpectations(t)
}

func TestLoginSuccess(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedResp := &dto.LoginResponse{
		Success:      true,
		AccessToken:  testAccessToken,
		RefreshToken: testRefreshToken,
		ExpiresIn:    3600,
	}

	mockAuth.On("Login", ctx, mock.MatchedBy(func(req *dto.LoginRequest) bool {
		return req.Email == testEmail && req.Password == testPassword
	})).Return(expectedResp, nil)

	req := &authpb.LoginRequest{
		Email:     testEmail,
		Password:  testPassword,
		IpAddress: testIPAddress,
		UserAgent: testUserAgent,
	}

	resp, err := grpcHandler.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, testAccessToken, resp.AccessToken)
	assert.Equal(t, testRefreshToken, resp.RefreshToken)
	assert.Equal(t, int64(3600), resp.ExpiresIn)
	mockAuth.AssertExpectations(t)
}

func TestLoginError(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedErr := errors.New("invalid credentials")
	mockAuth.On("Login", ctx, mock.AnythingOfType("*dto.LoginRequest")).Return(nil, expectedErr)

	req := &authpb.LoginRequest{
		Email:    testEmail,
		Password: "wrongpassword",
	}

	resp, err := grpcHandler.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockAuth.AssertExpectations(t)
}

func TestValidateTokenSuccess(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedResp := &dto.ValidateTokenResponse{
		Valid:  true,
		UserId: testUserID,
		Email:  testEmail,
	}

	mockAuth.On("ValidateToken", ctx, mock.MatchedBy(func(req *dto.ValidateTokenRequest) bool {
		return req.Token == "valid-token"
	})).Return(expectedResp, nil)

	req := &authpb.ValidateTokenRequest{
		Token: "valid-token",
	}

	resp, err := grpcHandler.ValidateToken(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, testUserID, resp.UserId)
	assert.Equal(t, testEmail, resp.Email)
	mockAuth.AssertExpectations(t)
}

func TestValidateTokenInvalidToken(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedErr := errors.New("invalid token")
	mockAuth.On("ValidateToken", ctx, mock.AnythingOfType("*dto.ValidateTokenRequest")).Return(nil, expectedErr)

	req := &authpb.ValidateTokenRequest{
		Token: "invalid-token",
	}

	resp, err := grpcHandler.ValidateToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockAuth.AssertExpectations(t)
}

func TestRefreshTokenSuccess(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedResp := &dto.RefreshTokenResponse{
		Success:      true,
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresIn:    7200,
	}

	mockAuth.On("RefreshToken", ctx, mock.MatchedBy(func(req *dto.RefreshTokenRequest) bool {
		return req.RefreshToken == "old-refresh-token"
	})).Return(expectedResp, nil)

	req := &authpb.RefreshTokenRequest{
		RefreshToken: "old-refresh-token",
	}

	resp, err := grpcHandler.RefreshToken(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "new-access-token", resp.AccessToken)
	assert.Equal(t, "new-refresh-token", resp.RefreshToken)
	assert.Equal(t, int64(7200), resp.ExpiresIn)
	mockAuth.AssertExpectations(t)
}

func TestRefreshTokenError(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedErr := errors.New("refresh token expired")
	mockAuth.On("RefreshToken", ctx, mock.AnythingOfType("*dto.RefreshTokenRequest")).Return(nil, expectedErr)

	req := &authpb.RefreshTokenRequest{
		RefreshToken: "expired-token",
	}

	resp, err := grpcHandler.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockAuth.AssertExpectations(t)
}

func TestLogoutSuccess(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedResp := &dto.LogoutResponse{
		Success: true,
		Message: "logged out successfully",
	}

	mockAuth.On("Logout", ctx, mock.MatchedBy(func(req *dto.LogoutRequest) bool {
		return req.Token == testTokenToRevoke && req.UserId == testUserID
	})).Return(expectedResp, nil)

	req := &authpb.LogoutRequest{
		Token:  testTokenToRevoke,
		UserId: testUserID,
	}

	resp, err := grpcHandler.Logout(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "logged out successfully", resp.Message)
	mockAuth.AssertExpectations(t)
}

func TestLogoutError(t *testing.T) {
	mockAuth := new(MockAuthService)
	grpcHandler, _ := setupTestGrpc(mockAuth)
	ctx := context.Background()

	expectedErr := errors.New("failed to logout")
	mockAuth.On("Logout", ctx, mock.AnythingOfType("*dto.LogoutRequest")).Return(nil, expectedErr)

	req := &authpb.LogoutRequest{
		Token:  testTokenToRevoke,
		UserId: testUserID,
	}

	resp, err := grpcHandler.Logout(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockAuth.AssertExpectations(t)
}
