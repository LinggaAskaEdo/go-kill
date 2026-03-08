package grpc

import (
	"context"

	orderpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/order"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
)

func (g *Grpc) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	reqData := &dto.CreateOrderRequest{
		UserID:            req.UserId,
		ShippingAddressID: req.ShippingAddressId,
		BillingAddressID:  req.BillingAddressId,
		PaymentMethod:     req.PaymentMethod,
		Items:             itemsToDTO(req.Items),
	}

	orderID, orderNumber, totalAmount, err := g.svc.Order.CreateOrder(ctx, reqData)
	if err != nil {
		return nil, err
	}

	return &orderpb.CreateOrderResponse{
		Success:     true,
		OrderId:     *orderID,
		OrderNumber: *orderNumber,
		TotalAmount: totalAmount,
	}, nil
}

func (g *Grpc) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	reqData := &dto.GetOrderRequest{
		OrderID: req.OrderId,
		UserID:  req.UserId,
	}

	order, err := g.svc.Order.GetOrder(ctx, reqData)
	if err != nil {
		return nil, err
	}

	return &orderpb.GetOrderResponse{
		Id:          order.ID,
		OrderNumber: order.OrderNumber,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		Items:       itemsToPB(order.Items),
	}, nil
}

func (g *Grpc) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	offset := (req.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	reqData := &dto.ListOrderRequest{
		UserID: req.UserId,
		Limit:  limit,
		Offset: offset,
	}

	orders, total, err := g.svc.Order.ListOrders(ctx, reqData)
	if err != nil {
		return nil, err
	}

	return &orderpb.ListOrdersResponse{
		Orders: ordersToPB(orders),
		Total:  total,
	}, nil
}

func (g *Grpc) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	return nil, nil
}
