package truelayer

import "fmt"

type StrError string

func (str StrError) Error() string {
	return string(str)
}

// ErrorResponse is a struct representation of the TrueLayer error response.
// This is used when returning an error from an API call.
type ErrorResponse struct {
	ErrorMessage     string            `json:"error"`
	ErrorDescription string            `json:"error_description"`
	ErrorDetails     map[string]string `json:"error_details"`
}

// Error is a method implemented by all `error` interfaces. Implementing this
// means that er can treat ErrorResponse as an error.
//
// TODO: Return a better and more helpful error message.
func (res *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", res.ErrorMessage, res.ErrorDescription)
}
