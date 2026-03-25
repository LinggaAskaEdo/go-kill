package grpc

import (
	"context"
	"errors"
	"testing"

	orderpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/order"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service/order"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testOrderID           = "order-123"
	testOrderNumber       = "ORD-12345"
	testUserID            = "user-123"
	testShippingAddressID = "addr-shipping"
	testBillingAddressID  = "addr-billing"
	testPaymentMethod     = "credit_card"
	testReason            = "Customer request"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, reqData *dto.CreateOrderRequest) (*string, *string, float64, error) {
	args := m.Called(mock.Anything, reqData)
	if args.Get(0) == nil {
		return nil, nil, 0, args.Error(3)
	}
	orderID := args.Get(0).(*string)
	orderNumber := args.Get(1).(*string)
	totalAmount := args.Get(2).(float64)
	return orderID, orderNumber, totalAmount, args.Error(3)
}

func (m *MockOrderService) GetOrder(ctx context.Context, reqData *dto.GetOrderRequest) (*entity.Order, error) {
	args := m.Called(mock.Anything, reqData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderService) ListOrders(ctx context.Context, reqData *dto.ListOrderRequest) ([]*entity.Order, int32, error) {
	args := m.Called(mock.Anything, reqData)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.Order), args.Get(1).(int32), args.Error(2)
}

func (m *MockOrderService) CancelOrder(ctx context.Context, reqData *dto.CancelOrderRequest) error {
	args := m.Called(mock.Anything, reqData)
	return args.Error(0)
}

var _ order.OrderServiceItf = (*MockOrderService)(nil)

func setupTestGrpc(mockOrder *MockOrderService) (*Grpc, *service.Service) {
	mockSvc := &service.Service{}
	mockSvc.Order = mockOrder

	grpcHandler := &Grpc{
		log: zerolog.Logger{},
		svc: mockSvc,
	}

	return grpcHandler, mockSvc
}

func TestCreateOrderSuccess(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	orderID := testOrderID
	orderNumber := testOrderNumber
	totalAmount := 199.99

	mockOrder.On("CreateOrder", ctx, mock.MatchedBy(func(req *dto.CreateOrderRequest) bool {
		return req.UserID == testUserID &&
			req.ShippingAddressID == testShippingAddressID &&
			req.BillingAddressID == testBillingAddressID &&
			req.PaymentMethod == testPaymentMethod
	})).Return(&orderID, &orderNumber, totalAmount, nil)

	req := &orderpb.CreateOrderRequest{
		UserId:            testUserID,
		ShippingAddressId: testShippingAddressID,
		BillingAddressId:  testBillingAddressID,
		PaymentMethod:     testPaymentMethod,
		Items:             []*orderpb.OrderItem{},
	}

	resp, err := grpcHandler.CreateOrder(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, testOrderID, resp.OrderId)
	assert.Equal(t, testOrderNumber, resp.OrderNumber)
	assert.Equal(t, totalAmount, resp.TotalAmount)
	mockOrder.AssertExpectations(t)
}

func TestCreateOrderError(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	expectedErr := errors.New("failed to create order")
	mockOrder.On("CreateOrder", ctx, mock.AnythingOfType("*dto.CreateOrderRequest")).Return(nil, nil, 0, expectedErr)

	req := &orderpb.CreateOrderRequest{
		UserId:            testUserID,
		ShippingAddressId: testShippingAddressID,
		BillingAddressId:  testBillingAddressID,
		PaymentMethod:     testPaymentMethod,
		Items:             []*orderpb.OrderItem{},
	}

	resp, err := grpcHandler.CreateOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockOrder.AssertExpectations(t)
}

func TestGetOrderSuccess(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	expectedOrder := &entity.Order{
		ID:          testOrderID,
		OrderNumber: testOrderNumber,
		Status:      entity.StatusPending,
		TotalAmount: 299.99,
		Items:       []*entity.OrderItem{},
	}

	mockOrder.On("GetOrder", ctx, mock.MatchedBy(func(req *dto.GetOrderRequest) bool {
		return req.OrderID == testOrderID && req.UserID == testUserID
	})).Return(expectedOrder, nil)

	req := &orderpb.GetOrderRequest{
		OrderId: testOrderID,
		UserId:  testUserID,
	}

	resp, err := grpcHandler.GetOrder(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testOrderID, resp.Id)
	assert.Equal(t, testOrderNumber, resp.OrderNumber)
	assert.Equal(t, string(entity.StatusPending), resp.Status)
	assert.Equal(t, 299.99, resp.TotalAmount)
	mockOrder.AssertExpectations(t)
}

func TestGetOrderError(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	expectedErr := errors.New("order not found")
	mockOrder.On("GetOrder", ctx, mock.AnythingOfType("*dto.GetOrderRequest")).Return(nil, expectedErr)

	req := &orderpb.GetOrderRequest{
		OrderId: "non-existent-order",
		UserId:  testUserID,
	}

	resp, err := grpcHandler.GetOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockOrder.AssertExpectations(t)
}

func TestListOrdersSuccess(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	orders := []*entity.Order{
		{
			ID:          testOrderID,
			OrderNumber: testOrderNumber,
			Status:      entity.StatusPending,
			TotalAmount: 100.00,
		},
	}

	mockOrder.On("ListOrders", ctx, mock.MatchedBy(func(req *dto.ListOrderRequest) bool {
		return req.UserID == testUserID && req.Limit == 20 && req.Offset == 0
	})).Return(orders, int32(1), nil)

	req := &orderpb.ListOrdersRequest{
		UserId: testUserID,
		Page:   1,
		Limit:  20,
	}

	resp, err := grpcHandler.ListOrders(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Orders, 1)
	assert.Equal(t, int32(1), resp.Total)
	mockOrder.AssertExpectations(t)
}

func TestListOrdersDefaultLimit(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	mockOrder.On("ListOrders", ctx, mock.MatchedBy(func(req *dto.ListOrderRequest) bool {
		return req.Limit == 20
	})).Return([]*entity.Order{}, int32(0), nil)

	req := &orderpb.ListOrdersRequest{
		UserId: testUserID,
		Page:   1,
		Limit:  0,
	}

	resp, err := grpcHandler.ListOrders(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockOrder.AssertExpectations(t)
}

func TestListOrdersMaxLimit(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	mockOrder.On("ListOrders", ctx, mock.MatchedBy(func(req *dto.ListOrderRequest) bool {
		return req.Limit == 100
	})).Return([]*entity.Order{}, int32(0), nil)

	req := &orderpb.ListOrdersRequest{
		UserId: testUserID,
		Page:   1,
		Limit:  200,
	}

	resp, err := grpcHandler.ListOrders(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockOrder.AssertExpectations(t)
}

func TestListOrdersError(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	expectedErr := errors.New("failed to list orders")
	mockOrder.On("ListOrders", ctx, mock.AnythingOfType("*dto.ListOrderRequest")).Return(nil, int32(0), expectedErr)

	req := &orderpb.ListOrdersRequest{
		UserId: testUserID,
		Page:   1,
		Limit:  20,
	}

	resp, err := grpcHandler.ListOrders(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockOrder.AssertExpectations(t)
}

func TestCancelOrderSuccess(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	mockOrder.On("CancelOrder", ctx, mock.MatchedBy(func(req *dto.CancelOrderRequest) bool {
		return req.OrderID == testOrderID && req.UserID == testUserID && req.Reason == testReason
	})).Return(nil)

	req := &orderpb.CancelOrderRequest{
		OrderId: testOrderID,
		UserId:  testUserID,
		Reason:  testReason,
	}

	resp, err := grpcHandler.CancelOrder(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	mockOrder.AssertExpectations(t)
}

func TestCancelOrderError(t *testing.T) {
	mockOrder := new(MockOrderService)
	grpcHandler, _ := setupTestGrpc(mockOrder)
	ctx := context.Background()

	expectedErr := errors.New("failed to cancel order")
	mockOrder.On("CancelOrder", ctx, mock.AnythingOfType("*dto.CancelOrderRequest")).Return(expectedErr)

	req := &orderpb.CancelOrderRequest{
		OrderId: testOrderID,
		UserId:  testUserID,
		Reason:  testReason,
	}

	resp, err := grpcHandler.CancelOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockOrder.AssertExpectations(t)
}
