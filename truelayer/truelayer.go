package truelayer

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const baseAuthURLSandbox = "https://auth.truelayer-sandbox.com"
const baseAuthURL = "https://auth.truelayer.com"

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

// AccessTokenResponse is the JSON Structure returned when requesting an
// AccessToken from TrueLayer.
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
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

// GetAuthenticationLink generates a link that can be used to authenticate
// against multiple providers with a specific scope of permissions.
//
// params
//   - providers - the allowed authentication providers
//   - permissions - the scope of permissions you want
//   - redirectURI - where to redirect the request to
//   - postCode - submit the code using `POST` over `GET`
//
// returns
//   - link - the authentication link
//   - err - error parsing the base URL - should not occur
func (t *TrueLayer) GetAuthenticationLink(providers []string, permissions []string, redirURI *url.URL, postCode bool) (link string, err error) {
	u, err := t.getBaseAuthURL()

	if err != nil {
		return link, err
	}

	q := u.Query()
	q.Add("client_id", t.clientID)
	q.Add("response_type", "code")

	q.Add("scope", strings.Join(permissions, " "))
	q.Add("providers", strings.Join(providers, " "))

	q.Add("redirect_uri", redirURI.String())

	if postCode {
		q.Add("response_mode", "form_post")
	}

	u.RawQuery = q.Encode()

	return u.String(), err
}

// GetAccessToken contacts the TrueLayer API and gets an access token to allow
// for authenticated requests to the TrueLayer Data API.
//
// params
//   - code - authentication code retrieved from the user
//
// returns
//   - token - access token
//   - err - any errors that have occurred
func (t *TrueLayer) GetAccessToken(code string, redirURI *url.URL) (token *AccessTokenResponse, err error) {
	u, err := t.getBaseAuthURL()
	u.Path = "/connect/token"

	body := url.Values{}
	body.Add("grant_type", "authorization_code")
	body.Add("client_id", t.clientID)
	body.Add("client_secret", t.clientSecret)
	body.Add("redirect_uri", redirURI.String())
	body.Add("code", code)

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(body.Encode()))
	if err != nil {
		return token, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := t.httpClient.Do(req)
	if err != nil {
		return token, err
	}

	defer res.Body.Close()

	token = &AccessTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(token)
	if err != nil {
		log.Println("error")
		return token, err
	}

	return token, err
}

// getBaseAuthURL parses the baseAuthURL for either the sandbox or non-sandbox
// TrueLayer environments and returns them. Using a utility method to reduce
// code duplication.
//
// returns
//   - the parsed url
//   - url parsing errors - should not occur as these are hard-coded values
func (t *TrueLayer) getBaseAuthURL() (*url.URL, error) {
	if t.sandbox {
		return url.Parse(baseAuthURLSandbox)
	}

	return url.Parse(baseAuthURL)
}
