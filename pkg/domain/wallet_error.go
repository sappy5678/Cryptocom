package domain

import "errors"

var (
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrWalletNotFound   = errors.New("wallet not found")
	ErrNotEnoughBalance = errors.New("not enough balance")
	ErrTransferToSelf   = errors.New("transfer to self")
	ErrUserIDRequired   = errors.New("userID is required")
)
