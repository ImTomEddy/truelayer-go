package truelayer

import (
	"fmt"
)

// ErrorResponse is a struct representation of the TrueLayer error response.
// This is used when returning an error from an API call.
type ErrorResponse struct {
	ErrorMessage     string            `json:"error"`
	ErrorDescription string            `json:"error_description"`
	ErrorDetails     map[string]string `json:"error_details"`
}

// Error is a method implemented by all `error` interfaces. Implementing this
// means that er can treat ErrorResponse as an error.
func (res *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", res.ErrorMessage, res.ErrorDescription)
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type AccountsResponse struct {
	Results []Account `json:"results"`
}

type BalanceResponse struct {
	Results []Balance `json:"results"`
}
