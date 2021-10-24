package truelayer

import (
	"net/url"
	"strings"
)

const baseURLSandbox = "https://auth.truelayer-sandbox.com/"
const baseURL = "https://auth.truelayer.com/"

type TrueLayer struct {
	clientID     string
	clientSecret string
	sandbox      bool
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
	return &TrueLayer{
		clientID:     clientID,
		clientSecret: clientSecret,
		sandbox:      sandbox,
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
	var u *url.URL

	if t.sandbox {
		u, err = url.Parse(baseURLSandbox)
	} else {
		u, err = url.Parse(baseURL)
	}

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
