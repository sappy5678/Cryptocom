package repository

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sappy5678/cryptocom/pkg/domain"

	"github.com/labstack/echo"
)

// postgresql error code define
// http://www.postgresql.org/docs/9.3/static/errcodes-appendix.html

func TimeToDBTime(t time.Time) time.Time {
	return t.UTC().Round(time.Microsecond)
}

// User represents the client for user table
type Wallet struct{}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

const existsQuery = `SELECT EXISTS(SELECT 1 FROM UserWallet WHERE userID = $1)`

func (w Wallet) Exists(ctx context.Context, db *sqlx.DB, user domain.User) (bool, error) {
	var exists bool
	if err := db.GetContext(ctx, &exists, existsQuery, user.ID); err != nil {
		return false, err
	}
	return exists, nil
}

const existsTransactionIDQuery = `SELECT EXISTS(SELECT 1 FROM UserWalletTransaction WHERE transactionID = $1)`

func (w Wallet) ExistsTransactionID(ctx context.Context, db *sqlx.DB, transactionID domain.TransactionID) (bool, error) {
	var exists bool
	if err := db.GetContext(ctx, &exists, existsTransactionIDQuery, transactionID.ID()); err != nil {
		return false, err
	}
	return exists, nil
}

const createWalletQuery = `INSERT INTO UserWallet (userID, balance) VALUES ($1, $2)`

// Create creates a new user on database
func (w Wallet) Create(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
	wallet := domain.Wallet{
		UserID:  user.ID,
		Balance: 0,
	}

	if exists, err := w.Exists(ctx, db, user); err != nil {
		return nil, err
	} else if exists {
		wallet, err := w.Get(ctx, db, user)
		return wallet, err
	}

	if _, err := db.ExecContext(ctx, createWalletQuery, wallet.UserID, wallet.Balance); err != nil {
		return nil, err
	}

	return &wallet, nil
}

const getWalletQuery = `SELECT ID, userID, balance FROM UserWallet WHERE userID = $1`

func (w Wallet) Get(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
	wallet := domain.Wallet{}

	if exists, err := w.Exists(ctx, db, user); err != nil {
		return nil, err
	} else if !exists {
		return nil, domain.ErrWalletNotFound
	}

	if err := db.GetContext(ctx, &wallet, getWalletQuery, user.ID); err != nil {
		return nil, err
	}
	return &wallet, nil
}

const depositQuery = `UPDATE UserWallet SET balance = balance + $2 WHERE userID = $1 RETURNING userID, balance`
const insertTransactionQuery = `INSERT INTO UserWalletTransaction (userID, transactionID, operationType, amount, passiveUserID, createdAt) VALUES ($1, $2, $3, $4, $5, $6)`

func (w Wallet) Deposit(ctx context.Context, db *sqlx.DB, now time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	// idempotent
	if exists, err := w.ExistsTransactionID(ctx, db, transactionID); err != nil {
		return nil, err
	} else if exists {
		return w.Get(ctx, db, user)
	}

	now = TimeToDBTime(now)

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
		return nil, domain.ErrWalletNotFound
	}
	if err := rows.StructScan(&wallet); err != nil {
		return nil, err
	}
	rows.Close()

	// insert transaction
	_, err = tx.ExecContext(ctx, insertTransactionQuery, user.ID, transactionID.ID(), domain.OperationTypeDeposit, amount, "", now)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &wallet, nil
}

const withdrawQuery = `UPDATE UserWallet SET balance = balance - $2 WHERE userID = $1 AND balance >= $2 RETURNING userID, balance`

func (w Wallet) Withdraw(ctx context.Context, db *sqlx.DB, now time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
	// check condition
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}
	if exists, err := w.Exists(ctx, db, user); err != nil {
		return nil, err
	} else if !exists {
		return nil, domain.ErrWalletNotFound
	}

	// idempotent
	if exists, err := w.ExistsTransactionID(ctx, db, transactionID); err != nil {
		return nil, err
	} else if exists {
		return w.Get(ctx, db, user)
	}

	now = TimeToDBTime(now)

	// start transaction
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
		return nil, domain.ErrNotEnoughBalance
	}
	if err := rows.StructScan(&wallet); err != nil {
		return nil, err
	}
	rows.Close()

	// insert transaction
	_, err = tx.ExecContext(ctx, insertTransactionQuery, user.ID, transactionID.ID(), domain.OperationTypeWithdraw, amount, "", now)
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

func (w Wallet) Transfer(ctx context.Context, db *sqlx.DB, now time.Time, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {
	// check condition
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}
	if user.ID == passiveUser.ID {
		return nil, domain.ErrTransferToSelf
	}
	if exists, err := w.Exists(ctx, db, user); err != nil {
		return nil, err
	} else if !exists {
		return nil, domain.ErrWalletNotFound
	}
	if exists, err := w.Exists(ctx, db, passiveUser); err != nil {
		return nil, err
	} else if !exists {
		return nil, domain.ErrWalletNotFound
	}

	// idempotent
	if exists, err := w.ExistsTransactionID(ctx, db, transactionID); err != nil {
		return nil, err
	} else if exists {
		return w.Get(ctx, db, user)
	}

	now = TimeToDBTime(now)

	// start transaction
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
		return nil, domain.ErrNotEnoughBalance
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
		return nil, domain.ErrWalletNotFound
	}

	// insert transaction
	_, err = tx.ExecContext(ctx, insertTransactionQuery, user.ID, transactionID.ID(),
		domain.OperationTypeTransferOut, amount, passiveUser.ID, now)
	if err != nil {
		return nil, err
	}

	// insert passive transaction
	_, err = tx.ExecContext(ctx, insertTransactionQuery, passiveUser.ID, transactionID.PassiveID(),
		domain.OperationTypeTransferIn, amount, user.ID, now)
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

	// default values
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if lastReturnedID == 0 {
		lastReturnedID = math.MaxInt64
	}
	if limit <= 0 {
		limit = 100
	}
	if exists, err := w.Exists(ctx, db, user); err != nil {
		return nil, err
	} else if !exists {
		return nil, domain.ErrWalletNotFound
	}
	transactions := []*domain.Transaction{}

	if err := db.SelectContext(ctx, &transactions, getTransactionsQuery, user.ID, createdAt,
		lastReturnedID, limit); err != nil {
		return nil, err
	}

	// remove timezone information
	for _, transaction := range transactions {
		transaction.CreatedAt = TimeToDBTime(transaction.CreatedAt)
	}
	return transactions, nil
}
