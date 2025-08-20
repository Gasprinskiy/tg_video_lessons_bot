package profile

import "fmt"

type RegisterStep string

const (
	RegisterStepFullName    RegisterStep = "full_name"
	RegisterStepBirthDate   RegisterStep = "birth_date"
	RegisterStepPhoneNumber RegisterStep = "phone_number"
)

func (s *RegisterStep) Scan(value any) error {
	if b, ok := value.([]byte); ok {
		*s = RegisterStep(string(b))
		return nil
	}
	if sVal, ok := value.(string); ok {
		*s = RegisterStep(sVal)
		return nil
	}
	return fmt.Errorf("не удалось записать %v в тип RegisterStep", value)
}

var (
	HelloMessage = "Assalomu alaykum! Academy of Arabic online o’quv platformasiga hush kelibsiz!\nIsmim Osiyo, men sizning online yordamchingizman."
	// HelloMessage              = "🕌 Ассаляму алейкум!\nДобро пожаловать в бота для изучения арабского языка!\nЗдесь ты сможешь шаг за шагом освоить арабский алфавит, слова, фразы и грамматику."
	SendPhoneNumber           = "📱 Telefon raqamni yuborish"
	RegistrationWasSuccessful = "✅ Yashasin! Muvaffaqiyatli ro’yhatdan o’tdingiz!"
	ProfileInfoMessage        = "ℹ️ <b>Shaxsiy ma’lumot:</b>\n\nIsm: <b>%s</b>\nSharif: <b>%s</b>\nTelefon raqam: <b>%s</b>\nObuna: <b>%s</b>"
)

var StepMessages = map[RegisterStep]string{
	RegisterStepFullName:    "Sizga nima deb murojaat qilay?\nIltimos ism sharifingizni to’liq yozib yuboring!",
	RegisterStepBirthDate:   "Judayam soz! Mamnun bo’ldim! Yana bir nechta qadam qoldi. Tug’ilgan kun,oy va yilingizni kiriting.",
	RegisterStepPhoneNumber: "Minnatdorman! Bittagina qadam qoldi, telefon raqamingizni yuboring!",
}

var StepValidationErrorMessages = map[RegisterStep]error{
	RegisterStepFullName:    ErrFullNameValidation,
	RegisterStepBirthDate:   ErrBirhDateInvalid,
	RegisterStepPhoneNumber: ErrPhoneNumberEmpty,
}
