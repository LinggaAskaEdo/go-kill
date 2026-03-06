package product

import (
	"context"

	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/repository/product"
)

type ProductServiceItf interface {
	CreateProduct(ctx context.Context, req dto.CreateProductRequest, qty int, rsv int) (*dto.Product, error)
	ListProduct(ctx context.Context) ([]*dto.Product, error)
	GetProduct(ctx context.Context, productID string) (*dto.Product, error)
	ListCategories(ctx context.Context) ([]*dto.Category, error)
	GetCategoriesByProduct(ctx context.Context, productID string) ([]*dto.Category, error)
	GetProductsByCategory(ctx context.Context, categoryID string) ([]*dto.Product, error)
	CheckInventory(ctx context.Context, productID string) (int32, int32, error)
	ReserveInventory(ctx context.Context, req []dto.CreateReserveInventory) error
	ReleaseInventory(ctx context.Context, req []dto.CreateReserveInventory) error
}

type productService struct {
	productRepository product.ProductRepositoryItf
}

func InitProductService(productRepository product.ProductRepositoryItf) ProductServiceItf {
	return &productService{
		productRepository: productRepository,
	}
}

func toProductEntity(product dto.CreateProductRequest) *entity.Product {
	return &entity.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Categories:  product.Categories,
	}
}

func toProducts(products []*entity.Product) []*dto.Product {
	if products == nil {
		return nil
	}

	result := make([]*dto.Product, len(products))
	for i, product := range products {
		result[i] = &dto.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			SKU:         product.SKU,
			IsActive:    product.IsActive,
			Categories:  product.Categories,
		}
	}

	return result
}

func toProduct(product *entity.Product) *dto.Product {
	return &dto.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Categories:  product.Categories,
	}
}

func toCategories(categories []*entity.Category) []*dto.Category {
	if categories == nil {
		return nil
	}

	result := make([]*dto.Category, len(categories))
	for i, category := range categories {
		result[i] = &dto.Category{
			ID:        category.ID,
			Name:      category.Name,
			Slug:      category.Slug,
			ParentID:  category.ParentID,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		}
	}

	return result
}
