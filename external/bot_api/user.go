package bot_api

import (
	"context"
	"fmt"
	"strings"
	"tg_video_lessons_bot/internal/entity/user"
	"tg_video_lessons_bot/tools/dump"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserBotApi struct {
	b *bot.Bot
}

func NewUserBotApi(b *bot.Bot) {
	api := UserBotApi{b}

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		StartHandler,
		[]bot.Middleware{
			// singleFlight,
			allreadyRegistered,
		}...,
	)

	api.b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		AnyHandler,
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

type MessageHandlerFunc = func(ctx context.Context, b *bot.Bot, update *models.Update)

var TempUserMap = make(map[int64]user.UserToRegiser, 10)
var RegisteredUsers = make(map[int64]user.UserToRegiser, 10)

var StepsHanders = map[user.RegisterStep]MessageHandlerFunc{
	user.RegisterStepFullName: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Введите ваше имя и фамилию",
		})
	},

	user.RegisterStepBirthDate: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Укажите дату рождения",
		})
	},
}

var SteptsValidation = map[user.RegisterStep]MessageHandlerFunc{
	user.RegisterStepFullName: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		from := update.Message.From
		text := strings.TrimSpace(update.Message.Text)

		cachedUser := TempUserMap[from.ID]

		splitted := strings.Split(text, " ")
		cachedUser.FirstName = splitted[0]
		cachedUser.LastName = splitted[1]

		cachedUser.RegisterStep = user.RegisterStepBirthDate

		TempUserMap[from.ID] = cachedUser

		StepsHanders[cachedUser.RegisterStep](ctx, b, update)
	},

	user.RegisterStepBirthDate: func(ctx context.Context, b *bot.Bot, update *models.Update) {
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
	},
}

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From

	cachedUser, exists := TempUserMap[from.ID]
	if !exists {
		cachedUser = user.UserToRegiser{
			ID:           from.ID,
			RegisterStep: user.RegisterStepFullName,
		}
		TempUserMap[from.ID] = cachedUser
	}

	step := cachedUser.RegisterStep
	fmt.Println("step: ", step)

	StepsHanders[step](ctx, b, update)
}

func AnyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From
	user := TempUserMap[from.ID]

	fmt.Println("user: ", dump.Struct(user))

	SteptsValidation[user.RegisterStep](ctx, b, update)
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
