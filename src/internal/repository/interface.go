package repository

import (
	"context"
	"tg_video_lessons_bot/internal/entity/profile"
)

type UserCache interface {
	SetRegisteredUserID(ctx context.Context, ID int64) error
	HasRegisteredUser(ctx context.Context, ID int64) (bool, error)
	SetUserToRegister(ctx context.Context, userToRegister profile.UserToRegiser) error
	GetUserToRegister(ctx context.Context, ID int64) (profile.UserToRegiser, error)
	DeleteUserToRegister(ctx context.Context, ID int64) error
}
