package user

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/query"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepositoryItf interface {
	RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error)
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
