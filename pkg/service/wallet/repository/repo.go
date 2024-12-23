package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sappy5678/cryptocom/pkg/domain"
)

// WalletRepository represents wallet repository interface
type WalletRepository interface {
	Create(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error)
	Get(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error)
	Withdraw(ctx context.Context, db *sqlx.DB, now time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error)
	Deposit(ctx context.Context, db *sqlx.DB, now time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error)
	GetTransactions(ctx context.Context, db *sqlx.DB, user domain.User, createdBefore time.Time, IDBefore int, limit int) ([]*domain.Transaction, error)
	Transfer(ctx context.Context, db *sqlx.DB, now time.Time, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error)
}
