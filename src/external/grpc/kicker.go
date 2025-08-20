package grpc_extrenal

import (
	"context"
	"tg_video_lessons_bot/external/grpc/proto/kicker"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
)

type KickerGrpcHandler struct {
	b   *bot.Bot
	ui  *uimport.UsecaseImport
	sm  transaction.SessionManager
	log *logger.Logger
	kicker.UnimplementedKickerServiceServer
}

func NewKickerGrpcHandler(
	b *bot.Bot,
	ui *uimport.UsecaseImport,
	sm transaction.SessionManager,
	log *logger.Logger,
) *KickerGrpcHandler {
	return &KickerGrpcHandler{
		b:   b,
		ui:  ui,
		sm:  sm,
		log: log,
	}
}

func (h *KickerGrpcHandler) KickExpiredSubsUsers(ctx context.Context, param *kicker.KickExpiredSubsUsersRequest) (*kicker.KickExpiredSubsUsersReply, error) {
	var result kicker.KickExpiredSubsUsersReply

	err := h.ui.Kicker.KickUsersByTGIDList(ctx, param.TgIds)
	if err != nil {
		return &result, err
	}

	result.Done = true
	return &result, nil
}
