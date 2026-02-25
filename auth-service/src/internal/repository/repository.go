package repository

import (
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/repository/auth"
	"github.com/linggaaskaedo/go-kill/common/component/query"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Auth auth.AuthRepositoryItf
}

func InitRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent, redis0 *redis.Client) *Repository {
	return &Repository{
		Auth: auth.InitAuthRepository(
			db0,
			queryLoader,
			redis0,
		),
	}
}
