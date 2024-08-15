package service

import (
	"auth-service/internal/lib/jwt"
	"auth-service/internal/lib/types"
	"auth-service/internal/storage"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authStorage storage.AuthStorageInterface
}

func (s *AuthService) RegisterNewUser(
	ctx context.Context, email string, pass string) (int64, error) {

	// Генерируем хэш и соль для пароля.
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to generate password hash", "err", err.Error())
		return 0, err
	}

	// Сохраняем пользователя в БД
	id, err := s.authStorage.SaveUser(ctx, email, passHash)
	if err != nil {
		slog.Error("failed to save user", "err", err)
		
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, types.ErrUserExists
			}
		}
		return 0, err
	}

	return id, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	appID int32,
	tokenTTL time.Duration,
) (string, error) {
	// Достаём пользователя из БД
	user, err := s.authStorage.GetUser(ctx, email)
	if err != nil {

		if err == pgx.ErrNoRows {
			slog.Warn("user not found", "err", err.Error())
			return "", types.ErrUserNotFound
		}

		slog.Error("failed to get user", "err", err.Error())
		return "", err
	}

	// Проверяем пароль пользователя на соответствие
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		slog.Info("invalid credentials", "err", err.Error())
		return "", types.ErrInvalidCredentials
	}

	// Получаем информацию о приложении
	app, err := s.authStorage.GetApp(ctx, appID)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Warn("app not found", "err", err.Error())
			return "", types.ErrAppNotFound
		}

		return "", err
	}

	slog.Info("user logged in successfully")

	// Создаём токен авторизации
	token, err := jwt.NewToken(user, app, tokenTTL)
	if err != nil {
		slog.Error("failed to generate token", "err", err.Error())
		return "", err
	}

	return token, nil
}
