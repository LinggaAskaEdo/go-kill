package user

import (
	"context"
	"strconv"
	"time"

	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/rs/zerolog"
)

func (u *userRepository) logActivityMongo(ctx context.Context, userID, activityType string, metadata map[string]any) error {
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

func (u *userRepository) getUserActivitiesMongo(ctx context.Context, userID string, page string, limit string) ([]*entity.UserActivity, int64, error) {
	var activities []*entity.UserActivity

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1 // default to first page
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10 // default page size
	}

	if limitInt > 100 {
		limitInt = 100
	}

	skip := (pageInt - 1) * limitInt

	collection := u.mongo0.Collection("user_activities")
	filter := bson.M{"user_id": userID}
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limitInt))

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

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("count_user_activities_mongo")
		return activities, 0, err
	}

	return activities, total, nil
}
