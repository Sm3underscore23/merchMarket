package service

import (
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/Sm3underscore23/merchStore/pkg/repository"
)

type GetUserInfoServece struct {
	repoUserProvider repository.UserProvider
}

func NewGetUserInfoServece(repoUserProvider repository.UserProvider) *GetUserInfoServece {
	return &GetUserInfoServece{
		repoUserProvider: repoUserProvider,
	}
}

func (s *GetUserInfoServece) GetUserInfo(id int) (models.UserInfoResponse, error) {
	var userInfo models.UserInfoResponse

	userInfo, err := s.repoUserProvider.GetUserInfo(id)

	if err != nil {
		return userInfo, err
	}

	return userInfo, nil

}
