package uimport

import (
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/internal/usecase"
	"tg_video_lessons_bot/rimport"
)

type UsecaseImport struct {
	SessionManager transaction.SessionManager
	Usecase
}

func NewUsecaseImport(
	ri *rimport.RepositoryImports,
	sessionManager transaction.SessionManager,
) *UsecaseImport {
	return &UsecaseImport{
		SessionManager: sessionManager,
		Usecase: Usecase{
			Profile: usecase.NewProfile(ri),
		},
	}
}
