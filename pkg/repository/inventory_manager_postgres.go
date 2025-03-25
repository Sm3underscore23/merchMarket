package repository

import (
	"fmt"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type InventoryManagerPostgres struct {
	db *sqlx.DB
}

func NewInventoryManagerPostgres(db *sqlx.DB) *InventoryManagerPostgres {
	return &InventoryManagerPostgres{db: db}
}

func (r *InventoryManagerPostgres) AddItemToInventory(userID, itemID int) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, item_id, quantity) VALUES ($1, $2, 1) ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = inventory.quantity + 1;",
		inventoryTable)
	_, err := r.db.Query(query, userID, itemID)

	if err != nil {
		return customerrors.ErrDatabase
	}
	return nil
}
