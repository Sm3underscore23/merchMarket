package repository

import (
	"fmt"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) AddP2PTransaction(fromUserID, toUserID, amount int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return customerrors.ErrTxStart
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`INSERT INTO %s (from_user, to_user, amount) VALUES ($1, $2, $3)`, transactionsTable)
	_, err = tx.Exec(query, fromUserID, toUserID, amount)
	if err != nil {
		return fmt.Errorf("add transaction error")
	}

	err = tx.Commit()
	if err != nil {
		return customerrors.ErrTxStop
	}

	return nil
}

