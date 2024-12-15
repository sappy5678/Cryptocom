package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sappy5678/cryptocom/pkg/domain"
)

type MockWalletRepository struct {
	CreateFunc          func(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error)
	GetFunc             func(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error)
	WithdrawFunc        func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error)
	DepositFunc         func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error)
	GetTransactionsFunc func(ctx context.Context, db *sqlx.DB, user domain.User, createdBefore time.Time, IDBefore int, limit int) ([]*domain.Transaction, error)
	TransferFunc        func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error)
}

func (m *MockWalletRepository) Create(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {

	return m.CreateFunc(ctx, db, user)
}

func (m *MockWalletRepository) Get(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {

	return m.GetFunc(ctx, db, user)
}

func (m *MockWalletRepository) Withdraw(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {

	return m.WithdrawFunc(ctx, db, time, user, transactionID, amount)
}

func (m *MockWalletRepository) Deposit(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {

	return m.DepositFunc(ctx, db, time, user, transactionID, amount)
}

func (m *MockWalletRepository) GetTransactions(ctx context.Context, db *sqlx.DB, user domain.User, createdBefore time.Time, IDBefore int, limit int) ([]*domain.Transaction, error) {

	return m.GetTransactionsFunc(ctx, db, user, createdBefore, IDBefore, limit)
}

func (m *MockWalletRepository) Transfer(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {

	return m.TransferFunc(ctx, db, time, user, transactionID, amount, passiveUser)
}
