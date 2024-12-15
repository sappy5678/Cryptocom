// Package user contains user application services
package wallet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sappy5678/cryptocom/pkg/domain"
)

// Create creates a new user account
func (w *Wallet) Create(ctx context.Context, user domain.User) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Create(ctx, w.db, user)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *Wallet) Get(ctx context.Context, user domain.User) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Get(ctx, w.db, user)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *Wallet) CreateTransactionID(ctx context.Context) domain.TransactionID {
	uuid := uuid.New()
	return domain.TransactionID(uuid.String())
}

func (w *Wallet) Withdraw(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Withdraw(ctx, w.db, time.Now(), user, transactionID, amount)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *Wallet) Deposit(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Deposit(ctx, w.db, time.Now(), user, transactionID, amount)

	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *Wallet) GetTransactions(ctx context.Context, user domain.User, createdAt time.Time, lastReturnedID int, limit int) ([]*domain.Transaction, error) {
	transactions, err := w.walletRepo.GetTransactions(ctx, w.db, user, createdAt, lastReturnedID, limit)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (w *Wallet) Transfer(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Transfer(ctx, w.db, time.Now(), user, transactionID, amount, passiveUser)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
