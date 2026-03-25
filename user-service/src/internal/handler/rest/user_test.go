package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testUserAuthID         = "auth-123"
	testUserEmail          = "test@example.com"
	testFirstName          = "John"
	testLastName           = "Doe"
	headerContentType      = "Content-Type"
	headerContentTypeValue = "application/json"
	headerAuthBearer       = "Bearer valid-token"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userpb.CreateUserResponse), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userpb.GetUserResponse), args.Error(1)
}

func (m *MockUserService) GetAddress(ctx context.Context, req *userpb.GetAddressRequest) (*userpb.GetAddressResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userpb.GetAddressResponse), args.Error(1)
}

func (m *MockUserService) LogActivity(ctx context.Context, req *userpb.LogActivityRequest) (*userpb.LogActivityResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userpb.LogActivityResponse), args.Error(1)
}

func (m *MockUserService) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authpb.ValidateTokenResponse), args.Error(1)
}

func (m *MockUserService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserRegResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserRegResp), args.Error(1)
}

func (m *MockUserService) GetMe(ctx context.Context, userAuthID string) (*dto.UserResp, error) {
	args := m.Called(ctx, userAuthID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResp), args.Error(1)
}

func (m *MockUserService) GetActivities(ctx context.Context, userAuthID string, page string, limit string) (*dto.UserActivity, error) {
	args := m.Called(ctx, userAuthID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserActivity), args.Error(1)
}

func (m *MockUserService) GetAddresses(ctx context.Context, userAuthID string, page string, limit string) (*dto.AddressesResp, error) {
	args := m.Called(ctx, userAuthID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AddressesResp), args.Error(1)
}

func (m *MockUserService) CreateAddress(ctx context.Context, userAuthID string, req dto.CreateUserAddress) (string, error) {
	args := m.Called(ctx, userAuthID, req)
	return args.Get(0).(string), args.Error(1)
}

var _ user.UserServiceItf = (*MockUserService)(nil)

func setupTestRest(mockUser *MockUserService) *rest {
	mockSvc := &service.Service{}
	mockSvc.User = mockUser

	return &rest{
		gin: nil,
		svc: mockSvc,
	}
}

func setupRouter(handler *rest) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST(pathUsersRegister, handler.handleRegister)
	r.GET(pathUsersMe, handler.authMiddleware(), handler.handleGetMe)
	r.GET("/api/v1/users/me/activities", handler.authMiddleware(), handler.handleGetActivities)
	r.GET(pathUsersMeAddresses, handler.authMiddleware(), handler.handleGetAddresses)
	r.POST(pathUsersMeAddresses, handler.authMiddleware(), handler.handleCreateAddress)

	return r
}

func TestHandleRegisterSuccess(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	resp := &dto.UserRegResp{
		ID:        "user-123",
		AuthID:    "auth-123",
		Email:     testUserEmail,
		FirstName: testFirstName,
		LastName:  testLastName,
	}

	mockUser.On("RegisterUser", mock.Anything, mock.AnythingOfType("dto.RegisterUserRequest")).Return(resp, nil)

	reqBody := dto.RegisterUserRequest{
		Email:     testUserEmail,
		Password:  "password123",
		FirstName: testFirstName,
		LastName:  testLastName,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, pathUsersRegister, bytes.NewBuffer(body))
	req.Header.Set(headerContentType, headerContentTypeValue)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.HttpSuccessResp
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Meta.StatusCode)

	mockUser.AssertExpectations(t)
}

func TestHandleRegisterInvalidBody(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	req, _ := http.NewRequest(http.MethodPost, pathUsersRegister, bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(headerContentType, headerContentTypeValue)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleRegisterError(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	mockUser.On("RegisterUser", mock.Anything, mock.AnythingOfType("dto.RegisterUserRequest")).Return(nil, errors.New("email already exists"))

	reqBody := dto.RegisterUserRequest{
		Email:     testUserEmail,
		Password:  "password123",
		FirstName: testFirstName,
		LastName:  testLastName,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, pathUsersRegister, bytes.NewBuffer(body))
	req.Header.Set(headerContentType, headerContentTypeValue)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUser.AssertExpectations(t)
}

func TestHandleGetMeSuccess(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	resp := &dto.UserResp{
		ID:        "user-123",
		Email:     testUserEmail,
		FirstName: testFirstName,
		LastName:  testLastName,
	}

	mockUser.On("ValidateToken", mock.Anything, mock.Anything).Return(&authpb.ValidateTokenResponse{Valid: true, UserId: testUserAuthID}, nil)
	mockUser.On("GetMe", mock.Anything, testUserAuthID).Return(resp, nil)

	req, _ := http.NewRequest(http.MethodGet, pathUsersMe, nil)
	req.Header.Set("Authorization", headerAuthBearer)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockUser.AssertExpectations(t)
}

func TestHandleGetMeUnauthorized(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	req, _ := http.NewRequest(http.MethodGet, pathUsersMe, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandleGetMeInvalidToken(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	mockUser.On("ValidateToken", mock.Anything, mock.Anything).Return(nil, errors.New("invalid token"))

	req, _ := http.NewRequest(http.MethodGet, pathUsersMe, nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUser.AssertExpectations(t)
}

func TestHandleGetActivitiesSuccess(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	resp := &dto.UserActivity{
		Success: true,
	}

	mockUser.On("ValidateToken", mock.Anything, mock.Anything).Return(&authpb.ValidateTokenResponse{Valid: true, UserId: testUserAuthID}, nil)
	mockUser.On("GetActivities", mock.Anything, testUserAuthID, "1", "20").Return(resp, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/me/activities", nil)
	req.Header.Set("Authorization", headerAuthBearer)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockUser.AssertExpectations(t)
}

func TestHandleGetAddressesSuccess(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	resp := &dto.AddressesResp{
		Success: true,
		Data:    []*dto.Address{},
	}

	mockUser.On("ValidateToken", mock.Anything, mock.Anything).Return(&authpb.ValidateTokenResponse{Valid: true, UserId: testUserAuthID}, nil)
	mockUser.On("GetAddresses", mock.Anything, testUserAuthID, "1", "20").Return(resp, nil)

	req, _ := http.NewRequest(http.MethodGet, pathUsersMeAddresses, nil)
	req.Header.Set("Authorization", headerAuthBearer)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockUser.AssertExpectations(t)
}

func TestHandleCreateAddressSuccess(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	mockUser.On("ValidateToken", mock.Anything, mock.Anything).Return(&authpb.ValidateTokenResponse{Valid: true, UserId: testUserAuthID}, nil)
	mockUser.On("CreateAddress", mock.Anything, testUserAuthID, mock.AnythingOfType("dto.CreateUserAddress")).Return("addr-123", nil)

	reqBody := dto.CreateUserAddress{
		AddressType:   "shipping",
		StreetAddress: "123 Main St",
		City:          "Jakarta",
		PostalCode:    "12345",
		Country:       "Indonesia",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, pathUsersMeAddresses, bytes.NewBuffer(body))
	req.Header.Set(headerContentType, headerContentTypeValue)
	req.Header.Set("Authorization", headerAuthBearer)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockUser.AssertExpectations(t)
}

func TestHandleCreateAddressInvalidBody(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	mockUser.On("ValidateToken", mock.Anything, mock.Anything).Return(&authpb.ValidateTokenResponse{Valid: true, UserId: testUserAuthID}, nil)

	req, _ := http.NewRequest(http.MethodPost, pathUsersMeAddresses, bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(headerContentType, headerContentTypeValue)
	req.Header.Set("Authorization", headerAuthBearer)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockUser.AssertExpectations(t)
}

func TestAuthMiddlewareNoHeader(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	req, _ := http.NewRequest(http.MethodGet, pathUsersMe, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No authorization header")
}

func TestAuthMiddlewareInvalidFormat(t *testing.T) {
	mockUser := new(MockUserService)
	handler := setupTestRest(mockUser)
	router := setupRouter(handler)

	req, _ := http.NewRequest(http.MethodGet, pathUsersMe, nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization header format")
}
