package transaction

import (
	"context"
	"tg_video_lessons_bot/internal/entity/global"
)

// CONTEXT NOT USED YET

type LoadDataFunc[T any] func(ctx context.Context) (T, error)
type ReturnErrFunc func(ctx context.Context) error

func RunInTx[T any](
	ctx context.Context,
	// log *logrus.Logger,
	sessionManager SessionManager,
	f LoadDataFunc[T],
) (T, error) {
	var zero T

	// создать сессию и открыть транзакцию
	ts := sessionManager.CreateSession()
	if err := ts.Start(); err != nil {
		// log.Errorln(fmt.Sprintf("ошибка открытия транзакции; ошибка: %v", err))
		return zero, global.ErrInternalError
	}
	defer ts.Rollback()

	// положить в контекст сессию и менеджер
	ctx = SetSessionManager(ctx, sessionManager)
	ctx = SetSession(ctx, ts)

	// выполнить функцию
	data, err := f(ctx)
	if err != nil {
		return data, err
	}

	return data, nil
}

func RunInTxCommit[T any](
	ctx context.Context,
	// log *logrus.Logger,
	sessionManager SessionManager,
	f LoadDataFunc[T],
) (T, error) {
	var zero T

	// создать сессию и открыть транзакцию
	ts := sessionManager.CreateSession()
	if err := ts.Start(); err != nil {
		// log.Errorln(fmt.Sprintf("ошибка открытия транзакции; ошибка: %v", err))
		return zero, global.ErrInternalError
	}
	defer ts.Rollback()

	// положить в контекст сессию и менеджер
	ctx = SetSessionManager(ctx, sessionManager)
	ctx = SetSession(ctx, ts)

	// выполнить функцию
	data, err := f(ctx)
	if err != nil {
		return data, err
	}

	// закоммитить транзакцию
	if err := ts.Commit(); err != nil {
		// log.Errorln(fmt.Sprintf("ошибка коммита транзакции; ошибка: %v", err))
		return zero, global.ErrInternalError
	}

	return data, nil
}

func RunInTxExec(
	ctx context.Context,
	// log *logrus.Logger,
	sessionManager SessionManager,
	f ReturnErrFunc,
) error {
	// выполнить функцию, результат которой не нужен
	_, err := RunInTxCommit(ctx, sessionManager, func(ctx context.Context) (struct{}, error) {
		err := f(ctx)
		return struct{}{}, err
	})

	return err
}
