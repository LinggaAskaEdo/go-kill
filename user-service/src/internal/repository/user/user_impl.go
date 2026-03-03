package user

import (
	"context"
	"database/sql"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (u *userRepository) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	userID, err := u.createUserSQL(ctx, req)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{Success: true, UserId: userID}, nil
}

func (u *userRepository) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	resp, err := u.getUserByIDSQL(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserResponse{
		Id:        resp.ID,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Found:     true,
	}, nil
}

func (u *userRepository) GetAddress(ctx context.Context, req *userpb.GetAddressRequest) (*userpb.GetAddressResponse, error) {
	resp, err := u.getUserAddressByIDSQL(ctx, req)
	if err != nil {
		return nil, err
	}

	return &userpb.GetAddressResponse{
		Id:            resp.ID,
		UserId:        resp.UserID,
		StreetAddress: resp.StreetAddress,
		City:          resp.City,
		State:         resp.State.String,
		PostalCode:    resp.PostalCode,
		Country:       resp.Country,
		Found:         true,
	}, nil
}

func (u *userRepository) LogActivity(ctx context.Context, req *userpb.LogActivityRequest) (*userpb.LogActivityResponse, error) {
	err := u.logActivityMongo(ctx, req.UserId, req.ActivityType, convertMetadata(req.Metadata))
	if err != nil {
		return nil, err
	}

	return &userpb.LogActivityResponse{Success: true}, nil
}

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
