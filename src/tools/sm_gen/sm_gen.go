package sm_gen

import (
	"log"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
)

func LoadDataIgnoreErrNoData[T any](
	sessionManager transaction.SessionManager,
	executor func(ts transaction.Session) (T, error),
	errMsg string,
) (T, error) {
	var data T

	ts := sessionManager.CreateSession()
	if err := ts.Start(); err != nil {
		log.Println("не удалось запустить транзакцию: ", err)
		return data, global.ErrInternalError
	}

	defer ts.Rollback()

	data, err := executor(ts)
	if err != nil {
		log.Printf("%s, ошибка: %v", errMsg, err)
		return data, global.ErrInternalError
	}

	return data, nil
}

func InTransactionSession(
	sessionManager transaction.SessionManager,
	executor func(ts transaction.Session) error,
) error {
	ts := sessionManager.CreateSession()
	if err := ts.Start(); err != nil {
		log.Println("не удалось запустить транзакцию: ", err)
		return global.ErrInternalError
	}

	defer ts.Commit()

	err := executor(ts)
	if err != nil {
		ts.Rollback()
	}
	return err
}
