# How to run

## Run via binary
(only work on `darwin/amd64` arch)
1. Locate your terminal to this repository
2. Perform `chmod +x main.go` (in some cases)
3. Perform `./main`

## Run via GO
1. Locate your terminal to this repository
2. Perform `go run app/main.go`

## Example Requests

### Init wallet
```
curl --location 'http://localhost:8000/api/v1/init' \
--form 'customer_xid="ea0212d3-abd6-406f-8c67-868e814a2436"'
```

### Enable wallet
```
curl --location --request POST 'http://localhost:8000/api/v1/wallet' \
--header 'Authorization: Token 9626dea978df941e8f7e72a2561b6bc66d15e125'
```

### View wallet balance
```
curl --location 'http://localhost:8000/api/v1/wallet' \
--header 'Authorization: Token 9626dea978df941e8f7e72a2561b6bc66d15e125'
```

### View wallet mutations
```
curl --location 'http://localhost:8000/api/v1/wallet/transactions' \
--header 'Authorization: Token 9626dea978df941e8f7e72a2561b6bc66d15e125'
```

### Deposit
```
curl --location 'http://localhost:8000/api/v1/wallet/deposits' \
--header 'Authorization: Token 9626dea978df941e8f7e72a2561b6bc66d15e125' \
--form 'amount="100000"' \
--form 'reference_id="50535246-dcb2-4929-8cc9-004ea06f5249"'
```

### Withdraw
```
curl --location 'http://localhost:8000/api/v1/wallet/withdrawals' \
--header 'Authorization: Token 9626dea978df941e8f7e72a2561b6bc66d15e125' \
--form 'amount="60000"' \
--form 'reference_id="4b01c9bb-3acd-47dc-87db-d9ac483d20b2"'
```

### Disable wallet
```
curl --location --request PATCH 'http://localhost:8000/api/v1/wallet' \
--header 'Authorization: Token 9626dea978df941e8f7e72a2561b6bc66d15e125' \
--form 'is_disabled="true"'
```