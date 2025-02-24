package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserProviderPostgres struct {
	db *sqlx.DB
}

func NewUserProviderPostgres(db *sqlx.DB) *UserProviderPostgres {
	return &UserProviderPostgres{db: db}
}

func (r *UserProviderPostgres) CreateUser(username, password_hash string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) VALUES ($1, $2) RETURNING id", usersTable)
	err := r.db.QueryRow(query, username, password_hash).Scan(&id)
	if err != nil {
		return 0, customerrors.ErrDatabase
	}
	return id, nil
}

func (r *UserProviderPostgres) GetUserId(username string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1", usersTable)
	err := r.db.QueryRow(query, username).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, customerrors.ErrUserNotFound // Пользователь не найден
		}
		return 0, customerrors.ErrDatabase
	}

	return id, nil
}

func (r *UserProviderPostgres) ChangeUserBalance(id, coins int) error {
	fmt.Println("ChangeUserBalance Repo SenderId", id)
	tx, err := r.db.Begin()
	if err != nil {
		return customerrors.ErrTxStart
	}
	defer tx.Rollback()

	var currentBalance int

	query := fmt.Sprintf("SELECT balance FROM %s WHERE id=$1 FOR UPDATE", usersTable)
	fmt.Println(id)
	err = tx.QueryRow(query, id).Scan(&currentBalance)
	if err != nil {
		fmt.Println(err)
		return customerrors.ErrGetBalance
	}

	newBalance := currentBalance + coins
	if newBalance < 0 {
		return customerrors.ErrChangeBalance
	}

	query = fmt.Sprintf("UPDATE %s SET balance=$1 WHERE id=$2", usersTable)
	_, err = tx.Exec(query, newBalance, id)
	if err != nil {
		return customerrors.ErrUpdateBalance
	}

	err = tx.Commit()
	if err != nil {
		return customerrors.ErrTxStop
	}

	return nil
}

func (r *UserProviderPostgres) GetUserInfo(id int) (models.UserInfoResponse, error) {
	var userInfo models.UserInfoResponse

	tx, err := r.db.Begin()
	if err != nil {
		return userInfo, customerrors.ErrTxStart
	}
	defer tx.Rollback()

	query := fmt.Sprintf("SELECT balance FROM %s WHERE id=$1", usersTable)
	err = tx.QueryRow(query, id).Scan(&userInfo.Balance)
	if err != nil {
		return userInfo, customerrors.ErrGetBalance
	}

	query = fmt.Sprintf("SELECT g.product_type, i.quantity FROM %s i JOIN goods g ON i.item_id = g.id WHERE i.user_id = $1;", inventoryTable)
	rows, err := tx.Query(query, id)
	if err != nil {
		return userInfo, customerrors.ErrDatabase
	}

	for rows.Next() {
		var inventoryItem models.InventoryItem
		err := rows.Scan(&inventoryItem.ItemType, &inventoryItem.Quantity)
		if err != nil {
			return userInfo, customerrors.ErrParseInventory
		}
		userInfo.Inventory = append(userInfo.Inventory, inventoryItem)
	}

	rows.Close()

	query = fmt.Sprintf("SELECT u.username AS receiver_username, t.amount FROM %s t JOIN users u ON t.to_user = u.id WHERE t.from_user = $1;", transactionsTable)
	rows, err = tx.Query(query, id)
	if err != nil {
		return userInfo, customerrors.ErrDatabase
	}

	for rows.Next() {
		var outgoingTransaction models.OutgoingTransaction
		err := rows.Scan(&outgoingTransaction.ToUser, &outgoingTransaction.Amount)
		if err != nil {
			return userInfo, customerrors.ErrParseTrx
		}
		userInfo.CoinHistory.Sent = append(userInfo.CoinHistory.Sent, outgoingTransaction)
	}

	rows.Close()

	query = fmt.Sprintf("SELECT u.username AS sender_username, t.amount FROM %s t JOIN users u ON t.from_user = u.id WHERE t.to_user = $1 ORDER BY t.created_at;", transactionsTable)
	rows, err = tx.Query(query, id)
	if err != nil {
		return userInfo, customerrors.ErrDatabase
	}

	for rows.Next() {
		var incomingTransaction models.IncomingTransaction
		err := rows.Scan(&incomingTransaction.FromUser, &incomingTransaction.Amount)
		if err != nil {
			return userInfo, customerrors.ErrParseTrx
		}
		userInfo.CoinHistory.Received = append(userInfo.CoinHistory.Received, incomingTransaction)
	}

	rows.Close()

	err = tx.Commit()
	if err != nil {
		return userInfo, customerrors.ErrTxStop
	}

	return userInfo, nil
}
