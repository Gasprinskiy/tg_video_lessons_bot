package usecase

import (
	"context"
	"fmt"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/entity/notify_message"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

type NotifyMessage struct {
	b    *bot.Bot
	ri   *rimport.RepositoryImports
	log  *logger.Logger
	conf *config.Config
}

func NewNotifyMessage(
	b *bot.Bot,
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	conf *config.Config,
) *NotifyMessage {
	return &NotifyMessage{b, ri, log, conf}
}

func (u *NotifyMessage) CreateChanelInviteLinkMessage(ctx context.Context, TGID int64) (global.InlineKeyboardMessage, error) {
	var message global.InlineKeyboardMessage

	lf := logrus.Fields{
		"tg_id": TGID,
	}
	ts := transaction.MustGetSession(ctx)

	user, err := u.ri.Repository.Profile.FindUserByTGID(ts, TGID)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln("не удалось найти пользователя по id: ", err)
		return message, global.ErrInternalError
	}

	_, err = u.b.UnbanChatMember(
		ctx,
		&bot.UnbanChatMemberParams{
			ChatID:       u.conf.BotChanelID,
			UserID:       TGID,
			OnlyIfBanned: true,
		},
	)

	if err != nil {
		u.log.Db.WithFields(lf).Errorln("не удалось разбанить пользователя по: ", err)
		return message, global.ErrInternalError
	}

	invite, err := u.b.CreateChatInviteLink(
		ctx,
		&bot.CreateChatInviteLinkParams{
			ChatID:      u.conf.BotChanelID,
			MemberLimit: 1,
		},
	)
	if err != nil {
		u.log.Db.WithFields(lf).Errorln("не удалось создать ссылку приглашение для канала: ", err)
		return message, global.ErrInternalError
	}

	message = global.NewInlineKeyboardMessage(
		fmt.Sprintf(notify_message.InviteMessage, user.FirstName),
		[]models.InlineKeyboardButton{
			{
				Text: "Вступить",
				URL:  invite.InviteLink,
			},
		},
	)

	return message, nil
}
