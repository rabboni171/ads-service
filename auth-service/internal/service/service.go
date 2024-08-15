package service

import "auth-service/internal/storage"

type Service struct {
	Auth AuthServiceInterface
}

func NewService(
	authStorage storage.AuthStorageInterface,
) *Service {
	return &Service{
		Auth: &AuthService{
			authStorage: authStorage,
		},
	}
}
