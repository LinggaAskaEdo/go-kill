package scheduler

import (
	"context"
	"errors"
	"testing"

	"github.com/linggaaskaedo/go-kill/common/component/scheduler"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service/product"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testBatchSize            = 5
	testCron                 = "0 * * * *"
	mockCreateProductRequest = "dto.CreateProductRequest"
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

func TestName(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)

	assert.Equal(t, "product_generator_job", job.Name())
}

func TestSchedule(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)

	assert.Equal(t, testCron, job.Schedule())
}

func TestRunDisabled(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   false,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)
	ctx := context.Background()

	err := job.Run(ctx)

	assert.NoError(t, err)
	mockProduct.AssertNotCalled(t, "CreateProduct")
}

func TestRunSuccess(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: 2,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)
	ctx := context.Background()

	categories := []*dto.Category{
		{ID: "cat-1", Name: "Electronics"},
		{ID: "cat-2", Name: "Accessories"},
	}

	mockProduct.On("ListCategories", ctx).Return(categories, nil)
	mockProduct.On("CreateProduct", ctx, mock.AnythingOfType(mockCreateProductRequest), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(&dto.Product{ID: "prod-1"}, nil)

	err := job.Run(ctx)

	assert.NoError(t, err)
	mockProduct.AssertNumberOfCalls(t, "CreateProduct", 2)
	mockProduct.AssertExpectations(t)
}

func TestRunCreateProductError(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: 1,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)
	ctx := context.Background()

	categories := []*dto.Category{
		{ID: "cat-1", Name: "Electronics"},
	}

	mockProduct.On("ListCategories", ctx).Return(categories, nil)
	mockProduct.On("CreateProduct", ctx, mock.AnythingOfType(mockCreateProductRequest), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, errors.New("failed to create product"))

	err := job.Run(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create product")
	mockProduct.AssertExpectations(t)
}

func TestRunNoCategories(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: 1,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)
	ctx := context.Background()

	mockProduct.On("ListCategories", ctx).Return([]*dto.Category{}, nil)
	mockProduct.On("CreateProduct", ctx, mock.AnythingOfType(mockCreateProductRequest), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(&dto.Product{ID: "prod-1"}, nil)

	err := job.Run(ctx)

	assert.NoError(t, err)
	mockProduct.AssertNumberOfCalls(t, "CreateProduct", 1)
	mockProduct.AssertExpectations(t)
}

func TestGenerateProductName(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)

	for i := 0; i < 10; i++ {
		name := job.generateProductName()
		assert.NotEmpty(t, name)
		assert.Contains(t, name, " ")
	}
}

func TestGenerateSKU(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)

	sku := job.generateSKU(1)
	assert.Contains(t, sku, "SKU-")
	assert.Contains(t, sku, "-000001")
}

func TestGeneratePrice(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)

	for i := 0; i < 10; i++ {
		price := job.generatePrice()
		assert.GreaterOrEqual(t, price, 5.0)
		assert.LessOrEqual(t, price, 1000.0)
	}
}

func TestGenerateDescription(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)

	desc := job.generateDescription("Test Product")
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "test product")
}

func TestSelectRandomCategoriesReturnsValidSubset(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)

	ids := []string{"cat-1", "cat-2", "cat-3"}

	categories := job.selectRandomCategories(ids)

	assert.NotNil(t, categories)
	assert.GreaterOrEqual(t, len(categories), 1)
	assert.Less(t, len(categories), len(ids))
}

func TestSelectRandomCategoriesEmptyIDs(t *testing.T) {
	mockProduct := new(MockProductService)
	cfg := scheduler.Config{
		Name:      "ProductGenerator",
		Enabled:   true,
		Cron:      testCron,
		BatchSize: testBatchSize,
	}

	job := NewProductGeneratorJob(zerolog.Logger{}, mockProduct, cfg)

	categories := job.selectRandomCategories(nil)

	assert.Nil(t, categories)
}
