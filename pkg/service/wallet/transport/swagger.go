package transport

import "github.com/sappy5678/cryptocom"

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*cryptocom.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []cryptocom.User `json:"users"`
		Page  int              `json:"page"`
	}
}
