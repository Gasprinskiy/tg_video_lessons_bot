package uimport

import (
	"tg_video_lessons_bot/internal/usecase"
	"tg_video_lessons_bot/rimport"
)

type UsecaseImport struct {
	Usecase
}

func NewUsecaseImport(ri *rimport.RepositoryImports) *UsecaseImport {
	return &UsecaseImport{
		Usecase: Usecase{
			Profile: usecase.NewProfile(ri),
		},
	}
}
