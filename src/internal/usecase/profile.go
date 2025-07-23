package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/chronos"
	"tg_video_lessons_bot/tools/sm_gen"
	"tg_video_lessons_bot/tools/sql_null"
	"tg_video_lessons_bot/tools/str"
	"time"

	"github.com/go-telegram/bot/models"
)

type Profile struct {
	ri *rimport.RepositoryImports
}

func NewProfile(
	ri *rimport.RepositoryImports,
) *Profile {
	return &Profile{ri}
}

func (u *Profile) HandlerStart(ctx context.Context, ID int64, userName string) ([]string, error) {
	messages := make([]string, 0, 2)

	cachedUser, err := u.findUserToRegister(ctx, ID)
	switch err {
	case nil:
		messages = append(messages, profile.StepMessages[cachedUser.Step])

	case global.ErrNoData:
		cachedUser = profile.NewDefaultUserToRegister(ID, userName, time.Now())

		err = u.setUserToRegister(ctx, cachedUser, false)
		if err != nil {
			return messages, global.ErrInternalError
		}

		messages = append(messages, []string{profile.HelloMessage, profile.StepMessages[cachedUser.Step]}...)
		err = nil

	default:
		err = global.ErrInternalError
	}

	return messages, err
}

func (u *Profile) HandlerFullName(ctx context.Context, ID int64, text string) (message string, err error) {
	text = strings.TrimSpace(text)

	cachedUser, err := u.findUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	splitted := str.SplitStringByEmptySpace(text)
	if len(splitted) < 2 {
		return message, profile.ErrFullNameValidation
	}

	cachedUser.FirstName = sql_null.NewString(str.CapFirstLowerRest(splitted[0]))
	cachedUser.LastName = sql_null.NewString(str.CapFirstLowerRest(splitted[1]))

	cachedUser.Step = profile.RegisterStepBirthDate

	err = u.setUserToRegister(ctx, cachedUser, true)
	if err != nil {
		return message, err
	}

	return profile.StepMessages[cachedUser.Step], nil
}

func (u *Profile) HandleBirthDate(ctx context.Context, ID int64, text string) (message global.ReplyMessage, err error) {

	parsed, err := time.Parse(chronos.DateMask, text)
	if err != nil {
		return message, profile.ErrBirhDateInvalid
	}

	if parsed.After(time.Now()) {
		return message, profile.ErrBirthDateInFuture
	}

	cachedUser, err := u.findUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	cachedUser.BirthDate = sql_null.NewNullTime(parsed)
	cachedUser.Step = profile.RegisterStepPhoneNumber

	err = u.setUserToRegister(ctx, cachedUser, true)
	if err != nil {
		return message, err
	}

	message = global.NewReplyMessage(
		profile.StepMessages[cachedUser.Step],
		[]models.KeyboardButton{
			{
				Text:           profile.SendPhoneNumber,
				RequestContact: true,
			},
		},
	)

	return message, nil
}

func (u *Profile) HandlePhoneNumber(ctx context.Context, ID int64, contact models.Contact) (message global.ReplyMessage, err error) {
	if contact.PhoneNumber == "" {
		return message, profile.ErrPhoneNumberEmpty
	}

	cachedUser, err := u.findUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	cachedUser.PhoneNumber = sql_null.NewString(contact.PhoneNumber)

	ts := transaction.MustGetSession(ctx)

	err = u.ri.Repository.Profile.MarkUserToRegiserAsRegistered(ts, ID)
	if err != nil {
		log.Println("не удалось удалить временные данные пользователя в базе: ", err)
		return message, global.ErrInternalError
	}

	err = u.ri.Profile.CreateNewUser(ts, cachedUser.NewUser(time.Now()))
	if err != nil {
		log.Println("не удалось сохранить данные пользователя: ", err)
		return message, global.ErrInternalError
	}

	err = u.ri.Repository.UserCache.DeleteUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось удалить временные данные пользователя в кеше: ", err)
		return message, global.ErrInternalError
	}

	err = u.ri.Repository.UserCache.SetRegisteredUserID(ctx, ID)
	if err != nil {
		log.Println("не удалось сохранить ID пользователя в кеш: ", err)
		return message, global.ErrInternalError
	}

	message = global.NewReplyMessage(
		profile.RegistrationWasSuccessful,
		[]models.KeyboardButton{
			{
				Text: global.TextCommandProfile[global.AppLangCode],
			},
		},
	)

	return message, nil
}

func (u *Profile) HandleStepsValidationMessages(ctx context.Context, ID int64) error {
	cachedUser, err := u.ri.Repository.UserCache.GetUserToRegister(ctx, ID)
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return global.ErrInternalError
	}

	return profile.StepValidationErrorMessages[cachedUser.Step]
}

func (u *Profile) HandlerProfileInfo(ctx context.Context, ID int64) (message string, err error) {
	var userData profile.User

	userData, err = u.ri.Repository.UserCache.GetTempUserData(ctx, ID)
	switch err {
	case nil:
	case global.ErrNoData:
		userData, err = sm_gen.LoadDataIgnoreErrNoData(ctx, func(ts transaction.Session) (profile.User, error) {
			return u.ri.Repository.Profile.FindUserByTGID(ts, ID)
		}, "не удалось найти пользователя по ID телеграм")
		if err != nil {
			return message, err
		}

		err = u.ri.UserCache.SetTempUserData(ctx, userData)
		if err != nil {
			log.Println("не удалось найти пользователя: ", err)
			return message, global.ErrInternalError
		}

	default:
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	message = fmt.Sprintf(
		profile.ProfileInfoMessage,
		userData.FirstName,
		userData.LastName,
		userData.PhoneNumber,
		userData.CalcAge(),
		userData.RegisterDate.Format(chronos.DateMask),
	)

	return message, nil
}

func (u *Profile) findUserToRegister(ctx context.Context, ID int64) (profile.UserToRegister, error) {
	cachedUser, err := u.ri.Repository.UserCache.GetUserToRegister(ctx, ID)
	switch err {
	case nil:
		return cachedUser, nil

	case global.ErrNoData:
		ts := transaction.MustGetSession(ctx)
		dbUser, err := u.ri.Repository.Profile.FindUserToRegiserByTGID(ts, ID)
		switch err {
		case nil:
			return dbUser, nil

		case global.ErrNoData:
			return dbUser, global.ErrNoData

		default:
			log.Println("не удалось найти пользователя в базе: ", err)
			return dbUser, global.ErrInternalError
		}

	default:
		log.Println("не удалось найти пользователя в кеше: ", err)
		return cachedUser, global.ErrInternalError
	}
}

func (u *Profile) setUserToRegister(ctx context.Context, user profile.UserToRegister, update bool) error {
	err := u.ri.Repository.UserCache.SetUserToRegister(ctx, user)
	if err != nil {
		log.Println("ошибка при кешировании пользователя: ", err)
	}

	dbFunc := u.ri.Repository.Profile.SaveUserToRegiser
	if update {
		dbFunc = u.ri.Repository.Profile.UpdateUserToRegister
	}

	ts := transaction.MustGetSession(ctx)
	err = dbFunc(ts, user)
	if err != nil {
		log.Println("ошибка при записи пользователя в базу: ", err)
		return global.ErrInternalError
	}

	return nil
}
