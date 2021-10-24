package truelayer

import (
	"encoding/json"
	"fmt"
	"time"
)

const EndpointDataV1Accounts = "/data/v1/accounts"
const EndpointDataV1Account = "/data/v1/accounts/%s"
const EndpointDataV1AccountBalance = "/data/v1/accounts/%s/balance"
const EndpointDataV1AccountTransactions = "/data/v1/accounts/%s/transactions"
const EndpointDataV1AccountPendingTransactions = "/data/v1/accounts/%s/pending-transactions"
const EndpointDataV1AccountStandingOrders = "/data/v1/accounts/%s/standing-orders"
const EndpointDataV1AccountDirectDebits = "/data/v1/accounts/%s/direct-debits"

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

type Balance struct {
	Currency        string    `json:"currency"`
	Available       float64   `json:"available"`
	Current         float64   `json:"current"`
	Overdraft       int       `json:"overdraft"`
	UpdateTimestamp time.Time `json:"update_timestamp"`
}

type Transaction struct {
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
	} `json:"running_balance,omitempty"`
	Meta struct {
		BankTransactionID           string `json:"bank_transaction_id"`
		ProviderTransactionCategory string `json:"provider_transaction_category"`
	} `json:"meta"`
}

// GetAccounts retrieves the account associated with the provided access token.
//
// params
//   - accessToken - access token to get the accounts from
//
// returns
//   - list of accounts
//   - errors from the api request
func (t *TrueLayer) GetAccounts(accessToken string) (*[]Account, error) {
	u, err := buildURL(t.getBaseURL(), EndpointDataV1Accounts)

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

	return &accountResp.Results, nil
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
func (t *TrueLayer) GetAccountBalance(accessToken string, accountID string) (*Balance, error) {
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

// GetAccountTransactions retrieves the specified account's transactions this
// account must be associated to the provided accessToken or an error will occur.
//
// params
//   - accessToken - access token to get the account from
//   - accountID - the account ID to get
//
// returns
//   - the transactions
//   - errors from the api request
func (t *TrueLayer) GetAccountTransactions(accessToken string, accountID string) (*[]Transaction, error) {
	u, err := buildURL(t.getBaseURL(), fmt.Sprintf(EndpointDataV1AccountTransactions, accountID))

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

	transactionsResp := AccountTransactionsResponse{}
	err = json.NewDecoder(res.Body).Decode(&transactionsResp)

	if err != nil {
		return nil, err
	}

	return &transactionsResp.Results, nil
}
