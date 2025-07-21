package repository

import (
	"context"
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/internal/transaction"
)

type UserCache interface {
	SetRegisteredUserID(ctx context.Context, ID int64) error
	HasRegisteredUser(ctx context.Context, ID int64) (bool, error)
	SetUserToRegister(ctx context.Context, userToRegister profile.UserToRegiser) error
	GetUserToRegister(ctx context.Context, ID int64) (profile.UserToRegiser, error)
	DeleteUserToRegister(ctx context.Context, ID int64) error
	SetTempUserData(ctx context.Context, user profile.User) error
	GetTempUserData(ctx context.Context, ID int64) (profile.User, error)
	DeleteTempUserData(ctx context.Context, ID int64) error
}

type Profile interface {
	CreateNewUser(ts transaction.Session, user profile.User) error
	FindUserByTGID(ts transaction.Session, ID int64) (profile.User, error)
	LoadAllActiveUserIDS(ts transaction.Session) ([]int64, error)
}
