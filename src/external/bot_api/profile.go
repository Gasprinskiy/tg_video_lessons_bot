package bot_api

import (
	"context"
	"fmt"
	"strings"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/tools/bot_tool"
	"tg_video_lessons_bot/tools/dump"
	"tg_video_lessons_bot/uimport"
	"time"

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
		[]bot.Middleware{
			// singleFlight,
			allreadyRegistered,
		}...,
	)

	api.b.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		profile.UserFullNameRegexp,
		api.FullNameHandler,
		[]bot.Middleware{
			// singleFlight,
			allreadyRegistered,
		}...,
	)

	api.b.RegisterHandlerRegexp(
		bot.HandlerTypeMessageText,
		profile.UserBirthDateRegexp,
		api.BirthDateHandler,
		[]bot.Middleware{
			// singleFlight,
			allreadyRegistered,
		}...,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		api.AnyHandler,
		[]bot.Middleware{
			// singleFlight,
			allreadyRegistered,
		}...,
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
	from := update.Message.From

	messages, err := e.ui.Usecase.Profile.HandlerStart(ctx, *from)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	for _, message := range messages {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)
	}
}

func (e *ProfileBotApi) FullNameHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From
	text := strings.TrimSpace(update.Message.Text)

	cachedUser := TempUserMap[from.ID]

	splitted := strings.Split(text, " ")
	cachedUser.FirstName = splitted[0]
	cachedUser.LastName = splitted[1]

	cachedUser.RegisterStep = profile.RegisterStepBirthDate

	TempUserMap[from.ID] = cachedUser

	StepsHanders[cachedUser.RegisterStep](ctx, b, update)
}

func (e *ProfileBotApi) AnyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From
	cachedUser := TempUserMap[from.ID]

	fmt.Println("cachedUser: ", dump.Struct(cachedUser))
	// SteptsValidation[cachedUser.RegisterStep](ctx, b, update)
}

func (e *ProfileBotApi) BirthDateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From
	text := update.Message.Text

	cachedUser := TempUserMap[from.ID]

	parsed, err := time.Parse("02.01.2006", text)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Некорректный формат. Пожалуйста, используй ДД.ММ.ГГГГ.",
		})
		return
	}

	cachedUser.BirthDate = parsed
	delete(TempUserMap, from.ID)

	fmt.Println("parsed: ", parsed)

	RegisteredUsers[from.ID] = cachedUser

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Регистрация прошла успешно",
	})
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
