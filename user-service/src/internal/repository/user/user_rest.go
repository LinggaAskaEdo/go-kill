package user

import (
	"context"
	"database/sql"

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
		_ = tx.Rollback()
		return user, err
	}

	tx, err = u.registerUserProfileSQL(ctx, tx, user.ID)
	if err != nil {
		_ = tx.Rollback()
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

	err = u.logActivityMongo(ctx, user.ID, "registration", metadata)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userRepository) GetMe(ctx context.Context, userAuthID string) (*entity.User, error) {
	user, err := u.getUserByAuthIDSQL(ctx, userAuthID)
	if err != nil {
		return nil, err
	}

	return u.getUserByIDSQL(ctx, user.ID)
}

func (u *userRepository) GetActivities(ctx context.Context, userAuthID string, page string, limit string) ([]*entity.UserActivity, int64, error) {
	user, err := u.getUserByAuthIDSQL(ctx, userAuthID)
	if err != nil {
		return nil, 0, err
	}

	return u.getUserActivitiesMongo(ctx, user.ID, page, limit)
}

func (u *userRepository) GetUserAddresses(ctx context.Context, userAuthID string) ([]*entity.UserAddress, error) {
	user, err := u.getUserByAuthIDSQL(ctx, userAuthID)
	if err != nil {
		return nil, err
	}

	return u.getUserAddressesSQL(ctx, user.ID)
}

func (u *userRepository) CreateAddress(ctx context.Context, userAuthID string, req dto.CreateUserAddress) (string, error) {
	user, err := u.getUserByAuthIDSQL(ctx, userAuthID)
	if err != nil {
		return "", err
	}

	return u.createUserAddressSQL(ctx, user.ID, req)
}
