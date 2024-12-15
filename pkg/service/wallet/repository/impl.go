package repository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sappy5678/cryptocom/pkg/domain"

	"github.com/labstack/echo"
)

// User represents the client for user table
type Wallet struct{}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

const createWalletQuery = `INSERT INTO UserWallet (userID, balance) VALUES ($1, $2)`

// Create creates a new user on database
func (w Wallet) Create(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
	wallet := domain.Wallet{
		UserID:  user.ID,
		Balance: 0,
	}

	if _, err := db.ExecContext(ctx, createWalletQuery, wallet.UserID, wallet.Balance); err != nil {
		return nil, err
	}

	return &wallet, nil
}

const getWalletQuery = `SELECT ID, userID, balance FROM UserWallet WHERE userID = $1`

func (w Wallet) Get(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
	wallet := domain.Wallet{}

	if err := db.GetContext(ctx, &wallet, getWalletQuery, user.ID); err != nil {
		return nil, err
	}
	return &wallet, nil
}

const depositQuery = `UPDATE UserWallet SET balance = balance + $2 WHERE userID = $1 RETURNING userID, balance`
const insertTransactionQuery = `INSERT INTO UserWalletTransaction (userID, transactionID, operationType, amount, passiveUserID) VALUES ($1, $2, $3, $4, $5)`

func (w Wallet) Deposit(ctx context.Context, db *sqlx.DB, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// update wallet balance and get the new balance
	rows, err := tx.QueryxContext(ctx, depositQuery, user.ID, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallet domain.Wallet
	if !rows.Next() {
		return nil, domain.ErrNotFound
	}
	if err := rows.StructScan(&wallet); err != nil {
		return nil, err
	}
	rows.Close()

	// insert transaction
	_, err = tx.ExecContext(ctx, insertTransactionQuery, user.ID, transactionID.ID(), domain.OperationTypeDeposit, amount, "")
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &wallet, nil
}

const withdrawQuery = `UPDATE UserWallet SET balance = balance - $2 WHERE userID = $1 AND balance >= $2 RETURNING userID, balance`

func (w Wallet) Withdraw(ctx context.Context, db *sqlx.DB, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// update wallet balance and get the new balance
	rows, err := tx.QueryxContext(ctx, withdrawQuery, user.ID, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallet domain.Wallet
	if !rows.Next() {
		return nil, domain.ErrInsufficientBalance
	}
	if err := rows.StructScan(&wallet); err != nil {
		return nil, err
	}
	rows.Close()

	// insert transaction
	_, err = tx.ExecContext(ctx, insertTransactionQuery, user.ID, transactionID.ID(), domain.OperationTypeWithdraw, amount, "")
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &wallet, nil
}

const transferQuery = `UPDATE UserWallet SET balance = balance - $2 WHERE userID = $1 AND balance >= $2 RETURNING userID, balance`
const passiveTransferQuery = `UPDATE UserWallet SET balance = balance + $2 WHERE userID = $1`
const insertTransferTransactionQuery = `INSERT INTO UserWalletTransaction (userID, transactionID, operationType, amount, passiveUserID) VALUES ($1, $2, $3, $4, $5)`

func (w Wallet) Transfer(ctx context.Context, db *sqlx.DB, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// update wallet balance and get the new balance
	rows, err := tx.QueryxContext(ctx, withdrawQuery, user.ID, amount)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, domain.ErrInsufficientBalance
	}
	var wallet domain.Wallet
	if err := rows.StructScan(&wallet); err != nil {
		return nil, err
	}
	rows.Close()

	// update passive wallet balance
	affected, err := tx.ExecContext(ctx, passiveTransferQuery, passiveUser.ID, amount)
	if err != nil {
		return nil, err
	}
	if affected, err := affected.RowsAffected(); err != nil || affected == 0 {
		if err != nil {
			return nil, err
		}
		return nil, domain.ErrNotFound
	}

	// insert transaction
	_, err = tx.ExecContext(ctx, insertTransferTransactionQuery, user.ID, transactionID.ID(), domain.OperationTypeTransferOut, amount, passiveUser.ID)
	if err != nil {
		return nil, err
	}

	// insert passive transaction
	_, err = tx.ExecContext(ctx, insertTransferTransactionQuery, passiveUser.ID, transactionID.PassiveID(), domain.OperationTypeTransferIn, amount, user.ID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &wallet, nil
}

const getTransactionsQuery = `SELECT ID, userID, transactionID, operationType, amount, passiveUserID, createdAt FROM UserWalletTransaction WHERE userID=$1 AND createdAt <= $2 AND ID < $3 ORDER BY createdAt DESC LIMIT $4`

func (w Wallet) GetTransactions(ctx context.Context, db *sqlx.DB, user domain.User, createdAt time.Time, lastReturnedID int, limit int) ([]*domain.Transaction, error) {
	transactions := []*domain.Transaction{}

	if err := db.SelectContext(ctx, &transactions, getTransactionsQuery, user.ID, createdAt, lastReturnedID, limit); err != nil {
		return nil, err
	}
	fmt.Println(getTransactionsQuery, user.ID, fmt.Sprintf("%+v", transactions))
	return transactions, nil
}
