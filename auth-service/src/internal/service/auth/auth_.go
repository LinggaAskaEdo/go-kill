package auth

import (
	"context"
	"time"

	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/model/dto"
	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (a *authService) CreateAuthUser(ctx context.Context, req *dto.CreateAuthUserRequest) (*dto.CreateAuthUserResponse, error) {
	authID, err := a.authRepository.CreateAuthUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.CreateAuthUserResponse{
		Success: true,
		AuthId:  authID,
	}, nil
}

func (a *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	userAuth, err := a.authRepository.FindAuthUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !userAuth.IsActive {
		zerolog.Ctx(ctx).Error().Msg("user_not_active")
		return nil, x.New("User is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAuth.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, x.Wrap(err, "Invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userAuth.ID,
		"email": userAuth.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
		"jti":   generateTokenID(),
	})

	accessToken, err := token.SignedString([]byte(a.authOptions.JwtSecret))
	if err != nil {
		return nil, x.Wrap(err, "Failed to generate token")
	}

	refreshToken := generateRefreshToken()

	err = a.authRepository.StoreSession(ctx, userAuth.ID, refreshToken, time.Now().Add(time.Hour*24*7), userAuth.Email, req.IpAddress)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	}, nil
}

func (a *authService) ValidateToken(ctx context.Context, req *dto.ValidateTokenRequest) (*dto.ValidateTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (any, error) { return []byte(a.authOptions.JwtSecret), nil })
	if err != nil {
		return nil, x.Wrap(err, "Failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, x.New("Failed to extract token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, x.New("Invalid token subject")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, x.New("Invalid token email")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, x.New("Invalid token ID")
	}

	exists := a.authRepository.FindTokenID(ctx, jti)
	if exists {
		return nil, x.New("token in the blacklist")
	}

	return &dto.ValidateTokenResponse{
		Valid:  true,
		UserId: sub,
		Email:  email,
	}, nil
}

func (a *authService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	user, err := a.authRepository.GetUserInfo(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	newRefreshToken := generateRefreshToken()

	err = a.authRepository.RotateRefreshToken(ctx, user.ID, req.RefreshToken, newRefreshToken, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
		"jti":   generateTokenID(),
	})

	accessToken, err := token.SignedString([]byte(a.authOptions.JwtSecret))
	if err != nil {
		return nil, x.Wrap(err, "Failed signed token")
	}

	return &dto.RefreshTokenResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600,
	}, nil
}

func (a *authService) Logout(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.authOptions.JwtSecret), nil
	})

	if err == nil && token != nil {
		if err := a.authRepository.BlacklistToken(ctx, token); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed_to_blacklist_token")
		}
	}

	if err := a.authRepository.ClearSession(ctx, req.UserId); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed_to_clear_session")
		return &dto.LogoutResponse{Success: false, Message: "Failed to logout"}, err
	}

	return &dto.LogoutResponse{Success: true, Message: "Logged out successfully"}, nil
}
