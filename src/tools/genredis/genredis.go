package genredis

import (
	"context"
	"encoding/json"
	"tg_video_lessons_bot/internal/entity/global"

	"github.com/redis/go-redis/v9"
)

func GetStruct[T any](ctx context.Context, db *redis.Client, key string) (T, error) {
	var data T

	stringResult, err := db.Get(ctx, key).Result()
	if err != nil {
		return data, HandleGetError(err)
	}

	if err = json.Unmarshal([]byte(stringResult), &data); err != nil {
		return data, err
	}

	return data, nil
}

func HandleGetError(err error) error {
	if err == redis.Nil {
		return global.ErrNoData
	}

	return err
}
