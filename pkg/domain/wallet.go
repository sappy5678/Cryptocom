package domain

import (
	"time"

	"github.com/labstack/echo"
)

type Wallet struct {
	ID        int       `json:"ID"`
	UserID    string    `json:"userID"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WalletService interface {
	CreateWallet(ctx echo.Context, userID string) (*Wallet, error)
	GetWallet(ctx echo.Context, userID string) (*Wallet, error)
	GetBalance(ctx echo.Context, userID string) (int, error)
	GetTransactions(ctx echo.Context, userID string) ([]Transaction, error)
	Transfer(ctx echo.Context, userID string, amount int) (*Wallet, error)
	Withdraw(ctx echo.Context, userID string, amount int) (*Wallet, error)
	Deposit(ctx echo.Context, userID string, amount int) (*Wallet, error)
}
