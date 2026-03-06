package service

import (
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/repository"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service/product"
)

type Service struct {
	Product product.ProductServiceItf
}

func InitService(repository *repository.Repository) *Service {
	return &Service{
		Product: product.InitProductService(
			repository.Product,
		),
	}
}
