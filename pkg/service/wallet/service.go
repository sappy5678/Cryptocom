package wallet

import (
	"github.com/jmoiron/sqlx"

	"github.com/sappy5678/cryptocom/pkg/domain"
	"github.com/sappy5678/cryptocom/pkg/service/wallet/repository"
)

// Service defines in domain

// New creates new wallet application service
// TODO: using WalletService
func New(db *sqlx.DB, walletRepo repository.WalletRepository) domain.WalletService {
	return &Wallet{db: db, walletRepo: walletRepo}
}

// Initialize initalizes Wallet application service with defaults
func Initialize(db *sqlx.DB) domain.WalletService {
	return New(db, repository.Wallet{})
}

// Wallet represents wallet application service
type Wallet struct {
	db         *sqlx.DB
	walletRepo repository.WalletRepository
}
