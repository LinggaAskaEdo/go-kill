package user

import (
	"context"
	"time"

	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (u *userRepository) logActivityMongo(ctx context.Context, userID, activityType string) error {
	ip, _ := ctx.Value("ip").(string)
	ua, _ := ctx.Value("user_agent").(string)

	metadata := map[string]any{
		"ip_address":          ip,
		"user_agent":          ua,
		"registration_method": "email",
	}

	activity := entity.UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		Metadata:     metadata,
		Timestamp:    time.Now(),
		CreatedAt:    time.Now(),
	}

	_, err := u.mongo0.Collection("user_activities").InsertOne(context.Background(), activity)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("log_activity_mongo")
		return err
	}

	return nil
}
