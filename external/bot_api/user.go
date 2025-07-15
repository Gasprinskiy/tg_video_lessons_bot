package bot_api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
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

	api.b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"level:",
		bot.MatchTypePrefix,
		LevelHandler,
		[]bot.Middleware{
			// singleFlight,
			allreadyCallbackQuery,
		}...,
	)
}

type Level int

const (
	LevelBeginer      Level = 0
	LevelIntermediate Level = 1
	LevelAdvanced     Level = 2
)

type RegisterStep string

const (
	RegisterStepFullName  RegisterStep = "full_name"
	RegisterStepLevel     RegisterStep = "level"
	RegisterStepBirthDate RegisterStep = "birth_date"
)

type UserToRegiser struct {
	ID        int64
	FirstName string
	LastName  string
	BirthDate time.Time
	RegisterStep
	Level
}

func (u UserToRegiser) HasFullName() bool {
	return u.FirstName != "" && u.LastName != ""
}

type MessageHandlerFunc = func(ctx context.Context, b *bot.Bot, update *models.Update)

var TempUserMap = make(map[int64]UserToRegiser, 10)
var RegisteredUsers = make(map[int64]UserToRegiser, 10)

var StepsHanders = map[RegisterStep]MessageHandlerFunc{
	RegisterStepFullName: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Введите ваше имя и фамилию",
		})
	},

	RegisterStepLevel: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Выбери свой уровень",
			ReplyMarkup: models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{Text: "Начинающий", CallbackData: fmt.Sprintf("level:%d", LevelBeginer)},
					},
					{
						{Text: "Средний", CallbackData: fmt.Sprintf("level:%d", LevelIntermediate)},
					},
					{
						{Text: "Продвинутый", CallbackData: fmt.Sprintf("level:%d", LevelAdvanced)},
					},
				},
			},
		})
	},

	RegisterStepBirthDate: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   "Укажите дату рождения",
		})
	},
}

var SteptsValidation = map[RegisterStep]MessageHandlerFunc{
	RegisterStepFullName: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		from := update.Message.From
		text := strings.TrimSpace(update.Message.Text)

		user := TempUserMap[from.ID]

		splitted := strings.Split(text, " ")
		user.FirstName = splitted[0]
		user.LastName = splitted[1]

		user.RegisterStep = RegisterStepLevel

		TempUserMap[from.ID] = user

		StepsHanders[user.RegisterStep](ctx, b, update)
	},

	RegisterStepLevel: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		from := update.Message.From
		user := TempUserMap[from.ID]
		data := update.CallbackQuery.Data

		levelStr := strings.Split(data, ":")[1]
		level, err := strconv.Atoi(levelStr)
		if err != nil {
			fmt.Println("не удалось конвертировать строку в число: ", err)
			return
		}

		user.Level = Level(level)
		user.RegisterStep = RegisterStepBirthDate
		TempUserMap[from.ID] = user

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Уровень знаний выбран ✅",
			ShowAlert:       false,
		})

		StepsHanders[user.RegisterStep](ctx, b, update)
	},

	RegisterStepBirthDate: func(ctx context.Context, b *bot.Bot, update *models.Update) {
		from := update.Message.From
		text := update.Message.Text

		user := TempUserMap[from.ID]

		parsed, err := time.Parse("02.01.2006", text)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Некорректный формат. Пожалуйста, используй ДД.ММ.ГГГГ.",
			})
			return
		}

		user.BirthDate = parsed
		delete(TempUserMap, from.ID)

		fmt.Println("parsed: ", parsed)

		RegisteredUsers[from.ID] = user

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Регистрация прошла успешно",
		})
	},
}

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var step RegisterStep

	from := update.Message.From

	user, exists := TempUserMap[from.ID]
	if !exists {
		step = RegisterStepFullName
		user = UserToRegiser{
			ID:           from.ID,
			RegisterStep: RegisterStepFullName,
		}
		TempUserMap[from.ID] = user
	} else {
		step = user.RegisterStep
	}
	fmt.Println("step: ", step)

	StepsHanders[step](ctx, b, update)
}

func LevelHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.CallbackQuery.From
	data := update.CallbackQuery.Data

	user := TempUserMap[from.ID]

	levelStr := strings.Split(data, ":")[1]
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		fmt.Println("не удалось конвертировать строку в число: ", err)
		return
	}

	user.Level = Level(level)
	user.RegisterStep = RegisterStepBirthDate
	TempUserMap[from.ID] = user

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		Text:            "Уровень знаний выбран ✅",
		ShowAlert:       false,
	})

	StepsHanders[user.RegisterStep](ctx, b, update)
}

func AnyHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	from := update.Message.From
	user := TempUserMap[from.ID]

	fmt.Println("user: ", dump.Struct(user))

	SteptsValidation[user.RegisterStep](ctx, b, update)
}

// middleware
func singleFlight(next bot.HandlerFunc) bot.HandlerFunc {
	sf := sync.Map{}
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		fmt.Println("CBQ")
		if update.CallbackQuery != nil {
			key := update.CallbackQuery.Message.Message.ID
			if _, loaded := sf.LoadOrStore(key, struct{}{}); loaded {
				return
			}
			defer sf.Delete(key)
			next(ctx, b, update)
		}
	}
}

func allreadyRegistered(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if _, exists := RegisteredUsers[update.Message.From.ID]; exists {
			return
		}
		next(ctx, b, update)
	}
}

func allreadyCallbackQuery(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if _, exists := RegisteredUsers[update.CallbackQuery.From.ID]; exists {
			return
		}
		next(ctx, b, update)
	}
}
