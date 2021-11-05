package truelayer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

const (
	EndpointDataV1Accounts                   = "/data/v1/accounts"
	EndpointDataV1Account                    = "/data/v1/accounts/%s"
	EndpointDataV1AccountBalance             = "/data/v1/accounts/%s/balance"
	EndpointDataV1AccountTransactions        = "/data/v1/accounts/%s/transactions"
	EndpointDataV1AccountPendingTransactions = "/data/v1/accounts/%s/transactions/pending"
	EndpointDataV1AccountStandingOrders      = "/data/v1/accounts/%s/standing_orders"
	EndpointDataV1AccountDirectDebits        = "/data/v1/accounts/%s/direct_debits"

	ErrToFromNil = StrError("to/from must not be nil")
)

type Account struct {
	UpdateTimestamp time.Time `json:"update_timestamp"`
	AccountID       string    `json:"account_id"`
	AccountType     string    `json:"account_type"`
	DisplayName     string    `json:"display_name"`
	Currency        string    `json:"currency"`
	AccountNumber   struct {
		Iban     string `json:"iban"`
		Number   string `json:"number"`
		SortCode string `json:"sort_code"`
		SwiftBic string `json:"swift_bic"`
	} `json:"account_number"`
	Provider struct {
		ProviderID string `json:"provider_id"`
	} `json:"provider"`
}

type AccountBalance struct {
	Currency        string    `json:"currency"`
	Available       float64   `json:"available"`
	Current         float64   `json:"current"`
	Overdraft       float64   `json:"overdraft"`
	UpdateTimestamp time.Time `json:"update_timestamp"`
}

type AccountTransaction struct {
	TransactionID                   string   `json:"transaction_id"`
	NormalisedProviderTransactionID string   `json:"normalised_provider_transaction_id"`
	ProviderTransactionID           string   `json:"provider_transaction_id"`
	Timestamp                       string   `json:"timestamp"`
	Description                     string   `json:"description"`
	Amount                          float64  `json:"amount"`
	Currency                        string   `json:"currency"`
	TransactionType                 string   `json:"transaction_type"`
	TransactionCategory             string   `json:"transaction_category"`
	TransactionClassification       []string `json:"transaction_classification"`
	MerchantName                    string   `json:"merchant_name"`
	RunningBalance                  struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	} `json:"running_balance"`
	Meta struct {
		BankTransactionID           string `json:"bank_transaction_id"`
		ProviderTransactionCategory string `json:"provider_transaction_category"`
	} `json:"meta"`
}

type AccountStandingOrder struct {
	Frequency string    `json:"frequency"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Currency  string    `json:"currency"`
	Meta      struct {
		ProviderAccountID string `json:"provider_account_id"`
	} `json:"meta"`
	NextPaymentDate    time.Time `json:"next_payment_date"`
	NextPaymentAmount  float64   `json:"next_payment_amount"`
	FirstPaymentDate   time.Time `json:"first_payment_date"`
	FirstPaymentAmount float64   `json:"first_payment_amount"`
	FinalPaymentDate   time.Time `json:"final_payment_date"`
	FinalPaymentAmount float64   `json:"final_payment_amount"`
	Reference          string    `json:"reference"`
	Payee              string    `json:"payee"`
}

type AccountDirectDebit struct {
	DirectDebitID            string    `json:"direct_debit_id"`
	Timestamp                time.Time `json:"timestamp"`
	Name                     string    `json:"name"`
	Status                   string    `json:"status"`
	PreviousPaymentTimestamp time.Time `json:"previous_payment_timestamp"`
	PreviousPaymentAmount    float64   `json:"previous_payment_amount"`
	Currency                 string    `json:"currency"`
	Meta                     struct {
		ProviderMandateIdentification string `json:"provider_mandate_identification"`
		ProviderAccountID             string `json:"provider_account_id"`
	} `json:"meta"`
}

type AccountOptions struct {
	To   *time.Time
	From *time.Time
}

// GetAccounts retrieves the account associated with the provided access token.
//
// params
//   - accessToken - access token to get the accounts from
//
// returns
//   - list of accounts
//   - errors from the api request
func (t *TrueLayer) GetAccounts(accessToken string) ([]Account, error) {
	u, err := buildURL(t.getBaseURL(), EndpointDataV1Accounts)

	if err != nil {
		return nil, err
	}

	return t.getAccounts(u, accessToken)
}

// GetAccountsAsync triggers an async request to TrueLayer to get a list of all
// accounts associated to the access token.
//
// params
//   - accessToken - access token to get the info from
//   - webhookURI - uri to access upon async job completion
//
// returns
//   - truelayer response
//   - errors from the api request
func (t *TrueLayer) GetAccountsAsync(accessToken string, webhookURI string) (*AsyncRequestResponse, error) {
	return t.doAsyncAccountRequest(EndpointDataV1Accounts, accessToken, webhookURI, nil)
}

// GetAccountsAsyncRequest takes the result from a Webhook request and sends a
// request to the correct endpoint to fetch the Accounts.
//
// params
//   - accessToken - the access token associated to the webhook request
//   - webhook - the webhook request to fetch data from
func (t *TrueLayer) GetAccountsAsyncRequest(accessToken string, webhook *WebhookRequest) ([]Account, error) {
	u, err := buildURL(t.getBaseURL(), fmt.Sprintf(EndpointDataV1Results, webhook.TaskID))

	if err != nil {
		return nil, err
	}

	return t.getAccounts(u, accessToken)
}

// getAccounts takes the Account data URL then does an authenticated GET request
// decoding the response and returning the correct data structure.
//
// params
//   - u - the URL to request
//   - accessToken - the account's associated access token
func (t *TrueLayer) getAccounts(u *url.URL, accessToken string) ([]Account, error) {
	res, err := t.doAuthorizedGetRequest(u, accessToken)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, parseErrorResponse(res)
	}

	accountResp := AccountsResponse{}
	err = json.NewDecoder(res.Body).Decode(&accountResp)

	if err != nil {
		return nil, err
	}

	return accountResp.Results, nil
}

// GetAccount retrieves the specified account based on accountID, this account
// must be associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the account from
//   - accountID - the account ID to get
//
// returns
//   - the account
//   - errors from the api request
func (t *TrueLayer) GetAccount(accessToken string, accountID string) (*Account, error) {
	u, err := buildURL(t.getBaseURL(), fmt.Sprintf(EndpointDataV1Account, accountID))

	if err != nil {
		return nil, err
	}

	res, err := t.doAuthorizedGetRequest(u, accessToken)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, parseErrorResponse(res)
	}

	accountResp := AccountsResponse{}
	err = json.NewDecoder(res.Body).Decode(&accountResp)

	if err != nil {
		return nil, err
	}

	return &accountResp.Results[0], nil
}

// GetAccountAsync triggers an async request to TrueLayer to get the specified
// account based on the accountID, this account must be associated to the
// provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the info from
//   - webhookURI - uri to access upon async job completion
//   - accountID - id of the account to get
//
// returns
//   - truelayer response
//   - errors from the api request
func (t *TrueLayer) GetAccountAsync(accessToken string, webhookURI string, accountID string) (*AsyncRequestResponse, error) {
	return t.doAsyncAccountRequest(fmt.Sprintf(EndpointDataV1Account, accountID), accessToken, webhookURI, nil)
}

// GetAccountBalance retrieves the specified account's balance this account must
// be associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the account from
//   - accountID - the account ID to get
//
// returns
//   - the balance
//   - errors from the api request
func (t *TrueLayer) GetAccountBalance(accessToken string, accountID string) (*AccountBalance, error) {
	u, err := buildURL(t.getBaseURL(), fmt.Sprintf(EndpointDataV1AccountBalance, accountID))

	if err != nil {
		return nil, err
	}

	res, err := t.doAuthorizedGetRequest(u, accessToken)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, parseErrorResponse(res)
	}

	balanceResp := AccountBalanceResponse{}
	err = json.NewDecoder(res.Body).Decode(&balanceResp)

	if err != nil {
		return nil, err
	}

	return &balanceResp.Results[0], nil
}

// GetAccountBalanceAsync triggers an async request to TrueLayer to get the
// specified account balance based on the accountID, this account must be
// associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the info from
//   - webhookURI - uri to access upon async job completion
//   - accountID - id of the account to get
//
// returns
//   - truelayer response
//   - errors from the api request
func (t *TrueLayer) GetAccountBalanceAsync(accessToken string, webhookURI string, accountID string) (*AsyncRequestResponse, error) {
	return t.doAsyncAccountRequest(fmt.Sprintf(EndpointDataV1AccountBalance, accountID), accessToken, webhookURI, nil)
}

// GetAccountTransactions retrieves the specified account's transactions this
// account must be associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the account from
//   - accountID - the account ID to get
//   - opts - options for the request
//
// returns
//   - the transactions
//   - errors from the api request
func (t *TrueLayer) GetAccountTransactions(accessToken string, accountID string, opts *AccountOptions) ([]AccountTransaction, error) {
	u, err := buildURL(t.getBaseURL(), fmt.Sprintf(EndpointDataV1AccountTransactions, accountID))

	if err != nil {
		return nil, err
	}

	return t.getAccountTransactions(u, accessToken, accountID, opts)
}

// GetAccountTransactionsAsync triggers an async request to TrueLayer to get the
// specified account transactions based on the accountID, this account must be
// associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the info from
//   - webhookURI - uri to access upon async job completion
//   - accountID - id of the account to get
//   - opts - options for the request
//
// returns
//   - truelayer response
//   - errors from the api request
func (t *TrueLayer) GetAccountTransactionsAsync(accessToken string, webhookURI string, accountID string, opts *AccountOptions) (*AsyncRequestResponse, error) {
	return t.doAsyncAccountRequest(fmt.Sprintf(EndpointDataV1AccountTransactions, accountID), accessToken, webhookURI, opts)
}

// GetAccountPendingTransactions retrieves the specified account's pending
// transactions this account must be associated to the provided accessToken or
// an error will occur.
//
// params
//   - accessToken - access token to get the account from
//   - accountID - the account ID to get
//   - opts - options for the request
//
// returns
//   - the transactions
//   - errors from the api request
func (t *TrueLayer) GetAccountPendingTransactions(accessToken string, accountID string, opts *AccountOptions) ([]AccountTransaction, error) {
	u, err := buildURL(t.getBaseURL(), fmt.Sprintf(EndpointDataV1AccountPendingTransactions, accountID))

	if err != nil {
		return nil, err
	}

	return t.getAccountTransactions(u, accessToken, accountID, opts)
}

// GetAccountPendingTransactionsAsync triggers an async request to TrueLayer to
// get the specified account pending transactions based on the accountID, this
// account must be associated to the provided accessToken or an error will
// occur.
//
// params
//   - accessToken - access token to get the info from
//   - webhookURI - uri to access upon async job completion
//   - accountID - id of the account to get
//   - opts - options for the request
//
// returns
//   - truelayer response
//   - errors from the api request
func (t *TrueLayer) GetAccountPendingTransactionsAsync(accessToken string, webhookURI string, accountID string, opts *AccountOptions) (*AsyncRequestResponse, error) {
	return t.doAsyncAccountRequest(fmt.Sprintf(EndpointDataV1AccountPendingTransactions, accountID), accessToken, webhookURI, opts)
}

// getAccountTransactions retrieves the specified account's transactions either
// pending or not depending on the passed URL.
//
// params
//   - url - the url to request
//     (EndpointDataV1AccountTransactions|EndpointDataV1AccountPendingTransactions)
//   - accessToken - access token to get the account from
//   - accountID - the account ID to get
//   - opts - options for the request
//
// returns
//   - the transactions
//   - errors from the api request
func (t *TrueLayer) getAccountTransactions(url *url.URL, accessToken string, accountID string, opts *AccountOptions) ([]AccountTransaction, error) {
	if opts != nil {
		if opts.From == nil || opts.To == nil {
			return nil, ErrToFromNil
		}

		q := url.Query()
		q.Add("to", opts.To.Format(time.RFC3339))
		q.Add("from", opts.From.Format(time.RFC3339))
		url.RawQuery = q.Encode()
	}

	res, err := t.doAuthorizedGetRequest(url, accessToken)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, parseErrorResponse(res)
	}

	transactionsResp := AccountTransactionsResponse{}
	err = json.NewDecoder(res.Body).Decode(&transactionsResp)

	if err != nil {
		return nil, err
	}

	return transactionsResp.Results, nil
}

// GetAccountStandingOrders retrieves the specified account's standing orders
// this account must be associated to the provided accessToken or an error will
// occur.
//
// params
//   - accessToken - access token to get the account from
//   - accountID - the account ID to get
//
// returns
//   - the standing orders
//   - errors from the api request
func (t *TrueLayer) GetAccountStandingOrders(accessToken string, accountID string) ([]AccountStandingOrder, error) {
	u, err := buildURL(t.getBaseURL(), fmt.Sprintf(EndpointDataV1AccountStandingOrders, accountID))

	if err != nil {
		return nil, err
	}

	res, err := t.doAuthorizedGetRequest(u, accessToken)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, parseErrorResponse(res)
	}

	standingOrderResp := AccountStandingOrderResponse{}
	err = json.NewDecoder(res.Body).Decode(&standingOrderResp)

	if err != nil {
		return nil, err
	}

	return standingOrderResp.Results, nil
}

// GetAccountStandingOrdersAsync triggers an async request to TrueLayer to get
// the specified account's standing orders based on the accountID, this account
// must be associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the info from
//   - webhookURI - uri to access upon async job completion
//   - accountID - id of the account to get
//
// returns
//   - truelayer response
//   - errors from the api request
func (t *TrueLayer) GetAccountStandingOrdersAsync(accessToken string, webhookURI string, accountID string) (*AsyncRequestResponse, error) {
	return t.doAsyncAccountRequest(fmt.Sprintf(EndpointDataV1AccountStandingOrders, accountID), accessToken, webhookURI, nil)
}

// GetAccountDirectDebits retrieves the specified account's direct debits this
//account must be associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the account from
//   - accountID - the account ID to get
//
// returns
//   - the direct debits
//   - errors from the api request
func (t *TrueLayer) GetAccountDirectDebits(accessToken string, accountID string) ([]AccountDirectDebit, error) {
	u, err := buildURL(t.getBaseURL(), fmt.Sprintf(EndpointDataV1AccountDirectDebits, accountID))

	if err != nil {
		return nil, err
	}

	res, err := t.doAuthorizedGetRequest(u, accessToken)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, parseErrorResponse(res)
	}

	directDebitResp := AccountDirectDebitResponse{}
	err = json.NewDecoder(res.Body).Decode(&directDebitResp)

	if err != nil {
		return nil, err
	}

	return directDebitResp.Results, nil
}

// GetAccountDirectDebitsAsync triggers an async request to TrueLayer to get the
// specified account's direct debits based on the accountID, this account must
// be associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the info from
//   - webhookURI - uri to access upon async job completion
//   - accountID - id of the account to get
//
// returns
//   - truelayer response
//   - errors from the api request
func (t *TrueLayer) GetAccountDirectDebitsAsync(accessToken string, webhookURI string, accountID string) (*AsyncRequestResponse, error) {
	return t.doAsyncAccountRequest(fmt.Sprintf(EndpointDataV1AccountStandingOrders, accountID), accessToken, webhookURI, nil)
}

// doAsyncAccountRequest starts the process of getting account information
// through TrueLayer's async function.
//
// params
//   - endpoint - api endpoint to access
//   - accessToken - access token to get the info from
//   - webhookURI - uri to access upon async job completion (optional)
//   - opts - any request options
//
// returns
//   - truelayer response
//   - errors from the api request
func (t *TrueLayer) doAsyncAccountRequest(endpoint string, accessToken string, webhookURI string, opts *AccountOptions) (*AsyncRequestResponse, error) {
	u, err := buildURL(t.getBaseURL(), endpoint)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("async", "true")
	q.Add("webhookURI", webhookURI)

	if opts != nil {
		if opts.From == nil || opts.To == nil {
			return nil, ErrToFromNil
		}

		q.Add("to", opts.To.Format(time.RFC3339))
		q.Add("from", opts.From.Format(time.RFC3339))
	}

	u.RawQuery = q.Encode()

	res, err := t.doAuthorizedGetRequest(u, accessToken)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, parseErrorResponse(res)
	}

	resp := AsyncRequestResponse{}
	err = json.NewDecoder(res.Body).Decode(&resp)

	return &resp, err
}
