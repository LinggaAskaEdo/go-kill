package user

import (
	"context"
	"database/sql"
	"time"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
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
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			zerolog.Ctx(ctx).Error().Err(rbErr).Msg("rollback_register_user")
		}
		return user, err
	}

	tx, err = u.registerUserProfileSQL(ctx, tx, user.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			zerolog.Ctx(ctx).Error().Err(rbErr).Msg("rollback_register_user_profile")
		}
		return user, err
	}

	if err = tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("commit_register_user")
		return user, x.Wrap(err, "commit_register_user")
	}

	ip, _ := ctx.Value("ip").(string)
	ua, _ := ctx.Value("user_agent").(string)

	metadata := map[string]any{
		"ip_address":          ip,
		"user_agent":          ua,
		"registration_method": "email",
	}

	mongoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = u.logActivityMongo(mongoCtx, user.ID, "registration", metadata); err != nil {
		zerolog.Ctx(ctx).Warn().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to log registration activity (non-fatal)")
	}

	return user, nil
}

func (u *userRepository) GetMe(ctx context.Context, userAuthID string) (*entity.User, error) {
	return u.getUserByAuthIDSQL(ctx, userAuthID)
}

func (u *userRepository) GetActivities(ctx context.Context, userID string, page string, limit string) ([]*entity.UserActivity, int64, error) {
	return u.getUserActivitiesMongo(ctx, userID, page, limit)
}

func (u *userRepository) GetUserAddresses(ctx context.Context, userID string, page string, limit string) ([]*entity.UserAddress, int64, error) {
	return u.getUserAddressesByUserIDSQL(ctx, userID, page, limit)
}

func (u *userRepository) CreateAddress(ctx context.Context, userID string, req dto.CreateUserAddress) (string, error) {
	return u.createUserAddressSQL(ctx, userID, req)
}
