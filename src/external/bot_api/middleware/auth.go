package middleware

import (
	"context"
	"log"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/repository"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/bot_tool"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AuthMiddleware struct {
	userCache repository.UserCache
}

func NewAuthMiddleware(
	ctx context.Context,
	//
	userCache repository.UserCache,
	userRepo repository.Profile,
	sm transaction.SessionManager,
) *AuthMiddleware {

	ts := sm.CreateSession()
	if err := ts.Start(); err != nil {
		log.Panic("не удалось запустить транзакцию: ", err)
	}

	defer ts.Rollback()

	userIDList, err := userRepo.LoadAllActiveUserIDS(ts)
	if err != nil && err != global.ErrNoData {
		log.Panic("ошибка при поиске активных пользовтелей")
	}

	for _, ID := range userIDList {
		err = userCache.SetRegisteredUserID(ctx, ID)
		if err != nil {
			log.Panicf("не удалось записать регистрированного пользователя c ID: %d в кеш, ошибка: %v", ID, err)
		}
	}

	return &AuthMiddleware{userCache}
}

func (m *AuthMiddleware) NotRegistered(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var ID int64

		if update.CallbackQuery != nil {
			ID = update.CallbackQuery.From.ID
		} else {
			ID = update.Message.From.ID
		}

		has, err := m.userCache.HasRegisteredUser(ctx, ID)
		if err != nil {
			log.Printf("не удалось получить пользователя по ID: %d, ошибка: %v", ID, err)
			bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
			return
		}

		if has {
			return
		}

		next(ctx, b, update)
	}
}

func (m *AuthMiddleware) IsRegsitered(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var ID int64

		if update.CallbackQuery != nil {
			ID = update.CallbackQuery.From.ID
		} else {
			ID = update.Message.From.ID
		}

		has, err := m.userCache.HasRegisteredUser(ctx, ID)
		if err != nil {
			log.Printf("не удалось получить пользователя по ID: %d, ошибка: %v", ID, err)
			bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
			return
		}

		if !has {
			bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrPermissionDenied])
			return
		}

		next(ctx, b, update)
	}
}

func (m *AuthMiddleware) IsContactShared(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message.Contact == nil {
			return
		}
		next(ctx, b, update)
	}
}
