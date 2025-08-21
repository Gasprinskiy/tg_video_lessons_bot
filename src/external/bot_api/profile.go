package bot_api

import (
	"context"
	"tg_video_lessons_bot/external/bot_api/middleware"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/repository"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/bot_api_gen"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/tools/str"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ProfileBotApi struct {
	b         *bot.Bot
	ui        *uimport.UsecaseImport
	m         *middleware.AuthMiddleware
	sm        transaction.SessionManager
	log       *logger.Logger
	userCache repository.UserCache
}

func NewProfileBotApi(
	b *bot.Bot,
	ui *uimport.UsecaseImport,
	m *middleware.AuthMiddleware,
	sm transaction.SessionManager,
	log *logger.Logger,
	userCache repository.UserCache,
) {
	api := ProfileBotApi{
		b,
		ui,
		m,
		sm,
		log,
		userCache,
	}

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		global.CommandStart,
		bot.MatchTypeExact,
		api.StartHandler,
		// middleware
		// api.m.NotRegistered,
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
		global.TextCommandProfile,
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
	registered, err := e.userCache.HasRegisteredUser(ctx, update.Message.From.ID)
	if registered && err == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   global.CommandStartMessage,
			ReplyMarkup: &models.ReplyKeyboardMarkup{
				Keyboard: global.MainMenuButtons,
			},
		})
		return
	}

	bot_api_gen.HanldeSendMultiplyMessageWitContextSession(
		ctx,
		b,
		update,
		e.sm,
		e.log,
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
		e.log,
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
		e.log,
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
		e.log,
		false,
		func(ctx context.Context) (global.ReplyMessage, error) {
			return e.ui.Usecase.Profile.HandlePhoneNumber(ctx, update.Message.From.ID, *update.Message.Contact)
		},
	)
}

func (e *ProfileBotApi) AnyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	bot_api_gen.HandleSendMessageByErrorMap(
		ctx,
		b,
		update,
		e.sm,
		e.log,
		func(ctx context.Context) error {
			return e.ui.Usecase.Profile.HandleStepsValidationMessages(ctx, update.Message.From.ID)
		},
	)
}

func (e *ProfileBotApi) HandlerProfile(ctx context.Context, b *bot.Bot, update *models.Update) {
	bot_api_gen.HanldeSendMessageWitContextSession(
		ctx,
		b,
		update,
		e.sm,
		e.log,
		func(ctx context.Context) (string, error) {
			return e.ui.Usecase.Profile.HandlerProfileInfo(ctx, update.Message.From.ID)
		},
	)
}
