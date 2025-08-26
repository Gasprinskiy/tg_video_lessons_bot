package postgres

import (
	"tg_video_lessons_bot/internal/entity/subscritions"
	"tg_video_lessons_bot/internal/repository"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/tools/sql_gen"
)

type subscritionsRepo struct{}

func NewSubscritions() repository.Subscritions {
	return &subscritionsRepo{}
}

func (r *subscritionsRepo) LoadSubscritionsList(ts transaction.Session) ([]subscritions.Subscrition, error) {
	sqlQuery := `
		SELECT
			st.sub_id,
			st.term_in_month,
			st.price
		FROM bot_subscription_types st
	`

	return sql_gen.Select[subscritions.Subscrition](SqlxTx(ts), sqlQuery)
}
