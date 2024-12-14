package repository

import (
	"context"
	"net/http"

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
func (w Wallet) Create(ctx context.Context, db *sqlx.DB, user domain.User) error {
	wallet := domain.Wallet{
		UserID:  user.ID,
		Balance: 0,
	}

	if _, err := db.ExecContext(ctx, createWalletQuery, wallet.UserID, wallet.Balance); err != nil {
		return err
	}

	return nil
}

const getWalletQuery = `SELECT ID, userID, balance FROM UserWallet WHERE userID = $1`

func (w Wallet) Get(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
	wallet := domain.Wallet{}

	if err := db.GetContext(ctx, &wallet, getWalletQuery, user.ID); err != nil {
		return nil, err
	}
	return &wallet, nil
}

const depositQuery = `UPDATE UserWallet SET balance = balance + $2 WHERE userID = $1`
const insertTransactionQuery = `INSERT INTO UserWalletTransaction (userID, operationType, amount, passiveUserID) VALUES ($1, $2, $3, $4)`

func (w Wallet) Deposit(ctx context.Context, db *sqlx.DB, user domain.User, amount int) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, depositQuery, user.ID, amount)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, insertTransactionQuery, user.ID, domain.OperationTypeDeposit, amount, "")
	if err != nil {
		return err
	}
	return tx.Commit()
}

const withdrawQuery = `UPDATE UserWallet SET balance = balance - $2 WHERE userID = $1 AND balance >= $2`

func (w Wallet) Withdraw(ctx context.Context, db *sqlx.DB, user domain.User, amount int) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, withdrawQuery, user.ID, amount)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, insertTransactionQuery, user.ID, domain.OperationTypeWithdraw, amount, "")
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (w Wallet) Transfer(ctx context.Context, db *sqlx.DB, user domain.User, amount int) error {
	return nil
}

func (w Wallet) GetTransactions(ctx context.Context, db *sqlx.DB, user domain.User) ([]*domain.Transaction, error) {
	return nil, nil
}
