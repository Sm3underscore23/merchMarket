package service

import (
	"github.com/Sm3underscore23/merchStore/pkg/repository"
)

type Authorization interface {
	GetUser(username, password string) (int, error)
	GenerateToken(id int) (string, error)
}

type UserInfo interface {
}

type Buy interface {
}

type SendCoins interface {
}

type Service struct {
	Authorization
	UserInfo
	Buy
	SendCoins
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
