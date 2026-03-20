package product

import (
	"context"
	"database/sql"
	"errors"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
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
		return tx, nil
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
	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("productID", productID).Int("qty", qty).Int("rsv", rsv).Msg("create_product_inventory_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "create_product_inventory_sql")
	}
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

	categories := make([]*entity.Category, 0, 10)
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

func (r *productRepository) getCategoriesByProductIDSQL(ctx context.Context, productID string) ([]*entity.Category, error) {
	query, _ := r.queryLoader.Get("GetCategoriesByProductID")
	rows, err := r.db0.QueryContext(ctx, query, productID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_categories_by_product_id_sql")
		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_categories_by_product_id_sql")
	}
	defer rows.Close()

	categories := make([]*entity.Category, 0, 10)
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
		); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("get_categories_by_product_id_sql")
			return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_categories_by_product_id_sql")
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_categories_by_product_id_sql")
		return nil, x.WrapWithCode(err, x.CodeSQLRead, "get_categories_by_product_id_sql")
	}

	return categories, nil
}

func (r *productRepository) getProductsByCategoryIDSQL(ctx context.Context, categoryID string) ([]*entity.Product, error) {
	query, _ := r.queryLoader.Get("GetProductsByCategoryID")
	rows, err := r.db0.QueryContext(ctx, query, categoryID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_products_by_category_id_sql")
		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_products_by_category_id_sql")
	}
	defer rows.Close()

	products := make([]*entity.Product, 0, 3)
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.SKU,
		); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("get_products_by_category_id_sql")
			return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_products_by_category_id_sql")
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_products_by_category_id_sql")
		return nil, x.WrapWithCode(err, x.CodeSQLRead, "get_products_by_category_id_sql")
	}

	return products, nil
}

func (r *productRepository) getInventoryByProductIDSQL(ctx context.Context, productID string) (int32, int32, error) {
	var quantity, reserved int32

	query, _ := r.queryLoader.Get("GetInventoryByProductID")
	err := r.db0.QueryRowContext(ctx, query, productID).Scan(&quantity, &reserved)
	if err != nil {
		zerolog.Ctx(ctx).Error().Str("id", productID).Msg("get_inventory_by_product_id_sql")
		return quantity, reserved, x.WrapWithCode(err, x.CodeSQLRead, "get_inventory_by_product_id_sql")
	}

	return quantity, reserved, nil
}

func (r *productRepository) createReserveInventorySQL(ctx context.Context, tx *sqlx.Tx, req []dto.CreateReserveInventory) (*sqlx.Tx, error) {
	query0, _ := r.queryLoader.Get("LockUpdateInventory")
	query1, _ := r.queryLoader.Get("UpdateReservedQuantity")

	for _, item := range req {
		var quantity, reserved int32

		// Lock and update inventory
		err := tx.QueryRowxContext(ctx, query0, item.ProductId).Scan(&quantity, &reserved)
		if err != nil {
			zerolog.Ctx(ctx).Error().Str("productID", item.ProductId).Msg("create_reserve_inventory_sql")
			return tx, x.WrapWithCode(err, x.CodeSQLUpdate, "create_reserve_inventory_sql")
		}

		available := quantity - reserved
		status := available < int32(item.Quantity)
		if status {
			zerolog.Ctx(ctx).Error().Bool("status", status).Msg("create_reserve_inventory_sql")
			return tx, x.New("Insufficient inventory")
		}

		// Update reserved quantity
		_, err = tx.ExecContext(ctx, query1, item.Quantity, item.ProductId)
		if err != nil {
			zerolog.Ctx(ctx).Error().Str("productID", item.ProductId).Int32("qty", item.Quantity).Msg("create_reserve_inventory_sql")
			return tx, x.WrapWithCode(err, x.CodeSQLUpdate, "create_reserve_inventory_sql")
		}

		r.setInventoryReservedCache(ctx, item.ProductId, item.Quantity)
	}

	return tx, nil
}

func (r *productRepository) createReleaseInventorySQLTx(ctx context.Context, tx *sqlx.Tx, req []dto.CreateReserveInventory) (*sqlx.Tx, error) {
	query, _ := r.queryLoader.Get("UpdateReleaseQuantity")

	for _, item := range req {
		var err error
		if tx != nil {
			_, err = tx.ExecContext(ctx, query, item.Quantity, item.ProductId)
		} else {
			_, err = r.db0.ExecContext(ctx, query, item.Quantity, item.ProductId)
		}
		if err != nil {
			zerolog.Ctx(ctx).Error().Str("productID", item.ProductId).Int32("qty", item.Quantity).Msg("create_release_inventory_sql")
			return tx, x.WrapWithCode(err, x.CodeSQLUpdate, "create_release_inventory_sql")
		}

		r.setInventoryReleaseCache(ctx, item.ProductId, item.Quantity)
	}

	return tx, nil
}
