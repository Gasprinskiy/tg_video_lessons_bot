package profile

import (
	"fmt"
	"math"
	"tg_video_lessons_bot/tools/chronos"
	"tg_video_lessons_bot/tools/sql_null"
	"time"
)

type User struct {
	UID           int                `json:"-" db:"u_id"`
	ID            int64              `json:"id" db:"tg_id"`
	FirstName     string             `json:"first_name" db:"first_name"`
	LastName      string             `json:"last_name" db:"last_name"`
	UserName      string             `json:"user_name" db:"tg_user_name"`
	BirthDate     time.Time          `json:"birth_date" db:"birth_date"`
	PhoneNumber   string             `json:"phone_number" db:"phone_number"`
	JoinDate      time.Time          `json:"join_date" db:"join_date"`
	RegisterDate  time.Time          `json:"register_date" db:"register_date"`
	HasPurchases  bool               `json:"-" db:"has_purchases"`
	PurchasesTime sql_null.NullTime  `json:"p_time" db:"p_time"`
	PurchasesTerm sql_null.NullInt64 `json:"term_in_month" db:"term_in_month"`
	KickTime      sql_null.NullTime  `json:"kick_time" db:"kick_time"`
}

func (u User) CalcAge() string {
	durationBetween := chronos.DurationBetween(u.BirthDate, chronos.NowTruncUTC())
	return fmt.Sprintf("%g", math.Round(durationBetween.Hours()/24/365.25))
}

func (u User) PurchasesStatus() string {
	if u.PurchasesTime.Valid {
		if u.KickTime.Valid {
			return "Faol emas"
		}

		now := chronos.NowTruncUTC()

		subDate := chronos.BeginingOfDate(u.PurchasesTime.Time)
		expireDate := subDate.AddDate(0, u.PurchasesTerm.GetInt(), 0)

		if now.After(expireDate) {
			return "Faol emas"
		}

		return fmt.Sprintf("%s %s", expireDate.Format(chronos.DateMask), "gacha faol")
	}

	return "Yo‘q"
}

type UserToRegister struct {
	ID          int64               `json:"id" db:"tg_id"`
	UserName    string              `json:"user_name" db:"tg_user_name"`
	JoinDate    time.Time           `json:"join_date" db:"join_date"`
	FirstName   sql_null.NullString `json:"first_name" db:"first_name"`
	LastName    sql_null.NullString `json:"last_name" db:"last_name"`
	BirthDate   sql_null.NullTime   `json:"birth_date" db:"birth_date"`
	PhoneNumber sql_null.NullString `json:"phone_number"`
	Step        RegisterStep        `json:"register_step"  db:"register_step"`
}

func (u UserToRegister) NewUser(registerDate time.Time) User {
	return User{
		ID:           u.ID,
		FirstName:    u.FirstName.String,
		LastName:     u.LastName.String,
		UserName:     u.UserName,
		BirthDate:    u.BirthDate.Time,
		PhoneNumber:  u.PhoneNumber.String,
		JoinDate:     u.JoinDate,
		RegisterDate: registerDate,
	}
}

func NewDefaultUserToRegister(ID int64, userName string, joinDate time.Time) UserToRegister {
	return UserToRegister{
		ID:       ID,
		UserName: userName,
		Step:     RegisterStepFullName,
		JoinDate: joinDate,
	}
}

type UserSubscrition struct {
	PurchasesTime sql_null.NullTime  `db:"p_time"`
	PurchasesTerm sql_null.NullInt64 `db:"term_in_month"`
	KickTime      sql_null.NullTime  `db:"kick_time"`
}

func (s UserSubscrition) IsActive() bool {
	return s.PurchasesTime.Valid && !s.KickTime.Valid
}
