package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sappy5678/cryptocom/pkg/domain"
)

// WalletRepository represents wallet repository interface
type WalletRepository interface {
	Create(ctx context.Context, db *sqlx.DB, user domain.User) error
	Get(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error)
}
