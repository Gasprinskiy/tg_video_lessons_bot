package uimport

import (
	"tg_video_lessons_bot/internal/usecase"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"
)

type UsecaseImport struct {
	Usecase
}

func NewUsecaseImport(
	ri *rimport.RepositoryImports,
	log *logger.Logger,
) *UsecaseImport {
	return &UsecaseImport{
		Usecase: Usecase{
			Profile: usecase.NewProfile(ri, log),
		},
	}
}
