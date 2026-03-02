package user

import (
	"context"
	"time"

	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

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

func (u *userRepository) getUserActivitiesMongo(ctx context.Context, userID string, page string, limit string) ([]entity.UserActivity, int64, error) {
	var activities []entity.UserActivity

	collection := u.mongo0.Collection("user_activities")
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(20)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_activities_mongo")
		return activities, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &activities); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_activities_mongo")
		return activities, 0, err
	}

	total, _ := collection.CountDocuments(context.Background(), filter)

	return activities, total, nil
}
