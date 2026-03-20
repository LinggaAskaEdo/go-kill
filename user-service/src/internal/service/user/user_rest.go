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
		AuthID:    authResp.AuthId,
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
	user, err := s.userRepository.GetMe(ctx, userAuthID)
	if err != nil {
		return nil, err
	}

	resp, total, err := s.userRepository.GetActivities(ctx, user.ID, page, limit)
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

func (s *userService) GetAddresses(ctx context.Context, userAuthID string, page string, limit string) (*dto.AddressesResp, error) {
	user, err := s.userRepository.GetMe(ctx, userAuthID)
	if err != nil {
		return nil, err
	}

	addresses, total, err := s.userRepository.GetUserAddresses(ctx, user.ID, page, limit)
	if err != nil {
		return nil, err
	}

	return &dto.AddressesResp{
		Success: true,
		Data:    toAddresses(addresses),
		Pagination: dto.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (s *userService) CreateAddress(ctx context.Context, userAuthID string, req dto.CreateUserAddress) (string, error) {
	user, err := s.userRepository.GetMe(ctx, userAuthID)
	if err != nil {
		return "", err
	}

	return s.userRepository.CreateAddress(ctx, user.ID, req)
}
