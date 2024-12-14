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

}

// User create request
// swagger:model userCreate
type createReq struct {
	UserID string `param:"userID" validate:"required"`
}

func (h HTTP) create(c echo.Context) error {
	r := createReq{}
	userID := c.Param("userID")
	if userID == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	r.UserID = userID

	err := h.svc.Create(c.Request().Context(), domain.User{
		ID: r.UserID,
	})

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
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
