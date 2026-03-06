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
