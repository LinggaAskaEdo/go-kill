package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func (r *orderRepository) createOrderSQL(ctx context.Context, tx *sqlx.Tx, order *entity.Order) (*sqlx.Tx, *entity.Order, error) {
	query, ok := r.queryLoader.Get("CreateOrder")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "CreateOrder").Msg("query_not_found")
		return tx, order, x.NewWithCode(x.CodeSQLQueryBuild, "query_CreateOrder_not_found")
	}
	_, err := tx.ExecContext(ctx, query, order.UserID, order.OrderNumber, order.Status, order.TotalAmount, order.ShippingAddressID, order.BillingAddressID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Str("userID", order.UserID).Str("orderID", order.OrderNumber).Msg("create_order_sql")
		return tx, order, x.WrapWithCode(err, x.CodeSQLCreate, "create_order_sql")
	}

	lastIDQuery, ok := r.queryLoader.Get("GetLastInsertID")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "GetLastInsertID").Msg("query_not_found")
		return tx, order, x.NewWithCode(x.CodeSQLQueryBuild, "query_GetLastInsertID_not_found")
	}
	var lastInsertID int64
	err = tx.QueryRowxContext(ctx, lastIDQuery).Scan(&lastInsertID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Str("userID", order.UserID).Str("orderID", order.OrderNumber).Msg("get_last_insert_id")
		return tx, order, x.WrapWithCode(err, x.CodeSQLCannotRetrieveLastInsertID, "get_last_insert_id")
	}
	order.ID = fmt.Sprintf("%d", lastInsertID)

	return tx, order, nil
}

func (r *orderRepository) createOrderItemsSQL(ctx context.Context, tx *sqlx.Tx, orderID string, productDetails []*dto.ProductDetails, createOrders *dto.CreateOrderRequest) (*sqlx.Tx, error) {
	query, ok := r.queryLoader.Get("CreateOrderItemsNamed")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "CreateOrderItemsNamed").Msg("query_not_found")
		return tx, x.NewWithCode(x.CodeSQLQueryBuild, "query_CreateOrderItemsNamed_not_found")
	}

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

	result, err := tx.NamedExecContext(ctx, query, items)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("orderID", orderID).Msg("batch_create_order_items_failed")
		return tx, x.NewWithCode(x.CodeSQLCreate, "create_order_items_batch_failed")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("orderID", orderID).Msg("rows_affected_error")
		return tx, x.NewWithCode(x.CodeSQLCannotRetrieveAffectedRows, "rows_affected_error")
	}
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
	query, ok := r.queryLoader.Get("CreatePayment")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "CreatePayment").Msg("query_not_found")
		return tx, x.NewWithCode(x.CodeSQLQueryBuild, "query_CreatePayment_not_found")
	}
	result, err := tx.ExecContext(ctx, query, payment.OrderID, payment.PaymentMethod, payment.Amount)
	if err != nil {
		zerolog.Ctx(ctx).Error().Str("orderID", payment.OrderID).Msg("create_payment_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "create_payment_sql")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		zerolog.Ctx(ctx).Error().Str("orderID", payment.OrderID).Msg("create_payment_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "create_payment_sql")
	}

	return tx, nil
}

func (r *orderRepository) createStatusHistorySQL(ctx context.Context, tx *sqlx.Tx, orderID, status, note string) (*sqlx.Tx, error) {
	query, ok := r.queryLoader.Get("CreateStatusHistory")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "CreateStatusHistory").Msg("query_not_found")
		return tx, x.NewWithCode(x.CodeSQLQueryBuild, "query_CreateStatusHistory_not_found")
	}
	result, err := tx.ExecContext(ctx, query, orderID, status, note)
	if err != nil {
		zerolog.Ctx(ctx).Error().Str("orderID", orderID).Msg("create_status_history_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "create_status_history_sql")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		zerolog.Ctx(ctx).Error().Str("orderID", orderID).Msg("create_status_history_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLCreate, "create_status_history_sql")
	}

	return tx, nil
}

func (r *orderRepository) getOrderSQL(ctx context.Context, orderID string, userID string) (*entity.Order, error) {
	var order entity.Order

	query, ok := r.queryLoader.Get("GetOrder")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "GetOrder").Msg("query_not_found")
		return nil, x.NewWithCode(x.CodeSQLQueryBuild, "query_GetOrder_not_found")
	}
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

	query, ok := r.queryLoader.Get("GetOrderItem")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "GetOrderItem").Msg("query_not_found")
		return orderItems, x.NewWithCode(x.CodeSQLQueryBuild, "query_GetOrderItem_not_found")
	}
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
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_item_sql_rows")
		return orderItems, x.WrapWithCode(err, x.CodeSQLRead, "get_order_item_sql_rows")
	}

	return orderItems, nil
}

func (r *orderRepository) getOrderLimitSQL(ctx context.Context, reqData *dto.ListOrderRequest) ([]*entity.Order, error) {
	orders := make([]*entity.Order, 0, int(reqData.Limit))

	query, ok := r.queryLoader.Get("GetOrderLimit")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "GetOrderLimit").Msg("query_not_found")
		return orders, x.NewWithCode(x.CodeSQLQueryBuild, "query_GetOrderLimit_not_found")
	}
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
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_order_limit_sql_rows")
		return nil, x.WrapWithCode(err, x.CodeSQLRead, "get_order_limit_sql_rows")
	}

	return orders, nil
}

func (r *orderRepository) getOrderTotalSQL(ctx context.Context, userID string) (int32, error) {
	var total int32

	query, ok := r.queryLoader.Get("GetOrderTotal")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "GetOrderTotal").Msg("query_not_found")
		return 0, x.NewWithCode(x.CodeSQLQueryBuild, "query_GetOrderTotal_not_found")
	}
	err := r.db0.QueryRowxContext(ctx, query, userID).Scan(&total)
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
	query, ok := r.queryLoader.Get("UpdateOrderStatus")
	if !ok {
		zerolog.Ctx(ctx).Error().Str("query", "UpdateOrderStatus").Msg("query_not_found")
		return tx, x.NewWithCode(x.CodeSQLQueryBuild, "query_UpdateOrderStatus_not_found")
	}
	_, err := tx.ExecContext(ctx, query, orderID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("update_order_status_sql")
		return tx, x.WrapWithCode(err, x.CodeSQLUpdate, "update_order_status_sql")
	}

	return tx, nil
}
