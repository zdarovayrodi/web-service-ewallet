# Run

```bash
docker build -t go-rest-api . &&
```

```bash
docker run -d -p 8080:8080 go-rest-api
```

# Test

```bash
go test -v ./db/models/
```

# Usage

## Get all wallets

```bash
curl http://localhost:8080/api/v1/wallets
```

## Get wallet with id (uuid)

```bash
curl http://localhost:8080/api/v1/wallets/00000000-0000-0000-0000-000000000000
```

## Create new wallet

Balance will be set to 100.

```bash
curl -X 'POST' \
  'http://localhost:8080/api/v1/wallets' \
  -H 'Content-Type: application/json' \
  -d ''
```

## Create new wallet with balance

```bash
curl -X 'POST' \
  'http://localhost:8080/api/v1/wallets' \
  -H 'Content-Type: application/json' \
  -d '{"balance": 10000}'
```

## Transfer

````bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"to": "DESTINATION-WALLET-UUID", "amount": AMOUNT}' \
     http://localhost:8080/api/v1/wallets/{SOURCE-WALLET-UUID}/send```
````
