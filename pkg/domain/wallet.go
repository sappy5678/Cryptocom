package domain

import (
	"context"
	"time"
)

type Wallet struct {
	ID      int    `json:"-"`
	UserID  string `json:"userID"`
	Balance int    `json:"balance"`
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
	TransactionID TransactionID `json:"transactionID"`
	UserID        string        `json:"userID"`
	Amount        int           `json:"amount"`
	OperationType OperationType `json:"operationType"`
	PassiveUserID string        `json:"passiveUserID"`
	CreatedAt     time.Time     `json:"createdAt"`
}

type TransactionID string

func (t TransactionID) ID() string {

	return string(t)
}
func (t TransactionID) PassiveID() string {

	return string(t) + "-passive"
}

type WalletService interface {
	Create(ctx context.Context, user User) (*Wallet, error)
	Get(ctx context.Context, user User) (*Wallet, error)
	CreateTransactionID(ctx context.Context) TransactionID
	GetTransactions(ctx context.Context, user User, createdAt time.Time, lastReturnedID int, limit int) ([]*Transaction, error)
	Transfer(ctx context.Context, user User, transactionID TransactionID, amount int, passiveUser User) (*Wallet, error)
	Withdraw(ctx context.Context, user User, transactionID TransactionID, amount int) (*Wallet, error)
	Deposit(ctx context.Context, user User, transactionID TransactionID, amount int) (*Wallet, error)
}
