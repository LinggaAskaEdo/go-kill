package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gin-gonic/gin"
)

const (
	testProductID   = "019d227d-6eac-749c-b935-263bddc5a630"
	testCategoryID  = "019d227d-6eac-749c-b935-263bddc5a631"
	testProductName = "Test Product"
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
	return args.Error(0)
}

var _ product.ProductServiceItf = (*MockProductService)(nil)

func setupTestRest(mockProduct *MockProductService) *rest {
	mockSvc := &service.Service{}
	mockSvc.Product = mockProduct

	return &rest{
		gin: nil,
		svc: mockSvc,
	}
}

func setupRouter(handler *rest) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/api/v1/products", handler.handleListProducts)
	r.GET("/api/v1/products/:id", handler.handleGetProduct)
	r.GET("/api/v1/categories", handler.handleListCategories)
	r.GET("/api/v1/products/:id/categories", handler.handleGetCategoriesByProduct)
	r.GET("/api/v1/categories/:id/products", handler.handleGetProductsByCategory)

	return r
}

func TestHandleListProducts_Success(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	products := []*dto.Product{
		{
			ID:   testProductID,
			Name: testProductName,
		},
	}

	mockProduct.On("ListProduct", mock.Anything).Return(products, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HttpSuccessResp
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Meta.StatusCode)

	mockProduct.AssertExpectations(t)
}

func TestHandleListProducts_Error(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	mockProduct.On("ListProduct", mock.Anything).Return(nil, errors.New("database error"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockProduct.AssertExpectations(t)
}

func TestHandleGetProduct_Success(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	product := &dto.Product{
		ID:   testProductID,
		Name: testProductName,
	}

	mockProduct.On("GetProduct", mock.Anything, testProductID).Return(product, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/"+testProductID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HttpSuccessResp
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Meta.StatusCode)

	mockProduct.AssertExpectations(t)
}

func TestHandleGetProduct_InvalidID(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/invalid-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestHandleGetProduct_NotFound(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	mockProduct.On("GetProduct", mock.Anything, testProductID).Return(nil, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/"+testProductID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
	mockProduct.AssertExpectations(t)
}

func TestHandleListCategories_Success(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	categories := []*dto.Category{
		{
			ID:   testCategoryID,
			Name: "Electronics",
		},
	}

	mockProduct.On("ListCategories", mock.Anything).Return(categories, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HttpSuccessResp
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Meta.StatusCode)

	mockProduct.AssertExpectations(t)
}

func TestHandleListCategories_Error(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	mockProduct.On("ListCategories", mock.Anything).Return(nil, errors.New("database error"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockProduct.AssertExpectations(t)
}

func TestHandleGetCategoriesByProduct_Success(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	categories := []*dto.Category{
		{
			ID:   testCategoryID,
			Name: "Electronics",
		},
	}

	mockProduct.On("GetCategoriesByProduct", mock.Anything, testProductID).Return(categories, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/"+testProductID+"/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HttpSuccessResp
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Meta.StatusCode)

	mockProduct.AssertExpectations(t)
}

func TestHandleGetCategoriesByProduct_InvalidID(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/products/invalid-id/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestHandleGetProductsByCategory_Success(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	products := []*dto.Product{
		{
			ID:   testProductID,
			Name: testProductName,
		},
	}

	mockProduct.On("GetProductsByCategory", mock.Anything, testCategoryID).Return(products, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/categories/"+testCategoryID+"/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.HttpSuccessResp
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Meta.StatusCode)

	mockProduct.AssertExpectations(t)
}

func TestHandleGetProductsByCategory_InvalidID(t *testing.T) {
	mockProduct := new(MockProductService)
	handler := setupTestRest(mockProduct)
	router := setupRouter(handler)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/categories/invalid-id/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}
