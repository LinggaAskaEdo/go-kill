package user

import (
	"context"
	"time"

	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/rs/zerolog"
)

func (u *userRepository) logActivityMongo(ctx context.Context, userID, activityType string, metadata map[string]any) {
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
	}
}
