package wallet

import (
	"context"
	"time"

	"github.com/sappy5678/cryptocom/pkg/domain"
)

// New creates new wallet logging service
func New(svc domain.WalletService, logger domain.Logger) *LogService {

	return &LogService{
		WalletService: svc,
		logger:        logger,
	}
}

// LogService represents wallet logging service
type LogService struct {
	domain.WalletService
	logger domain.Logger
}

const name = "wallet"

// Create logging
func (ls *LogService) Create(c context.Context, req domain.User) (wallet *domain.Wallet, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create wallet request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.WalletService.Create(c, req)
}

func (ls *LogService) Get(c context.Context, req domain.User) (wallet *domain.Wallet, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Get wallet request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.WalletService.Get(c, req)
}

func (ls *LogService) Withdraw(c context.Context, req domain.User, transactionID domain.TransactionID, amount int) (wallet *domain.Wallet, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Withdraw wallet request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.WalletService.Withdraw(c, req, transactionID, amount)
}

func (ls *LogService) Deposit(c context.Context, req domain.User, transactionID domain.TransactionID, amount int) (wallet *domain.Wallet, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Deposit wallet request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.WalletService.Deposit(c, req, transactionID, amount)
}

func (ls *LogService) GetTransactions(c context.Context, req domain.User, createdAt time.Time, lastReturnedID int, limit int) (transactions []*domain.Transaction, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Get transactions request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.WalletService.GetTransactions(c, req, createdAt, lastReturnedID, limit)
}

func (ls *LogService) Transfer(c context.Context, req domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (wallet *domain.Wallet, err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Transfer wallet request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.WalletService.Transfer(c, req, transactionID, amount, passiveUser)
}

func (ls *LogService) CreateTransactionID(c context.Context) domain.TransactionID {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create transaction ID request", nil,
			map[string]interface{}{
				"took": time.Since(begin),
			},
		)
	}(time.Now())

	return ls.WalletService.CreateTransactionID(c)
}
