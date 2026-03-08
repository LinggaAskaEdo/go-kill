package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/util"

	"github.com/rs/zerolog"
)

func (r *orderRepository) StoreOrder(ctx context.Context, productDetails []*dto.ProductDetails, createOrders *dto.CreateOrderRequest, totalAmount float64) (*string, *string, error) {
	// Step 4: Reserve inventory
	var inventoryItems []*productpb.InventoryItem
	for _, item := range createOrders.Items {
		inventoryItems = append(inventoryItems, &productpb.InventoryItem{ProductId: item.ProductID, Quantity: item.Quantity})
	}

	reserveResp, err := r.productClient.ReserveInventory(ctx, &productpb.ReserveInventoryRequest{Items: inventoryItems})
	if err != nil || !reserveResp.Success {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to reserve inventory")
		return nil, nil, x.New("Failed to reserve inventory", err)
	}

	// Step 5: Create order in MySQL
	tx, err := r.db0.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("tx_store_order")
		r.doReleaseInventory(ctx, inventoryItems)

		return nil, nil, err
	}

	// Insert order
	orderNumber := fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%1000000)
	order := &entity.Order{
		UserID:            createOrders.UserID,
		OrderNumber:       orderNumber,
		Status:            entity.StatusPending,
		TotalAmount:       totalAmount,
		ShippingAddressID: &createOrders.ShippingAddressID,
		BillingAddressID:  &createOrders.ShippingAddressID,
	}

	tx, order, err = r.createOrderSQL(ctx, tx, order)
	if err != nil {
		_ = tx.Rollback()
		r.doReleaseInventory(ctx, inventoryItems)

		return nil, nil, err
	}

	// Insert order items (one-to-many relationship)
	tx, err = r.createOrderItemsSQL(ctx, tx, order.ID, productDetails, createOrders)
	if err != nil {
		_ = tx.Rollback()
		r.doReleaseInventory(ctx, inventoryItems)

		return nil, nil, err
	}

	// Insert payment record
	payment := &entity.Payment{
		OrderID:       order.ID,
		PaymentMethod: createOrders.PaymentMethod,
		Amount:        totalAmount,
	}

	tx, err = r.createPaymentSQL(ctx, tx, payment)
	if err != nil {
		_ = tx.Rollback()
		r.doReleaseInventory(ctx, inventoryItems)

		return nil, nil, err
	}

	// Insert status history
	tx, err = r.createStatusHistorySQL(ctx, tx, order.ID, string(entity.StatusPending), "Order created")
	if err != nil {
		_ = tx.Rollback()
		r.doReleaseInventory(ctx, inventoryItems)

		return nil, nil, err
	}

	if err = tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("commit_store_order")

		return nil, nil, x.Wrap(err, "commit_store_order")
	}

	return &order.ID, &orderNumber, nil
}

func (r *orderRepository) doReleaseInventory(ctx context.Context, inventoryItems []*productpb.InventoryItem) {
	_, err := r.productClient.ReleaseInventory(ctx, &productpb.ReleaseInventoryRequest{Items: inventoryItems})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to release inventory")
	}
}

func (r *orderRepository) GetOrder(ctx context.Context, reqData *dto.GetOrderRequest) (*entity.Order, error) {
	order, err := r.getOrderSQL(ctx, reqData.OrderID, reqData.UserID)
	if err != nil {
		return nil, err
	}

	orderItems, err := r.getOrderItemSQL(ctx, reqData.OrderID)
	if err != nil {
		return nil, err
	}

	order.Items = orderItems

	return order, nil
}

func (r *orderRepository) ListOrders(ctx context.Context, reqData *dto.ListOrderRequest) ([]*entity.Order, int32, error) {
	orders, err := r.getOrderLimitSQL(ctx, reqData)
	if err != nil {
		return nil, 0, err
	}

	total, err := r.getOrderTotalSQL(ctx, reqData.UserID)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) CancelOrder(ctx context.Context, reqData *dto.CancelOrderRequest) error {
	order, err := r.getOrderSQL(ctx, reqData.OrderID, reqData.UserID)
	if err != nil {
		return err
	}

	if order.Status != "pending" && order.Status != "confirmed" {
		return x.New("Order cannot be cancelled")
	}

	orderItems, err := r.getOrderItemSQL(ctx, reqData.OrderID)
	if err != nil {
		return err
	}

	tx, err := r.db0.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("tx_cancel_order")
		return err
	}

	// Update order status
	tx, err = r.updateOrderStatusSQL(ctx, tx, reqData.OrderID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Add status history
	tx, err = r.createStatusHistorySQL(ctx, tx, order.ID, string(entity.StatusCancelled), reqData.Reason)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Release inventory
	_, err = r.productClient.ReleaseInventory(ctx, &productpb.ReleaseInventoryRequest{
		Items: util.ToInventoryItemPB(orderItems),
	})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to release inventory")
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("commit_cancel_order")
		return x.Wrap(err, "commit_cancel_order")
	}

	return nil
}
