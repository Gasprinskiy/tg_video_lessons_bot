package profile

import "errors"

var (
	ErrFullNameValidation = errors.New("ошибка при валидации имени и фамилии")
	ErrBirthDateInFuture  = errors.New("ошибка, дата рождения в бужущем")
	ErrBirhDateInvalid    = errors.New("ошибка при парсинге даты рождения")
	ErrPhoneNumberEmpty   = errors.New("номер телефона пустой")
	ErrPhoneNumberInvalid = errors.New("номер телефона введен не верно")
)
