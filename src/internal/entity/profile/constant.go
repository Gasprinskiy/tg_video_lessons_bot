package profile

type RegisterStep string

const (
	RegisterStepFullName    RegisterStep = "full_name"
	RegisterStepBirthDate   RegisterStep = "birth_date"
	RegisterStepPhoneNumber RegisterStep = "phone_number"
)

var (
	// HelloMessage = map[string]string{
	// 	global.LangCodeRU: "🕌 Ассаляму алейкум!\bДобро пожаловать в бота для изучения арабского языка!\bЗдесь ты сможешь шаг за шагом освоить арабский алфавит, слова, фразы и грамматику.",
	// 	global.LangCodeUZ: "🕌 Assalamu alaykum!\bArab tilini o‘rganish uchun botga xush kelibsiz!\bBu yerda siz harf, so‘z va grammatikani bosqichma-bosqich o‘rganasiz.",
	// }
	HelloMessage              = "🕌 Ассаляму алейкум!\nДобро пожаловать в бота для изучения арабского языка!\nЗдесь ты сможешь шаг за шагом освоить арабский алфавит, слова, фразы и грамматику."
	SendPhoneNumber           = "📱 Отправить номер телефона"
	RegistrationWasSuccessful = "✅ Регистрация прошла успешно. Добро пожаловать!"
	ProfileInfoMessage        = "ℹ️ <b>Информация о пользователе:</b>\n\nИмя: <b>%s</b>\nФамилия: <b>%s</b>\nНомер телефона: <b>%s</b>\nВозраст: <b>%d</b>"
)

var StepMessages = map[RegisterStep]string{
	RegisterStepFullName:    "Введите ваше <b>Имя</b> и <b>Фамилию</b>",
	RegisterStepBirthDate:   "Введите вашу дату рождения в формате <b>ДД.ММ.ГГГГ</b>",
	RegisterStepPhoneNumber: "Отправьте ваш номер телефона",
}

var StepValidationErrorMessages = map[RegisterStep]error{
	RegisterStepFullName:    ErrFullNameValidation,
	RegisterStepBirthDate:   ErrBirhDateInvalid,
	RegisterStepPhoneNumber: ErrPhoneNumberEmpty,
}
