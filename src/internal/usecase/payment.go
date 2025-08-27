package usecase

import (
	"context"
	"strconv"
	"strings"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/entity/contact"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/payment"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
			CallbackData: payment.SubscritionTypePrefix(sub.ID, sub.TermInMonth, sub.PriceInCents()),
		}

		buttons[i] = append(buttons[i], button)
	}

	message = global.NewInlineKeyboardMessage(
		payment.PickSubscritionTypeMessage,
		buttons,
	)

	return message, nil
}

func (u *Payment) CreatePaymentTypesMessage(queryData string, ID int64) (global.InlineKeyboardMessage, error) {
	var message global.InlineKeyboardMessage

	lf := logrus.Fields{
		"tg_id": ID,
	}

	splited := strings.Split(queryData, ":")

	subId, err := strconv.ParseInt(splited[1], 10, 64)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln("ошибка при пасинге id подписки")
		return message, global.ErrInternalError
	}

	term, err := strconv.ParseInt(splited[2], 10, 64)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln("ошибка при пасинге продолжительности подписки")
		return message, global.ErrInternalError
	}

	price, err := strconv.ParseFloat(splited[3], 64)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln("ошибка при пасинге цены подписки")
		return message, global.ErrInternalError
	}

	message = global.NewInlineKeyboardMessage(
		payment.PickPaymentMethodMessage,
		[][]models.InlineKeyboardButton{
			{
				{
					Text:         string(payment.PaymentMethodNamePayme),
					CallbackData: payment.PaymentMethodWithPrefix(payment.PaymentMethodNamePayme, subId, price),
				},
				// {
				// 	Text:         string(payment.PaymentMethodNameClick),
				// 	CallbackData: payment.PaymentMethodWithPrefix(payment.PaymentMethodNameClick, 1, 2),
				// },
			},
			{
				{
					Text: contact.ContactButton,
					URL:  contact.CreateUserLinkByUsername(u.conf.BotAdminUsername),
				},
			},
		},
	)
	message.CallbackQueryAnswerMessage = payment.SubscritionTypeName(int(term), price/100)

	return message, nil
}

func (u *Payment) CreatePaymentBill(ctx context.Context, queryData string, ID int64) (global.InlineKeyboardMessage, error) {
	var message global.InlineKeyboardMessage

	lf := logrus.Fields{
		"tg_id": ID,
	}

	splited := strings.Split(queryData, ":")

	paymentType := payment.PaymentMethodName(splited[1])

	subId, err := strconv.ParseInt(splited[2], 10, 64)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln("ошибка при пасинге id подписки:", err)
		return message, global.ErrInternalError
	}

	price, err := strconv.ParseFloat(splited[3], 64)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln("ошибка при пасинге цены подписки:", err)
		return message, global.ErrInternalError
	}

	tempId := uuid.NewString()

	message = global.NewInlineKeyboardMessage(
		payment.PaymentLinkMessage,
		[][]models.InlineKeyboardButton{
			{
				{
					Text: payment.PaymnetLinkButton,
					URL:  paymentType.GeneratePayLink("321sad", tempId, int(subId), price, u.conf.BotUserName, u.conf.IsDev),
				},
			},
		},
	)
	message.CallbackQueryAnswerMessage = string(paymentType)

	return message, nil
}
