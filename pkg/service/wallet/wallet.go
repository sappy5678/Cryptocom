// Package user contains user application services
package wallet

import (
	"context"

	"github.com/sappy5678/cryptocom/pkg/domain"
)

// Create creates a new user account
func (w Wallet) Create(ctx context.Context, user domain.User) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Create(ctx, w.db, user)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w Wallet) Get(ctx context.Context, user domain.User) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Get(ctx, w.db, user)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w Wallet) Withdraw(ctx context.Context, user domain.User, amount int) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Withdraw(ctx, w.db, user, amount)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w Wallet) Deposit(ctx context.Context, user domain.User, amount int) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Deposit(ctx, w.db, user, amount)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w Wallet) GetTransactions(ctx context.Context, user domain.User) ([]*domain.Transaction, error) {
	transactions, err := w.walletRepo.GetTransactions(ctx, w.db, user)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (w Wallet) Transfer(ctx context.Context, user domain.User, amount int, passiveUser domain.User) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Transfer(ctx, w.db, user, amount, passiveUser)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
