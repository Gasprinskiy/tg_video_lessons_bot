package usecase

import (
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/entity/contact"
	"tg_video_lessons_bot/internal/entity/global"

	"github.com/go-telegram/bot/models"
)

type Contact struct {
	conf *config.Config
}

func NewContact(conf *config.Config) *Contact {
	return &Contact{conf}
}

func (u *Contact) CreateContactsMessage() global.InlineKeyboardMessage {
	return global.NewInlineKeyboardMessage(
		contact.ContactMessage,
		[]models.InlineKeyboardButton{
			{
				Text: contact.ContactButton,
				URL:  contact.CreateUserLinkByUsername(u.conf.BotAdminUsername),
			},
		},
	)
}
