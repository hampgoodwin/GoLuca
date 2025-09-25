Just some grpcurl commands to get going


Accounts

```sh
grpcurl -protoset <(buf build proto -o -) -plaintext -d '{"name": "owners equity", "type": "ACCOUNT_TYPE_EQUITY", "basis": "BASIS_CREDIT"}' localhost:8080 goluca.account.v1.AccountService/CreateAccount

grpcurl -protoset <(buf build proto -o -) -plaintext -d '{"name": "cash", "type": "ACCOUNT_TYPE_ASSET", "basis": "BASIS_DEBIT"}' localhost:8080 goluca.account.v1.AccountService/CreateAccount

grpcurl -protoset <(buf build proto -o -) -plaintext localhost:8080 goluca.account.v1.AccountService/ListAccounts
```

Transactions

```sh
grpcurl -protoset <(buf build proto -o -) -plaintext -d '{"description": "fund", "entries": {"description": "fund", "debitAccount": "accountId", "creditAccount": "accountId", "amount":{"value": 100000, "currency": "usd"}}}' localhost:8080 goluca.transaction.v1.TransactionService/CreateTransaction

grpcurl -protoset <(buf build proto -o -) -plaintext localhost:8080 goluca.transaction.v1.TransactionService/ListTransactions
```
