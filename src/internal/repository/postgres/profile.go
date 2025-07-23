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
		INSERT INTO bot_users_profile (tg_id, first_name, last_name, tg_user_name, birth_date, phone_number, join_date, register_date)
		VALUES (:tg_id, :first_name, :last_name, :tg_user_name, :birth_date, :phone_number, :join_date, :register_date)
	`
	return sql_gen.ExecNamed(SqlxTx(ts), sqlQuery, user)
}

func (r *profileRepo) SaveUserToRegiser(ts transaction.Session, user profile.UserToRegister) error {
	sqlQuery := `
		INSERT INTO bot_users_to_register (tg_id, tg_user_name, register_step, join_date, registered)
		VALUES (:tg_id, :tg_user_name, :register_step, :join_date, false)
	`
	return sql_gen.ExecNamed(SqlxTx(ts), sqlQuery, user)
}

func (r *profileRepo) UpdateUserToRegister(ts transaction.Session, user profile.UserToRegister) error {
	sqlQuery := `
		UPDATE bot_users_to_register
		SET
			register_step = :register_step,
			first_name = :first_name,
			last_name = :last_name,
			birth_date = :birth_date
		WHERE tg_id = :tg_id
	`

	_, err := SqlxTx(ts).NamedExec(sqlQuery, user)
	return err
}

func (r *profileRepo) MarkUserToRegiserAsRegistered(ts transaction.Session, ID int64) error {
	sqlQuery := `
		UPDATE bot_users_to_register
		SET registered = true
		WHERE tg_id = $1
	`

	_, err := SqlxTx(ts).Exec(sqlQuery, ID)
	return err
}

func (r *profileRepo) FindUserToRegiserByTGID(ts transaction.Session, ID int64) (profile.UserToRegister, error) {
	sqlQuery := `
		SELECT
			utr.tg_id,
			utr.tg_user_name,
			utr.join_date,
			utr.register_step,
			utr.first_name,
			utr.last_name,
			utr.birth_date
		FROM bot_users_to_register utr
		WHERE utr.tg_id = $1
	`

	return sql_gen.Get[profile.UserToRegister](SqlxTx(ts), sqlQuery, ID)
}

func (r *profileRepo) FindUserByTGID(ts transaction.Session, ID int64) (profile.User, error) {
	sqlQuery := `
		SELECT
			up.tg_id,
			up.first_name,
			up.last_name,
			up.tg_user_name,
			up.birth_date,
			up.phone_number,
			up.register_date
		FROM bot_users_profile up
		WHERE up.tg_id = $1
	`

	return sql_gen.Get[profile.User](SqlxTx(ts), sqlQuery, ID)
}

func (r *profileRepo) LoadAllActiveUserIDS(ts transaction.Session) ([]int64, error) {
	sqlQuery := `
		SELECT
			up.tg_id
		FROM bot_users_profile up
	`

	return sql_gen.Select[int64](SqlxTx(ts), sqlQuery)
}
