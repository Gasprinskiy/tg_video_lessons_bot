package rimport

import "github.com/go-redis/redis/v9"

type RepositoryImports struct {
	RedisDB *redis.Client
	Repository
}

func NewRepositoryImports() *RepositoryImports {
	return &RepositoryImports{
		Repository: Repository{},
	}
}
