package bot_api_gen

import (
	"context"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/bot_tool"
	"tg_video_lessons_bot/tools/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

func HanldeSendMessageWitContextSession(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	sessionManager transaction.SessionManager,
	log *logger.Logger,
	executor func(ctx context.Context) (string, error),
) {
	lf := logrus.Fields{
		"tg_id": update.Message.From.ID,
	}

	ts := sessionManager.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		log.Db.WithFields(lf).Errorln("не удалось запустить транзакцию: ", err)
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
	}

	message, err := executor(ctx)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		if err := ts.Rollback(); err != nil {
			log.Db.WithFields(lf).Errorln("не удалось откатить транзакцию: ", err)
		}
		return
	}

	bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)

	if err := ts.Commit(); err != nil {
		log.Db.WithFields(lf).Errorln("не удалось зафиксировать транзакцию: ", err)
	}
}

func HanldeSendMultiplyMessageWitContextSession(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	sessionManager transaction.SessionManager,
	log *logger.Logger,
	executor func(ctx context.Context) ([]string, error),
) {
	lf := logrus.Fields{
		"tg_id": update.Message.From.ID,
	}

	ts := sessionManager.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		log.Db.WithFields(lf).Errorln("не удалось запустить транзакцию: ", err)
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
	}

	messages, err := executor(ctx)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		if err := ts.Rollback(); err != nil {
			log.Db.WithFields(lf).Errorln("не удалось откатить транзакцию: ", err)
		}
		return
	}

	for _, message := range messages {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, message)
	}

	if err := ts.Commit(); err != nil {
		log.Db.WithFields(lf).Errorln("не удалось зафиксировать транзакцию: ", err)
	}
}

func HanldeSendReplyMessageWitContextSession(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	sessionManager transaction.SessionManager,
	log *logger.Logger,
	closeAfterClick bool,
	executor func(ctx context.Context) (global.ReplyMessage, error),
) {
	lf := logrus.Fields{
		"tg_id": update.Message.From.ID,
	}

	ts := sessionManager.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		log.Db.WithFields(lf).Errorln("не удалось запустить транзакцию: ", err)
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
	}

	message, err := executor(ctx)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		if err := ts.Rollback(); err != nil {
			log.Db.WithFields(lf).Errorln("не удалось откатить транзакцию: ", err)
		}
		return
	}

	bot_tool.SendReplyKeyboardMessage(ctx, b, update, message, closeAfterClick)

	if err := ts.Commit(); err != nil {
		log.Db.WithFields(lf).Errorln("не удалось зафиксировать транзакцию: ", err)
	}
}

func HandleSendMessageByErrorMap(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	sessionManager transaction.SessionManager,
	log *logger.Logger,
	executor func(ctx context.Context) error,
) {
	lf := logrus.Fields{
		"tg_id": update.Message.From.ID,
	}

	ts := sessionManager.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		log.Db.WithFields(lf).Errorln("не удалось запустить транзакцию: ", err)
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[global.ErrInternalError])
	}

	err := executor(ctx)
	if err != nil {
		bot_tool.SendHTMLParseModeMessage(ctx, b, update, global.MessagesByError[err])
		return
	}

	if err := ts.Rollback(); err != nil {
		log.Db.WithFields(lf).Errorln("не удалось откатить транзакцию: ", err)
	}
}
