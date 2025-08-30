package chanel_kicker

const (
	KickReasonSubscritionExpired = iota + 1
	KickReasonMoneyBack
)

var KickMessageByReason = map[int]string{
	KickReasonSubscritionExpired: `Hurmatli %s, shaxsiy kanalga obunangiz muddati tugadi. Uni istalgan vaqtda "%s" bo'limida yangilashingiz mumkin.`,
	KickReasonMoneyBack:          `Hurmatli %s, obuna bekor qilindi. Istalgan vaqtda "%s" bo'limida qayta obuna bo'lishingiz mumkin`,
}
