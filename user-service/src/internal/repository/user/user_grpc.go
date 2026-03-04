package user

import (
	"context"

	userpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/user"
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
