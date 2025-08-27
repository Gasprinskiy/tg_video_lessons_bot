package payment

import (
	"encoding/base64"
	"fmt"
	"tg_video_lessons_bot/internal/entity/contact"
)

type PaymentMethodName string

const (
	PaymentMethodNamePayme PaymentMethodName = "Payme"
	PaymentMethodNameClick PaymentMethodName = "Click"
)

func (pm PaymentMethodName) GeneratePayLink(merchantID, tempID string, tgID int64, subID int, amount float64, botUserName string, isDev bool) string {
	if pm == PaymentMethodNamePayme {
		return pm.generatePaymeLink(merchantID, tempID, tgID, subID, amount, botUserName, isDev)
	}

	return "https://google.com"
}

func (pm PaymentMethodName) generatePaymeLink(merchantID, tempID string, tgID int64, subID int, amount float64, botUserName string, isDev bool) string {
	data := fmt.Sprintf(
		"m=%s;ac.temp_id=%s;ac.sub_id=%d;ac.u_id=%d;a=%f;c=%s",
		merchantID,
		tempID,
		subID,
		tgID,
		amount,
		contact.CreateUserLinkByUsername(botUserName),
	)

	// кодируем в base64
	encoded := base64.StdEncoding.EncodeToString([]byte(data))

	// выбираем URL (боевой или тестовый)
	baseURL := "https://checkout.paycom.uz/"
	if isDev {
		baseURL = "https://test.paycom.uz/"
	}

	return baseURL + encoded
}

var (
	PickSubTypePrefix = "sub_type"
	PickMethodPrefix  = "pay_method"
)

var (
	SubscritionTypeName = func(term int, price float64) string {
		return fmt.Sprintf("%d oy %d so‘m", term, int(price))
	}
	SubscritionTypePrefix = func(subID, term int, price float64) string {
		return fmt.Sprintf("%s:%d:%d:%f", PickSubTypePrefix, subID, term, price)
	}
	PaymentMethodWithPrefix = func(name PaymentMethodName, subID int64, price float64) string {
		return fmt.Sprintf("%s:%s:%d:%f", PickMethodPrefix, name, subID, price)
	}
)

const (
	PickPaymentMethodMessage   = "O‘zingiz to‘lov qiling yoki admin bilan bog‘laning"
	PickSubscritionTypeMessage = "Obuna turini tanlang"
	PaymentLinkMessage         = "To‘lash uchun havolaga o‘ting"
)

const (
	PaymnetLinkButton = "To‘lov havolasi"
)
