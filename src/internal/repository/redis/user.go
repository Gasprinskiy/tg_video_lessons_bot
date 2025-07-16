package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"tg_video_lessons_bot/internal/entity/user"
	"tg_video_lessons_bot/internal/repository"
	"tg_video_lessons_bot/tools/genredis"
	"time"

	"github.com/redis/go-redis/v9"
)

type userCache struct {
	db  *redis.Client
	ttl time.Duration
}

func NewUserCache(db *redis.Client, ttl time.Duration) repository.UserCache {
	return &userCache{db, ttl}
}

func (r userCache) SetRegisteredUserID(ctx context.Context, ID int64) error {
	return r.db.Set(ctx, fmt.Sprintf("%d:registered", ID), true, 0).Err()
}

func (r userCache) HasRegisteredUser(ctx context.Context, ID int64) (bool, error) {
	result, err := r.db.Get(ctx, fmt.Sprintf("%d:registered", ID)).Result()
	if err != nil {
		return false, genredis.HandleGetError(err)
	}

	return strconv.ParseBool(result)
}

func (r userCache) SetUserToRegister(ctx context.Context, userToRegister user.UserToRegiser) error {
	byteData, err := json.Marshal(userToRegister.User)
	if err != nil {
		return err
	}

	return r.db.Set(ctx, fmt.Sprintf("%d", userToRegister.ID), byteData, r.ttl).Err()
}

func (r userCache) GetUserToRegister(ctx context.Context, ID int64) (user.UserToRegiser, error) {
	return genredis.GetStruct[user.UserToRegiser](ctx, r.db, fmt.Sprintf("%d:registered", ID))

}
