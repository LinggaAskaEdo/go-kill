package user

import (
	"context"

	"github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository/user"
)

type UserServiceItf interface {
	RegisterUser(ctx context.Context, req dto.RegisterUserRequest, grpc *grpc.Grpc) (*entity.User, error)
}

type userService struct {
	userRepository user.UserRepositoryItf
}

func InitUserService(userRepository user.UserRepositoryItf) UserServiceItf {
	return &userService{
		userRepository: userRepository,
	}
}
