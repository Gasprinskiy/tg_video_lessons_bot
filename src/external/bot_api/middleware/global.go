package middleware

import (
	"context"
	"tg_video_lessons_bot/internal/repository"
	"tg_video_lessons_bot/internal/transaction"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type GlobalMiddleware struct {
	userCache repository.UserCache
	userRepo  repository.Profile
	sm        transaction.SessionManager
}

func NewGlobalMiddleware(
	userCache repository.UserCache,
	userRepo repository.Profile,
	sm transaction.SessionManager,
) *GlobalMiddleware {
	return &GlobalMiddleware{userCache, userRepo, sm}
}

func OnBotBlocked(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		isOldStatusMember := update.MyChatMember.OldChatMember.Type == models.ChatMemberTypeMember
		isNewStatusBanned := update.MyChatMember.NewChatMember.Type == models.ChatMemberTypeBanned
		isPrivate := update.MyChatMember.Chat.Type == models.ChatTypePrivate

		isUserBannedBot := isOldStatusMember && isNewStatusBanned && isPrivate

		if isUserBannedBot {

		}
		next(ctx, b, update)
	}
}
