package usecase

import (
	"context"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"
	"time"

	"github.com/go-telegram/bot"
	"github.com/sirupsen/logrus"
)

type KickerUsecase struct {
	b    *bot.Bot
	ri   *rimport.RepositoryImports
	log  *logger.Logger
	conf *config.Config
}

func NewKickerUsecase(
	b *bot.Bot,
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	conf *config.Config,
) *KickerUsecase {
	return &KickerUsecase{
		b,
		ri,
		log,
		conf,
	}
}

func (u *KickerUsecase) KickUsersByTGIDList(ctx context.Context, idList []int64) error {
	var errCount int

	for _, tgID := range idList {
		lf := logrus.Fields{
			"tg_id": tgID,
		}

		_, err := u.b.BanChatMember(ctx, &bot.BanChatMemberParams{
			ChatID: u.conf.BotChanelID,
			UserID: tgID,
		})

		if err != nil {
			errCount += 1
			u.log.Db.WithFields(lf).Errorln("не удалось кикнуть пользователя")
		}

		if err = u.ri.Repository.Profile.SetPurchaseKickTimeByTGID(transaction.MustGetSession(ctx), time.Now(), tgID); err != nil {
			errCount += 1
			u.log.Db.WithFields(lf).Errorln("не удалось обновить запись в базе у пользователя")
		}
	}

	if errCount == len(idList) {
		return global.ErrInternalError
	}

	return nil
}
