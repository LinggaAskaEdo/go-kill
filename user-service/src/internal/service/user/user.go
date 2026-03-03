package user

import (
	"context"

	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/repository/user"

	"google.golang.org/grpc"
)

type UserServiceItf interface {
	ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error)
	RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserRegResp, error)
	GetMe(ctx context.Context, userAuthID string) (*dto.UserResp, error)
	GetActivities(ctx context.Context, userAuthID string, page string, limit string) (*dto.UserActivity, error)
	GetAddresses(ctx context.Context, userAuthID string) ([]*dto.Address, error)
	CreateAddress(ctx context.Context, userAuthID string, req dto.CreateUserAddress) (string, error)
}

type userService struct {
	authClient     authpb.AuthServiceClient
	userRepository user.UserRepositoryItf
}

func InitUserService(clientConn *grpc.ClientConn, userRepository user.UserRepositoryItf) UserServiceItf {
	return &userService{
		authClient:     authpb.NewAuthServiceClient(clientConn),
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

func toAddresses(addresses []*entity.UserAddress) []*dto.Address {
	if addresses == nil {
		return nil
	}

	result := make([]*dto.Address, len(addresses))
	for i, address := range addresses {
		result[i] = &dto.Address{
			ID:            address.ID,
			AddressType:   address.AddressType,
			StreetAddress: address.StreetAddress,
			City:          address.City,
			State:         address.State,
			PostalCode:    address.PostalCode,
			Country:       address.Country,
			IsDefault:     address.IsDefault,
		}
	}

	return result
}
