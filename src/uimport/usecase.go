package uimport

import "tg_video_lessons_bot/internal/usecase"

type Usecase struct {
	Profile       *usecase.Profile
	NotifyMessage *usecase.NotifyMessage
	Kicker        *usecase.KickerUsecase
	Contact       *usecase.Contact
	Payment       *usecase.Payment
}
