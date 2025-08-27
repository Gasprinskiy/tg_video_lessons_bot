package rimport

import "tg_video_lessons_bot/internal/repository"

type Repository struct {
	UserCache        repository.UserCache
	PaymentBillCache repository.PaymentBillCache
	Profile          repository.Profile
	Subscritions     repository.Subscritions
}
