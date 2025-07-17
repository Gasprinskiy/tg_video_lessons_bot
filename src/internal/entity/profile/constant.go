package profile

import (
	"regexp"
)

type RegisterStep string

const (
	RegisterStepFullName  RegisterStep = "full_name"
	RegisterStepBirthDate RegisterStep = "birth_date"
)

var (
	UserFullNameRegexp  = regexp.MustCompile(`^(([А-ЯЁA-Z][а-яёa-z]+)\s+([А-ЯЁA-Z][а-яёa-z]+))$`)
	UserBirthDateRegexp = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
)

var (
	// HelloMessage = map[string]string{
	// 	global.LangCodeRU: "🕌 Ассаляму алейкум!\bДобро пожаловать в бота для изучения арабского языка!\bЗдесь ты сможешь шаг за шагом освоить арабский алфавит, слова, фразы и грамматику.",
	// 	global.LangCodeUZ: "🕌 Assalamu alaykum!\bArab tilini o‘rganish uchun botga xush kelibsiz!\bBu yerda siz harf, so‘z va grammatikani bosqichma-bosqich o‘rganasiz.",
	// }
	HelloMessage    = "🕌 Ассаляму алейкум!\nДобро пожаловать в бота для изучения арабского языка!\nЗдесь ты сможешь шаг за шагом освоить арабский алфавит, слова, фразы и грамматику."
	FullNameMessage = "Введите ваше <b>Имя</b> и <b>Фамилию</b>"
)

var StepMessages = map[RegisterStep]string{
	RegisterStepFullName:  "Введите ваше <b>Имя</b> и <b>Фамилию</b>",
	RegisterStepBirthDate: "Введите вашу дату рождения в формате <b>ДД.ММ.ГГГГ</b>",
}
