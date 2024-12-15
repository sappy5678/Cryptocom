package repository_test

import (
	"context"
	"math"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sappy5678/cryptocom/pkg/domain"
	"github.com/sappy5678/cryptocom/pkg/service/wallet/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	dbConnection *sqlx.DB
	pgdb         *embeddedpostgres.EmbeddedPostgres
}

func cleanTransaction(want *domain.Transaction) {
	want.ID = 0
	want.CreatedAt = time.Time{}
}

func (ts *TestSuite) SetupSuite() {
	ts.pgdb = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Password("password").
		Port(3000).Logger(nil))

	err := ts.pgdb.Start()
	assert.NoError(ts.T(), err)
	ts.dbConnection = sqlx.MustConnect("postgres", "postgres://postgres:password@localhost:3000/postgres?sslmode=disable")
	ts.dbConnection.Exec("CREATE DATABASE cryptocom;")
	err = ts.dbConnection.Close()
	assert.NoError(ts.T(), err)
}

func (ts *TestSuite) SetupTest() {
	ts.dbConnection = sqlx.MustConnect("postgres", "postgres://postgres:password@localhost:3000/cryptocom?sslmode=disable")
	driver, err := postgres.WithInstance(ts.dbConnection.DB, &postgres.Config{})
	assert.NoError(ts.T(), err)
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../../deploy/db/migrations",
		"postgres", driver)
	assert.NoError(ts.T(), err)
	err = m.Up()
	assert.NoError(ts.T(), err)
}

func (ts *TestSuite) TearDownTest() {
	driver, err := postgres.WithInstance(ts.dbConnection.DB, &postgres.Config{})
	assert.NoError(ts.T(), err)
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../../deploy/db/migrations",
		"postgres", driver)
	assert.NoError(ts.T(), err)
	err = m.Down()
	assert.NoError(ts.T(), err)
	err = ts.dbConnection.Close()
	assert.NoError(ts.T(), err)
}

func (ts *TestSuite) TearDownSuite() {
	err := ts.dbConnection.Close()
	assert.NoError(ts.T(), err)
	err = ts.pgdb.Stop()
	assert.NoError(ts.T(), err)
}

func (ts *TestSuite) TestCreate() {

	db := ts.dbConnection

	wallet := repository.Wallet{}
	ctx := context.Background()

	tests := []struct {
		name    string
		user    domain.User
		want    *domain.Wallet
		wantErr bool
	}{
		{
			name: "normal",
			user: domain.User{ID: "test-user-1"},
			want: &domain.Wallet{
				UserID:  "test-user-1",
				Balance: 0,
			},
		},
		{
			name: "duplicate",
			user: domain.User{ID: "test-user-1"},
			want: &domain.Wallet{
				UserID:  "test-user-1",
				Balance: 0,
			},
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := wallet.Create(ctx, db, tt.user)
			if tt.wantErr {
				assert.Error(ts.T(), err)
				return
			}
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.want.UserID, got.UserID)
			assert.Equal(ts.T(), tt.want.Balance, got.Balance)
		})
	}
}

func (ts *TestSuite) TestGet() {
	db := ts.dbConnection

	wallet := repository.Wallet{}
	ctx := context.Background()

	// create a wallet for test
	testUser := domain.User{ID: "test-user-2"}
	_, err := wallet.Create(ctx, db, testUser)
	assert.NoError(ts.T(), err)

	tests := []struct {
		name    string
		user    domain.User
		want    *domain.Wallet
		wantErr error
	}{
		{
			name: "get exist wallet",
			user: testUser,
			want: &domain.Wallet{
				UserID:  "test-user-2",
				Balance: 0,
			},
		},
		{
			name:    "get non-exist wallet",
			user:    domain.User{ID: "non-exist"},
			wantErr: domain.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := wallet.Get(ctx, db, tt.user)
			if tt.wantErr != nil {
				assert.ErrorIs(ts.T(), err, tt.wantErr)
				return
			}
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.want.UserID, got.UserID)
			assert.Equal(ts.T(), tt.want.Balance, got.Balance)
		})
	}
}

func (ts *TestSuite) TestDeposit() {
	db := ts.dbConnection

	wallet := repository.Wallet{}
	ctx := context.Background()

	// create a wallet for test
	testUser := domain.User{ID: "test-user-3"}
	_, err := wallet.Create(ctx, db, testUser)
	assert.NoError(ts.T(), err)

	mockNow := time.Now().UTC().Round(time.Microsecond)

	tests := []struct {
		name            string
		user            domain.User
		amount          int
		transactionID   domain.TransactionID
		want            *domain.Wallet
		wantTransaction []*domain.Transaction
		wantErr         error
	}{
		{
			name:          "normal deposit",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-1",
			want: &domain.Wallet{
				UserID:  "test-user-3",
				Balance: 100,
			},
			wantTransaction: []*domain.Transaction{
				{
					ID:            1,
					TransactionID: "test-tx-1",
					UserID:        "test-user-3",
					Amount:        100,
					OperationType: domain.OperationTypeDeposit,
					CreatedAt:     mockNow,
				},
			},
		},
		{
			name:          "test Idempotent",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-1",
			want: &domain.Wallet{
				UserID:  "test-user-3",
				Balance: 100,
			},
			wantTransaction: []*domain.Transaction{
				{
					ID:            1,
					TransactionID: "test-tx-1",
					UserID:        "test-user-3",
					Amount:        100,
					OperationType: domain.OperationTypeDeposit,
					CreatedAt:     mockNow,
				},
			},
		},
		{
			name:          "negative deposit",
			user:          testUser,
			amount:        -100,
			transactionID: "test-tx-2",
			wantErr:       domain.ErrInvalidAmount,
		},
		{
			name:          "deposit from non-exist user",
			user:          domain.User{ID: "non-exist"},
			amount:        100,
			transactionID: "test-tx-3",
			wantErr:       domain.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := wallet.Deposit(ctx, db, mockNow, tt.user, tt.transactionID, tt.amount)
			if tt.wantErr != nil {
				assert.ErrorIs(ts.T(), err, tt.wantErr)
				return
			}
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.want.UserID, got.UserID)
			assert.Equal(ts.T(), tt.want.Balance, got.Balance)

			// check transaction
			gotTransactions, err := wallet.GetTransactions(ctx, db, tt.user, time.Time{}, math.MaxInt64, 100)
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.wantTransaction, gotTransactions)
		})
	}
}

func (ts *TestSuite) TestWithdraw() {
	db := ts.dbConnection

	wallet := repository.Wallet{}
	ctx := context.Background()
	mockNow := time.Now().UTC().Round(time.Microsecond)

	// create a wallet for test, and deposit 1000 for test
	testUser := domain.User{ID: "test-user-4"}
	_, err := wallet.Create(ctx, db, testUser)
	assert.NoError(ts.T(), err)
	_, err = wallet.Deposit(ctx, db, mockNow, testUser, "test-tx-0", 1000)
	assert.NoError(ts.T(), err)

	tests := []struct {
		name            string
		user            domain.User
		amount          int
		transactionID   domain.TransactionID
		want            *domain.Wallet
		wantTransaction []*domain.Transaction
		wantErr         error
	}{
		{
			name:          "normal Withdraw",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-1",
			want: &domain.Wallet{
				UserID:  "test-user-4",
				Balance: 900,
			},
			wantTransaction: []*domain.Transaction{
				{
					ID:            2,
					TransactionID: "test-tx-1",
					UserID:        "test-user-4",
					Amount:        100,
					OperationType: domain.OperationTypeWithdraw,
					CreatedAt:     mockNow,
				},
				{
					ID:            1,
					TransactionID: "test-tx-0",
					UserID:        "test-user-4",
					Amount:        1000,
					OperationType: domain.OperationTypeDeposit,
					CreatedAt:     mockNow,
				},
			},
		},
		{
			name:          "test Idempotent",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-1",
			want: &domain.Wallet{
				UserID:  "test-user-4",
				Balance: 900,
			},
			wantTransaction: []*domain.Transaction{
				{
					ID:            2,
					TransactionID: "test-tx-1",
					UserID:        "test-user-4",
					Amount:        100,
					OperationType: domain.OperationTypeWithdraw,
					CreatedAt:     mockNow,
				},
				{
					ID:            1,
					TransactionID: "test-tx-0",
					UserID:        "test-user-4",
					Amount:        1000,
					OperationType: domain.OperationTypeDeposit,
					CreatedAt:     mockNow,
				},
			},
		},
		{
			name:          "negative withdraw",
			user:          testUser,
			amount:        -100,
			transactionID: "test-tx-2",
			wantErr:       domain.ErrInvalidAmount,
		},
		{
			name:          "not enough balance",
			user:          testUser,
			amount:        1000,
			transactionID: "test-tx-3",
			wantErr:       domain.ErrNotEnoughBalance,
		},
		{
			name:          "withdraw from non-exist user",
			user:          domain.User{ID: "non-exist"},
			amount:        100,
			transactionID: "test-tx-4",
			wantErr:       domain.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := wallet.Withdraw(ctx, db, mockNow, tt.user, tt.transactionID, tt.amount)
			if tt.wantErr != nil {
				assert.ErrorIs(ts.T(), err, tt.wantErr)
				return
			}
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.want.UserID, got.UserID)
			assert.Equal(ts.T(), tt.want.Balance, got.Balance)

			// check transaction
			gotTransactions, err := wallet.GetTransactions(ctx, db, tt.user, time.Now(), math.MaxInt64, 100)
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.wantTransaction, gotTransactions)
		})
	}
}

func (ts *TestSuite) TestTransfer() {
	db := ts.dbConnection

	wallet := repository.Wallet{}
	ctx := context.Background()
	mockNow := repository.TimeToDBTime(time.Now())

	// create a wallet for test, and deposit 1000 for test
	testUser := domain.User{ID: "test-user-5"}
	_, err := wallet.Create(ctx, db, testUser)
	assert.NoError(ts.T(), err)
	_, err = wallet.Deposit(ctx, db, repository.TimeToDBTime(mockNow.Add(10*time.Second)), testUser, "test-tx-0", 1000)
	assert.NoError(ts.T(), err)

	passiveUser := domain.User{ID: "test-user-6"}
	_, err = wallet.Create(ctx, db, passiveUser)
	assert.NoError(ts.T(), err)

	tests := []struct {
		name                   string
		user                   domain.User
		amount                 int
		transactionID          domain.TransactionID
		passiveUser            domain.User
		want                   *domain.Wallet
		wantPassive            *domain.Wallet
		wantTransaction        []*domain.Transaction
		wantPassiveTransaction []*domain.Transaction
		wantErr                error
	}{
		{
			name:          "normal transfer",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-1",
			passiveUser:   passiveUser,
			want: &domain.Wallet{
				UserID:  "test-user-5",
				Balance: 900,
			},
			wantPassive: &domain.Wallet{
				UserID:  "test-user-6",
				Balance: 100,
			},
			wantTransaction: []*domain.Transaction{
				{
					ID:            1,
					TransactionID: "test-tx-0",
					UserID:        "test-user-5",
					Amount:        1000,
					OperationType: domain.OperationTypeDeposit,
					CreatedAt:     mockNow,
				},
				{
					ID:            2,
					TransactionID: "test-tx-1",
					UserID:        "test-user-5",
					Amount:        100,
					OperationType: domain.OperationTypeTransferOut,
					PassiveUserID: passiveUser.ID,
					CreatedAt:     mockNow,
				},
			},
			wantPassiveTransaction: []*domain.Transaction{
				{
					TransactionID: "test-tx-1-passive",
					UserID:        "test-user-6",
					Amount:        100,
					OperationType: domain.OperationTypeTransferIn,
					PassiveUserID: testUser.ID,
					CreatedAt:     mockNow,
				},
			},
		},
		{
			name:          "test Idempotent",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-1",
			passiveUser:   passiveUser,
			want: &domain.Wallet{
				UserID:  "test-user-5",
				Balance: 900,
			},
			wantPassive: &domain.Wallet{
				UserID:  "test-user-6",
				Balance: 100,
			},
			wantTransaction: []*domain.Transaction{
				{
					ID:            1,
					TransactionID: "test-tx-0",
					UserID:        "test-user-5",
					Amount:        1000,
					OperationType: domain.OperationTypeDeposit,
					CreatedAt:     mockNow,
				},
				{
					ID:            2,
					TransactionID: "test-tx-1",
					UserID:        "test-user-5",
					Amount:        100,
					OperationType: domain.OperationTypeTransferOut,
					PassiveUserID: passiveUser.ID,
					CreatedAt:     mockNow,
				},
			},
			wantPassiveTransaction: []*domain.Transaction{
				{
					TransactionID: "test-tx-1-passive",
					UserID:        "test-user-6",
					Amount:        100,
					OperationType: domain.OperationTypeTransferIn,
					PassiveUserID: testUser.ID,
					CreatedAt:     mockNow,
				},
			},
		},
		{
			name:          "transfer two times",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-2",
			passiveUser:   passiveUser,
			want: &domain.Wallet{
				UserID:  "test-user-5",
				Balance: 800,
			},
			wantPassive: &domain.Wallet{
				UserID:  "test-user-6",
				Balance: 200,
			},
			wantTransaction: []*domain.Transaction{
				{
					ID:            1,
					TransactionID: "test-tx-0",
					UserID:        "test-user-5",
					Amount:        1000,
					OperationType: domain.OperationTypeDeposit,
					CreatedAt:     mockNow,
				},
				{
					ID:            3,
					TransactionID: "test-tx-1",
					UserID:        "test-user-5",
					Amount:        100,
					OperationType: domain.OperationTypeTransferOut,
					PassiveUserID: passiveUser.ID,
					CreatedAt:     mockNow,
				},
				{
					ID:            5,
					TransactionID: "test-tx-2",
					UserID:        "test-user-5",
					Amount:        100,
					OperationType: domain.OperationTypeTransferOut,
					PassiveUserID: passiveUser.ID,
					CreatedAt:     mockNow,
				},
			},
			wantPassiveTransaction: []*domain.Transaction{
				{
					TransactionID: "test-tx-1-passive",
					UserID:        "test-user-6",
					Amount:        100,
					OperationType: domain.OperationTypeTransferIn,
					PassiveUserID: testUser.ID,
					CreatedAt:     mockNow,
				},
				{
					TransactionID: "test-tx-2-passive",
					UserID:        "test-user-6",
					Amount:        100,
					OperationType: domain.OperationTypeTransferIn,
					PassiveUserID: testUser.ID,
					CreatedAt:     mockNow,
				},
			},
		},
		{
			name:          "negative transfer",
			user:          testUser,
			amount:        -100,
			transactionID: "test-tx-3",
			passiveUser:   passiveUser,
			wantErr:       domain.ErrInvalidAmount,
		},
		{
			name:          "not enough balance",
			user:          testUser,
			amount:        10000,
			transactionID: "test-tx-4",
			passiveUser:   passiveUser,
			wantErr:       domain.ErrNotEnoughBalance,
		},
		{
			name:          "transfer to self",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-5",
			passiveUser:   testUser,
			wantErr:       domain.ErrTransferToSelf,
		},
		{
			name:          "transfer to non-exist user",
			user:          testUser,
			amount:        100,
			transactionID: "test-tx-6",
			passiveUser:   domain.User{ID: "non-exist"},
			wantErr:       domain.ErrWalletNotFound,
		},
		{
			name:          "transfer from non-exist user",
			user:          domain.User{ID: "non-exist"},
			amount:        100,
			transactionID: "test-tx-7",
			passiveUser:   passiveUser,
			wantErr:       domain.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := wallet.Transfer(ctx, db, mockNow, tt.user, tt.transactionID, tt.amount, tt.passiveUser)
			if tt.wantErr != nil {
				assert.ErrorIs(ts.T(), err, tt.wantErr)
				return
			}
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.want.UserID, got.UserID)
			assert.Equal(ts.T(), tt.want.Balance, got.Balance)

			// check passive wallet
			gotPassive, err := wallet.Get(ctx, db, tt.passiveUser)
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.wantPassive.UserID, gotPassive.UserID)
			assert.Equal(ts.T(), tt.wantPassive.Balance, gotPassive.Balance)

			// check transaction
			gotTransactions, err := wallet.GetTransactions(ctx, db, tt.user, time.Now(), math.MaxInt64, 100)
			assert.NoError(ts.T(), err)
			for i, want := range tt.wantTransaction {
				cleanTransaction(want)
				cleanTransaction(gotTransactions[i])
			}
			assert.ElementsMatch(ts.T(), tt.wantTransaction, gotTransactions)

			// check passive transaction
			gotPassiveTransactions, err := wallet.GetTransactions(ctx, db, tt.passiveUser, time.Now(), math.MaxInt64, 100)
			assert.NoError(ts.T(), err)
			for i, want := range tt.wantPassiveTransaction {
				cleanTransaction(want)
				cleanTransaction(gotPassiveTransactions[i])
			}
			assert.ElementsMatch(ts.T(), tt.wantPassiveTransaction, gotPassiveTransactions)
		})
	}
}

func (ts *TestSuite) TestGetTransactions() {
	db := ts.dbConnection

	wallet := repository.Wallet{}
	ctx := context.Background()
	mockNow := time.Now()

	// create a wallet for test, and create transactions
	testUser := domain.User{ID: "test-user-7"}
	testPassiveUser := domain.User{ID: "test-user-8"}
	_, err := wallet.Create(ctx, db, testUser)
	assert.NoError(ts.T(), err)
	_, err = wallet.Create(ctx, db, testPassiveUser)
	assert.NoError(ts.T(), err)
	_, err = wallet.Deposit(ctx, db, mockNow, testUser, "test-tx-1", 1000)
	assert.NoError(ts.T(), err)
	_, err = wallet.Withdraw(ctx, db, mockNow, testUser, "test-tx-2", 100)
	assert.NoError(ts.T(), err)
	_, err = wallet.Transfer(ctx, db, mockNow, testUser, "test-tx-3", 100, testPassiveUser)
	assert.NoError(ts.T(), err)
	_, err = wallet.Transfer(ctx, db, mockNow, testPassiveUser, "test-tx-4", 100, testUser)
	assert.NoError(ts.T(), err)

	tests := []struct {
		name    string
		user    domain.User
		want    []*domain.Transaction
		wantErr error
	}{
		{
			name: "normal",
			user: domain.User{ID: "test-user-7"},
			want: []*domain.Transaction{
				{
					ID:            6,
					TransactionID: "test-tx-4-passive",
					UserID:        "test-user-7",
					Amount:        100,
					OperationType: domain.OperationTypeTransferIn,
					PassiveUserID: testPassiveUser.ID,
					CreatedAt:     mockNow,
				},
				{
					ID:            2,
					TransactionID: "test-tx-3",
					UserID:        "test-user-7",
					Amount:        100,
					PassiveUserID: testPassiveUser.ID,
					OperationType: domain.OperationTypeTransferOut,
					CreatedAt:     mockNow,
				},
				{
					ID:            1,
					TransactionID: "test-tx-2",
					UserID:        "test-user-7",
					Amount:        100,
					OperationType: domain.OperationTypeWithdraw,
					CreatedAt:     mockNow,
				},
				{
					ID:            0,
					TransactionID: "test-tx-1",
					UserID:        "test-user-7",
					Amount:        1000,
					OperationType: domain.OperationTypeDeposit,
					CreatedAt:     mockNow,
				},
			},
			wantErr: nil,
		},
		{
			name:    "no transaction",
			user:    domain.User{ID: "not-exist"},
			wantErr: domain.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func() {

			got, err := wallet.GetTransactions(ctx, db, tt.user, time.Now(), math.MaxInt64, 100)
			if tt.wantErr != nil {
				assert.ErrorIs(ts.T(), err, tt.wantErr, tt.name+": error is not equal")
				return
			}
			assert.NoError(ts.T(), err, tt.name+": error is not nil")
			for i, want := range tt.want {
				cleanTransaction(want)
				cleanTransaction(got[i])
			}
			assert.ElementsMatch(ts.T(), tt.want, got)
		})
	}
}

func (ts *TestSuite) TestExistsTransactionID() {
	db := ts.dbConnection

	wallet := repository.Wallet{}
	ctx := context.Background()
	mockNow := repository.TimeToDBTime(time.Now())
	// create a wallet for test, and create transactions
	testUser := domain.User{ID: "test-user-9"}
	_, err := wallet.Create(ctx, db, testUser)
	assert.NoError(ts.T(), err)
	_, err = wallet.Deposit(ctx, db, mockNow, testUser, "test-tx-1", 1000)
	assert.NoError(ts.T(), err)

	tests := []struct {
		name          string
		transactionID domain.TransactionID
		want          bool
		wantErr       error
	}{
		{
			name:          "normal",
			transactionID: "test-tx-1",
			want:          true,
		},
		{
			name:          "not exist",
			transactionID: "test-tx-2",
			want:          false,
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := wallet.ExistsTransactionID(ctx, db, tt.transactionID)
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.want, got)
		})
	}
}

func (ts *TestSuite) TestExists() {
	db := ts.dbConnection

	wallet := repository.Wallet{}
	ctx := context.Background()
	// create a wallet for test
	testUser := domain.User{ID: "test-user-10"}
	_, err := wallet.Create(ctx, db, testUser)
	assert.NoError(ts.T(), err)

	tests := []struct {
		name    string
		user    domain.User
		want    bool
		wantErr error
	}{
		{
			name: "normal",
			user: testUser,
			want: true,
		},
		{
			name: "not exist",
			user: domain.User{ID: "not-exist"},
			want: false,
		},
	}

	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := wallet.Exists(ctx, db, tt.user)
			assert.NoError(ts.T(), err)
			assert.Equal(ts.T(), tt.want, got)
		})
	}
}

func TestWalletSuite(t *testing.T) {
	ts := new(TestSuite)
	suite.Run(t, ts)
}
