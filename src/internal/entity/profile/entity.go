package profile

import (
	"fmt"
	"math"
	"tg_video_lessons_bot/tools/chronos"
	"time"
)

type User struct {
	ID          int64     `json:"id" db:"tg_id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	UserName    string    `json:"user_name" db:"tg_user_name"`
	BirthDate   time.Time `json:"birth_date" db:"birth_date"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
}

type UserToRegiser struct {
	User
	RegisterStep `json:"register_step"`
}

func (u User) CalcAge() string {
	durationBetween := chronos.DurationBetween(u.BirthDate, chronos.NowTruncUTC())
	return fmt.Sprintf("%g", math.Round(durationBetween.Hours()/24/365.25))
}
