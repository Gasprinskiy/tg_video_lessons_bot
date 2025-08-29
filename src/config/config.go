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
	PostgresURL      string
	BotToken         string
	BotChanelID      int
	BotUserName      string
	BotAdminUsername string
	RedisAddr        string
	RedisPass        string
	RedisTtl         time.Duration
	RedisPaymentTtl  time.Duration
	PaymeMerchantID  string
	GrpcPort         string
	IsDev            bool
}

// NewConfig загружает переменные из .env и возвращает структуру Config
func NewConfig() *Config {
	redisTtl, err := strconv.Atoi(os.Getenv("REDIS_TTL"))
	if err != nil {
		log.Panic("не удалось получить время жизни кеша: ", err)
	}

	redisPaymentTtl, err := strconv.Atoi(os.Getenv("REDIS_PAYMENT_TTL"))
	if err != nil {
		log.Panic("не удалось получить время жизни кеша платежных чеков: ", err)
	}

	botChanelID, err := strconv.Atoi(os.Getenv("BOT_CHANEL_ID"))
	if err != nil {
		log.Panic("не удалось id платного канала: ", err)
	}

	isDev, err := strconv.ParseBool(os.Getenv("IS_DEV"))
	if err != nil {
		isDev = false
	}

	return &Config{
		PostgresURL:      os.Getenv("POSTGRES_URL"),
		BotToken:         os.Getenv("BOT_TOKEN"),
		BotChanelID:      botChanelID,
		BotAdminUsername: os.Getenv("BOT_ADMIN_USERNAME"),
		RedisPass:        os.Getenv("REDIS_PASSWORD"),
		BotUserName:      os.Getenv("BOT_USERNAME"),
		PaymeMerchantID:  os.Getenv("PAYME_MERHCANT_ID"),
		RedisAddr:        fmt.Sprintf("redis:%s", os.Getenv("REDIS_PORT")),
		RedisTtl:         time.Minute * time.Duration(redisTtl),
		RedisPaymentTtl:  time.Minute * time.Duration(redisPaymentTtl),
		GrpcPort:         fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")),
		IsDev:            isDev,
	}
}
