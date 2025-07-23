package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/external/bot_api"
	"tg_video_lessons_bot/external/bot_api/middleware"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/uimport"
	"time"

	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	loc, err := time.LoadLocation("Asia/Tashkent")
	if err != nil {
		log.Fatalf("cannot load time location: %v", err)
	}
	time.Local = loc

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := config.NewConfig()

	// подключение к redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
	})
	defer rdb.Close()
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Panic("ошибка при пинге redis: ", err)
	}

	// подключение к postgres
	pgdb, err := sqlx.Connect("postgres", config.PostgresURL)
	if err != nil {
		log.Fatalln("не удалось подключиться к базе postgres: ", err)
	}
	defer pgdb.Close()

	if err := pgdb.Ping(); err != nil {
		log.Fatal("ошибка при пинге postgres : ", err)
	}

	// инициализация логгера
	hook := logger.NewPostgresHook(pgdb)
	logger, err := logger.InitLogger(hook)
	if err != nil {
		log.Fatalln("Не удалось инициализировать логгер:", err)
	}

	// инициализация session manager
	sessionManager := transaction.NewSQLSessionManager(pgdb)

	// инициализация бота
	b, err := bot.New(config.BotToken)
	if err != nil {
		log.Panic("ошибка при чтении токена бота: ", err)
	}
	defer b.Close(ctx)

	// инициализация репо
	ri := rimport.NewRepositoryImports(config, rdb)

	// инициализация usecase
	ui := uimport.NewUsecaseImport(ri, logger)

	// инициализация middleware
	mid := middleware.NewAuthMiddleware(
		ctx,
		ri.Repository.UserCache,
		ri.Repository.Profile,
		sessionManager,
	)

	bot_api.NewPrfileBotApi(b, ui, mid, sessionManager, logger)

	log.Println("бот запущен")

	b.Start(ctx)
}
