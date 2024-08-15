package services

import (
	"ads-service/internal/lib/types"
	"ads-service/internal/models"
	"ads-service/internal/repository"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type adService struct {
	repository *repository.Repository
}

var _ AdServiceInterface = (*adService)(nil)

func (s *adService) GetAll(priceSort string, dateSort string, page int, userId int,
	) ([]*models.Ad, error) {
	ads, err := s.repository.Ad.GetAll(priceSort, dateSort, page, userId)
	if err != nil {
		slog.Error("service.ad.GetAll", "err", err.Error())
		return nil, err
	}

	return ads, nil
}

func (s *adService) Create(ad *models.Ad) (int, error) {
	id, err := s.repository.Ad.Create(ad)
	if err != nil {
		slog.Error("service.ad.Create", "err", err.Error())
		return -1, err
	}

	return id, nil
}

func (s *adService) GetOne(id int) (*models.Ad, error) {
	ad, err := s.repository.Ad.GetOne(id)

	if err == pgx.ErrNoRows {
		return nil, types.ErrAdNotFound
	}

	if err != nil {
		slog.Error("service.ad.GetOne", "err", err.Error())
		return nil, err
	}

	return ad, nil
}
