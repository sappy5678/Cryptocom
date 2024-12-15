package transport

import (
	"net/http"
	"time"

	"github.com/sappy5678/cryptocom/pkg/domain"

	"github.com/labstack/echo"
)

// HTTP represents user http service
type HTTP struct {
	svc domain.WalletService
}

// NewHTTP creates new user http service
func NewHTTP(svc domain.WalletService, r *echo.Group) {
	h := HTTP{svc}
	ur := r.Group("/user/:userID/wallet")
	// swagger:route POST /v1/users users userCreate
	// Creates new user account.
	// responses:
	//  200: userResp
	//  400: errMsg
	//  401: err
	//  403: errMsg
	//  500: err
	ur.POST("/create", h.create)

	ur.GET("", h.get)

	ur.POST("/transactionID", h.createTransactionID)

	ur.POST("/deposit", h.deposit)
	ur.POST("/withdraw", h.withdraw)
	ur.GET("/transactions", h.getTransactions)
	ur.POST("/transfer", h.transfer)
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
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID

	wallet, err := h.svc.Create(c.Request().Context(), domain.User{
		ID: r.UserID,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, wallet)
}

func (h HTTP) get(c echo.Context) error {
	r := createReq{}
	userID := c.Param("userID")
	if userID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID
	wallet, err := h.svc.Get(c.Request().Context(), domain.User{
		ID: r.UserID,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}

type createTransactionIDResp struct {
	TransactionID domain.TransactionID `json:"transactionID"`
}

func (h HTTP) createTransactionID(c echo.Context) error {
	transactionID := h.svc.CreateTransactionID(c.Request().Context())
	return c.JSON(http.StatusOK, createTransactionIDResp{TransactionID: transactionID})
}

type depositReq struct {
	UserID        string
	TransactionID string `json:"transactionID" validate:"required"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
}

func (h HTTP) deposit(c echo.Context) error {
	r := depositReq{}
	if err := c.Bind(&r); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusBadRequest)
	}
	userID := c.Param("userID")
	if userID == "" {
		c.Logger().Error("userID is required")
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID
	wallet, err := h.svc.Deposit(c.Request().Context(), domain.User{
		ID: r.UserID,
	}, domain.TransactionID(r.TransactionID), r.Amount)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}

type withdrawReq struct {
	UserID        string
	TransactionID string `json:"transactionID" validate:"required"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
}

func (h HTTP) withdraw(c echo.Context) error {
	r := withdrawReq{}
	if err := c.Bind(&r); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	userID := c.Param("userID")
	if userID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID
	wallet, err := h.svc.Withdraw(c.Request().Context(), domain.User{
		ID: r.UserID,
	}, domain.TransactionID(r.TransactionID), r.Amount)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}

type getTransactionsReq struct {
	UserID         string
	CreatedAt      time.Time `json:"createdAt"`
	LastReturnedID int       `json:"lastReturnedID"`
	Limit          int       `json:"limit"`
}

func (h HTTP) getTransactions(c echo.Context) error {
	r := getTransactionsReq{}
	if err := c.Bind(&r); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	userID := c.Param("userID")
	if userID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID
	transactions, err := h.svc.GetTransactions(c.Request().Context(), domain.User{
		ID: r.UserID,
	}, r.CreatedAt, r.LastReturnedID, r.Limit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, transactions)
}

type transferReq struct {
	UserID        string
	PassiveUserID string `json:"passiveUserID" validate:"required"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
	TransactionID string `json:"transactionID" validate:"required"`
}

func (h HTTP) transfer(c echo.Context) error {
	r := transferReq{}
	if err := c.Bind(&r); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	userID := c.Param("userID")
	if userID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID
	wallet, err := h.svc.Transfer(c.Request().Context(), domain.User{ID: r.UserID},
		domain.TransactionID(r.TransactionID), r.Amount, domain.User{ID: r.PassiveUserID})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}
