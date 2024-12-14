// Package user contains user application services
package wallet

import (
	"context"

	"github.com/sappy5678/cryptocom/pkg/domain"
)

// Create creates a new user account
func (w Wallet) Create(ctx context.Context, user domain.User) error {
	if err := w.walletRepo.Create(ctx, w.db, user); err != nil {
		return err
	}
	return nil
}

func (w Wallet) Get(ctx context.Context, user domain.User) (*domain.Wallet, error) {
	wallet, err := w.walletRepo.Get(ctx, w.db, user)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w Wallet) GetTransactions(ctx context.Context, user domain.User) ([]*domain.Transaction, error) {
	return nil, nil
}

func (w Wallet) Transfer(ctx context.Context, user domain.User, amount int) (*domain.Wallet, error) {
	return nil, nil
}

func (w Wallet) Withdraw(ctx context.Context, user domain.User, amount int) (*domain.Wallet, error) {
	return nil, nil
}

func (w Wallet) Deposit(ctx context.Context, user domain.User, amount int) (*domain.Wallet, error) {
	return nil, nil
}
