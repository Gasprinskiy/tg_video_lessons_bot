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
	BotChanelID int
	RedisAddr   string
	RedisPass   string
	RedisTtl    time.Duration
	GrpcPort    string
}

// NewConfig загружает переменные из .env и возвращает структуру Config
func NewConfig() *Config {
	redisTtl, err := strconv.Atoi(os.Getenv("REDIS_TTL"))
	if err != nil {
		log.Panic("не удалось получить время жизни кеша: ", err)
	}

	botChanelID, err := strconv.Atoi(os.Getenv("BOT_CHANEL_ID"))
	if err != nil {
		log.Panic("не удалось id платного канала: ", err)
	}

	return &Config{
		PostgresURL: os.Getenv("POSTGRES_URL"),
		BotToken:    os.Getenv("BOT_TOKEN"),
		BotChanelID: botChanelID,
		RedisPass:   os.Getenv("REDIS_PASSWORD"),
		RedisAddr:   fmt.Sprintf("redis:%s", os.Getenv("REDIS_PORT")),
		RedisTtl:    time.Minute * time.Duration(redisTtl),
		GrpcPort:    fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")),
	}
}
