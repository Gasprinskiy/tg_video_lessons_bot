package bot_api

import (
	"context"
	"fmt"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/tools/bot_tool"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ProfileBotApi struct {
	b  *bot.Bot
	ui *uimport.UsecaseImport
}

func NewPrfileBotApi(b *bot.Bot, ui *uimport.UsecaseImport) {
	api := ProfileBotApi{b, ui}

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		api.StartHandler,
		// middleware
		allreadyRegistered,
	)

	api.b.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		profile.UserFullNameRegexp,
		api.FullNameHandler,
		// middleware
		allreadyRegistered,
	)

	api.b.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		profile.UserBirthDateRegexp,
		api.BirthDateHandler,
		// middleware
		allreadyRegistered,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypeExact,
		api.AnyHandler,
		// middleware
		allreadyRegistered,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		api.AnyHandler,
		// middleware
		allreadyRegistered,
	)

	// api.b.RegisterHandler(
	// 	bot.HandlerTypeCallbackQueryData,
	// 	"level:",
	// 	bot.MatchTypePrefix,
	// 	LevelHandler,
	// 	[]bot.Middleware{
	// 		// singleFlight,
	// 		allreadyCallbackQuery,
	// 	}...,
	// )
}

var TempUserMap = make(map[int64]profile.UserToRegiser, 10)
var RegisteredUsers = make(map[int64]profile.UserToRegiser, 10)

type MessageHandlerFunc = func(ctx context.Context, b *bot.Bot, update *models.Update)

var StepsHanders = map[profile.RegisterStep]MessageHandlerFunc{
	profile.RegisterStepFullName: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Введите ваше имя и фамилию",
		})
	},

	profile.RegisterStepBirthDate: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Укажите дату рождения",
		})
	},
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

	bot_tool.SendReplyKeyboardMessage(ctx, b, update, message, false)
}

func (e *ProfileBotApi) PhoneNumberHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	message, err := e.ui.Usecase.Profile.HandlePhoneNumber(ctx, update.Message.From.ID, *update.Message.Contact)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)
}

func (e *ProfileBotApi) AnyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("update.Message.Text: ", update.Message.Text)
	// SteptsValidation[cachedUser.RegisterStep](ctx, b, update)
}

// middleware
func allreadyRegistered(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if _, exists := RegisteredUsers[update.Message.From.ID]; exists {
			return
		}
		next(ctx, b, update)
	}
}
