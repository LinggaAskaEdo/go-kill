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
	err := r.redis0.Get(ctx, cacheKey).Scan(&product)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Msg("Product not found in cache: " + err.Error())
		return nil
	}

	return &product
}

func (r *productRepository) setInventoryReservedCache(ctx context.Context, productID string, qty int32) {
	if err := r.redis0.HIncrBy(ctx, fmt.Sprintf("inventory:%s", productID), "reserved", int64(qty)).Err(); err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("failed to update inventory reserved cache")
	}
}

func (r *productRepository) setInventoryReleaseCache(ctx context.Context, productID string, qty int32) {
	if err := r.redis0.HIncrBy(ctx, fmt.Sprintf("inventory:%s", productID), "reserved", -int64(qty)).Err(); err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("failed to update inventory release cache")
	}
}

func (r *productRepository) setProductCache(ctx context.Context, product *entity.Product) error {
	cacheKey := fmt.Sprintf("product:%s", product.ID)
	return r.redis0.Set(ctx, cacheKey, product, 0).Err()
}
