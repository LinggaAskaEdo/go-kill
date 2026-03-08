package grpc

import (
	orderpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/order"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/service"

	"github.com/rs/zerolog"
)

type Grpc struct {
	orderpb.OrderServiceServer
	log zerolog.Logger
	svc *service.Service
}

func InitGrpcHandler(log zerolog.Logger, svc *service.Service) *Grpc {
	return &Grpc{
		log: log,
		svc: svc,
	}
}

func itemsToDTO(original []*orderpb.OrderItem) []*dto.OrderItem {
	result := make([]*dto.OrderItem, len(original))

	for i, item := range original {
		if item != nil {
			result[i] = &dto.OrderItem{
				ProductID: item.ProductId,
				Quantity:  item.Quantity,
			}
		}
	}

	return result
}

func itemsToPB(original []*entity.OrderItem) []*orderpb.OrderItemDetail {
	result := make([]*orderpb.OrderItemDetail, len(original))

	for i, item := range original {
		if item != nil {
			result[i] = &orderpb.OrderItemDetail{
				Id:          item.ID,
				ProductId:   item.ProductID,
				ProductName: item.ProductName,
				Quantity:    int32(item.Quantity),
				UnitPrice:   item.UnitPrice,
				Subtotal:    item.Subtotal,
			}
		}
	}

	return result
}

func ordersToPB(original []*entity.Order) []*orderpb.GetOrderResponse {
	result := make([]*orderpb.GetOrderResponse, len(original))

	for i, item := range original {
		if item != nil {
			result[i] = &orderpb.GetOrderResponse{
				Id:          item.ID,
				OrderNumber: item.OrderNumber,
				Status:      string(item.Status),
				TotalAmount: item.TotalAmount,
			}
		}
	}

	return result
}
