package repository

import (
	merchstore "github.com/Sm3underscore23/merchStore"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(username, password_hash string) error
	GetUser(username, password_hash string) (int, error)
}

type UserInfo interface {
	GetUserInfo(user merchstore.User) (interface{}, error)
}

type Buy interface {
}

type SendCoins interface {
}

type Repository struct {
	Authorization
	UserInfo
	Buy
	SendCoins
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
