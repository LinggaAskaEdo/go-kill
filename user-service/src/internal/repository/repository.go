package repository

import (
	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository/user"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	User user.UserRepositoryItf
}

func InitRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent, mongo0 *mongo.Database) *Repository {
	return &Repository{
		User: user.InitUserRepository(
			db0,
			queryLoader,
			mongo0,
		),
	}
}
