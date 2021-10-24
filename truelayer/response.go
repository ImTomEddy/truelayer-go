package truelayer

import (
	"fmt"
	"time"
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

type Account struct {
	UpdateTimestamp time.Time `json:"update_timestamp"`
	AccountID       string    `json:"account_id"`
	AccountType     string    `json:"account_type"`
	DisplayName     string    `json:"display_name"`
	Currency        string    `json:"currency"`
	AccountNumber   struct {
		Iban     string `json:"iban"`
		Number   string `json:"number"`
		SortCode string `json:"sort_code"`
		SwiftBic string `json:"swift_bic"`
	} `json:"account_number"`
	Provider struct {
		ProviderID string `json:"provider_id"`
	} `json:"provider"`
}
