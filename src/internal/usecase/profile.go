package usecase

import (
	"context"
	"fmt"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/rimport"

	"github.com/go-telegram/bot/models"
)

type Profile struct {
	ri *rimport.RepositoryImports
}

func NewProfile(ri *rimport.RepositoryImports) *Profile {
	return &Profile{ri}
}

func (u *Profile) HandlerStart(ctx context.Context, from models.User) ([]string, error) {
	messages := make([]string, 0, 2)

	cachedUser, err := u.ri.Repository.UserCache.GetUserToRegister(ctx, from.ID)
	switch err {
	case nil:
		messages = append(messages, profile.StepMessages[cachedUser.RegisterStep])

	case global.ErrNoData:
		cachedUser = profile.UserToRegiser{
			User: profile.User{
				ID: from.ID,
			},
			RegisterStep: profile.RegisterStepFullName,
		}

		err = u.ri.Repository.UserCache.SetUserToRegister(ctx, cachedUser)
		if err != nil {
			fmt.Println("ошибка при кешировании пользователя: ", err)
			return messages, global.ErrInternalError
		}

		messages = append(messages, []string{profile.HelloMessage, profile.StepMessages[cachedUser.RegisterStep]}...)
		err = nil

	default:
		fmt.Println("ошибка при получении данных пользователя: ", err)
		err = global.ErrInternalError
	}

	return messages, err
}
