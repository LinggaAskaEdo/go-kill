package grpc

import (
	"context"
	"errors"
	"testing"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service/user"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testUserID    = "user-123"
	testAuthID    = "auth-123"
	testUserEmail = "test@example.com"
	testAddressID = "addr-123"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	args := m.Called(mock.Anything, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userpb.CreateUserResponse), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	args := m.Called(mock.Anything, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userpb.GetUserResponse), args.Error(1)
}

func (m *MockUserService) GetAddress(ctx context.Context, req *userpb.GetAddressRequest) (*userpb.GetAddressResponse, error) {
	args := m.Called(mock.Anything, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userpb.GetAddressResponse), args.Error(1)
}

func (m *MockUserService) LogActivity(ctx context.Context, req *userpb.LogActivityRequest) (*userpb.LogActivityResponse, error) {
	args := m.Called(mock.Anything, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userpb.LogActivityResponse), args.Error(1)
}

func (m *MockUserService) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	args := m.Called(mock.Anything, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authpb.ValidateTokenResponse), args.Error(1)
}

func (m *MockUserService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserRegResp, error) {
	args := m.Called(mock.Anything, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserRegResp), args.Error(1)
}

func (m *MockUserService) GetMe(ctx context.Context, userAuthID string) (*dto.UserResp, error) {
	args := m.Called(mock.Anything, userAuthID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResp), args.Error(1)
}

func (m *MockUserService) GetActivities(ctx context.Context, userAuthID string, page string, limit string) (*dto.UserActivity, error) {
	args := m.Called(mock.Anything, userAuthID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserActivity), args.Error(1)
}

func (m *MockUserService) GetAddresses(ctx context.Context, userAuthID string, page string, limit string) (*dto.AddressesResp, error) {
	args := m.Called(mock.Anything, userAuthID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AddressesResp), args.Error(1)
}

func (m *MockUserService) CreateAddress(ctx context.Context, userAuthID string, req dto.CreateUserAddress) (string, error) {
	args := m.Called(mock.Anything, userAuthID, req)
	return args.Get(0).(string), args.Error(1)
}

var _ user.UserServiceItf = (*MockUserService)(nil)

func setupTestGrpc(mockUser *MockUserService) (*Grpc, *service.Service) {
	mockSvc := &service.Service{}
	mockSvc.User = mockUser

	grpcHandler := &Grpc{
		log: zerolog.Logger{},
		svc: mockSvc,
	}

	return grpcHandler, mockSvc
}

func TestCreateUserSuccess(t *testing.T) {
	mockUser := new(MockUserService)
	grpcHandler, _ := setupTestGrpc(mockUser)
	ctx := context.Background()

	expectedResp := &userpb.CreateUserResponse{
		UserId:  testUserID,
		Success: true,
	}

	mockUser.On("CreateUser", mock.Anything, mock.Anything).Return(expectedResp, nil)

	req := &userpb.CreateUserRequest{
		Email: testUserEmail,
	}

	resp, err := grpcHandler.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testUserID, resp.GetUserId())
	assert.True(t, resp.GetSuccess())
	mockUser.AssertExpectations(t)
}

func TestCreateUserError(t *testing.T) {
	mockUser := new(MockUserService)
	grpcHandler, _ := setupTestGrpc(mockUser)
	ctx := context.Background()

	expectedErr := errors.New("failed to create user")
	mockUser.On("CreateUser", mock.Anything, mock.Anything).Return(nil, expectedErr)

	req := &userpb.CreateUserRequest{
		Email: testUserEmail,
	}

	resp, err := grpcHandler.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockUser.AssertExpectations(t)
}

func TestGetUserSuccess(t *testing.T) {
	mockUser := new(MockUserService)
	grpcHandler, _ := setupTestGrpc(mockUser)
	ctx := context.Background()

	expectedResp := &userpb.GetUserResponse{
		Id:    testUserID,
		Email: testUserEmail,
	}

	mockUser.On("GetUser", mock.Anything, mock.Anything).Return(expectedResp, nil)

	req := &userpb.GetUserRequest{
		UserId: testUserID,
	}

	resp, err := grpcHandler.GetUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testUserID, resp.GetId())
	assert.Equal(t, testUserEmail, resp.GetEmail())
	mockUser.AssertExpectations(t)
}

func TestGetUserError(t *testing.T) {
	mockUser := new(MockUserService)
	grpcHandler, _ := setupTestGrpc(mockUser)
	ctx := context.Background()

	expectedErr := errors.New("user not found")
	mockUser.On("GetUser", mock.Anything, mock.Anything).Return(nil, expectedErr)

	req := &userpb.GetUserRequest{
		UserId: "non-existent-user",
	}

	resp, err := grpcHandler.GetUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockUser.AssertExpectations(t)
}

func TestGetAddressSuccess(t *testing.T) {
	mockUser := new(MockUserService)
	grpcHandler, _ := setupTestGrpc(mockUser)
	ctx := context.Background()

	expectedResp := &userpb.GetAddressResponse{
		Id: testAddressID,
	}

	mockUser.On("GetAddress", mock.Anything, mock.Anything).Return(expectedResp, nil)

	req := &userpb.GetAddressRequest{
		UserId:    testUserID,
		AddressId: testAddressID,
	}

	resp, err := grpcHandler.GetAddress(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testAddressID, resp.GetId())
	mockUser.AssertExpectations(t)
}

func TestGetAddressError(t *testing.T) {
	mockUser := new(MockUserService)
	grpcHandler, _ := setupTestGrpc(mockUser)
	ctx := context.Background()

	expectedErr := errors.New("address not found")
	mockUser.On("GetAddress", mock.Anything, mock.Anything).Return(nil, expectedErr)

	req := &userpb.GetAddressRequest{
		UserId:    testUserID,
		AddressId: "non-existent-addr",
	}

	resp, err := grpcHandler.GetAddress(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockUser.AssertExpectations(t)
}

func TestLogActivitySuccess(t *testing.T) {
	mockUser := new(MockUserService)
	grpcHandler, _ := setupTestGrpc(mockUser)
	ctx := context.Background()

	expectedResp := &userpb.LogActivityResponse{
		Success: true,
	}

	mockUser.On("LogActivity", mock.Anything, mock.Anything).Return(expectedResp, nil)

	req := &userpb.LogActivityRequest{
		UserId: testUserID,
	}

	resp, err := grpcHandler.LogActivity(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.GetSuccess())
	mockUser.AssertExpectations(t)
}

func TestLogActivityError(t *testing.T) {
	mockUser := new(MockUserService)
	grpcHandler, _ := setupTestGrpc(mockUser)
	ctx := context.Background()

	expectedErr := errors.New("failed to log activity")
	mockUser.On("LogActivity", mock.Anything, mock.Anything).Return(nil, expectedErr)

	req := &userpb.LogActivityRequest{
		UserId: testUserID,
	}

	resp, err := grpcHandler.LogActivity(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockUser.AssertExpectations(t)
}
