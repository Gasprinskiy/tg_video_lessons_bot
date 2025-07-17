package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/chronos"
	"tg_video_lessons_bot/tools/dump"
	"time"

	"github.com/go-telegram/bot/models"
)

type Profile struct {
	ri *rimport.RepositoryImports
}

func NewProfile(ri *rimport.RepositoryImports) *Profile {
	return &Profile{ri}
}

func (u *Profile) HandlerStart(ctx context.Context, ID int64, userName string) ([]string, error) {
	messages := make([]string, 0, 2)

	cachedUser, err := u.ri.Repository.UserCache.GetUserToRegister(ctx, ID)
	switch err {
	case nil:
		messages = append(messages, profile.StepMessages[cachedUser.RegisterStep])

	case global.ErrNoData:
		cachedUser = profile.UserToRegiser{
			User: profile.User{
				ID:       ID,
				UserName: userName,
			},
			RegisterStep: profile.RegisterStepFullName,
		}

		err = u.ri.Repository.UserCache.SetUserToRegister(ctx, cachedUser)
		if err != nil {
			log.Println("ошибка при кешировании пользователя: ", err)
			return messages, global.ErrInternalError
		}

		messages = append(messages, []string{profile.HelloMessage, profile.StepMessages[cachedUser.RegisterStep]}...)
		err = nil

	default:
		log.Println("ошибка при получении данных пользователя: ", err)
		err = global.ErrInternalError
	}

	return messages, err
}

func (u *Profile) HandlerFullName(ctx context.Context, ID int64, text string) (message string, err error) {
	text = strings.TrimSpace(text)

	cachedUser, err := u.ri.Repository.UserCache.GetUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	splitted := strings.Split(text, " ")
	if len(splitted) != 2 {
		return message, profile.ErrFullNameValidation
	}

	cachedUser.FirstName = strings.TrimSpace(splitted[0])
	cachedUser.LastName = strings.TrimSpace(splitted[1])

	cachedUser.RegisterStep = profile.RegisterStepBirthDate

	err = u.ri.Repository.UserCache.SetUserToRegister(ctx, cachedUser)
	if err != nil {
		log.Println("не удалось обновить данные пользователя: ", err)
		return message, global.ErrInternalError
	}

	return profile.StepMessages[cachedUser.RegisterStep], nil
}

func (u *Profile) HandleBirthDate(ctx context.Context, ID int64, text string) (message global.ReplyMessage, err error) {

	parsed, err := time.Parse(chronos.DateMask, text)
	if err != nil {
		return message, profile.ErrBirhDateInvalid
	}

	if parsed.After(time.Now()) {
		return message, profile.ErrBirthDateInFuture
	}

	cachedUser, err := u.ri.Repository.UserCache.GetUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	cachedUser.BirthDate = parsed
	cachedUser.RegisterStep = profile.RegisterStepPhoneNumber

	err = u.ri.Repository.UserCache.SetUserToRegister(ctx, cachedUser)
	if err != nil {
		log.Println("не удалось обновить данные пользователя: ", err)
		return message, global.ErrInternalError
	}

	message = global.NewReplyMessage(
		profile.StepMessages[cachedUser.RegisterStep],
		[]models.KeyboardButton{
			{
				Text:           profile.SendPhoneNumber,
				RequestContact: true,
			},
		},
	)

	return message, nil
}

func (u *Profile) HandlePhoneNumber(ctx context.Context, ID int64, contact models.Contact) (message string, err error) {
	if contact.PhoneNumber == "" {
		return message, profile.ErrPhoneNumberEmpty
	}

	cachedUser, err := u.ri.Repository.UserCache.GetUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	cachedUser.PhoneNumber = contact.PhoneNumber

	err = u.ri.Repository.UserCache.DeleteUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось удалить временные данные пользователя: ", err)
		return message, global.ErrInternalError
	}

	fmt.Println("cachedUser: ", dump.Struct(cachedUser.User))

	err = u.ri.Repository.UserCache.SetRegisteredUserID(ctx, ID)
	if err != nil {
		log.Println("не удалось сохранить ID пользователя в кеш: ", err)
		return message, global.ErrInternalError
	}

	return profile.RegistrationWasSuccessful, nil
}
