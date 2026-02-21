package user

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/query"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
)

type UserRepositoryItf interface {
	RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error)
}

type userRepository struct {
	db0         *sqlx.DB
	queryLoader *query.QueryComponent
}

func InitUserRepository(db0 *sqlx.DB, queryLoader *query.QueryComponent) UserRepositoryItf {
	return &userRepository{
		db0:         db0,
		queryLoader: queryLoader,
	}
}
