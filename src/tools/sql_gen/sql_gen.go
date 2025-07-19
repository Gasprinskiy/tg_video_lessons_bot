package sql_gen

import (
	"database/sql"
	"tg_video_lessons_bot/internal/entity/global"

	"github.com/jmoiron/sqlx"
)

func ExecNamed[T any](tx *sqlx.Tx, sqlQuery string, data T) error {
	_, err := tx.NamedExec(sqlQuery, data)
	return err
}

func Get[T any](tx *sqlx.Tx, sqlQuery string, params ...any) (T, error) {
	var data T

	err := tx.Get(&data, sqlQuery, params...)

	return data, HandleError(err)
}

func Select[T any](tx *sqlx.Tx, sqlQuery string, params ...any) ([]T, error) {
	var data []T

	err := tx.Select(&data, sqlQuery, params...)

	return data, HandleError(err)
}

func HandleError(err error) error {
	if err == sql.ErrNoRows {
		return global.ErrNoData
	}

	return err
}
