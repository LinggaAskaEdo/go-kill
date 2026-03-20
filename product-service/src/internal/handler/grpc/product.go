package grpc

import (
	"context"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
)

func (g *Grpc) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	if req.ProductId == "" {
		return nil, x.New("product ID is required")
	}

	resp, err := g.svc.Product.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &productpb.GetProductResponse{
		Id:          resp.ID,
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
		Sku:         resp.SKU,
		IsActive:    resp.IsActive,
	}, nil
}

func (g *Grpc) CheckInventory(ctx context.Context, req *productpb.CheckInventoryRequest) (*productpb.CheckInventoryResponse, error) {
	if req.ProductId == "" {
		return nil, x.New("product ID is required")
	}

	quantity, reserved, err := g.svc.Product.CheckInventory(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	available := quantity - reserved
	isAvailable := available >= int32(req.Quantity)

	return &productpb.CheckInventoryResponse{
		Available:        isAvailable,
		CurrentQuantity:  int32(quantity),
		ReservedQuantity: int32(reserved),
	}, nil
}

func (g *Grpc) ReserveInventory(ctx context.Context, req *productpb.ReserveInventoryRequest) (*productpb.ReserveInventoryResponse, error) {
	if req.Items == nil {
		return nil, x.New("Item is empty")
	}

	err := g.svc.Product.ReserveInventory(ctx, convertItems(req.Items))
	if err != nil {
		return nil, err
	}

	return &productpb.ReserveInventoryResponse{Success: true}, nil
}

func (g *Grpc) ReleaseInventory(ctx context.Context, req *productpb.ReleaseInventoryRequest) (*productpb.ReleaseInventoryResponse, error) {
	if req.Items == nil {
		return nil, x.New("Item is empty")
	}

	err := g.svc.Product.ReleaseInventory(ctx, convertItems(req.Items))
	if err != nil {
		return nil, err
	}

	return &productpb.ReleaseInventoryResponse{Success: true}, nil
}
