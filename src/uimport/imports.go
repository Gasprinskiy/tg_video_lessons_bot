package uimport

import (
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/internal/usecase"
	"tg_video_lessons_bot/rimport"
)

type UsecaseImport struct {
	Usecase
}

func NewUsecaseImport(
	ri *rimport.RepositoryImports,
	sessionManager transaction.SessionManager,
) *UsecaseImport {
	return &UsecaseImport{
		Usecase: Usecase{
			Profile: usecase.NewProfile(ri, sessionManager),
		},
	}
}
