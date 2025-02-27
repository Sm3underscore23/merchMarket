package service

import (
	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/pkg/repository"
)

type SendCoinsService struct {
	userProviderRepo repository.UserProvider
	trxRepo          repository.Transaction
}

func NewSendCoinsService(
	userProviderRepo repository.UserProvider,
	trxRepo repository.Transaction,
) *SendCoinsService {
	return &SendCoinsService{
		userProviderRepo: userProviderRepo,
		trxRepo:          trxRepo,
	}
}

func (s *SendCoinsService) SendCoins(toUserUsername string, fromUserId, amount int) error {
	toUserId, err := s.userProviderRepo.GetUserId(toUserUsername)
	if err != nil {
		return err
	}

	if fromUserId == toUserId {
		return customerrors.ErrSendCoinsToYousel
	}

	err = s.userProviderRepo.ChangeUserBalance(fromUserId, amount*-1)
	if err != nil {
		return err
	}

	err = s.userProviderRepo.ChangeUserBalance(toUserId, amount)
	if err != nil {
		return err
	}

	err = s.trxRepo.AddP2PTransaction(fromUserId, toUserId, amount)
	if err != nil {
		return err
	}

	return nil
}
