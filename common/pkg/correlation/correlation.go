package correlation

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/pkg/preference"

	"github.com/gin-gonic/gin"
)

func AttachKeyValCtx(ctx context.Context, pairs ...any) context.Context {
	if len(pairs)%2 != 0 {
		panic("AttachValues: odd number of arguments, expected key-value pairs")
	}

	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]
		value := pairs[i+1]
		ctx = context.WithValue(ctx, key, value)
	}

	return ctx
}

func GetCtxKeyVal(c *gin.Context, key preference.CtxKey) string {
	if val := c.Request.Context().Value(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}

	return ""
}

func WithReqID(ctx context.Context, key preference.CtxKey, id string) context.Context {
	return context.WithValue(ctx, key, id)
}

func GetReqID(ctx context.Context, key preference.CtxKey) string {
	if id, ok := ctx.Value(key).(string); ok {
		return id
	}

	return ""
}
