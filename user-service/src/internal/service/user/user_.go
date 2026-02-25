package user

import (
	"context"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (s *userService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest, grpc *grpc.Grpc) (*entity.User, error) {
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

	return s.userRepository.RegisterUser(ctx, user)
}
