package user

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/query"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepositoryItf interface {
	// gRPC
	CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error)
	GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error)
	GetAddress(ctx context.Context, req *userpb.GetAddressRequest) (*userpb.GetAddressResponse, error)
	LogActivity(ctx context.Context, req *userpb.LogActivityRequest) (*userpb.LogActivityResponse, error)

	// REST
	RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetMe(ctx context.Context, userAuthID string) (*entity.User, error)
	GetActivities(ctx context.Context, userID string, page string, limit string) ([]*entity.UserActivity, int64, error)
	GetUserAddresses(ctx context.Context, userID string, page string, limit string) ([]*entity.UserAddress, int64, error)
	CreateAddress(ctx context.Context, userID string, req dto.CreateUserAddress) (string, error)
}

type userRepository struct {
	db0         *sqlx.DB
	queryLoader *query.QueryComponent
	mongo0      *mongo.Database
}

func InitUserRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent, mongo0 *mongo.Database) UserRepositoryItf {
	return &userRepository{
		db0:         db0,
		queryLoader: queryLoader,
		mongo0:      mongo0,
	}
}

func convertMetadata(m map[string]string) map[string]any {
	result := make(map[string]any)

	for k, v := range m {
		result[k] = v
	}

	return result
}
