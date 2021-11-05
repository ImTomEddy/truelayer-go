# truelayer-go
An unofficial TrueLayer SDK for Go

## Contents
- [truelayer-go](#truelayer-go)
  - [Contents](#contents)
  - [Usage](#usage)
  - [Supported Providers](#supported-providers)
  - [Supported Features](#supported-features)

This SDK only has plans to support the TrueLayer Data API

## Usage
An example of how to use the truelayer-go SDK can be found in the 
[examples](examples/) directory.

The README will be updated to include more in-depth usage instructions soon.

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