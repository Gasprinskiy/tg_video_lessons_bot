package user

import "regexp"

type RegisterStep string

const (
	RegisterStepFullName  RegisterStep = "full_name"
	RegisterStepBirthDate RegisterStep = "birth_date"
)

var (
	UserFullNameRegexp  = regexp.MustCompile(`^(([А-ЯЁA-Z][а-яёa-z]+)\s+([А-ЯЁA-Z][а-яёa-z]+))$`)
	UserBirthDateRegexp = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
)
