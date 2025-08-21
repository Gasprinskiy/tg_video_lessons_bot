package bot_api

import (
	"context"
	"tg_video_lessons_bot/external/bot_api/middleware"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/tools/bot_tool"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ContactBotApi struct {
	b  *bot.Bot
	ui *uimport.UsecaseImport
	m  *middleware.AuthMiddleware
}

func NewContactBotApi(
	b *bot.Bot,
	ui *uimport.UsecaseImport,
	m *middleware.AuthMiddleware,
) {
	api := ContactBotApi{
		b,
		ui,
		m,
	}

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		global.TextCommandContact,
		bot.MatchTypeExact,
		api.HandleContactInfo,
		// middleware
		api.m.IsRegsitered,
	)
}

func (h *ContactBotApi) HandleContactInfo(ctx context.Context, b *bot.Bot, update *models.Update) {
	message := h.ui.Contact.CreateContactsMessage()
	bot_tool.SendInlineKeyboardMarkupMessage(ctx, h.b, update, message)
}
