package storage

import (
	"auth-service/internal/config"
	"context"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	Auth AuthStorageInterface
}

func NewStorage(ctx context.Context) *Storage {
	// подключение к postgres
	conn, err := pgx.Connect(context.Background(), config.Cfg.PGUrl)
	if err != nil {
		slog.Error("Unable to connect to database",
			"err", err.Error())
	}

	// запуск миграций postgres
	m, err := migrate.New("file://"+config.Cfg.MigrationsPath, config.Cfg.PGUrl)
	if err != nil {
		slog.Error("new migrations",
			"err", err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migrations up",
			"err", err.Error())
	}

	var storage = &Storage{
		Auth: &AuthStorage{
			conn: conn,
		},
	}

	return storage
}
