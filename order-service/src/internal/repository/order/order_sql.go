package order

import (
	"context"
	"database/sql"
	"errors"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func (r *orderRepository) createOrderSQL(ctx context.Context, tx *sqlx.Tx, order *entity.Order) (*sqlx.Tx, *entity.Order, error) {
	query, _ := r.queryLoader.Get("CreateOrder")
	err := tx.QueryRowContext(ctx, query, order.UserID, order.OrderNumber, order.Status, order.TotalAmount, order.ShippingAddressID, order.BillingAddressID).Scan(&order.ID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Str("userID", order.UserID).Str("orderID", order.OrderNumber).Msg("create_product_sql")
		return tx, order, x.WrapWithCode(err, x.CodeSQLCreate, "create_product_sql")
	}

	return tx, order, nil
}

func (r *orderRepository) createOrderItemsSQL(ctx context.Context, tx *sqlx.Tx, orderID string, productDetails []*dto.ProductDetails, createOrders *dto.CreateOrderRequest) (*sqlx.Tx, error) {
	query, _ := r.queryLoader.Get("CreateOrderItemsNamed")

	// Build the batch data slice
	items := make([]map[string]any, len(createOrders.Items))
	for i, item := range createOrders.Items {
		subtotal := productDetails[i].Price * float64(item.Quantity)
		items[i] = map[string]any{
			"order_id":     orderID,
			"product_id":   item.ProductID,
			"product_name": productDetails[i].Name,
			"quantity":     item.Quantity,
			"unit_price":   productDetails[i].Price,
			"subtotal":     subtotal,
		}
	}

	// Execute batch insert
	result, err := tx.NamedExecContext(ctx, query, items)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("orderID", orderID).Msg("batch_create_order_items_failed")
		return tx, x.NewWithCode(x.CodeSQLCreate, "create_order_items_batch_failed")
	}

	// Verify rows affected match expected count
	rowsAffected, _ := result.RowsAffected()
	if int(rowsAffected) != len(items) {
		zerolog.Ctx(ctx).Error().
			Str("orderID", orderID).
			Int64("expected", int64(len(items))).
			Int64("actual", rowsAffected).
			Msg("rows_affected_mismatch")

		return tx, x.NewWithCode(x.CodeSQLCreate, "create_order_items_sql")
	}

	return tx, nil
}

func (r *orderRepository) createPaymentSQL(ctx context.Context, tx *sqlx.Tx, payment *entity.Payment) (*sqlx.Tx, error) {
	query, _ := r.queryLoader.Get("CreatePayment")
	rows, err := tx.MustExecContext(ctx, query, payment.OrderID, payment.PaymentMethod, payment.Amount).RowsAffected()
	if rows == 0 || err != nil {
		zerolog.Ctx(ctx).Error().Str("orderID", payment.OrderID).Msg("create_payment_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "create_payment_sql")
	}

	return tx, nil
}

func (r *orderRepository) createStatusHistorySQL(ctx context.Context, tx *sqlx.Tx, orderID, status, note string) (*sqlx.Tx, error) {
	query, _ := r.queryLoader.Get("CreateStatusHistory")
	rows, err := tx.MustExecContext(ctx, query, orderID, status, note).RowsAffected()
	if rows == 0 || err != nil {
		zerolog.Ctx(ctx).Error().Str("orderID", orderID).Msg("create_status_history_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "create_status_history_sql")
	}

	return tx, nil
}

func (r *orderRepository) getOrderSQL(ctx context.Context, orderID string, userID string) (*entity.Order, error) {
	var order entity.Order

	query, _ := r.queryLoader.Get("GetOrder")
	err := r.db0.QueryRowxContext(ctx, query, orderID, userID).StructScan(&order)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_sql")

		if errors.Is(err, sql.ErrNoRows) {
			return nil, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_order_sql")
		}

		return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_order_sql")
	}

	return &order, nil
}

func (r *orderRepository) getOrderItemSQL(ctx context.Context, orderID string) ([]*entity.OrderItem, error) {
	orderItems := make([]*entity.OrderItem, 0, 3)

	query, _ := r.queryLoader.Get("GetOrderItem")
	rows, err := r.db0.QueryContext(ctx, query, orderID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_item_sql")

		if errors.Is(err, sql.ErrNoRows) {
			return orderItems, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_order_item_sql")
		}

		return orderItems, x.WrapWithCode(err, x.CodeSQLRowScan, "get_order_item_sql")
	}
	defer rows.Close()

	for rows.Next() {
		var orderItem entity.OrderItem
		if err := rows.Scan(
			&orderItem.ID,
			&orderItem.ProductID,
			&orderItem.ProductName,
			&orderItem.Quantity,
			&orderItem.UnitPrice,
			&orderItem.Subtotal,
		); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_item_sql_row_scan")
			return orderItems, x.WrapWithCode(err, x.CodeSQLRowScan, "get_order_item_sql_row_scan")
		}

		orderItems = append(orderItems, &orderItem)
		// order.Items = append(order.Items, &orderItem)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_item_sql_rows")
		return orderItems, x.WrapWithCode(err, x.CodeSQLRead, "get_order_item_sql_rows")
	}

	return orderItems, nil
}

func (r *orderRepository) getOrderLimitSQL(ctx context.Context, reqData *dto.ListOrderRequest) ([]*entity.Order, error) {
	orders := make([]*entity.Order, 0, 3)

	query, _ := r.queryLoader.Get("GetOrderLimit")
	rows, err := r.db0.QueryContext(ctx, query, reqData.UserID, reqData.Limit, reqData.Offset)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_limit_sql")

		if errors.Is(err, sql.ErrNoRows) {
			return orders, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_order_limit_sql")
		}

		return orders, x.WrapWithCode(err, x.CodeSQLRowScan, "get_order_limit_sql")
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&order.Status,
			&order.TotalAmount,
		); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_limit_sql_row_scan")
			return nil, x.WrapWithCode(err, x.CodeSQLRowScan, "get_order_limit_sql_row_scan")
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_limti_sql_rows")
		return nil, x.WrapWithCode(err, x.CodeSQLRead, "get_order_limit_sql_rows")
	}

	return orders, nil
}

func (r *orderRepository) getOrderTotalSQL(ctx context.Context, UserID string) (int32, error) {
	var total int32

	query, _ := r.queryLoader.Get("GetOrderTotal")
	err := r.db0.QueryRowxContext(ctx, query, UserID).Scan(&total)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_total_sql")

		if errors.Is(err, sql.ErrNoRows) {
			return 0, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_order_total_sql")
		}

		return 0, x.WrapWithCode(err, x.CodeSQLRowScan, "get_order_total_sql")
	}

	return total, nil
}

func (r *orderRepository) updateOrderStatusSQL(ctx context.Context, tx *sqlx.Tx, orderID string) (*sqlx.Tx, error) {
	query, _ := r.queryLoader.Get("UpdateOrderStatus")
	_, err := tx.ExecContext(ctx, query, orderID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_update_order_status_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLUpdate, "get_update_order_status_sql")
	}

	return tx, nil
}
