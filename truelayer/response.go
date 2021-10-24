package truelayer

import "fmt"

// ErrorResponse is a struct representation of the TrueLayer error response.
// This is used when returning an error from an API call.
type ErrorResponse struct {
	ErrorMessage     string            `json:"error"`
	ErrorDescription string            `json:"error_description"`
	ErrorDetails     map[string]string `json:"error_details"`
}

func (res *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", res.ErrorMessage, res.ErrorDescription)
}

// AccessTokenResponse is the JSON Structure returned when requesting an
// AccessToken from TrueLayer.
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}
