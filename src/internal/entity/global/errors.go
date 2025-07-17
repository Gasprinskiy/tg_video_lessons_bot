package global

import (
	"errors"
	"tg_video_lessons_bot/internal/entity/profile"
)

var (
	// ErrNoData данные не найдены"
	ErrNoData = errors.New("данные не найдены")
	// ErrInternalError внутряя ошибка
	ErrInternalError = errors.New("произошла внутреняя ошибка")
)

var MessagesByError = map[error]string{
	ErrInternalError:              "Внутреняя ошибка бота, попробуйте позже или свяжитесь с поддержкой.",
	profile.ErrFullNameValidation: "Имя и фамилия введены не верно, попробуйте еще раз",
	profile.ErrBirthDateInFuture:  "Дата рождения не может быть в будущем",
	profile.ErrBirhDateInvalid:    "Дата рождения введено не верно, попробуйте еще раз",
	profile.ErrPhoneNumberEmpty:   "Номер телефона не отправлен или отправлен не верно, попробуйте еще раз",
}
