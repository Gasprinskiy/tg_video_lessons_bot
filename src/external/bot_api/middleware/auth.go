package middleware

import (
	"context"
	"log"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AuthMiddleware struct {
	userCache repository.UserCache
}

func NewAuthMiddleware(userCache repository.UserCache) *AuthMiddleware {
	return &AuthMiddleware{userCache}
}

func (m *AuthMiddleware) AllreadyRegistered(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var ID int64

		if update.CallbackQuery != nil {
			ID = update.CallbackQuery.From.ID
		} else {
			ID = update.Message.From.ID
		}

		has, err := m.userCache.HasRegisteredUser(ctx, ID)
		if err != nil && err != global.ErrNoData {
			log.Println("не удалось получить пользователя по ID: ", err)
			return
		}

		if has {
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
