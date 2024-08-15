package storage

import (
	"auth-service/internal/models"
	"context"
)

type AuthStorageInterface interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
	GetUser(ctx context.Context, email string) (*models.User, error)
	GetApp(ctx context.Context, appID int32) (*models.App, error)
}