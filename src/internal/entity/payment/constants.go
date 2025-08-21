package payment

import (
	"fmt"
	"strings"
)

type PaymentMethodName string

const (
	PaymentMethodNamePayme PaymentMethodName = "Payme"
	PaymentMethodNameClick PaymentMethodName = "Click"
)

var PaymentMethodWithPrefix = func(name PaymentMethodName) string {
	return fmt.Sprintf("pay_method:%s", strings.ToLower(string(name)))
}

const (
	PickPaymentMethodMessage = "O‘zingiz to‘lov qiling yoki admin bilan bog‘laning"
)
