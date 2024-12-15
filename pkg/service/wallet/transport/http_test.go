package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/sappy5678/cryptocom/pkg/domain"
	"github.com/sappy5678/cryptocom/pkg/service/wallet"
	"github.com/sappy5678/cryptocom/pkg/service/wallet/transport"
	"github.com/sappy5678/cryptocom/pkg/utl/server"
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

var mockError = errors.New("error")
var mockErrorWalletService = &wallet.MockWalletService{
	CreateFunc: func(ctx context.Context, user domain.User) (*domain.Wallet, error) {
		return nil, mockError
	},
	GetFunc: func(ctx context.Context, user domain.User) (*domain.Wallet, error) {
		return nil, mockError
	},
	WithdrawFunc: func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
		return nil, mockError
	},
	DepositFunc: func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int) (*domain.Wallet, error) {
		return nil, mockError
	},
	GetTransactionsFunc: func(ctx context.Context, user domain.User, createdAt time.Time, lastReturnedID int, limit int) ([]*domain.Transaction, error) {
		return nil, mockError
	},
	CreateTransactionIDFunc: func(ctx context.Context) domain.TransactionID {
		return domain.TransactionID("test-transaction-id")
	},
	TransferFunc: func(ctx context.Context, user domain.User, transactionID domain.TransactionID, amount int, passiveUser domain.User) (*domain.Wallet, error) {
		return nil, mockError
	},
}

func TestCreate(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name        string
		userID      string
		wantStatus  int
		wantResp    *domain.Wallet
		wantErrResp *domain.ErrorRespond
		svc         domain.WalletService
	}{
		{
			name:       "success",
			userID:     "1",
			wantStatus: http.StatusOK,
			wantResp: &domain.Wallet{
				UserID:  "1",
				Balance: 0,
			},
			svc: mockWalletService,
		},
		{
			name:       "userID is empty",
			userID:     "",
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrUserIDRequired.Error(),
			},
			svc: mockWalletService,
		},
		{
			name:       "error",
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: mockError.Error(),
			},
			svc: mockErrorWalletService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("v1")
			transport.NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/user/" + tt.userID + "/wallet/create"
			req, err := http.NewRequest(http.MethodPut, path, bytes.NewBufferString("{}"))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(domain.Wallet)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			} else {
				response := new(domain.ErrorRespond)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantErrResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestGet(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name        string
		userID      string
		wantStatus  int
		wantResp    *domain.Wallet
		wantErrResp *domain.ErrorRespond
		svc         domain.WalletService
	}{
		{
			name:       "normal",
			userID:     "1",
			wantStatus: http.StatusOK,
			wantResp: &domain.Wallet{
				UserID: "1",
			},
			svc: mockWalletService,
		},
		{
			name:       "missing userID",
			userID:     "",
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrUserIDRequired.Error(),
			},
			svc: mockWalletService,
		},
		{
			name:       "error",
			userID:     "1",
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: mockError.Error(),
			},
			svc: mockErrorWalletService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("v1")
			transport.NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/user/" + tt.userID + "/wallet"
			res, err := http.Get(path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(domain.Wallet)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			} else {
				response := new(domain.ErrorRespond)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantErrResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestCreateTransactionID(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name       string
		userID     string
		wantStatus int
		wantResp   *transport.CreateTransactionIDResp
		svc        domain.WalletService
	}{
		{
			name:       "normal",
			userID:     "1",
			wantStatus: http.StatusOK,
			wantResp: &transport.CreateTransactionIDResp{
				TransactionID: domain.TransactionID("test-transaction-id"),
			},
			svc: mockWalletService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("v1")
			transport.NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/user/" + tt.userID + "/wallet/transactionID"
			res, err := http.Post(path, "application/json", bytes.NewBufferString("{}"))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			response := new(transport.CreateTransactionIDResp)
			if err := json.NewDecoder(res.Body).Decode(response); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.wantResp, response)
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestDeposit(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := []struct {
		name        string
		userID      string
		req         transport.DepositReq
		wantStatus  int
		wantResp    *domain.Wallet
		wantErrResp *domain.ErrorRespond
		svc         domain.WalletService
	}{
		{
			name:   "normal",
			userID: "1",
			req: transport.DepositReq{
				TransactionID: "txn-1",
				Amount:        100,
			},
			wantStatus: http.StatusOK,
			wantResp: &domain.Wallet{
				UserID: "1",
			},
			svc: mockWalletService,
		},
		{
			name:   "missing userID",
			userID: "",
			req: transport.DepositReq{
				TransactionID: "txn-1",
				Amount:        100,
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrUserIDRequired.Error(),
			},
			svc: mockWalletService,
		},
		{
			name:   "error",
			userID: "1",
			req: transport.DepositReq{
				TransactionID: "txn-1",
				Amount:        100,
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: mockError.Error(),
			},
			svc: mockErrorWalletService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("v1")
			transport.NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/user/" + tt.userID + "/wallet/deposit"
			reqBody, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(domain.Wallet)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			} else {
				response := new(domain.ErrorRespond)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantErrResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestWithdraw(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name        string
		userID      string
		req         transport.WithdrawReq
		wantStatus  int
		wantResp    *domain.Wallet
		wantErrResp *domain.ErrorRespond
		svc         domain.WalletService
	}{
		{
			name:   "success",
			userID: "1",
			req: transport.WithdrawReq{
				TransactionID: "txn-1",
				Amount:        100,
			},
			wantStatus: http.StatusOK,
			wantResp: &domain.Wallet{
				UserID: "1",
			},
			svc: mockWalletService,
		},
		{
			name:   "missing userID",
			userID: "",
			req: transport.WithdrawReq{
				TransactionID: "txn-1",
				Amount:        100,
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrUserIDRequired.Error(),
			},
			svc: mockWalletService,
		},
		{
			name:   "error",
			userID: "1",
			req: transport.WithdrawReq{
				TransactionID: "txn-1",
				Amount:        100,
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: mockError.Error(),
			},
			svc: mockErrorWalletService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("v1")
			transport.NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/user/" + tt.userID + "/wallet/withdraw"
			reqBody, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(domain.Wallet)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			} else {
				response := new(domain.ErrorRespond)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantErrResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestGetTransactions(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name        string
		userID      string
		req         transport.GetTransactionsReq
		wantStatus  int
		wantResp    []*domain.Transaction
		wantErrResp *domain.ErrorRespond
		svc         domain.WalletService
	}{
		{
			name:   "success",
			userID: "1",
			req: transport.GetTransactionsReq{
				CreatedBeforeStr: time.Now().Format(time.RFC3339),
				IDBefore:         0,
				Limit:            10,
			},
			wantStatus: http.StatusOK,
			wantResp:   []*domain.Transaction{},
			svc:        mockWalletService,
		},
		{
			name:   "missing userID",
			userID: "",
			req: transport.GetTransactionsReq{
				CreatedBeforeStr: time.Now().Format(time.RFC3339),
				IDBefore:         0,
				Limit:            10,
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrUserIDRequired.Error(),
			},
			svc: mockWalletService,
		},
		{
			name:   "error",
			userID: "1",
			req: transport.GetTransactionsReq{
				CreatedBeforeStr: time.Now().Format(time.RFC3339),
				IDBefore:         0,
				Limit:            10,
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: mockError.Error(),
			},
			svc: mockErrorWalletService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("v1")
			transport.NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/user/" + tt.userID + "/wallet/transactions"
			query := url.Values{}
			query.Add("createdBefore", tt.req.CreatedBeforeStr)
			query.Add("IDBefore", strconv.Itoa(tt.req.IDBefore))
			query.Add("limit", strconv.Itoa(tt.req.Limit))

			req, err := http.NewRequest(http.MethodGet, path+"?"+query.Encode(), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				var response []*domain.Transaction
				if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			} else {
				response := new(domain.ErrorRespond)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantErrResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestTransfer(t *testing.T) {
	defer goleak.VerifyNone(t)
	tests := []struct {
		name        string
		userID      string
		req         transport.TransferReq
		wantStatus  int
		wantResp    *domain.Wallet
		wantErrResp *domain.ErrorRespond
		svc         domain.WalletService
	}{
		{
			name:   "success",
			userID: "1",
			req: transport.TransferReq{
				PassiveUserID: "2",
				Amount:        100,
				TransactionID: "txn-1",
			},
			wantStatus: http.StatusOK,
			wantResp: &domain.Wallet{
				UserID: "1",
			},
			svc: mockWalletService,
		},
		{
			name:   "missing userID",
			userID: "",
			req: transport.TransferReq{
				PassiveUserID: "2",
				Amount:        100,
				TransactionID: "txn-1",
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: domain.ErrUserIDRequired.Error(),
			},
			svc: mockWalletService,
		},
		{
			name:   "error",
			userID: "1",
			req: transport.TransferReq{
				PassiveUserID: "2",
				Amount:        100,
				TransactionID: "txn-1",
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   nil,
			wantErrResp: &domain.ErrorRespond{
				Error: mockError.Error(),
			},
			svc: mockErrorWalletService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			rg := r.Group("v1")
			transport.NewHTTP(tt.svc, rg)
			ts := httptest.NewServer(r)
			defer ts.Close()
			path := ts.URL + "/v1/user/" + tt.userID + "/wallet/transfer"
			reqBody, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if tt.wantResp != nil {
				response := new(domain.Wallet)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantResp, response)
			} else {
				response := new(domain.ErrorRespond)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.wantErrResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
