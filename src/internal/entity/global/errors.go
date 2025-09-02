package global

import (
	"errors"
	"tg_video_lessons_bot/internal/entity/payment"
	"tg_video_lessons_bot/internal/entity/profile"
)

var (
	// ErrNoData данные не найдены"
	ErrNoData = errors.New("данные не найдены")
	// ErrInternalError внутряя ошибка
	ErrInternalError = errors.New("произошла внутреняя ошибка")
	//
	ErrPermissionDenied = errors.New("отказано в доступе")
)

// var MessagesByError = map[error]string{
// 	ErrInternalError:              "Внутреняя ошибка бота, попробуйте позже или свяжитесь с поддержкой.",
// 	ErrPermissionDenied:           "Отказано в доступе выполнения команды",
// 	profile.ErrFullNameValidation: "Имя и фамилия введены не верно, попробуйте еще раз",
// 	profile.ErrBirthDateInFuture:  "Дата рождения не может быть в будущем",
// 	profile.ErrBirhDateInvalid:    "Дата рождения введено не верно, попробуйте еще раз",
// 	profile.ErrPhoneNumberEmpty:   "Номер телефона не отправлен или отправлен не верно, попробуйте еще раз",
// 	profile.ErrPhoneNumberInvalid: "Номер телефона введен не верно, попробуйте еще раз",
// }

var MessagesByError = map[error]string{
	ErrInternalError:                "Ichki xato paydo bo'ldi, keyinroq urinib ko‘ring yoki qo‘llab-quvvatlash xizmati bilan bog‘laning.",
	ErrPermissionDenied:             "Rad etdi",
	profile.ErrFullNameValidation:   "Ism va familiya noto‘g‘ri kiritilgan, qaytadan urinib ko‘ring",
	profile.ErrBirthDateInFuture:    "Tug‘ilgan sana kelajakda bo‘lishi mumkin emas",
	profile.ErrBirhDateInvalid:      `Tug‘ilgan sana noto‘g‘ri kiritilgan, qaytadan urinib ko‘ring, "KK.OY.YYYY" formatida kiriting`,
	profile.ErrPhoneNumberEmpty:     "Telefon raqami yuborilmagan yoki noto‘g‘ri yuborilgan, qaytadan urinib ko‘ring",
	profile.ErrPhoneNumberInvalid:   "Telefon raqami noto‘g‘ri kiritilgan, qaytadan urinib ko‘ring",
	payment.ErrAllreadyHasActiveSub: "Sizning obunangiz faol",
}
