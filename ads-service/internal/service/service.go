package services

import (
	"ads-service/internal/models"
	"ads-service/internal/repository"
)

//go:generate mockery --name AdServiceInterface --output ./mocks
type AdServiceInterface interface {
	Create(ad *models.Ad) (int, error)
	GetOne(id int) (*models.Ad, error)
	GetAll(priceSort string, dateSort string, page int, userId int) ([]*models.Ad, error)
}

type Service struct {
	Ad AdServiceInterface
}

func NewService(
	repository *repository.Repository,
) *Service {
	return &Service{
		Ad: &adService{repository: repository},
	}
}
