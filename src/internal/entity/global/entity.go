package global

import "github.com/go-telegram/bot/models"

type ReplyMessage struct {
	Message    string
	ButtonList []models.KeyboardButton
}

func NewReplyMessage(message string, buttonList []models.KeyboardButton) ReplyMessage {
	return ReplyMessage{
		Message:    message,
		ButtonList: buttonList,
	}
}
