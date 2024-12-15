package wallet_test

import (
	"context"
	"testing"
	"time"

	"github.com/sappy5678/cryptocom/pkg/domain"
	"github.com/sappy5678/cryptocom/pkg/service/wallet"
	wl "github.com/sappy5678/cryptocom/pkg/service/wallet/logging"
	"github.com/sappy5678/cryptocom/pkg/utl/zlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

var mockWalletService = &wallet.MockWalletService{
	CreateFunc: func(ctx context.Context, user domain.User) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: user.ID, Balance: 0}, nil
	},
	GetFunc: func(ctx context.Context, user domain.User) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: user.ID, Balance: 0}, nil
	},
	WithdrawFunc: func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: user.ID, Balance: 0}, nil
	},
	DepositFunc: func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: user.ID, Balance: 0}, nil
	},
	GetTransactionsFunc: func(ctx context.Context, user domain.User, createdAt time.Time, lastReturnedID int, limit int) ([]*domain.Transaction, error) {
		return []*domain.Transaction{}, nil
	},
	CreateTransactionIDFunc: func(ctx context.Context) domain.TransactionID {
		return domain.TransactionID("test-transaction-id")
	},
	TransferFunc: func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {
		return &domain.Wallet{UserID: user.ID, Balance: 0}, nil
	},
}

func TestCreate(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := zlog.New()
	svc := wl.New(mockWalletService, log)
	r1, e1 := svc.Create(context.Background(), domain.User{ID: "test-user-id"})
	r2, e2 := mockWalletService.Create(context.Background(), domain.User{ID: "test-user-id"})

	assert.Equal(t, r1, r2)
	assert.Equal(t, e1, e2)
}

func TestGet(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := zlog.New()
	svc := wl.New(mockWalletService, log)
	r1, e1 := svc.Get(context.Background(), domain.User{ID: "test-user-id"})
	r2, e2 := mockWalletService.Get(context.Background(), domain.User{ID: "test-user-id"})

	assert.Equal(t, r1, r2)
	assert.Equal(t, e1, e2)
}

func TestWithdraw(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := zlog.New()
	svc := wl.New(mockWalletService, log)
	r1, e1 := svc.Withdraw(context.Background(), domain.User{ID: "test-user-id"}, "txn-1", 100)
	r2, e2 := mockWalletService.Withdraw(context.Background(), domain.User{ID: "test-user-id"}, "txn-1", 100)

	assert.Equal(t, r1, r2)
	assert.Equal(t, e1, e2)
}

func TestDeposit(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := zlog.New()
	svc := wl.New(mockWalletService, log)
	r1, e1 := svc.Deposit(context.Background(), domain.User{ID: "test-user-id"}, "txn-1", 100)
	r2, e2 := mockWalletService.Deposit(context.Background(), domain.User{ID: "test-user-id"}, "txn-1", 100)

	assert.Equal(t, r1, r2)
	assert.Equal(t, e1, e2)
}

func TestGetTransactions(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := zlog.New()
	svc := wl.New(mockWalletService, log)
	r1, e1 := svc.GetTransactions(context.Background(), domain.User{ID: "test-user-id"}, time.Now(), 0, 10)
	r2, e2 := mockWalletService.GetTransactions(context.Background(), domain.User{ID: "test-user-id"}, time.Now(), 0, 10)

	assert.Equal(t, r1, r2)
	assert.Equal(t, e1, e2)
}

func TestCreateTransactionID(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := zlog.New()
	svc := wl.New(mockWalletService, log)
	r1 := svc.CreateTransactionID(context.Background())
	r2 := mockWalletService.CreateTransactionID(context.Background())

	assert.Equal(t, r1, r2)
}

func TestTransfer(t *testing.T) {
	defer goleak.VerifyNone(t)

	log := zlog.New()
	svc := wl.New(mockWalletService, log)
	r1, e1 := svc.Transfer(context.Background(), domain.User{ID: "test-user-id"}, "txn-1", 100, domain.User{ID: "test-user-id-2"})
	r2, e2 := mockWalletService.Transfer(context.Background(), domain.User{ID: "test-user-id"}, "txn-1", 100, domain.User{ID: "test-user-id-2"})

	assert.Equal(t, r1, r2)
	assert.Equal(t, e1, e2)
}
