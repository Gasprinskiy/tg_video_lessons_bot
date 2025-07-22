package bot_api

import (
	"context"
	"tg_video_lessons_bot/external/bot_api/middleware"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/bot_api_gen"
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
	sm transaction.SessionManager
}

func NewPrfileBotApi(
	b *bot.Bot,
	ui *uimport.UsecaseImport,
	m *middleware.AuthMiddleware,
	sm transaction.SessionManager,
) {
	api := ProfileBotApi{b, ui, m, sm}

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		global.CommandStart,
		bot.MatchTypeExact,
		api.StartHandler,
		// middleware
		api.m.NotRegistered,
	)

	api.b.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		str.FullNameRegexp,
		api.FullNameHandler,
		// middleware
		api.m.NotRegistered,
	)

	api.b.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		str.BirthDateRegexp,
		api.BirthDateHandler,
		// middleware
		api.m.NotRegistered,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		global.TextCommandProfile[global.AppLangCode],
		bot.MatchTypeExact,
		api.HandlerProfile,
		// middleware
		api.m.IsRegsitered,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypeExact,
		api.PhoneNumberHandler,
		// middleware
		api.m.NotRegistered,
		api.m.IsContactShared,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		api.AnyHandler,
		// middleware
		api.m.NotRegistered,
	)
}

func (e *ProfileBotApi) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	bot_api_gen.HanldeSendMultiplyMessageWitContextSession(
		ctx,
		b,
		update,
		e.sm,
		func(ctx context.Context) ([]string, error) {
			return e.ui.Usecase.Profile.HandlerStart(ctx, update.Message.From.ID, update.Message.From.Username)
		},
	)
}

func (e *ProfileBotApi) FullNameHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	bot_api_gen.HanldeSendMessageWitContextSession(
		ctx,
		b,
		update,
		e.sm,
		func(ctx context.Context) (string, error) {
			return e.ui.Usecase.Profile.HandlerFullName(ctx, update.Message.From.ID, update.Message.Text)
		},
	)
}

func (e *ProfileBotApi) BirthDateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	bot_api_gen.HanldeSendReplyMessageWitContextSession(
		ctx,
		b,
		update,
		e.sm,
		true,
		func(ctx context.Context) (global.ReplyMessage, error) {
			return e.ui.Usecase.Profile.HandleBirthDate(ctx, update.Message.From.ID, update.Message.Text)
		},
	)
}

func (e *ProfileBotApi) PhoneNumberHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	bot_api_gen.HanldeSendReplyMessageWitContextSession(
		ctx,
		b,
		update,
		e.sm,
		false,
		func(ctx context.Context) (global.ReplyMessage, error) {
			return e.ui.Usecase.Profile.HandlePhoneNumber(ctx, update.Message.From.ID, *update.Message.Contact)
		},
	)
}

func (e *ProfileBotApi) AnyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := e.ui.Usecase.Profile.HandleStepsValidationMessages(ctx, update.Message.From.ID)
	bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
}

func (e *ProfileBotApi) HandlerProfile(ctx context.Context, b *bot.Bot, update *models.Update) {
	bot_api_gen.HanldeSendMessageWitContextSession(
		ctx,
		b,
		update,
		e.sm,
		func(ctx context.Context) (string, error) {
			return e.ui.Usecase.Profile.HandlerProfileInfo(ctx, update.Message.From.ID)
		},
	)
}
