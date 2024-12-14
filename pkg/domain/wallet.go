package domain

import (
	"context"
	"time"
)

type Wallet struct {
	ID        int       `json:"ID"`
	UserID    string    `json:"userID"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type OperationType int

const (
	OperationTypeDummy       OperationType = 0
	OperationTypeDeposit     OperationType = 1
	OperationTypeWithdraw    OperationType = 2
	OperationTypeTransferIn  OperationType = 3
	OperationTypeTransferOut OperationType = 4
)

type Transaction struct {
	ID            int           `json:"ID"`
	UserID        string        `json:"userID"`
	Amount        int           `json:"amount"`
	OperationType OperationType `json:"operationType"`
	PassiveUserID string        `json:"passiveUserID"`
	CreatedAt     time.Time     `json:"createdAt"`
}

type WalletService interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, user User) (*Wallet, error)
	GetTransactions(ctx context.Context, user User) ([]*Transaction, error)
	Transfer(ctx context.Context, user User, amount int) (*Wallet, error)
	Withdraw(ctx context.Context, user User, amount int) (*Wallet, error)
	Deposit(ctx context.Context, user User, amount int) (*Wallet, error)
}
