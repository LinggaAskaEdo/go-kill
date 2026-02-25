package auth

import (
	"context"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (a *authRepository) CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error) {
	emailExists, err := a.checkUserExist(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if emailExists {
		zerolog.Ctx(ctx).Error().Msg("user_exist")
		return nil, x.NewWithCode(x.CodeSQLConflict, "Email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("hashed_password")
		return nil, x.Wrap(err, "hashed_password")
	}

	authID, err := a.saveUser(ctx, req.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return &authpb.CreateAuthUserResponse{
		Success: true,
		AuthId:  authID,
	}, nil
}
