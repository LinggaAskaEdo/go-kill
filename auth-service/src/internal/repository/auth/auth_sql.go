package auth

import (
	"context"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"

	"github.com/rs/zerolog"
)

func (a *authRepository) checkUserExist(ctx context.Context, email string) (bool, error) {
	var emailExists bool

	query, _ := a.queryLoader.Get("CheckUser")
	err := a.db0.QueryRowContext(ctx, query, email).Scan(&emailExists)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("check_user_exist")
		return emailExists, x.WrapWithCode(err, x.CodeSQLRowScan, "check_user_exist")
	}

	return emailExists, nil
}

func (a *authRepository) saveUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	var authID string

	query, _ := a.queryLoader.Get("SaveUSer")
	err := a.db0.QueryRowContext(ctx, query, email, string(hashedPassword)).Scan(&authID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("save_user")
		return authID, x.WrapWithCode(err, x.CodeSQLCreate, "save_user")
	}

	return authID, nil
}
