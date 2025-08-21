package usecase

import (
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/entity/contact"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/payment"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Payment struct {
	b    *bot.Bot
	ri   *rimport.RepositoryImports
	log  *logger.Logger
	conf *config.Config
}

func NewPayment(
	b *bot.Bot,
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	conf *config.Config,
) *Payment {
	return &Payment{
		b,
		ri,
		log,
		conf,
	}
}

func (u *Payment) CreatePaymentTypesMessage() global.InlineKeyboardMessage {
	return global.NewInlineKeyboardMessage(
		payment.PickPaymentMethodMessage,
		[][]models.InlineKeyboardButton{
			{
				{
					Text:         string(payment.PaymentMethodNamePayme),
					CallbackData: payment.PaymentMethodWithPrefix(payment.PaymentMethodNamePayme),
				},
				{
					Text:         string(payment.PaymentMethodNameClick),
					CallbackData: payment.PaymentMethodWithPrefix(payment.PaymentMethodNameClick),
				},
			},
			{
				{
					Text: contact.ContactButton,
					URL:  contact.CreateUserLinkByUsername(u.conf.BotAdminUsername),
				},
			},
		},
	)
}
