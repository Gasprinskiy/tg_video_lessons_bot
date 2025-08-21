package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"tg_video_lessons_bot/config"
	"tg_video_lessons_bot/external/bot_api"
	"tg_video_lessons_bot/external/bot_api/middleware"
	grpc_extrenal "tg_video_lessons_bot/external/grpc"
	"tg_video_lessons_bot/external/grpc/proto/kicker"
	"tg_video_lessons_bot/external/grpc/proto/notify_message"
	"tg_video_lessons_bot/internal/transaction"
	"tg_video_lessons_bot/rimport"
	"tg_video_lessons_bot/tools/logger"
	"tg_video_lessons_bot/uimport"

	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	var wg sync.WaitGroup

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
	ui := uimport.NewUsecaseImport(b, ri, logger, config)

	// инициализация middleware
	mid := middleware.NewAuthMiddleware(
		ctx,
		ri.UserCache,
		ri.Profile,
		sessionManager,
	)

	bot_api.NewPayemntBotApi(b, ui, mid, sessionManager, logger)
	bot_api.NewContactBotApi(b, ui, mid)
	bot_api.NewProfileBotApi(b, ui, mid, sessionManager, logger, ri.UserCache)

	wg.Add(2)
	// запуск gin
	go func() {
		defer wg.Done()

		lis, err := net.Listen("tcp", config.GrpcPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()

		notifyMessageHandler := grpc_extrenal.NewMessageGrpcHandler(b, ui, sessionManager, logger)
		kickerHadner := grpc_extrenal.NewKickerGrpcHandler(b, ui, sessionManager, logger)

		notify_message.RegisterBotServiceServer(grpcServer, notifyMessageHandler)
		kicker.RegisterKickerServiceServer(grpcServer, kickerHadner)

		log.Printf("gRPC server started on: %s", config.GrpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// запуск api бота
	go func() {
		defer wg.Done()

		log.Println("bot started")

		b.Start(ctx)
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()

	wg.Wait()
}
