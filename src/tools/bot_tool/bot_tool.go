package bot_tool

import (
	"context"
	"tg_video_lessons_bot/internal/entity/global"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SendHTMLParseModeMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	message string,
) {
	var ID int64

	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "🔴",
		})
		ID = update.CallbackQuery.From.ID
	} else {
		ID = update.Message.From.ID
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    ID,
		Text:      message,
		ParseMode: "HTML",
	})
}

func SendHTMLParseModeMessageDeleteMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	message string,
) {
	var (
		fromID    int64
		messageID int
	)

	if update.CallbackQuery != nil {
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "🔴",
		})
		fromID = update.CallbackQuery.From.ID
		messageID = update.CallbackQuery.Message.Message.ID
	} else {
		fromID = update.Message.From.ID
		messageID = update.Message.ID
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    fromID,
		Text:      message,
		ParseMode: "HTML",
	})

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    fromID,
		MessageID: messageID,
	})
}

func SendReplyKeyboardMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	replyMessage global.ReplyMessage,
	closeAfterClick bool,
) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   replyMessage.Message,
		ReplyMarkup: &models.ReplyKeyboardMarkup{
			Keyboard:        replyMessage.ButtonList,
			ResizeKeyboard:  true,
			OneTimeKeyboard: closeAfterClick,
		},
	})
}

func SendInlineKeyboardMarkupMessage(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	replyMessage global.InlineKeyboardMessage,
) {
	var ID int64

	if update.CallbackQuery != nil {
		ID = update.CallbackQuery.From.ID
	} else {
		ID = update.Message.From.ID
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: ID,
		Text:   replyMessage.Message,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: replyMessage.ButtonList,
		},
	})

}
