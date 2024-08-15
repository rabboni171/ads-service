package logger

import (
	"io"
	"log/slog"
	"os"
)

func MustLoad(cfgPath string) {
	// Создаем файл для записи логов logs.txt
	outfile, err := os.Create("logs.txt")
	if err != nil {
		panic("creating file for logs: " + err.Error())
	}

	var logger *slog.Logger

	switch cfgPath {
	case "local":
		logger = slog.New(slog.NewTextHandler(io.MultiWriter(outfile, os.Stdout),
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		logger = slog.New(slog.NewJSONHandler(io.MultiWriter(outfile, os.Stdout),
			&slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("logger, unknown env path: " + cfgPath)
	}

	// Заменяем стандартный логгер log и slog на настроенный нами
	slog.SetDefault(logger)
}
