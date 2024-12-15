package transport

import (
	"net/http"
	"time"

	"github.com/sappy5678/cryptocom/pkg/domain"

	"github.com/labstack/echo"
)

// HTTP represents user http service
type HTTP struct {
	Service domain.WalletService
}

// NewHTTP creates new user http service
func NewHTTP(svc domain.WalletService, r *echo.Group) {
	h := HTTP{Service: svc}
	ur := r.Group("/user/:userID/wallet")

	// Create wallet
	// PUT /v1/users/{userID}/wallet/create
	ur.PUT("/create", h.create)

	// Get wallet
	// GET /v1/users/{userID}/wallet
	ur.GET("", h.get)

	// Get transactions
	// GET /v1/users/{userID}/wallet/transactions
	ur.GET("/transactions", h.getTransactions)

	// Create transactionID
	// POST /v1/users/{userID}/wallet/transactionID
	ur.POST("/transactionID", h.createTransactionID)

	// Deposit
	// PUT /v1/users/{userID}/wallet/deposit
	ur.PUT("/deposit", h.deposit)

	// Withdraw
	// PUT /v1/users/{userID}/wallet/withdraw
	ur.PUT("/withdraw", h.withdraw)

	// Transfer
	// PUT /v1/users/{userID}/wallet/transfer
	ur.PUT("/transfer", h.transfer)
}

// User create request
// swagger:model userCreate
type createReq struct {
	UserID string
}

func (h HTTP) create(c echo.Context) error {
	r := createReq{}
	userID := c.Param("userID")
	if userID == "" {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: domain.ErrUserIDRequired.Error()})

		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}
	r.UserID = userID

	wallet, err := h.Service.Create(c.Request().Context(), domain.User{
		ID: r.UserID,
	})

	if err != nil {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})

		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, wallet)
}

func (h HTTP) get(c echo.Context) error {
	r := createReq{}
	userID := c.Param("userID")

	if userID == "" {

		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: domain.ErrUserIDRequired.Error()})
		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	r.UserID = userID
	wallet, err := h.Service.Get(c.Request().Context(), domain.User{
		ID: r.UserID,
	})

	if err != nil {

		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})
		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, wallet)
}

type CreateTransactionIDResp struct {
	TransactionID domain.TransactionID `json:"transactionID"`
}

func (h HTTP) createTransactionID(c echo.Context) error {
	transactionID := h.Service.CreateTransactionID(c.Request().Context())

	return c.JSON(http.StatusOK, CreateTransactionIDResp{TransactionID: transactionID})
}

type DepositReq struct {
	UserID        string
	TransactionID string `json:"transactionID" validate:"required"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
}

func (h HTTP) deposit(c echo.Context) error {
	r := DepositReq{}

	if err := c.Bind(&r); err != nil {
		c.Logger().Error(err)

		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})
		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	userID := c.Param("userID")
	if userID == "" {
		c.Logger().Error("userID is required")
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: domain.ErrUserIDRequired.Error()})
		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	r.UserID = userID
	wallet, err := h.Service.Deposit(c.Request().Context(), domain.User{
		ID: r.UserID,
	}, domain.TransactionID(r.TransactionID), r.Amount)

	if err != nil {
		c.Logger().Error(err)

		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})
		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, wallet)
}

type WithdrawReq struct {
	UserID        string
	TransactionID string `json:"transactionID" validate:"required"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
}

func (h HTTP) withdraw(c echo.Context) error {
	r := WithdrawReq{}

	if err := c.Bind(&r); err != nil {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})
		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	userID := c.Param("userID")

	if userID == "" {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: domain.ErrUserIDRequired.Error()})
		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	r.UserID = userID
	wallet, err := h.Service.Withdraw(c.Request().Context(), domain.User{
		ID: r.UserID,
	}, domain.TransactionID(r.TransactionID), r.Amount)

	if err != nil {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})

		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, wallet)
}

type GetTransactionsReq struct {
	UserID           string
	CreatedBeforeStr string `query:"createdBefore"`
	IDBefore         int    `query:"IDBefore"`
	Limit            int    `query:"limit"`
}

func (h HTTP) getTransactions(c echo.Context) error {
	r := GetTransactionsReq{}

	if err := c.Bind(&r); err != nil {

		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})
		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}
	userID := c.Param("userID")

	if userID == "" {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: domain.ErrUserIDRequired.Error()})

		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	createdBefore := time.Time{}
	if r.CreatedBeforeStr != "" {
		var err error
		createdBefore, err = time.Parse(time.RFC3339, r.CreatedBeforeStr)
		if err != nil {
			err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})
			if err != nil {
				c.Logger().Error(err)
			}

			return err
		}
	}

	r.UserID = userID
	transactions, err := h.Service.GetTransactions(c.Request().Context(), domain.User{
		ID: r.UserID,
	}, createdBefore, r.IDBefore, r.Limit)

	if err != nil {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})

		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, transactions)
}

type TransferReq struct {
	UserID        string
	PassiveUserID string `json:"passiveUserID" validate:"required"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
	TransactionID string `json:"transactionID" validate:"required"`
}

func (h HTTP) transfer(c echo.Context) error {
	r := TransferReq{}
	if err := c.Bind(&r); err != nil {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})

		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}
	userID := c.Param("userID")
	if userID == "" {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: domain.ErrUserIDRequired.Error()})

		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}
	r.UserID = userID
	wallet, err := h.Service.Transfer(c.Request().Context(), domain.User{ID: r.UserID},
		domain.TransactionID(r.TransactionID), r.Amount, domain.User{ID: r.PassiveUserID})
	if err != nil {
		err := c.JSON(http.StatusBadRequest, domain.ErrorRespond{Error: err.Error()})

		if err != nil {
			c.Logger().Error(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, wallet)
}
