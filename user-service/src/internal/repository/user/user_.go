package user

import (
	"context"
	"database/sql"

	x "github.com/linggaaskaedo/go-kill/common/errors"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (u *userRepository) RegisterUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx, err := u.db0.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("tx_register_user")
		return nil, err
	}

	tx, user, err = u.registerUserSQL(ctx, tx, user)
	if err != nil {
		_ = tx.Rollback()
		zerolog.Ctx(ctx).Error().Err(err).Msg("sql_register_user")
		return user, x.Wrap(err, "sql_register_user")
	}

	if err = tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("commit_register_user")
		return user, x.Wrap(err, "commit_register_user")
	}

	return user, nil
}
