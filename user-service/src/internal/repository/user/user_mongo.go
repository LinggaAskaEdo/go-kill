package user

import (
	"context"
	"strconv"
	"time"

	x "github.com/linggaaskaedo/go-kill/common/pkg/errors"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/model/entity"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
	collectionName  = "user_activities"
)

func (u *userRepository) logActivityMongo(ctx context.Context, userID, activityType string, metadata map[string]any) error {
	activity := entity.UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		Metadata:     metadata,
		Timestamp:    time.Now(),
	}

	_, err := u.mongo0.Collection(collectionName).InsertOne(ctx, activity)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("log_activity_mongo")
		return x.WrapWithCode(err, x.CodeSQLCreate, "log_activity_mongo")
	}

	return nil
}

func (u *userRepository) getUserActivitiesMongo(ctx context.Context, userID string, page string, limit string) ([]*entity.UserActivity, int64, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = defaultPage
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = defaultPageSize
	}

	if limitInt > maxPageSize {
		limitInt = maxPageSize
	}

	skip := (pageInt - 1) * limitInt

	collection := u.mongo0.Collection(collectionName)
	filter := bson.M{"user_id": userID}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limitInt))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_activities_mongo")
		return nil, 0, x.WrapWithCode(err, x.CodeSQLRead, "get_user_activities_mongo")
	}
	defer cursor.Close(ctx)

	var activities []*entity.UserActivity
	if err = cursor.All(ctx, &activities); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("get_user_activities_mongo")
		return nil, 0, x.WrapWithCode(err, x.CodeSQLRead, "get_user_activities_mongo")
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("count_user_activities_mongo")
		return nil, 0, x.WrapWithCode(err, x.CodeSQLRead, "count_user_activities_mongo")
	}

	return activities, total, nil
}
