package bot_api

import (
	"context"
	"tg_video_lessons_bot/external/bot_api/middleware"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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
		api.HandlePaymentMethods,
		// middleware
		api.m.IsRegsitered,
	)
}

func (e *PayemntBotApi) HandlePaymentMethods(ctx context.Context, b *bot.Bot, update *models.Update) {

}
