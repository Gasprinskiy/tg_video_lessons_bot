package grpc_extrenal

import (
	"context"
	"tg_video_lessons_bot/external/grpc/proto/kicker"
	"tg_video_lessons_bot/internal/entity/global"
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

func (h *KickerGrpcHandler) KickUsers(ctx context.Context, param *kicker.KickUsersRequest) (*kicker.KickExpiredSubsUsersReply, error) {
	var result kicker.KickExpiredSubsUsersReply

	ts := h.sm.CreateSession()
	ctx = transaction.SetSession(ctx, ts)

	if err := ts.Start(); err != nil {
		h.log.Db.Errorln("не удалось запустить транзакцию: ", err)
		return &result, global.ErrInternalError
	}

	defer ts.Rollback()

	err := h.ui.Kicker.KickUsersByTGIDList(ctx, param.Params)
	if err != nil {
		return &result, err
	}

	if err = ts.Commit(); err != nil {
		h.log.Db.Errorln("не удалось зафиксировать транзакцию: ", err)
		return &result, global.ErrInternalError
	}

	result.Done = true
	return &result, nil
}
