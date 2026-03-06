package product

import (
	"context"

	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
)

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest, qty int, rsv int) (*dto.Product, error) {
	product := toProductEntity(req)

	product, err := s.productRepository.CreateProduct(ctx, product, qty, rsv)
	if err != nil {
		return nil, err
	}

	return toProduct(product), nil
}

func (s *productService) ListProduct(ctx context.Context) ([]*dto.Product, error) {
	products, err := s.productRepository.GetListProduct(ctx)
	if err != nil {
		return nil, err
	}

	return toProducts(products), nil
}

func (s *productService) GetProduct(ctx context.Context, productID string) (*dto.Product, error) {
	product, err := s.productRepository.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	return toProduct(product), nil
}

func (s *productService) ListCategories(ctx context.Context) ([]*dto.Category, error) {
	categories, err := s.productRepository.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	return toCategories(categories), nil
}

func (s *productService) GetCategoriesByProduct(ctx context.Context, productID string) ([]*dto.Category, error) {
	categories, err := s.productRepository.GetCategoriesByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	return toCategories(categories), nil
}

func (s *productService) GetProductsByCategory(ctx context.Context, categoryID string) ([]*dto.Product, error) {
	products, err := s.productRepository.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return toProducts(products), nil
}

func (s *productService) CheckInventory(ctx context.Context, productID string) (int32, int32, error) {
	return s.productRepository.CheckInventory(ctx, productID)
}

func (s *productService) ReserveInventory(ctx context.Context, req []dto.CreateReserveInventory) error {
	return s.productRepository.ReserveInventory(ctx, req)
}

func (s *productService) ReleaseInventory(ctx context.Context, req []dto.CreateReserveInventory) error {
	return s.productRepository.ReleaseInventory(ctx, req)
}
