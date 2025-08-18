package grpc_extrenal

import (
	"context"
	messageproto "tg_video_lessons_bot/external/grpc/proto/message"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
)

type MessageGrpcHandler struct {
	b   *bot.Bot
	ui  *uimport.UsecaseImport
	sm  transaction.SessionManager
	log *logger.Logger
	messageproto.UnimplementedBotServiceServer
}

func NewMessageGrpcHandler(
	b *bot.Bot,
	ui *uimport.UsecaseImport,
	sm transaction.SessionManager,
	log *logger.Logger,
) *MessageGrpcHandler {
	return &MessageGrpcHandler{
		b:   b,
		ui:  ui,
		sm:  sm,
		log: log,
	}
}

func (h MessageGrpcHandler) SendInviteLink(ctx context.Context, request *messageproto.SendInviteLinkRequest) (*messageproto.SendInviteLinkReply, error) {
	_, err := h.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: request.TgIs,
		Text:   "FUCK YOU PAL",
	})

	if err != nil {
		return &messageproto.SendInviteLinkReply{
			Sent: false,
		}, global.ErrInternalError
	}

	return &messageproto.SendInviteLinkReply{
		Sent: true,
	}, nil
}
