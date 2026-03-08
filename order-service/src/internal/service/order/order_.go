package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"

	"github.com/openpcc/openpcc/uuidv7"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

func (s *orderService) CreateOrder(ctx context.Context, reqData *dto.CreateOrderRequest) (*string, *string, float64, error) {
	// Step 1: Validate user
	userResp, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{UserId: reqData.UserID})
	if err != nil || !userResp.Found {
		zerolog.Ctx(ctx).Error().Err(err).Msg("create_order")
		return nil, nil, 0, x.New("User not found", err)
	}

	// Step 2: Validate shipping address
	addressResp, err := s.userClient.GetAddress(ctx, &userpb.GetAddressRequest{
		AddressId: reqData.ShippingAddressID,
		UserId:    reqData.UserID,
	})
	if err != nil || !addressResp.Found {
		zerolog.Ctx(ctx).Error().Err(err).Msg("create_order")
		return nil, nil, 0, x.New("Invalid shipping address", err)
	}

	// Step 3: Validate products and calculate total
	productDetails, totalAmount, err := s.validateAndCollectProducts(ctx, reqData.Items)
	if err != nil {
		return nil, nil, 0, err
	}

	// Step 4: Create order in MySQL
	orderID, orderNumber, err := s.orderRepository.StoreOrder(ctx, productDetails, reqData, totalAmount)
	if err != nil {
		return nil, nil, 0, err
	}

	// Step 5: Publish event to Kafka (non‑critical, errors are logged)
	if err := s.publishOrderEvent(ctx, orderID, orderNumber, reqData, userResp, productDetails, totalAmount); err != nil {
		return nil, nil, 0, err
	}

	return orderID, orderNumber, totalAmount, nil
}

func (s *orderService) validateAndCollectProducts(ctx context.Context, items []*dto.OrderItem) ([]*dto.ProductDetails, float64, error) {
	if len(items) == 0 {
		return nil, 0, x.New("no items in order", nil)
	}

	productDetails := make([]*dto.ProductDetails, len(items))
	g, ctx := errgroup.WithContext(ctx)

	for i, item := range items {
		g.Go(func() error {
			details, err := s.fetchAndCheckProduct(ctx, item)
			if err != nil {
				return err
			}

			productDetails[i] = details
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	totalAmount := 0.0
	for i, item := range items {
		totalAmount += productDetails[i].Price * float64(item.Quantity)
	}

	return productDetails, totalAmount, nil
}

func (s *orderService) fetchAndCheckProduct(ctx context.Context, item *dto.OrderItem) (*dto.ProductDetails, error) {
	prodResp, err := s.productClient.GetProduct(ctx, &productpb.GetProductRequest{ProductId: item.ProductID})
	if err != nil {
		return nil, x.New(fmt.Sprintf("failed to get product %s", item.ProductID), err)
	}
	if !prodResp.Found {
		return nil, x.New(fmt.Sprintf("product %s not found", item.ProductID), nil)
	}

	invResp, err := s.productClient.CheckInventory(ctx, &productpb.CheckInventoryRequest{
		ProductId: item.ProductID,
		Quantity:  item.Quantity,
	})
	if err != nil {
		return nil, x.New(fmt.Sprintf("inventory check failed for %s", prodResp.Name), err)
	}
	if !invResp.Available {
		return nil, x.New(fmt.Sprintf("insufficient inventory for %s", prodResp.Name), nil)
	}

	return &dto.ProductDetails{
		ID:    prodResp.Id,
		Name:  prodResp.Name,
		Price: prodResp.Price,
	}, nil
}

func (s *orderService) publishOrderEvent(
	ctx context.Context,
	orderID, orderNumber *string,
	reqData *dto.CreateOrderRequest,
	userResp *userpb.GetUserResponse,
	productDetails []*dto.ProductDetails,
	totalAmount float64,
) error {
	orderItems := make([]dto.OrderItem, len(reqData.Items))

	for i, item := range reqData.Items {
		orderItems[i] = dto.OrderItem{
			ProductID:   item.ProductID,
			ProductName: productDetails[i].Name,
			Quantity:    int32(item.Quantity),
			UnitPrice:   productDetails[i].Price,
		}
	}

	event := dto.OrderEvent{
		EventID:   uuidv7.MustNew().String(),
		EventType: "order.created",
		Version:   "1.0",
		Timestamp: time.Now(),
		Source:    "order-service",
		Data: dto.OrderData{
			OrderID:     *orderID,
			OrderNumber: *orderNumber,
			UserID:      reqData.UserID,
			UserEmail:   userResp.Email,
			TotalAmount: totalAmount,
			Status:      "pending",
			Items:       orderItems,
		},
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return x.New("Failed to marshal event", err)
	}

	partition, offset, err := s.kafkaProducer.SendMessage(s.orderOptions.TopicOrderCreated, []byte(event.Data.UserID), eventBytes)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to send Kafka message")
	} else {
		zerolog.Ctx(ctx).Debug().Int32("partition", partition).Int64("offset", offset).Msg("Kafka message sent")
	}

	return nil
}

func (s *orderService) GetOrder(ctx context.Context, reqData *dto.GetOrderRequest) (*entity.Order, error) {
	return s.orderRepository.GetOrder(ctx, reqData)
}

func (s *orderService) ListOrders(ctx context.Context, reqData *dto.ListOrderRequest) ([]*entity.Order, int32, error) {
	return s.orderRepository.ListOrders(ctx, reqData)
}
