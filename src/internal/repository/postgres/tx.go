package postgres

import (
	"tg_video_lessons_bot/internal/transaction"

	"github.com/jmoiron/sqlx"
)

func SqlxTx(ts transaction.Session) *sqlx.Tx {
	tx := ts.Tx()
	return tx.(*sqlx.Tx)
}
