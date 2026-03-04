package user

import (
	"context"

	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
)

func (s *userService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	return s.userRepository.CreateUser(ctx, req)
}

func (s *userService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	return s.userRepository.GetUser(ctx, req)
}

func (s *userService) GetAddress(ctx context.Context, req *userpb.GetAddressRequest) (*userpb.GetAddressResponse, error) {
	return s.userRepository.GetAddress(ctx, req)
}

func (s *userService) LogActivity(ctx context.Context, req *userpb.LogActivityRequest) (*userpb.LogActivityResponse, error) {
	return s.userRepository.LogActivity(ctx, req)
}
