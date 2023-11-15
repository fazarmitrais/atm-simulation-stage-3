# atm-simulation-stage-3

This is the backend side of ATM simulation.
The service can be accessed using Controller API.

## How to run the app
Open terminal or command prompt, then go to the app root directory.
Run this command : go run main.go

## endpoints' CURL examples
### PIN validation (login)
curl --location 'http://localhost:8080/api/v1/account/validate' \
--header 'Content-Type: application/json' \
--data '{
    "accountNumber": "112244",
    "pin": "932012"
}'

### Balance check
curl --location 'http://localhost:8080/api/v1/account/balance' \

### Withdraw
curl --location 'http://localhost:8080/api/v1/account/withdraw' \
--header 'Content-Type: application/json' \
--data '{
    "amount": 20.02
}'

### Transfer
curl --location 'http://localhost:8080/api/v1/account/transfer' \
--header 'Content-Type: application/json' \
--data '{
    "toAccountNumber": "112244",
    "referenceNumber": "123",
    "amount": 20
}'

### Exit (logout)
curl --location 'http://localhost:8080/api/v1/account/exit' \
--header 'Content-Type: application/json' \

