package user

import (
	"context"
	"database/sql"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
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
		return user, x.Wrap(err, "sql_register_user")
	}

	if err = tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("commit_register_user")
		return user, x.Wrap(err, "commit_register_user")
	}

	u.logActivityMongo(ctx, user.ID, "registration", map[string]any{
		"ip_address":          ctx.Value("ip").(string),
		"user_agent":          ctx.Value("user_agent").(string),
		"registration_method": "email",
	})

	return user, nil
}
