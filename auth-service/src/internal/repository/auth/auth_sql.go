package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/entity"
	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"

	"github.com/rs/zerolog"
)

func (a *authRepository) checkUserExistSql(ctx context.Context, email string) (bool, error) {
	var emailExists bool

	query, _ := a.queryLoader.Get("CheckUser")
	err := a.db0.QueryRowContext(ctx, query, email).Scan(&emailExists)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("check_user_exist_sql")

		if err == sql.ErrNoRows {
			return emailExists, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "check_user_exist_sql")
		}

		return emailExists, x.WrapWithCode(err, x.CodeSQLRowScan, "check_user_exist_sql")
	}

	return emailExists, nil
}

func (a *authRepository) saveUserSql(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	var authID string

	query, _ := a.queryLoader.Get("SaveUSer")
	err := a.db0.QueryRowContext(ctx, query, email, string(hashedPassword)).Scan(&authID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("save_user_sql")
		return authID, x.WrapWithCode(err, x.CodeSQLCreate, "save_user_sql")
	}

	return authID, nil
}

func (a *authRepository) getUserByEmailSql(ctx context.Context, email string) (*entity.UserAuth, error) {
	var userAuth entity.UserAuth

	query, _ := a.queryLoader.Get("GetUserByEmail")
	err := a.db0.QueryRowxContext(ctx, query, email).StructScan(&userAuth)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_by_email_sql")

		if err == sql.ErrNoRows {
			return &userAuth, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_user_by_email_sql")
		}

		return &userAuth, x.WrapWithCode(err, x.CodeSQLRead, "get_user_by_email_sql")
	}

	return &userAuth, nil
}

func (a *authRepository) storeRefreshTokenSql(ctx context.Context, userID string, refreshToken string, expired time.Time) error {
	query, _ := a.queryLoader.Get("StoreRefreshToken")
	_, err := a.db0.ExecContext(ctx, query, userID, refreshToken, expired)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("store_refresh_token_sql")
		return x.WrapWithCode(err, x.CodeSQLCreate, "store_refresh_token_sql")
	}

	return nil
}

func (a *authRepository) getUserByIDSql(ctx context.Context, userID string) (*entity.UserAuth, error) {
	var userAuth entity.UserAuth

	query, _ := a.queryLoader.Get("GetUserWithID")
	err := a.db0.QueryRowContext(ctx, query, userID).Scan(&userAuth)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_by_id_sql")

		if err == sql.ErrNoRows {
			return &userAuth, x.WrapWithCode(err, x.CodeSQLRecordDoesNotExist, "get_user_by_id_sql")
		}

		return &userAuth, x.WrapWithCode(err, x.CodeSQLRowScan, "get_user_by_id_sql")
	}

	return &userAuth, nil
}

func (a *authRepository) deleteRefreshTokenSql(ctx context.Context, userID string) error {
	query, _ := a.queryLoader.Get("DeleteRefreshToken")
	_, err := a.db0.ExecContext(ctx, query, userID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("delete_refresh_token_sql")
		return x.WrapWithCode(err, x.CodeSQLRowScan, "delete_refresh_token_sql")
	}

	return nil
}
