package repository

import (
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

func (r *AuthPostgres) CheckPassword_hash(id int, password_hash string) error {
	var storedHash string

	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE id=$1", usersTable)
	err := r.db.QueryRow(query, id).Scan(&storedHash)

	if err != nil {
		return err
	}

	if storedHash != password_hash {
		return customerrors.ErrWrongPassword // Неверный пароль
	}

	return nil // Успешная аутентификация
}
