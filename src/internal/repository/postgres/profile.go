package postgres

import (
	"tg_video_lessons_bot/internal/entity/profile"
	"tg_video_lessons_bot/internal/repository"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/sql_gen"
	"time"
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
		WITH purchase AS (
			SELECT
				bp.p_id,
				bp.u_id,
				bp.p_time,
				bp.kick_time,
				st.term_in_month
			FROM bot_users_purchases bp
				JOIN bot_subscription_types st ON (st.sub_id = bp.sub_id)
			ORDER BY bp.p_id DESC
			LIMIT 1
		)
		SELECT
			up.tg_id,
			up.first_name,
			up.last_name,
			up.tg_user_name,
			up.birth_date,
			up.phone_number,
			up.register_date,
			p.p_id IS NOT NULL as has_purchases,
			p.p_time,
			p.term_in_month
		FROM bot_users_profile up
			LEFT JOIN purchase p ON (p.u_id = up.u_id)
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

func (r *profileRepo) SetPurchaseKickTimeByTGID(ts transaction.Session, date time.Time, reason, uid int) error {
	sqlQuery := `
		UPDATE
			bot_users_purchases
		SET
			kick_time = $2,
			kick_reason = $3
		WHERE u_id = $1
	`

	_, err := SqlxTx(ts).Exec(sqlQuery, uid, date, reason)
	return err
}

func (r *profileRepo) BulkSearchUsersByTGID(ts transaction.Session, tgIDList []int64) ([]profile.User, error) {
	sqlQuery := `
		SELECT
			up.u_id,
			up.tg_id,
			up.first_name,
			up.last_name,
			up.tg_user_name,
			up.birth_date,
			up.phone_number,
			up.register_date
		FROM bot_users_profile up
		WHERE up.tg_id IN (:TG_ID_LIST)
	`

	return sql_gen.SelectNamedIn[profile.User](SqlxTx(ts), sqlQuery, map[string]any{
		"TG_ID_LIST": tgIDList,
	})
}

func (r *profileRepo) GetUserLastSubscrition(ts transaction.Session, ID int64) (profile.UserSubscrition, error) {
	sqlQuery := `
		SELECT
			bp.p_time,
			bp.kick_time,
			st.term_in_month
		FROM bot_users_purchases bp
			JOIN bot_subscription_types st ON (st.sub_id = bp.sub_id)
		WHERE bp.u_id = (SELECT u.u_id FROM bot_users_profile u WHERE u.tg_id = $1)
		ORDER BY bp.p_id DESC
		LIMIT 1
	`

	return sql_gen.Get[profile.UserSubscrition](SqlxTx(ts), sqlQuery, ID)
}
