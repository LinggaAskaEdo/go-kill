package grpc

import (
	"context"

	orderpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/order"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/util"
)

func (g *Grpc) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	reqData := &dto.CreateOrderRequest{
		UserID:            req.UserId,
		ShippingAddressID: req.ShippingAddressId,
		BillingAddressID:  req.BillingAddressId,
		PaymentMethod:     req.PaymentMethod,
		Items:             util.ToOrderItemDTO(req.Items),
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
		Items:       util.ToOrderItemDetailPB(order.Items),
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
		Orders: util.ToGetOrderResponsePB(orders),
		Total:  total,
	}, nil
}

func (g *Grpc) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	reqData := &dto.CancelOrderRequest{
		OrderID: req.OrderId,
		UserID:  req.UserId,
		Reason:  req.Reason,
	}

	err := g.svc.Order.CancelOrder(ctx, reqData)
	if err != nil {
		return nil, err
	}

	return &orderpb.CancelOrderResponse{Success: true}, nil
}
