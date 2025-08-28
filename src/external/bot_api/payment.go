package bot_api

import (
	"context"
	"tg_video_lessons_bot/external/bot_api/middleware"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/payment"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/bot_tool"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

type PayemntBotApi struct {
	b   *bot.Bot
	ui  *uimport.UsecaseImport
	m   *middleware.AuthMiddleware
	sm  transaction.SessionManager
	log *logger.Logger
}

func NewPayemntBotApi(
	b *bot.Bot,
	ui *uimport.UsecaseImport,
	m *middleware.AuthMiddleware,
	sm transaction.SessionManager,
	log *logger.Logger,
) {
	api := PayemntBotApi{
		b,
		ui,
		m,
		sm,
		log,
	}

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		global.TextCommandBuySub,
		bot.MatchTypeExact,
		api.HandleBuySubscription,
		// middleware
		api.m.IsRegsitered,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		payment.PickSubTypePrefix,
		bot.MatchTypePrefix,
		api.HandleSubscriptionType,
		// middleware
		api.m.IsRegsitered,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		payment.PickMethodPrefix,
		bot.MatchTypePrefix,
		api.HandlePaymentType,
		// middleware
		api.m.IsRegsitered,
	)
}

func (h *PayemntBotApi) HandleBuySubscription(ctx context.Context, b *bot.Bot, update *models.Update) {
	lf := logrus.Fields{
		"tg_id": update.Message.From.ID,
	}

	ts := h.sm.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		h.log.Db.WithFields(lf).Errorln("не удалось запустить транзакцию: ", err)
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
	}

	defer ts.Rollback()

	message, err := h.ui.Payment.CreatePickSubsTypeMessage(ctx)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	bot_tool.SendInlineKeyboardMarkupMessage(ctx, h.b, update, message)
}

func (h *PayemntBotApi) HandleSubscriptionType(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := h.ui.Payment.CreatePaymentTypesMessage(update.CallbackQuery.Data, update.CallbackQuery.From.ID)
	if err != nil {
		bot_tool.SendHTMLParseModeMessageDeleteMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            message.CallbackQueryAnswerMessage,
	})

	bot_tool.SendInlineKeyboardMarkupMessage(ctx, h.b, update, message)

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.From.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
}

func (h *PayemntBotApi) HandlePaymentType(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := h.ui.Payment.CreatePaymentBill(ctx, update.CallbackQuery.Data, update.CallbackQuery.From.ID)
	if err != nil {
		bot_tool.SendHTMLParseModeMessageDeleteMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            message.CallbackQueryAnswerMessage,
	})

	bot_tool.SendInlineKeyboardMarkupMessage(ctx, h.b, update, message)

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.From.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	})
}
