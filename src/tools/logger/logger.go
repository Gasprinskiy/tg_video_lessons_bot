package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var FileLogger *logrus.Logger
var CombinedLogger *logrus.Logger

type Logger struct {
	File *logrus.Logger
	Db   *logrus.Logger
}

func InitLogger(pgHook logrus.Hook) (logger *Logger, err error) {
	formatter := &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableQuote:    true,
	}

	// 1. Простой файловый логгер
	fileLog := logrus.New()
	file, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return logger, err
	}
	fileLog.SetOutput(io.MultiWriter(os.Stdout, file))
	fileLog.SetFormatter(formatter)
	fileLog.SetLevel(logrus.InfoLevel)

	// 2. Комбинированный логгер
	combinedLog := logrus.New()
	combinedLog.SetOutput(io.MultiWriter(os.Stdout, file))
	combinedLog.SetFormatter(formatter)
	combinedLog.SetLevel(logrus.InfoLevel)
	combinedLog.AddHook(pgHook)

	return &Logger{
		File: fileLog,
		Db:   combinedLog,
	}, nil
}
