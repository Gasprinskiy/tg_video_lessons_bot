package rimport

import (
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/internal/repository/postgres"
	"tg_video_lessons_bot/internal/repository/redis_cache"

	"github.com/redis/go-redis/v9"
)

type RepositoryImports struct {
	Repository
}

func NewRepositoryImports(config *config.Config, rdb *redis.Client) *RepositoryImports {
	return &RepositoryImports{
		Repository: Repository{
			UserCache:    redis_cache.NewUserCache(rdb, config.RedisTtl),
			Profile:      postgres.NewProfile(),
			Subscritions: postgres.NewSubscritions(),
		},
	}
}
