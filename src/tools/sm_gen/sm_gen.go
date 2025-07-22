package sm_gen

import (
	"context"
	"log"
	"tg_video_lessons_bot/internal/entity/global"
	"tg_video_lessons_bot/internal/transaction"
)

func LoadDataIgnoreErrNoData[T any](
	ctx context.Context,
	executor func(ts transaction.Session) (T, error),
	errMsg string,
) (T, error) {
	ts := transaction.MustGetSession(ctx)

	data, err := executor(ts)
	if err != nil {
		log.Printf("%s, ошибка: %v", errMsg, err)
		return data, global.ErrInternalError
	}

	return data, nil
}
