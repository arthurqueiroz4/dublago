package user

import (
	"context"
	"errors"
	"log/slog"
	"tradutor-dos-crias/database"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserService struct{}

var (
	createUserSql = `INSERT INTO "user" (email, name, sso_id) VALUES ($1, $2, $3)`
	getByEmailSql = `SELECT id, email, name, sso_id FROM "user" WHERE email = $1`
)

func (us *UserService) Create(user *User) error {
	_, err := database.Pool.Exec(context.Background(), createUserSql,
		user.Email,
		user.Name,
		user.SsoId)
	if err != nil {
		return handleError(user, err)
	}
	return nil
}

func (us *UserService) GetByEmail(email string) (*User, error) {
	var user User
	err := database.Pool.QueryRow(context.Background(), getByEmailSql, email).Scan(&user.ID, &user.Email, &user.Name, &user.SsoId)
	if err != nil {
		return nil, errors.New("Error on QueryRow [GetByEmail]")
	}
	return &user, nil
}

func handleError(user *User, err error) error {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			slog.Error("User already exists", "email", user.Email)
			return ErrUserAlreadyExists
		}
	}

	slog.Error("Error on creating user", "error", err)
	return err
}
