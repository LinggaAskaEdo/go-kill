package product

import (
	"context"
	"database/sql"
	"errors"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func (r *productRepository) createProductSQL(ctx context.Context, tx *sqlx.Tx, product *entity.Product) (*sqlx.Tx, *entity.Product, error) {
	query, _ := r.queryLoader.Get("CreateProduct")
	row := tx.QueryRowContext(ctx, query, product.Name, product.Description, product.Price, product.SKU, product.IsActive).Scan(&product.ID)
	if err := row; err != nil {
		zerolog.Ctx(ctx).Error().Str("id", product.ID).Msg("create_product_sql")
		return tx, nil, x.WrapWithCode(err, x.CodeSQLCreate, "create_product_sql")
	}

	return tx, product, nil
}

func (r *productRepository) createProductCategoriesSQL(ctx context.Context, tx *sqlx.Tx, productID string, categoryIDs []string) (*sqlx.Tx, error) {
	if len(categoryIDs) == 0 {
		zerolog.Ctx(ctx).Error().Str("productID", productID).Any("categoryIDs", categoryIDs).Msg("create_product_categories_sql")
		return tx, x.NewWithCode(x.CodeSQLCreate, "create_product_categories_sql")
	}

	query, _ := r.queryLoader.Get("CreateProductCategories")
	result := tx.MustExecContext(ctx, query, productID, categoryIDs)
	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("productID", productID).Strs("categoryIDs", categoryIDs).Msg("create_product_categories_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "create_product_categories_sql")
	}

	if rows != int64(len(categoryIDs)) {
		zerolog.Ctx(ctx).Error().Str("productID", productID).Strs("categoryID", categoryIDs).Msg("create_product_categories_sql")
		return tx, x.NewWithCode(x.CodeSQLCannotRetrieveAffectedRows, "create_product_categories_sql")
	}

	return tx, nil
}

func (r *productRepository) createProductInventorySQL(ctx context.Context, tx *sqlx.Tx, productID string, qty, rsv int) (*sqlx.Tx, error) {
	query, _ := r.queryLoader.Get("CreateProductInventory")
	result := tx.MustExecContext(ctx, query, productID, qty, rsv)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		zerolog.Ctx(ctx).Error().Str("productID", productID).Int("qty", qty).Int("rsv", rsv).Msg("create_product_inventory_sql")
		return tx, x.NewWithCode(x.CodeSQLCreate, "create_product_inventory_sql")
	}

	return tx, nil
}

func (r *productRepository) getListProductSQL(ctx context.Context) ([]*entity.Product, error) {
	query, _ := r.queryLoader.Get("GetListProducts")
	rows, err := r.db0.QueryContext(ctx, query)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_list_product_sql")

		if errors.Is(err, sql.ErrNoRows) {
			return nil, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_list_product_sql")
		}

		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_list_product_sql")
	}
	defer rows.Close()

	products := make([]*entity.Product, 0, 10) // adjust capacity as needed
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.SKU,
			&product.IsActive,
		); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("get_list_product_sql_row_scan")
			return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_list_product_sql_row_scan")
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_list_product_sql_rows")
		return nil, x.WrapWithCode(err, x.CodeSQLRead, "get_list_product_sql_rows")
	}

	return products, nil
}

func (r *productRepository) getProductByIDSQL(ctx context.Context, productID string) (*entity.Product, error) {
	var product entity.Product

	query, _ := r.queryLoader.Get("GetProductByID")
	err := r.db0.QueryRowxContext(ctx, query, productID).StructScan(&product)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_product_by_id_sql")

		if errors.Is(err, sql.ErrNoRows) {
			return nil, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_product_by_id_sql")
		}

		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_product_by_id_sql")
	}

	return &product, nil
}

func (r *productRepository) getCategoriesSQL(ctx context.Context) ([]*entity.Category, error) {
	query, _ := r.queryLoader.Get("GetListCategories")
	rows, err := r.db0.QueryContext(ctx, query)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_list_categories_sql")
		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_list_categories_sql")
	}
	defer rows.Close()

	categories := make([]*entity.Category, 0, 10) // adjust capacity as needed
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
		); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("get_list_categories_sql_row_scan")
			return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_list_categories_sql_row_scan")
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_list_categories_sql_rows")
		return nil, x.WrapWithCode(err, x.CodeSQLRead, "get_list_categories_sql_rows")
	}

	return categories, nil
}
