package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"tg_video_lessons_bot/external/bot_api"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Ошибка загрузки .env файла: ", err)
		return
	}

	botToken := os.Getenv("BOT_TOKEN")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(botToken)
	if err != nil {
		log.Panic("ошибка при чтении токена бота: ", err)
		return
	}

	defer b.Close(ctx)

	bot_api.NewUserBotApi(b)

	b.Start(ctx)

}
