package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) GetUser(username, password_hash string) (int, error) {

	var id int
	var storedHash string

	query := fmt.Sprintf("SELECT id, password_hash FROM %s WHERE username=$1", usersTable)
	err := r.db.QueryRow(query, username).Scan(&id, &storedHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, customerrors.ErrUserNotFound // Пользователь не найден
		}
		return 0, err
	}

	if storedHash != password_hash {
		return 0, customerrors.ErrWrongPasswod
	}

	return id, nil // Успешная аутентификация
}

func (r *AuthPostgres) CreateUser(username, password_hash string) error {
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) values($1, $2) RETURNING id", usersTable)
	row := r.db.QueryRow(query, username, password_hash)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}
