package payment

import "errors"

var (
	ErrAllreadyHasActiveSub = errors.New("уже есть активная подписка")
)
