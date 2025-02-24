package service

import "github.com/Sm3underscore23/merchStore/pkg/repository"

type BuyService struct {
	productRepo   repository.Product
	userProvider  repository.UserProvider
	inventoryRepo repository.InventoryManager
}

func NewProductService(
	productRepo repository.Product,
	userProvider repository.UserProvider,
	inventoryRepo repository.InventoryManager,
) *BuyService {
	return &BuyService{
		productRepo:   productRepo,
		userProvider:  userProvider,
		inventoryRepo: inventoryRepo,
	}
}

func (s *BuyService) Buy(userId int, productType string) error {
	productId, productPrice, err := s.productRepo.GetProductIdAndPrice(productType)
	if err != nil {
		// customerrors.ErrProductNotFound
		// customerrors.ErrDatabase
		return err
	}

	err = s.userProvider.ChangeUserBalance(userId, productPrice*-1)
	if err != nil {
		// ErrTxStart
		// ErrGetBalance
		// ErrChangeBalance
		// ErrUpdateBalance
		// ErrTxStop
		return err
	}

	err = s.inventoryRepo.AddItemToInventory(userId, productId)
	if err != nil {
		// customerrors.ErrDatabase
		return err
	}

	return nil
}
