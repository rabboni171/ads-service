package repository

import (
	"ads-service/internal/config"
	"ads-service/internal/models"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:generate mockery --name AdRepoInterface --output ./mocks
type AdRepoInterface interface {
	Create(ad *models.Ad) (int, error)
	GetOne(id int) (*models.Ad, error)
	GetAll(priceSort string, dateSort string, page int, userId int) ([]*models.Ad, error)
}


type Repository struct {
	Ad AdRepoInterface
}

func NewRepository(ctx context.Context) *Repository {
	// подключение к postgres
	conn, err := pgx.Connect(context.Background(), config.Cfg.PgUrl)
	if err != nil {

		slog.Error("Unable to connect to database",
			"err", err.Error())
	}
	
	// запуск миграций
	m, err := migrate.New(config.Cfg.MigrationsDir, config.Cfg.PgUrl)
	if err != nil {
		slog.Error("new migrations",
			"err", err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migrations up",
			"err", err.Error())
	}

	var repository = &Repository{
		Ad: NewAdRepo(ctx, conn),
	}

	return repository
}
