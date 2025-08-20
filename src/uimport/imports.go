package uimport

import (
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/usecase"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"

	"github.com/go-telegram/bot"
)

type UsecaseImport struct {
	Usecase
}

func NewUsecaseImport(
	b *bot.Bot,
	ri *rimport.RepositoryImports,
	log *logger.Logger,
	conf *config.Config,
) *UsecaseImport {
	return &UsecaseImport{
		Usecase: Usecase{
			Profile:       usecase.NewProfile(ri, log),
			NotifyMessage: usecase.NewNotifyMessage(b, ri, log, conf),
			Kicker:        usecase.NewKickerUsecase(b, ri, log, conf),
		},
	}
}
