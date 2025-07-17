package global

import "errors"

var (
	// ErrNoData данные не найдены"
	ErrNoData = errors.New("данные не найдены")
	// ErrInternalError внутряя ошибка
	ErrInternalError = errors.New("произошла внутреняя ошибка")
)

var MessagesByError = map[error]string{
	ErrInternalError: "Внутреняя ошибка бота, попробуйте позже или свяжитесь с поддержкой.",
}
