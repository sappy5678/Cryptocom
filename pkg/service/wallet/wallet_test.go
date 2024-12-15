package wallet_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/sappy5678/cryptocom/pkg/domain"
	"github.com/sappy5678/cryptocom/pkg/service/wallet"
	"github.com/sappy5678/cryptocom/pkg/service/wallet/repository"
)

var mockWalletRepository = &repository.MockWalletRepository{
	CreateFunc: func(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: "1"}, nil
	},
	GetFunc: func(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: "1"}, nil
	},
	WithdrawFunc: func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: "1"}, nil
	},
	DepositFunc: func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: "1"}, nil
	},
	GetTransactionsFunc: func(ctx context.Context, db *sqlx.DB, user domain.User, createdBefore time.Time, IDBefore int, limit int) ([]*domain.Transaction, error) {
		return []*domain.Transaction{{UserID: "1"}}, nil
	},
	TransferFunc: func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: "1"}, nil
	},
}

var mockErrorWalletRepository = &repository.MockWalletRepository{
	CreateFunc: func(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
		return nil, errors.New("error")
	},
	GetFunc: func(ctx context.Context, db *sqlx.DB, user domain.User) (*domain.Wallet, error) {
		return nil, errors.New("error")
	},
	WithdrawFunc: func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
		return nil, errors.New("error")
	},
	DepositFunc: func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
		return nil, errors.New("error")
	},
	GetTransactionsFunc: func(ctx context.Context, db *sqlx.DB, user domain.User, createdBefore time.Time, IDBefore int, limit int) ([]*domain.Transaction, error) {
		return nil, errors.New("error")
	},
	TransferFunc: func(ctx context.Context, db *sqlx.DB, time time.Time, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {
		return nil, errors.New("error")
	},
}

func TestNew(t *testing.T) {
	cases := []struct {
		name    string
		db      *sqlx.DB
		wantErr bool
	}{
		{
			name:    "create wallet service",
			db:      &sqlx.DB{},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.New(tt.db, nil)
			assert.NotNil(t, svc)
		})
	}
}

func TestInitialize(t *testing.T) {
	cases := []struct {
		name    string
		db      *sqlx.DB
		wantErr bool
	}{
		{
			name:    "initialize wallet service",
			db:      &sqlx.DB{},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.Initialize(tt.db)
			assert.NotNil(t, svc)
		})
	}
}

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		db       *sqlx.DB
		mockRepo repository.WalletRepository
		wantErr  bool
	}{
		{
			name:     "create wallet service",
			db:       &sqlx.DB{},
			mockRepo: mockWalletRepository,
			wantErr:  false,
		},
		{
			name:     "create wallet service error",
			db:       &sqlx.DB{},
			mockRepo: mockErrorWalletRepository,
			wantErr:  true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.New(tt.db, tt.mockRepo)
			_, err := svc.Create(context.Background(), domain.User{ID: "1"})
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		name     string
		db       *sqlx.DB
		mockRepo repository.WalletRepository
		wantErr  bool
	}{
		{
			name:     "get wallet success",
			db:       &sqlx.DB{},
			mockRepo: mockWalletRepository,
			wantErr:  false,
		},
		{
			name:     "get wallet error",
			db:       &sqlx.DB{},
			mockRepo: mockErrorWalletRepository,
			wantErr:  true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.New(tt.db, tt.mockRepo)
			_, err := svc.Get(context.Background(), domain.User{ID: "1"})
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	cases := []struct {
		name     string
		db       *sqlx.DB
		mockRepo repository.WalletRepository
		wantErr  bool
	}{
		{
			name:     "withdraw success",
			db:       &sqlx.DB{},
			mockRepo: mockWalletRepository,
			wantErr:  false,
		},
		{
			name:     "withdraw error",
			db:       &sqlx.DB{},
			mockRepo: mockErrorWalletRepository,
			wantErr:  true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.New(tt.db, tt.mockRepo)
			_, err := svc.Withdraw(context.Background(), domain.User{ID: "1"}, "txn-1", 100)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDeposit(t *testing.T) {
	cases := []struct {
		name     string
		db       *sqlx.DB
		mockRepo repository.WalletRepository
		wantErr  bool
	}{
		{
			name:     "deposit success",
			db:       &sqlx.DB{},
			mockRepo: mockWalletRepository,
			wantErr:  false,
		},
		{
			name:     "deposit error",
			db:       &sqlx.DB{},
			mockRepo: mockErrorWalletRepository,
			wantErr:  true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.New(tt.db, tt.mockRepo)
			_, err := svc.Deposit(context.Background(), domain.User{ID: "1"}, "txn-1", 100)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	cases := []struct {
		name     string
		db       *sqlx.DB
		mockRepo repository.WalletRepository
		wantErr  bool
	}{
		{
			name:     "get transactions success",
			db:       &sqlx.DB{},
			mockRepo: mockWalletRepository,
			wantErr:  false,
		},
		{
			name:     "get transactions error",
			db:       &sqlx.DB{},
			mockRepo: mockErrorWalletRepository,
			wantErr:  true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.New(tt.db, tt.mockRepo)
			_, err := svc.GetTransactions(context.Background(), domain.User{ID: "1"}, time.Now(), 0, 10)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	cases := []struct {
		name     string
		db       *sqlx.DB
		mockRepo repository.WalletRepository
		wantErr  bool
	}{
		{
			name:     "transfer success",
			db:       &sqlx.DB{},
			mockRepo: mockWalletRepository,
			wantErr:  false,
		},
		{
			name:     "transfer error",
			db:       &sqlx.DB{},
			mockRepo: mockErrorWalletRepository,
			wantErr:  true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.New(tt.db, tt.mockRepo)
			_, err := svc.Transfer(context.Background(), domain.User{ID: "1"}, "txn-1", 100, domain.User{ID: "2"})
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCreateTransactionID(t *testing.T) {
	cases := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "create transaction id success",
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			svc := wallet.New(&sqlx.DB{}, nil)
			txnID := svc.CreateTransactionID(context.Background())
			assert.Equal(t, string(txnID), txnID.ID())
			assert.Equal(t, string(txnID)+"-passive", txnID.PassiveID())
		})
	}
}
