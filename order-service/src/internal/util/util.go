package util

import (
	orderpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/order"
	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"
)

func ToOrderItemDTO(original []*orderpb.OrderItem) []*dto.OrderItem {
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

func ToOrderItemDetailPB(original []*entity.OrderItem) []*orderpb.OrderItemDetail {
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

func ToGetOrderResponsePB(original []*entity.Order) []*orderpb.GetOrderResponse {
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

func ToInventoryItemPB(original []*entity.OrderItem) []*productpb.InventoryItem {
	result := make([]*productpb.InventoryItem, len(original))

	for i, item := range original {
		if item != nil {
			result[i] = &productpb.InventoryItem{
				ProductId: item.ProductID,
				Quantity:  int32(item.Quantity),
			}
		}
	}

	return result
}
