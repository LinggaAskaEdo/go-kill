package user

import (
	"context"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (s *userService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest, grpc *grpc.Grpc) (*dto.UserRegResp, error) {
	authResp, err := grpc.CreateAuthUser(ctx, &authpb.CreateAuthUserRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("grpc_create_auth_user")
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

func (s *userService) GetMe(ctx context.Context, userID string) (*dto.UserResp, error) {
	user, err := s.userRepository.GetMe(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toUserResp(user), nil
}

func (s *userService) GetActivities(ctx context.Context, userID string, page string, limit string) (dto.UserActivity, error) {
	resp, total, err := s.userRepository.GetActivities(ctx, userID, page, limit)
	if err != nil {
		return dto.UserActivity{}, err
	}

	return dto.UserActivity{
		Success: true,
		Data:    resp,
		Pagination: dto.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}
