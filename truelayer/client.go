package truelayer

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type TrueLayer struct {
	clientID     string
	clientSecret string
	sandbox      bool
	httpClient   httpClient
}

// httpClient is an interface to define the methods required from any kind of
// HTTP Client that will be used by the TrueLayer Client.
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// New creates a new instance of the TrueLayer go client. This is done to allow
// for mocking within user implementation to allow for greater test coverage.
//
// params
//   - clientID - TrueLayer client_id
//   - clientSecret - TrueLayer client_secret
//   - sandbox - true if using the sandbox environment
//
// returns
//   - instance of TrueLayer client
func New(clientID, clientSecret string, sandbox bool) *TrueLayer {
	return NewWithHTTPClient(clientID, clientSecret, sandbox, &http.Client{})
}

// NewWithHTTPClient creates a new instance of the TrueLayer go client, with a
// custom HTTP Client. This is done to allow for mocking within user
// implementation to allow for greater test coverage.
//
// params
//   - clientID - TrueLayer client_id
//   - clientSecret - TrueLayer client_secret
//   - sandbox - true if using the sandbox environment
//   - httpClient - custom HTTP client for the TrueLayer client to use
//
// returns
//   - instance of TrueLayer client
func NewWithHTTPClient(clientID, clientSecret string, sandbox bool, httpClient httpClient) *TrueLayer {
	return &TrueLayer{
		clientID:     clientID,
		clientSecret: clientSecret,
		sandbox:      sandbox,
		httpClient:   httpClient,
	}
}

// buildURL takes a base URL as well as a path and combines them into a url.URL
// object.
//
// params
//   - baseURL - the base url to use
//   - path - the path to combine
//
// returns
//   - combined url
//   - parsing errors
func buildURL(baseURL string, path string) (*url.URL, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	u.Path = path

	return u, nil
}

// parseErrorResponse takes a http.Response object and decodes the body into an
// ErrorResponse object which implements the `error` interface.
//
// params
//   - res - the http response to decode
//
// returns
//   - err - the decoded error or the error returned from decoding
func parseErrorResponse(res *http.Response) (err error) {
	respErr := &ErrorResponse{}
	err = json.NewDecoder(res.Body).Decode(respErr)
	if err != nil {
		return err
	}
	return respErr
}
