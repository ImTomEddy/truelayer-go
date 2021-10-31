package truelayer

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type AccountsResponse struct {
	Results []Account `json:"results"`
}

type AccountBalanceResponse struct {
	Results []AccountBalance `json:"results"`
}

type AccountTransactionsResponse struct {
	Results []AccountTransaction `json:"results"`
}

type AccountStandingOrderResponse struct {
	Results []AccountStandingOrder `json:"results"`
}

type AccountDirectDebitResponse struct {
	Results []AccountDirectDebit `json:"results"`
}
