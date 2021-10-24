package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/ImTomEddy/truelayer-go/truelayer"
	"github.com/ImTomEddy/truelayer-go/truelayer/providers"
	env "github.com/Netflix/go-env"
)

type Config struct {
	Host         string `env:"HOST,default=http://localhost:3000"`
	RedirectPath string `env:"REDIRECT_PATH,default=/callback"`
	ClientID     string `env:"TRUELAYER_CLIENT_ID,required=true"`
	ClientSecret string `env:"TRUELAYER_CLIENT_SECRET,required=true"`
	Sandbox      bool   `env:"TRUELAYER_SANDBOX,default=true"`
}

func main() {
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)

	if err != nil {
		log.Fatal(err)
	}

	redirectURL, err := url.Parse(config.Host)
	if err != nil {
		log.Fatal(err)
	}
	redirectURL.Path = config.RedirectPath

	t := truelayer.New(config.ClientID, config.ClientSecret, config.Sandbox)
	link, _ := t.GetAuthenticationLink([]string{providers.UKMock, providers.UKOAuthAll, providers.UKOpenBankingAll}, []string{truelayer.PermissionAll}, redirectURL, false)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, link, http.StatusMovedPermanently)
	})

	http.HandleFunc(config.RedirectPath, handleCallback(t, redirectURL))
	http.ListenAndServe(":"+redirectURL.Port(), nil)
}

func handleCallback(t *truelayer.TrueLayer, redirectURL *url.URL) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			log.Println("No code")
			return
		}

		token, err := t.GetAccessToken(code, redirectURL)
		if err != nil {
			log.Println(err)
			return
		}

		b, err := json.Marshal(token)
		if err != nil {
			return
		}

		rw.Write(b)
	}
}
