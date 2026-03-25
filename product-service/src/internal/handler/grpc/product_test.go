package grpc

import (
	"context"
	"errors"
	"testing"

	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service/product"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testProductID          = "product-123"
	testProductName        = "Test Product"
	testProductDescription = "Test Description"
	testProductPrice       = 99.99
	testProductSKU         = "SKU-12345"
	mockInventoryType      = "[]dto.CreateReserveInventory"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest, qty int, rsv int) (*dto.Product, error) {
	args := m.Called(ctx, req, qty, rsv)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Product), args.Error(1)
}

func (m *MockProductService) ListProduct(ctx context.Context) ([]*dto.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.Product), args.Error(1)
}

func (m *MockProductService) GetProduct(ctx context.Context, productID string) (*dto.Product, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.Product), args.Error(1)
}

func (m *MockProductService) ListCategories(ctx context.Context) ([]*dto.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.Category), args.Error(1)
}

func (m *MockProductService) GetCategoriesByProduct(ctx context.Context, productID string) ([]*dto.Category, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.Category), args.Error(1)
}

func (m *MockProductService) GetProductsByCategory(ctx context.Context, categoryID string) ([]*dto.Product, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.Product), args.Error(1)
}

func (m *MockProductService) CheckInventory(ctx context.Context, productID string) (int32, int32, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).(int32), args.Get(1).(int32), args.Error(2)
}

func (m *MockProductService) ReserveInventory(ctx context.Context, req []dto.CreateReserveInventory) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockProductService) ReleaseInventory(ctx context.Context, req []dto.CreateReserveInventory) error {
	args := m.Called(ctx, req)
	_ = len(req)
	return args.Error(0)
}

var _ product.ProductServiceItf = (*MockProductService)(nil)

func setupTestGrpc(mockProduct *MockProductService) (*Grpc, *service.Service) {
	mockSvc := &service.Service{}
	mockSvc.Product = mockProduct

	grpcHandler := &Grpc{
		log: zerolog.Logger{},
		svc: mockSvc,
	}

	return grpcHandler, mockSvc
}

func TestGetProductSuccess(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	expectedResp := &dto.Product{
		ID:          testProductID,
		Name:        testProductName,
		Description: testProductDescription,
		Price:       testProductPrice,
		SKU:         testProductSKU,
		IsActive:    true,
	}

	mockProduct.On("GetProduct", ctx, testProductID).Return(expectedResp, nil)

	req := &productpb.GetProductRequest{
		ProductId: testProductID,
	}

	resp, err := grpcHandler.GetProduct(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, testProductID, resp.Id)
	assert.Equal(t, testProductName, resp.Name)
	assert.Equal(t, testProductDescription, resp.Description)
	assert.Equal(t, testProductPrice, resp.Price)
	assert.Equal(t, testProductSKU, resp.Sku)
	assert.True(t, resp.IsActive)
	mockProduct.AssertExpectations(t)
}

func TestGetProductError(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	expectedErr := errors.New("product not found")
	mockProduct.On("GetProduct", ctx, testProductID).Return(nil, expectedErr)

	req := &productpb.GetProductRequest{
		ProductId: testProductID,
	}

	resp, err := grpcHandler.GetProduct(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, expectedErr, err)
	mockProduct.AssertExpectations(t)
}

func TestGetProductMissingProductID(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	req := &productpb.GetProductRequest{
		ProductId: "",
	}

	resp, err := grpcHandler.GetProduct(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "product ID is required")
}

func TestCheckInventorySuccess(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	mockProduct.On("CheckInventory", ctx, testProductID).Return(int32(100), int32(20), nil)

	req := &productpb.CheckInventoryRequest{
		ProductId: testProductID,
		Quantity:  10,
	}

	resp, err := grpcHandler.CheckInventory(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Available)
	assert.Equal(t, int32(100), resp.CurrentQuantity)
	assert.Equal(t, int32(20), resp.ReservedQuantity)
	mockProduct.AssertExpectations(t)
}

func TestCheckInventoryInsufficientStock(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	mockProduct.On("CheckInventory", ctx, testProductID).Return(int32(5), int32(3), nil)

	req := &productpb.CheckInventoryRequest{
		ProductId: testProductID,
		Quantity:  10,
	}

	resp, err := grpcHandler.CheckInventory(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Available)
	mockProduct.AssertExpectations(t)
}

func TestCheckInventoryError(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	expectedErr := errors.New("failed to check inventory")
	mockProduct.On("CheckInventory", ctx, testProductID).Return(int32(0), int32(0), expectedErr)

	req := &productpb.CheckInventoryRequest{
		ProductId: testProductID,
		Quantity:  10,
	}

	resp, err := grpcHandler.CheckInventory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockProduct.AssertExpectations(t)
}

func TestCheckInventoryMissingProductID(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	req := &productpb.CheckInventoryRequest{
		ProductId: "",
	}

	resp, err := grpcHandler.CheckInventory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "product ID is required")
}

func TestReserveInventorySuccess(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	mockProduct.On("ReserveInventory", ctx, mock.AnythingOfType(mockInventoryType)).Return(nil)

	req := &productpb.ReserveInventoryRequest{
		Items: []*productpb.InventoryItem{
			{ProductId: testProductID, Quantity: 5},
		},
	}

	resp, err := grpcHandler.ReserveInventory(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	mockProduct.AssertExpectations(t)
}

func TestReserveInventoryError(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	expectedErr := errors.New("failed to reserve inventory")
	mockProduct.On("ReserveInventory", ctx, mock.AnythingOfType(mockInventoryType)).Return(expectedErr)

	req := &productpb.ReserveInventoryRequest{
		Items: []*productpb.InventoryItem{
			{ProductId: testProductID, Quantity: 5},
		},
	}

	resp, err := grpcHandler.ReserveInventory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockProduct.AssertExpectations(t)
}

func TestReserveInventoryEmptyItems(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	req := &productpb.ReserveInventoryRequest{
		Items: nil,
	}

	resp, err := grpcHandler.ReserveInventory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Item is empty")
}

func TestReleaseInventorySuccess(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	mockProduct.On("ReleaseInventory", ctx, mock.AnythingOfType(mockInventoryType)).Return(nil)

	req := &productpb.ReleaseInventoryRequest{
		Items: []*productpb.InventoryItem{
			{ProductId: testProductID, Quantity: 10},
		},
	}

	resp, err := grpcHandler.ReleaseInventory(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	mockProduct.AssertExpectations(t)
}

func TestReleaseInventoryError(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	expectedErr := errors.New("failed to release inventory")
	mockProduct.On("ReleaseInventory", ctx, mock.AnythingOfType(mockInventoryType)).Return(expectedErr)

	req := &productpb.ReleaseInventoryRequest{
		Items: []*productpb.InventoryItem{
			{ProductId: testProductID, Quantity: 10},
		},
	}

	resp, err := grpcHandler.ReleaseInventory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockProduct.AssertExpectations(t)
}

func TestReleaseInventoryEmptyItems(t *testing.T) {
	mockProduct := new(MockProductService)
	grpcHandler, _ := setupTestGrpc(mockProduct)
	ctx := context.Background()

	req := &productpb.ReleaseInventoryRequest{
		Items: nil,
	}

	resp, err := grpcHandler.ReleaseInventory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Item is empty")
}

func TestConvertItems(t *testing.T) {
	items := []*productpb.InventoryItem{
		{ProductId: "product-1", Quantity: 10},
		{ProductId: "product-2", Quantity: 20},
		nil,
		{ProductId: "product-3", Quantity: 30},
	}

	result := convertItems(items)

	assert.Len(t, result, 3)
	assert.Equal(t, "product-1", result[0].ProductId)
	assert.Equal(t, int32(10), result[0].Quantity)
	assert.Equal(t, "product-2", result[1].ProductId)
	assert.Equal(t, int32(20), result[1].Quantity)
	assert.Equal(t, "product-3", result[2].ProductId)
	assert.Equal(t, int32(30), result[2].Quantity)
}

func TestConvertItemsEmpty(t *testing.T) {
	result := convertItems(nil)
	assert.Len(t, result, 0)
}
