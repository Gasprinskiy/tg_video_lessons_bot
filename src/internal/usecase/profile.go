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
	"tg_video_lessons_bot/tools/gennull"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/tools/sql_null"
	"tg_video_lessons_bot/tools/str"
	"time"

	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

type Profile struct {
	ri  *rimport.RepositoryImports
	log *logger.Logger
}

func NewProfile(
	ri *rimport.RepositoryImports,
	log *logger.Logger,
) *Profile {
	return &Profile{ri, log}
}

func (u *Profile) logPrefix() string {
	return "[profile]"
}

func (u *Profile) HandlerStart(ctx context.Context, ID int64, userName string) ([]string, error) {
	messages := make([]string, 0, 2)

	cachedUser, err := u.findUserToRegister(ctx, ID, gennull.GenericNull[profile.RegisterStep]{})
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

	cachedUser, err := u.findUserToRegister(ctx, ID, gennull.NewGenericNull(profile.RegisterStepFullName))
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	if cachedUser.Step != profile.RegisterStepFullName {
		return message, profile.StepValidationErrorMessages[cachedUser.Step]
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

	cachedUser, err := u.findUserToRegister(ctx, ID, gennull.NewGenericNull(profile.RegisterStepBirthDate))
	if err != nil {
		log.Println("не удалось найти пользователя: ", err)
		return message, global.ErrInternalError
	}

	if cachedUser.Step != profile.RegisterStepBirthDate {
		return message, profile.StepValidationErrorMessages[cachedUser.Step]
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
	lf := logrus.Fields{
		"tg_id": ID,
	}

	if contact.PhoneNumber == "" {
		u.log.File.WithFields(lf).Warningln(u.logPrefix(), "пользователь не поделился номером телефона")
		return message, profile.ErrPhoneNumberEmpty
	}

	cachedUser, err := u.findUserToRegister(ctx, ID, gennull.NewGenericNull(profile.RegisterStepPhoneNumber))
	if err != nil {
		return message, global.ErrInternalError
	}

	if cachedUser.Step != profile.RegisterStepPhoneNumber {
		return message, profile.StepValidationErrorMessages[cachedUser.Step]
	}

	lf["register_step"] = cachedUser.Step

	cachedUser.PhoneNumber = sql_null.NewString(contact.PhoneNumber)

	ts := transaction.MustGetSession(ctx)

	err = u.ri.Repository.Profile.MarkUserToRegiserAsRegistered(ts, ID)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось изменить флаг registered не зарегестрированного пользователя:", err)
		return message, global.ErrInternalError
	}

	err = u.ri.Profile.CreateNewUser(ts, cachedUser.NewUser(time.Now()))
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось создать данные в таблице зарегестрированных пользователей:", err)
		return message, global.ErrInternalError
	}

	err = u.ri.Repository.UserCache.DeleteUserToRegister(ctx, ID)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось удалить временные данные пользователя в кеше:", err)
		return message, global.ErrInternalError
	}

	err = u.ri.Repository.UserCache.SetRegisteredUserID(ctx, ID)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось сохранить ID пользователя в кеш:", err)
		return message, global.ErrInternalError
	}

	message = global.NewReplyMessage(
		profile.RegistrationWasSuccessful,
		[]models.KeyboardButton{
			{
				Text: global.TextCommandProfile,
			},
		},
	)

	return message, nil
}

func (u *Profile) HandleStepsValidationMessages(ctx context.Context, ID int64) error {
	cachedUser, err := u.findUserToRegister(ctx, ID, gennull.GenericNull[profile.RegisterStep]{})
	if err != nil {
		return global.ErrInternalError
	}

	return profile.StepValidationErrorMessages[cachedUser.Step]
}

func (u *Profile) HandlerProfileInfo(ctx context.Context, ID int64) (message string, err error) {
	lf := logrus.Fields{
		"tg_id": ID,
	}

	var userData profile.User

	userData, err = u.ri.Repository.UserCache.GetTempUserData(ctx, ID)
	switch err {
	case nil:
	case global.ErrNoData:
		ts := transaction.MustGetSession(ctx)

		userData, err = u.ri.Repository.Profile.FindUserByTGID(ts, ID)
		if err != nil {
			return message, err
		}

		err = u.ri.UserCache.SetTempUserData(ctx, userData)
		if err != nil {
			u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти данные пользователя в базе:", err)
			return message, global.ErrInternalError
		}

	default:
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти данные пользователя из кеша:", err)
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

func (u *Profile) findUserToRegister(ctx context.Context, ID int64, step gennull.GenericNull[profile.RegisterStep]) (profile.UserToRegister, error) {
	lf := logrus.Fields{
		"tg_id":         ID,
		"register_step": step.Value,
	}

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
			u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти пользователя в таблице не зарегестрированных:", err)
			return dbUser, global.ErrInternalError
		}

	default:
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось найти не зарегестрированного пользователя в кеше:", err)
		return cachedUser, global.ErrInternalError
	}
}

func (u *Profile) setUserToRegister(ctx context.Context, user profile.UserToRegister, update bool) error {
	lf := logrus.Fields{
		"tg_id":         user.ID,
		"register_step": user.Step,
	}

	err := u.ri.Repository.UserCache.SetUserToRegister(ctx, user)
	if err != nil {
		u.log.Db.WithFields(lf).Warningln(u.logPrefix(), "не удалось записать не зарегестрированного пользователя в кеш:", err)
	}

	dbFunc := u.ri.Repository.Profile.SaveUserToRegiser
	if update {
		dbFunc = u.ri.Repository.Profile.UpdateUserToRegister
	}

	ts := transaction.MustGetSession(ctx)
	err = dbFunc(ts, user)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln(u.logPrefix(), "не удалось обновить данные пользователя в таблице не зарегестрированных:", err)
		return global.ErrInternalError
	}

	return nil
}
