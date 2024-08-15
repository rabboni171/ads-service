package service

import (
	"context"
	"time"
)


type AuthServiceInterface interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int32,
		tokenTTL time.Duration,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}
