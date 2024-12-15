package transport

import (
	"net/http"

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

type depositReq struct {
	UserID string
	Amount int `json:"amount" validate:"required,gt=0"`
}

func (h HTTP) deposit(c echo.Context) error {
	r := depositReq{}
	if err := c.Bind(&r); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	userID := c.Param("userID")
	if userID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID
	wallet, err := h.svc.Deposit(c.Request().Context(), domain.User{
		ID: r.UserID,
	}, r.Amount)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}

type withdrawReq struct {
	UserID string
	Amount int `json:"amount" validate:"required,gt=0"`
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
	}, r.Amount)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}

type getTransactionsReq struct {
	UserID string
}

func (h HTTP) getTransactions(c echo.Context) error {
	r := getTransactionsReq{}
	userID := c.Param("userID")
	if userID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID
	transactions, err := h.svc.GetTransactions(c.Request().Context(), domain.User{
		ID: r.UserID,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, transactions)
}

type transferReq struct {
	UserID        string
	PassiveUserID string `json:"passiveUserID" validate:"required"`
	Amount        int    `json:"amount" validate:"required,gt=0"`
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
	wallet, err := h.svc.Transfer(c.Request().Context(), domain.User{
		ID: r.UserID,
	}, r.Amount, domain.User{
		ID: r.PassiveUserID,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, wallet)
}
