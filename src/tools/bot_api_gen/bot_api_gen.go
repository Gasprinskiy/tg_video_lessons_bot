package bot_api_gen

import (
	"context"
	"log"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/bot_tool"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HanldeSendMessageWitContextSession(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	sessionManager transaction.SessionManager,
	executor func(ctx context.Context) (string, error),
) {
	ts := sessionManager.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		log.Println("не удалось запустить транзакцию: ", err)
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
	}

	message, err := executor(ctx)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		if err := ts.Rollback(); err != nil {
			log.Println("не удалось откатить транзакцию: ", err)
		}
		return
	}

	bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)

	if err := ts.Commit(); err != nil {
		log.Println("не удалось зафиксировать транзакцию: ", err)
	}
}

func HanldeSendMultiplyMessageWitContextSession(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	sessionManager transaction.SessionManager,
	executor func(ctx context.Context) ([]string, error),
) {
	ts := sessionManager.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		log.Println("не удалось запустить транзакцию: ", err)
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
	}

	messages, err := executor(ctx)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		if err := ts.Rollback(); err != nil {
			log.Println("не удалось откатить транзакцию: ", err)
		}
		return
	}

	for _, message := range messages {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)
	}

	if err := ts.Commit(); err != nil {
		log.Println("не удалось зафиксировать транзакцию: ", err)
	}
}

func HanldeSendReplyMessageWitContextSession(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	sessionManager transaction.SessionManager,
	closeAfterClick bool,
	executor func(ctx context.Context) (global.ReplyMessage, error),
) {
	ts := sessionManager.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		log.Println("не удалось запустить транзакцию: ", err)
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
	}

	message, err := executor(ctx)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		if err := ts.Rollback(); err != nil {
			log.Println("не удалось откатить транзакцию: ", err)
		}
		return
	}

	bot_tool.SendReplyKeyboardMessage(ctx, b, update, message, closeAfterClick)

	if err := ts.Commit(); err != nil {
		log.Println("не удалось зафиксировать транзакцию: ", err)
	}
}
