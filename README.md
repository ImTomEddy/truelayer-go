# truelayer-go
An unofficial TrueLayer SDK for Go

## Contents
- [truelayer-go](#truelayer-go)
  - [Contents](#contents)
  - [Usage](#usage)
    - [Synchronous](#synchronous)
    - [Asynchronous](#asynchronous)
  - [Supported Providers](#supported-providers)
  - [Supported Features](#supported-features)

This SDK only has plans to support the TrueLayer Data API

## Usage
The README will be updated to include more in-depth usage instructions soon.

### Synchronous
A Synchronous example of the truelayer-go SDK can be found within the
[examples](examples/) directory. To run this example you will need to export two
environment variables `TRUELAYER_CLIENT_ID` and `TRUELAYER_CLIENT_SECRET`. These
variables should be the sandbox mode variables - the example wasn't built using 
the live mode.

With the example application running you can go to `http://localhost:3000` which
will redirect you to the TrueLayer authentication page. Once authenticated 
(user: `john`, password: `doe`) you will be redirected back to localhost using a
`POST` request.

### Asynchronous
TrueLayer recommends using the asynchronous API over the synchronous API this is
because it can help mitigate issues that are beyond TrueLayer's control. It also
means that the user has to impliment less retry logic.

Read more at https://docs.truelayer.com/docs/asynchronous-calls-and-webhooks

The current working asynchronous implementation is using webhooks. There is also
the ability to poll an asynchronous request but this is yet to be implemented.


## Supported Providers
truelayer-go doesn't inherently limit the providers that can be used however, 
the SDK does provide hard-coded provider values to make it easier to manage.

The hard-coded providers can be found under 
[truelayer/providers](truelayer/providers/).

## Supported Features
- [x] Authentication
  - [x] Access Token
  - [x] Refresh Token
- [ ] Data API
  - [ ] Accounts
    - [ ] Async Support
      - [ ] Webhook
      - [ ] Polling
    - [ ] Correlation ID
    - [ ] PSU-IP
    - [x] Routes
      - [x] Get Accounts
      - [x] Get Account
      - [x] Get Account Balance
      - [x] Get Account Transactions
      - [x] Get Account Pending Transactions
      - [x] Get Standing Orders
      - [x] Get Direct Debits
  - [ ] Batch
    - [ ] Initiate Batch Call
    - [ ] Get Data In Batch
  - [ ] Cards
    - [ ] Get Cards
    - [ ] Get Card
    - [ ] Get Card Balance
    - [ ] Get Card Transactions
    - [ ] Get Card Pending Transactions
  - [ ] Meta
    - [ ] Get Connection Metadata
    - [ ] Get User Info