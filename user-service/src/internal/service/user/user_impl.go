package user

import (
	"context"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (s *userService) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	authResp, err := s.authClient.ValidateToken(ctx, req)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("validate_token")
		return nil, err
	}

	return authResp, nil
}

func (s *userService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserRegResp, error) {
	authResp, err := s.authClient.CreateAuthUser(ctx, &authpb.CreateAuthUserRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("register_user")
		return nil, err
	}

	user := &entity.User{
		AutdID:    authResp.AuthId,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err = s.userRepository.RegisterUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return toUserRegResp(user), nil
}

func (s *userService) GetMe(ctx context.Context, userAuthID string) (*dto.UserResp, error) {
	user, err := s.userRepository.GetMe(ctx, userAuthID)
	if err != nil {
		return nil, err
	}

	return toUserResp(user), nil
}

func (s *userService) GetActivities(ctx context.Context, userAuthID string, page string, limit string) (*dto.UserActivity, error) {
	resp, total, err := s.userRepository.GetActivities(ctx, userAuthID, page, limit)
	if err != nil {
		return nil, err
	}

	return &dto.UserActivity{
		Success: true,
		Data:    resp,
		Pagination: dto.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (s *userService) GetAddresses(ctx context.Context, userAuthID string) ([]*dto.Address, error) {
	addresses, err := s.userRepository.GetUserAddresses(ctx, userAuthID)
	if err != nil {
		return nil, err
	}

	return toAddresses(addresses), nil
}

func (s *userService) CreateAddress(ctx context.Context, userAuthID string, req dto.CreateUserAddress) (string, error) {
	var result string

	result, err := s.userRepository.CreateAddress(ctx, userAuthID, req)
	if err != nil {
		return result, err
	}

	return result, nil
}
