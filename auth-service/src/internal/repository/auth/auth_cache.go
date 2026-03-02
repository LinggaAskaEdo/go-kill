package auth

import (
	"context"
	"fmt"
	"time"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
)

func (a *authRepository) storeSessionCache(ctx context.Context, userID string, refreshToken string, email string, ipAddress string) error {
	sessionKey := fmt.Sprintf("session:%s", userID)
	sessionData := fmt.Sprintf(`{"user_id":"%s","email":"%s","ip":"%s"}`, userID, email, ipAddress)

	if err := a.redis0.Set(ctx, sessionKey, sessionData, time.Hour*1).Err(); err != nil {
		return x.WrapWithCode(err, x.CodeCacheSetHashKey, "set_cache_session_user")
	}

	if err := a.redis0.Set(ctx, fmt.Sprintf("refresh:%s", refreshToken), userID, time.Hour*24*7).Err(); err != nil {
		return x.WrapWithCode(err, x.CodeCacheSetHashKey, "set_cache_refresh_token_user")
	}

	return nil
}

func (a *authRepository) findTokenIDCache(ctx context.Context, tokenID string) bool {
	exists, _ := a.redis0.Exists(ctx, fmt.Sprintf("blacklist:%s", tokenID)).Result()

	return exists > 0
}

func (a *authRepository) findRefreshTokenCache(ctx context.Context, refreshToken string) (string, error) {
	userID, err := a.redis0.Get(ctx, fmt.Sprintf("refresh:%s", refreshToken)).Result()
	if err != nil {
		return userID, x.WrapWithCode(err, x.CodeCacheGetSimpleKey, "Invalid refresh token")
	}

	return userID, nil
}

func (a *authRepository) blacklistTokenCache(ctx context.Context, tokenID string) error {
	if err := a.redis0.Set(ctx, fmt.Sprintf("blacklist:%s", tokenID), "revoked", time.Hour*1).Err(); err != nil {
		return x.WrapWithCode(err, x.CodeCacheSetSimpleKey, "blacklist_token_cache")
	}

	return nil
}

func (a *authRepository) deleteSessionCache(ctx context.Context, userID string) error {
	if err := a.redis0.Del(ctx, fmt.Sprintf("session:%s", userID)).Err(); err != nil {
		return x.WrapWithCode(err, x.CodeCacheDeleteSimpleKey, "delete_session_cache")
	}

	return nil
}
