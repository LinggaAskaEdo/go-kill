package user

import (
	"context"

	x "github.com/linggaaskaedo/go-kill/common/errors"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func (u *userRepository) registerUserSQL(ctx context.Context, tx *sqlx.Tx, user *entity.User) (*sqlx.Tx, *entity.User, error) {
	query, _ := u.queryLoader.Get("RegisterUser")
	row := tx.QueryRowContext(ctx, query, user.AutdID, user.Email, user.FirstName, user.LastName).Scan(&user.ID)
	if err := row; err != nil {
		return tx, user, x.Wrap(err, "register_user_sql")
	}

	query, _ = u.queryLoader.Get("RegisterUserProfile")
	result := tx.MustExecContext(ctx, query, user.ID)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		zerolog.Ctx(ctx).Error().Str("id", user.ID).Msg("register_user_profile_err")
		return tx, user, x.NewWithCode(x.CodeSQLCreate, "register_user_profile_err")
	}

	return tx, user, nil
}
