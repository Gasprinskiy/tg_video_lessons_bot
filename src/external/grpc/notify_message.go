package grpc_extrenal

import (
	"context"
	"tg_video_lessons_bot/external/grpc/proto/notify_message"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

type NotifyMessageGrpcHandler struct {
	b   *bot.Bot
	ui  *uimport.UsecaseImport
	sm  transaction.SessionManager
	log *logger.Logger
	notify_message.UnimplementedBotServiceServer
}

func NewMessageGrpcHandler(
	b *bot.Bot,
	ui *uimport.UsecaseImport,
	sm transaction.SessionManager,
	log *logger.Logger,
) *NotifyMessageGrpcHandler {
	return &NotifyMessageGrpcHandler{
		b:   b,
		ui:  ui,
		sm:  sm,
		log: log,
	}
}

func (h NotifyMessageGrpcHandler) SendInviteLink(ctx context.Context, request *notify_message.SendInviteLinkRequest) (*notify_message.SendInviteLinkReply, error) {
	var reply notify_message.SendInviteLinkReply

	lf := logrus.Fields{
		"tg_id": request.TgId,
	}

	ts := h.sm.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		h.log.Db.WithFields(lf).Errorln("не удалось запустить транзакцию: ", err)
		return &reply, global.ErrInternalError
	}

	defer ts.Rollback()

	replyMessage, err := h.ui.NotifyMessage.CreateChanelInviteLinkMessage(ctx, request.TgId)
	if err != nil {
		return &reply, global.ErrInternalError
	}

	_, err = h.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: request.TgId,
		Text:   replyMessage.Message,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: replyMessage.ButtonList,
		},
	})

	if err != nil {
		return &reply, global.ErrInternalError
	}

	reply.Sent = true
	return &reply, nil
}
