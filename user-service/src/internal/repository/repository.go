package repository

import (
	"github.com/linggaaskaedo/go-kill/common/query"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository/user"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User user.UserRepositoryItf
}

func InitRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent) *Repository {
	return &Repository{
		User: user.InitUserRepository(
			db0,
			queryLoader,
		),
	}
}
