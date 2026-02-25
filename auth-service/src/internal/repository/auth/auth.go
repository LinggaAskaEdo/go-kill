package auth

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/query"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthRepositoryItf interface {
	CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error)
}

type authRepository struct {
	db0         *sqlx.DB
	queryLoader *query.QueryComponent
	redis0      *redis.Client
}

func InitAuthRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent, redis0 *redis.Client) AuthRepositoryItf {
	return &authRepository{
		db0:         db0,
		queryLoader: queryLoader,
		redis0:      redis0,
	}
}
