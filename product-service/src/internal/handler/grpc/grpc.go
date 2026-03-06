package grpc

import (
	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service"

	"github.com/rs/zerolog"
)

type Grpc struct {
	productpb.ProductServiceServer
	log zerolog.Logger
	svc *service.Service
}

func InitGrpcHandler(log zerolog.Logger, svc *service.Service) *Grpc {
	return &Grpc{
		log: log,
		svc: svc,
	}
}

func convertItems(original []*productpb.InventoryItem) []dto.CreateReserveInventory {
	result := make([]dto.CreateReserveInventory, len(original))

	for i, item := range original {
		if item != nil {
			result[i] = dto.CreateReserveInventory{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			}
		}
	}

	return result
}
