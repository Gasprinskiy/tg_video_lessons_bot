package redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"tg_video_lessons_bot/internal/entity/profile"
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

func (r *userCache) SetRegisteredUserID(ctx context.Context, ID int64) error {
	return r.db.Set(ctx, fmt.Sprintf("%d:registered", ID), true, 0).Err()
}

func (r *userCache) HasRegisteredUser(ctx context.Context, ID int64) (bool, error) {
	result, err := r.db.Get(ctx, fmt.Sprintf("%d:registered", ID)).Result()
	switch err {
	case nil:
		return strconv.ParseBool(result)

	case redis.Nil:
		return false, nil

	default:
		return false, err
	}
}

func (r *userCache) SetUserToRegister(ctx context.Context, userToRegister profile.UserToRegiser) error {
	byteData, err := json.Marshal(userToRegister)
	if err != nil {
		return err
	}

	return r.db.Set(ctx, fmt.Sprintf("%d", userToRegister.ID), byteData, r.ttl).Err()
}

func (r *userCache) GetUserToRegister(ctx context.Context, ID int64) (profile.UserToRegiser, error) {
	return genredis.GetStruct[profile.UserToRegiser](ctx, r.db, fmt.Sprintf("%d", ID))
}

func (r *userCache) DeleteUserToRegister(ctx context.Context, ID int64) error {
	_, err := r.db.Del(ctx, fmt.Sprintf("%d", ID)).Result()
	return err
}

func (r *userCache) SetTempUserData(ctx context.Context, user profile.User) error {
	byteData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.db.Set(ctx, fmt.Sprintf("%d:temp", user.ID), byteData, r.ttl).Err()
}

func (r *userCache) GetTempUserData(ctx context.Context, ID int64) (profile.User, error) {
	return genredis.GetStruct[profile.User](ctx, r.db, fmt.Sprintf("%d:temp", ID))
}

func (r *userCache) DeleteTempUserData(ctx context.Context, ID int64) error {
	_, err := r.db.Del(ctx, fmt.Sprintf("%d:temp", ID)).Result()
	return err
}
