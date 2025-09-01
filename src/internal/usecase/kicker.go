package usecase

import (
	"context"
	"fmt"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/external/grpc/proto/kicker"
	"tg_video_lessons_bot/internal/entity/chanel_kicker"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/tools/slice"
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

func (u *KickerUsecase) KickUsersByTGIDList(ctx context.Context, list []*kicker.KickUserParam) error {
	var errCount int

	ts := transaction.MustGetSession(ctx)

	tgIdList := slice.Map(list, func(item *kicker.KickUserParam) int64 {
		return item.TgId
	})

	usersMap := make(map[int64]profile.User, len(tgIdList))

	kickUserList, err := u.ri.Profile.BulkSearchUsersByTGID(ts, tgIdList)
	switch err {
	case nil:
		for _, user := range kickUserList {
			usersMap[user.ID] = user
		}

	case global.ErrNoData:
		return err

	default:
		u.log.Db.Errorln("не удалось найти список пользователей:", err)
		return global.ErrInternalError
	}

	for _, row := range list {
		lf := logrus.Fields{
			"tg_id":     row.TgId,
			"reason_id": row.ReasonId,
		}

		userToKick, exists := usersMap[row.TgId]
		if !exists {
			u.log.Db.WithFields(lf).Errorln("не удалось найти пользователя")
			continue
		}
		_, err := u.b.BanChatMember(ctx, &bot.BanChatMemberParams{
			ChatID: u.conf.BotChanelID,
			UserID: row.TgId,
		})
		if err != nil {
			errCount += 1
			u.log.Db.WithFields(lf).Errorln("не удалось кикнуть пользователя")
		}

		if err = u.ri.Repository.Profile.SetPurchaseKickTimeByTGID(transaction.MustGetSession(ctx), time.Now(), int(row.ReasonId), userToKick.UID); err != nil {
			errCount += 1
			u.log.Db.WithFields(lf).Errorln("не удалось обновить запись в базе у пользователя")
		}

		kickMessage, exists := chanel_kicker.KickMessageByReason[int(row.ReasonId)]
		if !exists {
			u.log.Db.WithFields(lf).Warningln("не верно указана причина")
			continue
		}

		_, err = u.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: row.TgId,
			Text:   fmt.Sprintf(kickMessage, userToKick.FirstName, global.TextCommandBuySub),
		})
		if err != nil {
			u.log.Db.WithFields(lf).Warningln("не удалось отправить сообщение с причиной кика:", err)
		}
	}

	if errCount == len(list) {
		return global.ErrInternalError
	}

	return nil
}
