package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/ImTomEddy/truelayer-go/truelayer"
	"github.com/ImTomEddy/truelayer-go/truelayer/providers"
	env "github.com/Netflix/go-env"
)

type Config struct {
	Host         string `env:"HOST,default=http://localhost:3000"`
	RedirectPath string `env:"REDIRECT_PATH,default=/"`
	ClientID     string `env:"TRUELAYER_CLIENT_ID,required=true"`
	ClientSecret string `env:"TRUELAYER_CLIENT_SECRET,required=true"`
	Sandbox      bool   `env:"TRUELAYER_SANDBOX,default=true"`
}

type TemplateData struct {
	AccountID      string
	Accounts       []truelayer.Account
	Balance        *truelayer.AccountBalance
	Transactions   []truelayer.AccountTransaction
	StandingOrders []truelayer.AccountStandingOrder
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
	link, _ := t.GetAuthenticationLink([]string{providers.UKMock, providers.UKOAuthAll, providers.UKOpenBankingAll}, []string{truelayer.PermissionAll}, redirectURL, true)

	http.HandleFunc("/", handle(t, link, redirectURL))

	http.ListenAndServe(":"+redirectURL.Port(), nil)
}

func handle(t *truelayer.TrueLayer, redirectURL string, callbackURL *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Recieved Request")
		if r.Method != http.MethodPost {
			log.Println("Not Post - Redirecting")
			http.Redirect(rw, r, redirectURL, http.StatusMovedPermanently)
			return
		}

		log.Println("Post - Reading Body")
		if r.Body == nil {
			log.Println("no body")
			rw.Write([]byte("no request body"))
			return
		}

		log.Println("Parsing Form")
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
			rw.Write([]byte(err.Error()))
			return
		}

		log.Println("Getting Code")
		code := r.PostForm.Get("code")
		if code == "" {
			log.Println("code is empty")
			rw.Write([]byte("code is empty"))
		}

		log.Println("Getting Access Token")
		token, err := t.GetAccessToken(code, callbackURL)
		if err != nil {
			log.Println(err.Error())
			rw.Write([]byte(err.Error()))
			return
		}

		log.Println("Getting Accounts")
		accounts, err := t.GetAccounts(token.AccessToken)
		if err != nil {
			log.Println(err.Error())
			rw.Write([]byte(err.Error()))
			return
		}

		log.Println("Getting Balance")
		balance, err := t.GetAccountBalance(token.AccessToken, accounts[0].AccountID)
		if err != nil {
			log.Println(err.Error())
			rw.Write([]byte(err.Error()))
			return
		}

		log.Println("Getting Transactions")
		transactions, err := t.GetAccountTransactions(token.AccessToken, accounts[0].AccountID)
		if err != nil {
			log.Println(err.Error())
			rw.Write([]byte(err.Error()))
			return
		}

		log.Println("Generating HTML")
		data := TemplateData{
			AccountID:    accounts[0].AccountID,
			Accounts:     accounts,
			Balance:      balance,
			Transactions: transactions[:10],
		}

		temp := template.Must(template.ParseFiles("examples/template.html"))
		temp.Execute(rw, data)
	}
}
