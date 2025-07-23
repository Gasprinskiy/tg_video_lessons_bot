package logger

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PostgresHook struct {
	DB *sqlx.DB
}

func NewPostgresHook(db *sqlx.DB) *PostgresHook {
	return &PostgresHook{DB: db}
}

func (hook *PostgresHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *PostgresHook) Fire(entry *logrus.Entry) error {
	timestamp := entry.Time
	level := entry.Level.String()
	message := entry.Message

	data := entry.Data

	// Выделяем tg_id
	tgIDValue, ok := data["tg_id"]
	if !ok {
		return fmt.Errorf("tg_id отсутствует в логах")
	}

	var tgID int64
	switch v := tgIDValue.(type) {
	case int64:
		tgID = v
	case int:
		tgID = int64(v)
	case float64:
		tgID = int64(v)
	default:
		return fmt.Errorf("tg_id имеет неподдерживаемый тип: %T", tgIDValue)
	}

	// Копируем всё кроме tg_id в additional_fields
	additionalFields := make(logrus.Fields)
	for k, v := range data {
		if k != "tg_id" {
			additionalFields[k] = v
		}
	}

	additionalJSON, err := json.Marshal(additionalFields)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать additional_fields: %v", err)
	}

	_, err = hook.DB.Exec(`
		INSERT INTO bot_logs (type, message, tg_id, additional_fields, time)
		VALUES ($1, $2, $3, $4, $5)
	`, level, message, tgID, additionalJSON, timestamp)

	return err
}
