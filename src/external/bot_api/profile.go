package bot_api

import (
	"context"
	"tg_video_lessons_bot/external/bot_api/middleware"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/tools/bot_tool"
	"tg_video_lessons_bot/tools/str"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ProfileBotApi struct {
	b  *bot.Bot
	ui *uimport.UsecaseImport
	m  *middleware.AuthMiddleware
}

func NewPrfileBotApi(
	b *bot.Bot,
	ui *uimport.UsecaseImport,
	m *middleware.AuthMiddleware,
) {
	api := ProfileBotApi{b, ui, m}

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		global.CommandStart,
		bot.MatchTypeExact,
		api.StartHandler,
		// middleware
		api.m.AllreadyRegistered,
	)

	api.b.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		str.FullNameRegexp,
		api.FullNameHandler,
		// middleware
		api.m.AllreadyRegistered,
	)

	api.b.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		str.BirthDateRegexp,
		api.BirthDateHandler,
		// middleware
		api.m.AllreadyRegistered,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypeExact,
		api.PhoneNumberHandler,
		// middleware
		api.m.AllreadyRegistered,
		api.m.IsContactShared,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		api.AnyHandler,
		// middleware
		api.m.AllreadyRegistered,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		global.CommandProfile,
		bot.MatchTypeExact,
		api.HandlerProfile,
	)
}

func (e *ProfileBotApi) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	messages, err := e.ui.Usecase.Profile.HandlerStart(ctx, update.Message.From.ID, update.Message.From.Username)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	for _, message := range messages {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)
	}
}

func (e *ProfileBotApi) FullNameHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := e.ui.Usecase.Profile.HandlerFullName(ctx, update.Message.From.ID, update.Message.Text)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)
}

func (e *ProfileBotApi) BirthDateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := e.ui.Usecase.Profile.HandleBirthDate(ctx, update.Message.From.ID, update.Message.Text)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	bot_tool.SendReplyKeyboardMessage(ctx, b, update, message, true)
}

func (e *ProfileBotApi) PhoneNumberHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := e.ui.Usecase.Profile.HandlePhoneNumber(ctx, update.Message.From.ID, *update.Message.Contact)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	bot_tool.SendReplyKeyboardMessage(ctx, b, update, message, false)
}

func (e *ProfileBotApi) AnyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := e.ui.Usecase.Profile.HandleStepsValidationMessages(ctx, update.Message.From.ID)
	bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
}

func (e *ProfileBotApi) HandlerProfile(ctx context.Context, b *bot.Bot, update *models.Update) {}
