package repository

import (
	"context"
	"tg_video_lessons_bot/internal/entity/user"
)

type UserCache interface {
	SetRegisteredUserID(ctx context.Context, ID int64) error
	HasRegisteredUser(ctx context.Context, ID int64) (bool, error)
	SetUserToRegister(ctx context.Context, userToRegister user.UserToRegiser) error
	GetUserToRegister(ctx context.Context, ID int64) (user.UserToRegiser, error)
}
