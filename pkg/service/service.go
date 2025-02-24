package service

import (
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/Sm3underscore23/merchStore/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go


type Authorization interface {
	AuthUser(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Info interface {
	GetUserInfo(id int) (models.UserInfoResponse, error)
}

type SendCoins interface {
	SendCoins(toUserUsername string, fromUserId, amount int) error
}

type Buy interface {
	Buy(userId int, productType string) error
}

type Service struct {
	Authorization
	Info
	Buy
	SendCoins
}

func NewService(repos *repository.Repository, authConfig models.AuthConfig) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, repos.UserProvider, authConfig),
		Info:          NewGetUserInfoServece(repos.UserProvider),
		Buy:           NewProductService(repos.Product, repos.UserProvider, repos.InventoryManager),
		SendCoins:     NewSendCoinsService(repos.UserProvider, repos.Transaction),
	}
}
