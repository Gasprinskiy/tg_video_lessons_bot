package usecase

import (
	"context"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/entity/contact"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/payment"
	"tg_video_lessons_bot/internal/transaction"
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

func (u *Payment) logPrefix() string {
	return "[payment]"
}

func (u *Payment) CreatePickSubsTypeMessage(ctx context.Context) (global.InlineKeyboardMessage, error) {
	var message global.InlineKeyboardMessage

	ts := transaction.MustGetSession(ctx)

	subsList, err := u.ri.Subscritions.LoadSubscritionsList(ts)
	if err != nil {
		u.log.Db.Errorln(u.logPrefix(), "не удалось загрузить типы подписок, err", err)
		return message, global.ErrInternalError
	}

	buttons := make([][]models.InlineKeyboardButton, len(subsList))

	for i, sub := range subsList {
		button := models.InlineKeyboardButton{
			Text:         payment.SubscritionTypeName(sub.TermInMonth, sub.Price),
			CallbackData: payment.SubscritionTypePrefix(sub.ID, sub.Price),
		}

		buttons[i] = append(buttons[i], button)
	}

	message = global.NewInlineKeyboardMessage(
		payment.PickSubscritionTypeMessage,
		buttons,
	)

	return message, nil
}

func (u *Payment) CreatePaymentTypesMessage() global.InlineKeyboardMessage {
	return global.NewInlineKeyboardMessage(
		payment.PickPaymentMethodMessage,
		[][]models.InlineKeyboardButton{
			{
				{
					Text:         string(payment.PaymentMethodNamePayme),
					CallbackData: payment.PaymentMethodWithPrefix(payment.PaymentMethodNamePayme, 1, 2),
				},
				{
					Text:         string(payment.PaymentMethodNameClick),
					CallbackData: payment.PaymentMethodWithPrefix(payment.PaymentMethodNameClick, 1, 2),
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
