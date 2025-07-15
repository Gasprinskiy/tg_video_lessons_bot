package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"tg_video_lessons_bot/external/bot_api"

	"github.com/go-telegram/bot"
	"github.com/redis/go-redis/v9"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Panic("Ошибка загрузки .env файла: ", err)
	// 	return
	// }

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("redis:%s", os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Panic("ошибка при подключении пинге redis: ", err)
	}

	defer rdb.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic("ошибка при чтении токена бота: ", err)
		return
	}

	defer b.Close(ctx)

	bot_api.NewUserBotApi(b)

	b.Start(ctx)

	fmt.Println("бот запущен")
}
