package product

import (
	"context"
	"database/sql"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (r *productRepository) CreateProduct(ctx context.Context, product *entity.Product, qty int, rsv int) (*entity.Product, error) {
	tx, err := r.db0.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("tx_create_product")
		return nil, err
	}

	tx, product, err = r.createProductSQL(ctx, tx, product)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	tx, err = r.createProductCategoriesSQL(ctx, tx, product.ID, product.Categories)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	tx, err = r.createProductInventorySQL(ctx, tx, product.ID, qty, rsv)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("commit_create_product")
		return nil, x.Wrap(err, "commit_create_product")
	}

	return product, nil
}

func (r *productRepository) GetListProduct(ctx context.Context) ([]*entity.Product, error) {
	return r.getListProductSQL(ctx)
}

func (r *productRepository) GetProduct(ctx context.Context, productID string) (*entity.Product, error) {
	product := r.getProductCache(ctx, productID)
	if product == nil {
		product, err := r.getProductByIDSQL(ctx, productID)
		if err != nil {
			return nil, err
		}

		return product, nil
	}

	return product, nil
}

func (r *productRepository) ListCategories(ctx context.Context) ([]*entity.Category, error) {
	return r.getCategoriesSQL(ctx)
}

func (r *productRepository) GetCategoriesByProduct(ctx context.Context, productID string) ([]*entity.Category, error) {
	return r.getCategoriesByProductIDSQL(ctx, productID)
}

func (r *productRepository) GetProductsByCategory(ctx context.Context, categoryID string) ([]*entity.Product, error) {
	return r.getProductsByCategoryIDSQL(ctx, categoryID)
}
func (r *productRepository) CheckInventory(ctx context.Context, productID string) (int32, int32, error) {
	return r.getInventoryByProductIDSQL(ctx, productID)
}

func (r *productRepository) ReserveInventory(ctx context.Context, req []dto.CreateReserveInventory) error {
	tx, err := r.db0.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("tx_reserve_inventory")
		return err
	}

	tx, err = r.createReserveInventorySQL(ctx, tx, req)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("commit_reserve_inventory")
		return x.Wrap(err, "commit_reserve_inventory")
	}

	return nil
}

func (r *productRepository) ReleaseInventory(ctx context.Context, req []dto.CreateReserveInventory) error {
	return r.createReleaseInventorySQL(ctx, req)
}
