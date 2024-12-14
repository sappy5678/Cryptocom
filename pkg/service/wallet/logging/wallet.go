package user

import (
	"context"
	"time"

	"github.com/sappy5678/cryptocom"

	"github.com/sappy5678/cryptocom/pkg/domain"
)

// New creates new user logging service
func New(svc domain.WalletService, logger cryptocom.Logger) *LogService {
	return &LogService{
		WalletService: svc,
		logger:        logger,
	}
}

// LogService represents user logging service
type LogService struct {
	domain.WalletService
	logger cryptocom.Logger
}

const name = "wallet"

// Create logging
func (ls *LogService) Create(c context.Context, req domain.User) (err error) {
	defer func(begin time.Time) {
		ls.logger.Log(
			c,
			name, "Create user request", err,
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
			name, "Get user request", err,
			map[string]interface{}{
				"req":  req,
				"took": time.Since(begin),
			},
		)
	}(time.Now())
	return ls.WalletService.Get(c, req)
}
