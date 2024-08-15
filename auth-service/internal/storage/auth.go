package storage

import (
	"auth-service/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
)

type AuthStorage struct {
	conn *pgx.Conn
}
var _ AuthStorageInterface = (*AuthStorage)(nil)

func (a *AuthStorage) GetApp(ctx context.Context, appID int32) (*models.App, error) {
	app := &models.App{}
	err := a.conn.QueryRow(ctx,
		`SELECT id, name, secret 
		FROM apps 
		WHERE id=$1;`,
		appID,
	).Scan(&app.ID, &app.Name, &app.Secret)

	return app, err
}

func (a *AuthStorage) GetUser(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := a.conn.QueryRow(ctx,
		`SELECT id, email, pass_hash 
		FROM users 
		WHERE email=$1;`,
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash)

	return user, err
}

func (a *AuthStorage) SaveUser(ctx context.Context, email string, passHash []byte,
) (int64, error) {
	var id int64
	err := a.conn.QueryRow(ctx,
		`INSERT INTO users(email, pass_hash) 
		VALUES($1, $2)
		RETURNING id;`,
		email, passHash,
	).Scan(&id)

	return id, err
}
