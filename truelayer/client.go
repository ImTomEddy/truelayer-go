package truelayer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TrueLayer struct {
	clientID     string
	clientSecret string
	sandbox      bool
	httpClient   httpClient
}

const (
	baseURL        = "https://api.truelayer.com"
	baseSandboxURL = "https://api.truelayer-sandbox.com"

	EndpointDataV1Results = "/data/v1/results/%s"

	ErrRequestBodyNil = StrError("request body is nil")
)

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

// doAuthorizedGetRequest executes a request with an Authorization header with
// the provided accessToken to the provided URL.
//
// params
//   - url - the URL to make a request
//   - accessToken - the access token to use
//
// returns
//   - the http response
//   - any errors that have occurred
func (t *TrueLayer) doAuthorizedGetRequest(url *url.URL, accessToken string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := t.httpClient.Do(req)

	return res, err
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

// getBaseURL parses the baseAuthURL for either the sandbox or non-sandbox
// TrueLayer environments and returns them. Using a utility method to reduce
// code duplication.
//
// returns
//   - the base url
func (t *TrueLayer) getBaseURL() string {
	if t.sandbox {
		return baseSandboxURL
	}

	return baseURL
}

// HandleAsyncWebhookRequest will take an HTTP request and return the
// WebhookRequest object or an error. This is used as a part of the async flow
// for the TrueLayer api.
//
// params
//   - req - the http request to handle
//
// returns
//   - the webhook request
//   - error if an error occurs
func (t *TrueLayer) HandleAsyncWebhookRequest(req *http.Request) (*WebhookRequest, error) {
	if req.Body == nil {
		return nil, ErrRequestBodyNil
	}

	return nil, nil
}

// HandleAsyncWebhookRequestBody will take an io.ReadCloser and return the
// WebhookRequest object or an error. This is used as a part of the async flow
// for the TrueLayer api.
//
// params
//   - body - the readcloser to decode
//
// returns
//   - the webhook request
//   - error if an error occurs
func (t *TrueLayer) HandleAsyncWebhookRequestBody(body io.ReadCloser) (*WebhookRequest, error) {
	req := &WebhookRequest{}

	err := json.NewDecoder(body).Decode(req)

	if err != nil {
		return nil, err
	}

	if req.Status == "Failed" {
		return req, &ErrorResponse{
			ErrorMessage:     req.Error,
			ErrorDescription: req.ErrorDescription,
		}
	}

	return req, nil
}
