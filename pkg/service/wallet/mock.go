package wallet

import (
	"context"
	"time"

	"github.com/sappy5678/cryptocom/pkg/domain"
)

type MockWalletService struct {
	CreateFunc              func(ctx context.Context, user domain.User) (*domain.Wallet, error)
	GetFunc                 func(ctx context.Context, user domain.User) (*domain.Wallet, error)
	WithdrawFunc            func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error)
	DepositFunc             func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error)
	GetTransactionsFunc     func(ctx context.Context, user domain.User, createdAt time.Time, lastReturnedID int, limit int) ([]*domain.Transaction, error)
	TransferFunc            func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error)
	CreateTransactionIDFunc func(ctx context.Context) domain.TransactionID
}

func (m *MockWalletService) Create(ctx context.Context, user domain.User) (*domain.Wallet, error) {
	return m.CreateFunc(ctx, user)
}

func (m *MockWalletService) Get(ctx context.Context, user domain.User) (*domain.Wallet, error) {
	return m.GetFunc(ctx, user)
}

func (m *MockWalletService) Withdraw(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
	return m.WithdrawFunc(ctx, user, transactionID, amount)
}

func (m *MockWalletService) Deposit(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
	return m.DepositFunc(ctx, user, transactionID, amount)
}

func (m *MockWalletService) GetTransactions(ctx context.Context, user domain.User, createdAt time.Time, lastReturnedID int, limit int) ([]*domain.Transaction, error) {
	return m.GetTransactionsFunc(ctx, user, createdAt, lastReturnedID, limit)
}

func (m *MockWalletService) Transfer(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {
	return m.TransferFunc(ctx, user, transactionID, amount, passiveUser)
}

func (m *MockWalletService) CreateTransactionID(ctx context.Context) domain.TransactionID {
	return m.CreateTransactionIDFunc(ctx)
}
