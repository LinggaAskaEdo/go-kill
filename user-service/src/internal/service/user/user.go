package user

import (
	"context"

	"github.com/linggaaskaedo/go-kill/user-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository/user"
)

type UserServiceItf interface {
	RegisterUser(ctx context.Context, req dto.RegisterUserRequest, grpc *grpc.Grpc) (*dto.UserRegResp, error)
	GetMe(ctx context.Context, userID string) (*dto.UserResp, error)
	GetActivities(ctx context.Context, userID string, page string, limit string) (dto.UserActivity, error)
}

type userService struct {
	userRepository user.UserRepositoryItf
}

func InitUserService(userRepository user.UserRepositoryItf) UserServiceItf {
	return &userService{
		userRepository: userRepository,
	}
}

func toUserRegResp(u *entity.User) *dto.UserRegResp {
	if u == nil {
		return nil
	}

	return &dto.UserRegResp{
		ID:        u.ID,
		AutdID:    u.AutdID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toUserResp(u *entity.User) *dto.UserResp {
	if u == nil {
		return nil
	}

	return &dto.UserResp{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}
