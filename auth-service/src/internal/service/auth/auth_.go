package auth

import (
	"context"
	"time"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

func (a *authService) CreateAuthUser(ctx context.Context, req *authpb.CreateAuthUserRequest) (*authpb.CreateAuthUserResponse, error) {
	authID, err := a.authRepository.CreateAuthUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &authpb.CreateAuthUserResponse{
		Success: true,
		AuthId:  authID,
	}, nil
}

func (a *authService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	userAuth, err := a.authRepository.FindAuthUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !userAuth.IsActive {
		zerolog.Ctx(ctx).Error().Msg("user_not_active")
		return nil, x.New("User is not active")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(userAuth.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, x.Wrap(err, "Invalid email or password")
	}

	// Generate JWT token
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

	// Generate refresh token
	refreshToken := generateRefreshToken()

	//  Store refresh token in DB and session in Redis
	err = a.authRepository.StoreSession(ctx, userAuth.ID, hashToken(refreshToken), time.Now().Add(time.Hour*24*7), userAuth.Email, req.IpAddress)
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		Success:      true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	}, nil
}

func (a *authService) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (any, error) { return []byte(a.authOptions.JwtSecret), nil })
	if err != nil || !token.Valid {
		return nil, x.Wrap(err, "Failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, x.New("Failed to claims token")
	}

	// Check if token is blacklisted
	tokenID := claims["jti"].(string)
	exists := a.authRepository.FindTokenID(ctx, tokenID)
	if exists {
		return nil, x.New("token in the blacklist")
	}

	return &authpb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims["sub"].(string),
		Email:  claims["email"].(string),
	}, nil
}

func (a *authService) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	user, err := a.authRepository.GetUserInfo(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Generate new access token
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

	return &authpb.RefreshTokenResponse{
		Success:     true,
		AccessToken: accessToken,
		ExpiresIn:   3600,
	}, nil
}

func (a *authService) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	// Parse token to get JTI
	token, _ := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.authOptions.JwtSecret), nil
	})

	if token != nil {
		_ = a.authRepository.BlacklistToken(ctx, token)
	}

	_ = a.authRepository.ClearSession(ctx, req.UserId)

	return &authpb.LogoutResponse{Success: true, Message: "Logged out successfully"}, nil
}
