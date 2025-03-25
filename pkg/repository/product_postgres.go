package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type ProductPostgres struct {
	db *sqlx.DB
}

func NewProductPostgres(db *sqlx.DB) *ProductPostgres {
	return &ProductPostgres{db: db}
}

func (r *ProductPostgres) GetProductIdAndPrice(productType string) (int, int, error) {
	var id, price int
	query := fmt.Sprintf("SELECT id, price FROM %s WHERE product_type=$1", goodsTable)
	err := r.db.QueryRow(query, productType).Scan(&id, &price)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, customerrors.ErrProductNotFound
		}
		return 0, 0, customerrors.ErrDatabase
	}
	return id, price, nil
	// Получить id и ценну продукта, если он есть (productType)(id, price)
	// Поменять баланс пользователя ChangeBalane
	// Записать транзакцию в покупки 
}
