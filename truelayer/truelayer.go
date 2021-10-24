package truelayer

import (
	"encoding/json"
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

	q := t.getURLValuesWithClientInfo(u.Query(), false)
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

	body := t.getNewURLValuesWithClientInfo(true)
	body.Add("grant_type", "authorization_code")
	body.Add("redirect_uri", redirURI.String())
	body.Add("code", code)

	res, err := t.doRequestWithFormURLEncodedBody(http.MethodPost, u.String(), body)
	if err != nil {
		return token, err
	}

	defer res.Body.Close()

	token = &AccessTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(token)
	if err != nil {
		return token, err
	}

	return token, err
}

// RefreshAccessToken takes a refresh token and returns a refreshed access token
// allowing for continued use without user disruptuon.
//
// params
//   - refreshToken - user refresh token
//
// returns
//   - token - access token
//   - err - any errors that have occurred
func (t *TrueLayer) RefreshAccessToken(refreshToken string) (token *AccessTokenResponse, err error) {
	u, err := t.getBaseAuthURL()
	u.Path = "/connect/token"

	body := t.getNewURLValuesWithClientInfo(true)
	body.Add("grant_type", "refresh_token")
	body.Add("refresh_token", refreshToken)

	res, err := t.doRequestWithFormURLEncodedBody(http.MethodPost, u.String(), body)
	if err != nil {
		return token, err
	}

	defer res.Body.Close()

	token = &AccessTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(token)
	if err != nil {
		return token, err
	}

	return token, err
}

// doRequestWithFormURLEncodedBody creates a HTTP request object with the
// Content-Type header set to `application/x-www-form-urlencoded` which is used
// in authentication requests.
//
// params
//   - method - http request method
//   - url - url to request
//   - body - url.Values to be encoded for the request body
//
// returns
//   - response from http request
//   - any errors from creating the requests or executing the request
func (t *TrueLayer) doRequestWithFormURLEncodedBody(method, url string, body url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return t.httpClient.Do(req)
}

// getNewURLValuesWithClientInfo creates a new url.Values object and injects it
// with the TrueLayer client information.
//
// params
//   - withSecret - inject client secret
//
// returns
//   - a new url.Values with client info
func (t *TrueLayer) getNewURLValuesWithClientInfo(withSecret bool) url.Values {
	return t.getURLValuesWithClientInfo(url.Values{}, withSecret)
}

// getURLValuesWithClientInfo injects TrueLayer client information into the
// provided url.Values object.
//
// params
//   - values - existing url.Values object
//   - withSecret - inject client secret
//
// returns
//   - url.Values with client info
func (t *TrueLayer) getURLValuesWithClientInfo(values url.Values, withSecret bool) url.Values {
	values.Add("client_id", t.clientID)

	if withSecret {
		values.Add("client_secret", t.clientSecret)
	}

	return values
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
