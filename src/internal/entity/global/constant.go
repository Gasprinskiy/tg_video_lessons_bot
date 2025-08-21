package global

import (
	"github.com/go-telegram/bot/models"
)

const (
	LangCodeRU = "ru"
	LangCodeUZ = "uz"
)

const AppLangCode = LangCodeUZ

const (
	CommandStart = "/start"
)

var (
	TextCommandProfile = "Shaxsiy ma’lumot"
	TextCommandContact = "Adminga bog’lanish"
	TextCommandBuySub  = "Obunani rasmiylashtirish"
)

var MainMenuButtons = [][]models.KeyboardButton{
	{
		{
			Text: TextCommandProfile,
		},
	},
	{
		{
			Text: TextCommandBuySub,
		},
	},
	{
		{
			Text: TextCommandContact,
		},
	},
}
