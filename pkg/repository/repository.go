package repository

import (
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CheckPassword_hash(id int, password_hash string) error
}

type UserProvider interface {
	GetUserId(username string) (int, error)
	ChangeUserBalance(id, coins int) error
	CreateUser(username, password_hash string) (int, error)
	GetUserInfo(id int) (models.UserInfoResponse, error)
}

type InventoryManager interface {
	AddItemToInventory(userID, itemID int) error
	// GetUserInventory(userID int) ([]InventoryItem, error)
}

type Product interface {
	GetProductIdAndPrice(productType string) (int, int, error)
}

type Transaction interface {
	AddP2PTransaction(fromUserID, toUserID, amount int) error
}

type Repository struct {
	Authorization
	UserProvider
	InventoryManager
	Product
	Transaction
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:    NewAuthPostgres(db),
		UserProvider:     NewUserProviderPostgres(db),
		InventoryManager: NewInventoryManagerPostgres(db),
		Product:          NewProductPostgres(db),
		Transaction:      NewTransactionPostgres(db),
	}
}
