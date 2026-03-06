package product

import (
	"context"
	"fmt"

	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (r *productRepository) getProductCache(ctx context.Context, productID string) *entity.Product {
	var product entity.Product

	cacheKey := fmt.Sprintf("product:%s", productID)
	err := r.redis0.Get(context.Background(), cacheKey).Scan(&product)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Msg("Product not found in cache: " + err.Error())
		return nil
	}

	return &product
}
