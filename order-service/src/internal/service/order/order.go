package order

import (
	"context"

	productpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/product"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/order-service/src/internal/repository/order"

	"google.golang.org/grpc"
)

type OrderServiceItf interface {
	CreateOrder(ctx context.Context, reqData *dto.CreateOrderRequest) (*string, *string, float64, error)
	GetOrder(ctx context.Context, reqData *dto.GetOrderRequest) (*entity.Order, error)
	ListOrders(ctx context.Context, reqData *dto.ListOrderRequest) ([]*entity.Order, int32, error)
	CancelOrder(ctx context.Context, reqData *dto.CancelOrderRequest) error
}

type KafkaProducer interface {
	SendMessage(topic string, key, value []byte) (partition int32, offset int64, err error)
}

type orderService struct {
	orderRepository order.OrderRepositoryItf
	userClient      userpb.UserServiceClient
	productClient   productpb.ProductServiceClient
	kafkaProducer   KafkaProducer
	orderOptions    Options
}

type Options struct {
	TopicOrderCreated  string `yaml:"topic_order_created"`
	TopicOrderCanceled string `yaml:"topic_order_canceled"`
}

func InitOrderService(orderRepository order.OrderRepositoryItf, userClientConn *grpc.ClientConn, productClientConn *grpc.ClientConn, kafkaProducer KafkaProducer) OrderServiceItf {
	return &orderService{
		orderRepository: orderRepository,
		userClient:      userpb.NewUserServiceClient(userClientConn),
		productClient:   productpb.NewProductServiceClient(productClientConn),
		kafkaProducer:   kafkaProducer,
	}
}
