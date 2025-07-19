package postgres

import (
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/internal/repository"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/sql_gen"
)

type profileRepo struct{}

func NewProfile() repository.Profile {
	return &profileRepo{}
}

func (r *profileRepo) CreateNewUser(ts transaction.Session, user profile.User) error {
	sqlQuery := `
		INSERT INTO bot_users_profile (tg_id, first_name, last_name, tg_user_name, birth_date, phone_number)
		VALUES (:tg_id, :first_name, :last_name, :tg_user_name, :birth_date, :phone_number)
	`
	return sql_gen.ExecNamed(SqlxTx(ts), sqlQuery, user)
}

func (r *profileRepo) FindUserByTGID(ts transaction.Session, ID int64) (profile.User, error) {
	sqlQuery := `
		SELECT
			bu.tg_id,
			bu.first_name,
			bu.last_name,
			bu.tg_user_name,
			bu.birth_date,
			bu.phone_number
		FROM bot_users_profile bu
		WHERE bu.tg_id = $1
	`

	return sql_gen.Get[profile.User](SqlxTx(ts), sqlQuery, ID)
}

func (r *profileRepo) LoadAllActiveUserIDS(ts transaction.Session) ([]int64, error) {
	sqlQuery := `
		SELECT
			bu.tg_id
		FROM bot_users_profile bu
		WHERE bu.active = true
	`

	return sql_gen.Select[int64](SqlxTx(ts), sqlQuery)
}
