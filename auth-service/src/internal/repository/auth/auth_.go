package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/dto"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/entity"
	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (a *authRepository) CreateAuthUser(ctx context.Context, req *dto.CreateAuthUserRequest) (string, error) {
	var authID string

	emailExists, err := a.checkUserExistSql(ctx, req.Email)
	if err != nil {
		return authID, err
	}

	if emailExists {
		zerolog.Ctx(ctx).Error().Msg("user_exist")
		return authID, x.New("Email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("hashed_password")
		return authID, x.Wrap(err, "hashed_password")
	}

	authID, err = a.saveUserSql(ctx, req.Email, hashedPassword)
	if err != nil {
		return authID, err
	}

	return authID, nil
}

func (a *authRepository) FindAuthUserByEmail(ctx context.Context, email string) (*entity.UserAuth, error) {
	userAuth, err := a.getUserByEmailSql(ctx, email)
	if err != nil {
		return nil, err
	}

	return userAuth, nil
}

func (a *authRepository) StoreSession(ctx context.Context, userID string, refreshToken string, expired time.Time, email string, ipAddress string) error {
	err := a.storeRefreshTokenSql(ctx, userID, refreshToken, expired)
	if err != nil {
		return err
	}

	err = a.storeSessionCache(ctx, userID, refreshToken, email, ipAddress)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) FindTokenID(ctx context.Context, tokenID string) bool {
	return a.findTokenIDCache(ctx, tokenID)
}

func (a *authRepository) GetUserInfo(ctx context.Context, refreshToken string) (*entity.UserAuth, error) {
	userID, err := a.findRefreshTokenCache(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	result, err := a.getUserByIDSql(ctx, userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *authRepository) BlacklistToken(ctx context.Context, token *jwt.Token) error {
	claims, _ := token.Claims.(jwt.MapClaims)
	tokenID := claims["jti"].(string)

	// Blacklist token
	err := a.blacklistTokenCache(ctx, tokenID)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) ClearSession(ctx context.Context, userID string) error {
	// Delete session
	err := a.deleteSessionCache(ctx, userID)
	if err != nil {
		return err
	}

	// Delete refresh tokens from DB
	err = a.deleteRefreshTokenSql(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (a *authRepository) RotateRefreshToken(ctx context.Context, userID string, oldRefreshToken string, newRefreshToken string, expired time.Time) error {
	err := a.deleteRefreshTokenByTokenSql(ctx, oldRefreshToken)
	if err != nil {
		return err
	}

	err = a.storeRefreshTokenSql(ctx, userID, newRefreshToken, expired)
	if err != nil {
		return err
	}

	hashedOldToken := hashToken(oldRefreshToken)
	hashedNewToken := hashToken(newRefreshToken)
	a.redis0.Del(ctx, fmt.Sprintf("refresh:%s", hashedOldToken))
	a.redis0.Set(ctx, fmt.Sprintf("refresh:%s", hashedNewToken), userID, time.Hour*24*7)

	return nil
}
