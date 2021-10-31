package truelayer

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

const (
	authBaseURLSandbox = "https://auth.truelayer-sandbox.com"
	authBaseURL        = "https://auth.truelayer.com"
	authTokenEndpoint  = "/connect/token"
)

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
	u, err := buildURL(t.getAuthBaseURL(), "")

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
//   - err - any errors that have occurred including API errors
func (t *TrueLayer) GetAccessToken(code string, redirURI *url.URL) (token *AccessTokenResponse, err error) {
	body := t.getNewURLValuesWithClientInfo(true)
	body.Add("grant_type", "authorization_code")
	body.Add("redirect_uri", redirURI.String())
	body.Add("code", code)

	return t.authDoTokenRequest(body)
}

// RefreshAccessToken takes a refresh token and returns a refreshed access token
// allowing for continued use without user disruptuon.
//
// params
//   - refreshToken - user refresh token
//
// returns
//   - token - access token
//   - err - any errors that have occurred including API errors
func (t *TrueLayer) RefreshAccessToken(refreshToken string) (token *AccessTokenResponse, err error) {
	body := t.getNewURLValuesWithClientInfo(true)
	body.Add("grant_type", "refresh_token")
	body.Add("refresh_token", refreshToken)

	return t.authDoTokenRequest(body)
}

// authDoTokenRequest builds and executes authentication requests for the
// TrueLayer api.
//
// params
//   - body - request payload
//
// returns
//   - token - access token
//   - err - any errors that have occurred including API errors
func (t *TrueLayer) authDoTokenRequest(body url.Values) (token *AccessTokenResponse, err error) {
	u, err := buildURL(t.getAuthBaseURL(), authTokenEndpoint)

	if err != nil {
		return nil, err
	}

	res, err := t.doRequestWithFormURLEncodedBody(http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, parseErrorResponse(res)
	}

	token = &AccessTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(token)
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

// getAuthBaseURL parses the baseAuthURL for either the sandbox or non-sandbox
// TrueLayer environments and returns them. Using a utility method to reduce
// code duplication.
//
// returns
//   - the base url
func (t *TrueLayer) getAuthBaseURL() string {
	if t.sandbox {
		return authBaseURLSandbox
	}

	return authBaseURL
}
