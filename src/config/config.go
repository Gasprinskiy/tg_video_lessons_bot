package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Config структура для хранения переменных окружения
type Config struct {
	PostgresURL string
	BotToken    string
	RedisAddr   string
	RedisPass   string
	RedisTtl    time.Duration
}

// NewConfig загружает переменные из .env и возвращает структуру Config
func NewConfig() *Config {
	redisTtl, err := strconv.Atoi(os.Getenv("REDIS_TTL"))
	if err != nil {
		log.Panic("не удалось получить время жизни кеша: ", err)
	}

	return &Config{
		PostgresURL: os.Getenv("POSTGRES_URL"),
		BotToken:    os.Getenv("BOT_TOKEN"),
		RedisPass:   os.Getenv("REDIS_PASSWORD"),
		RedisAddr:   fmt.Sprintf("redis:%s", os.Getenv("REDIS_PORT")),
		RedisTtl:    time.Minute * time.Duration(redisTtl),
	}
}
