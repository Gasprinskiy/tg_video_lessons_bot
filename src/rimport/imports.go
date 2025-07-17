package rimport

import (
	"context"
	"log"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/repository/redis_cache"

	"github.com/redis/go-redis/v9"
)

type RepositoryImports struct {
	RedisDB *redis.Client
	Config  *config.Config
	Repository
}

func NewRepositoryImports(config *config.Config) *RepositoryImports {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Panic("ошибка при пинге redis: ", err)
	}

	return &RepositoryImports{
		RedisDB: rdb,
		Repository: Repository{
			UserCache: redis_cache.NewUserCache(rdb, config.RedisTtl),
		},
	}
}
