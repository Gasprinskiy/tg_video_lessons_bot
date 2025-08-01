package repository

import (
	"context"
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/internal/transaction"
)

type UserCache interface {
	SetRegisteredUserID(ctx context.Context, ID int64) error
	HasRegisteredUser(ctx context.Context, ID int64) (bool, error)
	SetUserToRegister(ctx context.Context, userToRegister profile.UserToRegister) error
	GetUserToRegister(ctx context.Context, ID int64) (profile.UserToRegister, error)
	DeleteUserToRegister(ctx context.Context, ID int64) error
	SetTempUserData(ctx context.Context, user profile.User) error
	GetTempUserData(ctx context.Context, ID int64) (profile.User, error)
	DeleteTempUserData(ctx context.Context, ID int64) error
}

type Profile interface {
	CreateNewUser(ts transaction.Session, user profile.User) error
	SaveUserToRegiser(ts transaction.Session, user profile.UserToRegister) error
	UpdateUserToRegister(ts transaction.Session, user profile.UserToRegister) error
	FindUserToRegiserByTGID(ts transaction.Session, ID int64) (profile.UserToRegister, error)
	MarkUserToRegiserAsRegistered(ts transaction.Session, ID int64) error
	FindUserByTGID(ts transaction.Session, ID int64) (profile.User, error)
	LoadAllActiveUserIDS(ts transaction.Session) ([]int64, error)
}
