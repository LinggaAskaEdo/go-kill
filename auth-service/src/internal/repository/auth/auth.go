package auth

import (
	"context"
	"time"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/common/component/query"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthRepositoryItf interface {
	CreateAuthUser(ctx context.Context, req *dto.CreateAuthUserRequest) (string, error)
	FindAuthUserByEmail(ctx context.Context, email string) (*entity.UserAuth, error)
	StoreSession(ctx context.Context, userID string, refreshToken string, expired time.Time, email string, ipAddress string) error
	RotateRefreshToken(ctx context.Context, userID string, oldRefreshToken string, newRefreshToken string, expired time.Time) error
	FindTokenID(ctx context.Context, tokenID string) bool
	GetUserInfo(ctx context.Context, refreshToken string) (*entity.UserAuth, error)
	BlacklistToken(ctx context.Context, token *jwt.Token) error
	ClearSession(ctx context.Context, userID string) error
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
