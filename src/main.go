package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/external/bot_api"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Panic("Ошибка загрузки .env файла: ", err)
	// 	return
	// }

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := config.NewConfig()

	// инициализация бота
	b, err := bot.New(config.BotToken)
	if err != nil {
		log.Panic("ошибка при чтении токена бота: ", err)
	}
	defer b.Close(ctx)

	// инициализация репо
	ri := rimport.NewRepositoryImports(config)
	defer ri.RedisDB.Close()

	// инициализация usecase
	ui := uimport.NewUsecaseImport(ri)

	bot_api.NewPrfileBotApi(b, ui)

	b.Start(ctx)

	fmt.Println("бот запущен")
}
