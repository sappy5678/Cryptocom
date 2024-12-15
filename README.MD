
# How to run

```
docker compose up --build -d
```
and the server will run on port 8080


# DB Design
## User Wallet
1. Single Wallet per User
   - Simplified design for demonstration purposes
2. Authentication Omitted
   - Focus on core wallet functionality
   - Authentication can be added later as a separate service
   - Allows for easier testing and development
3. Integer-based Balance
   - Uses bigint (10^6 = 1 dollar) to avoid floating point precision issues
   - Ensures accurate calculations for all transactions
   - Use constraint to ensure balance is non-negative

## User Wallet Transaction
1. Transaction Types
   - Clear distinction between different operations:
     - Type 1: Deposit 
     - Type 2: Withdrawal 
     - Type 3: TransferIn
     - Type 4: TransferOut
   - Write TransferIn and TransferOut at the same time
     - Simplifies transaction history queries for specific user
     - Enables straightforward reporting and analytics for specific user
2. Passive Wallet ID Design
   - Clearly identifies the passive user in transfer transactions
   - Enables bi-directional transaction tracking
   - Simplifies transaction history queries
3. API Design Choices
   - Transaction IDs for idempotency

## Index Design
1. Index: UserWallet(userID)
   - Optimizes wallet lookups by userID
   - Ensures uniqueness constraint for one wallet per user
2. Composite Index: UserWalletTransaction(userID, createdAt, ID)
   - Efficiently supports transaction history queries
   - Enables pagination with createdAt or ID
   - Optimizes filtering by time range for specific user
3. Secondary Index: UserWalletTransaction(transactionID)
   - Supports fast transaction ID lookups
   - Helps enforce idempotency by checking existing transactions
   - Enables quick transaction status verification

# API Design
## Key Features
1. Idempotency
   - Any transaction need to get transactionID from server first
   - All write operations require transactionID for idempotency
   - Prevents duplicate transactions
   - Safe for retry operations
2. Pagination and Limit
   - Cursor-based pagination using createdBefore and IDBefore
   - Efficient for large transaction histories
   - Consistent ordering by IDBefore
3. Error Handling
   - Standardized error responses
   - Clear validation messages
   - Proper HTTP status codes
4. Transaction Atomicity
   - Balance update and transaction record are atomic
   - Consistent wallet balances
   - No lost or duplicate transactions
5. Simplify API Design
   - Use POST for all write operations
   - Use GET for read operations

## API Design
1. Create Wallet
   - PUT /api/v1/users/{userID}/wallet
   - Creates a new wallet for specified user
   - Returns wallet details with initial balance of 0
   - If wallet already exists, return the existing wallet
2. Get Wallet
   - GET /api/v1/users/{userID}/wallet  
   - Retrieves wallet information for specified user
   - Returns wallet balance and details
   - If wallet not found, return error
3. Create Transaction ID
   - POST /api/v1/users/{userID}/wallet/transactionID
   - Generates unique transaction ID for subsequent operations
   - Response:
     ```json
     {
       "transactionID": "unique-transaction-id"
     }
     ```
   - Must obtain transactionID before any write operation
   - Each transactionID can only be used once
   - Ensures idempotency and prevents duplicate transactions
4. Deposit
   - PUT /api/v1/users/{userID}/wallet/deposit
   - Deposits funds into user's wallet
   - Request body:
     ```json
     {
       "amount": 1000000,
       "transactionID": "unique-transaction-id"
     }
     ```
   - Amount is in base units (10^6 = 1 dollar)
   - TransactionID ensures idempotency, so it is safe to retry
   - Using PUT instead of POST to let client know it is a idempotent operation
5. Withdraw
   - PUT /api/v1/users/{userID}/wallet/withdraw
   - Withdraws funds from user's wallet
   - Request body:
     ```json
     {
       "amount": 1000000,
       "transactionID": "unique-transaction-id"
     }
     ```
   - TransactionID ensures idempotency, so it is safe to retry
   - Using PUT instead of POST to let client know it is a idempotent operation
   - If balance is insufficient, return error
6. Transfer
   - PUT /api/v1/users/{userID}/wallet/transfer
   - Transfers funds between wallets
   - Request body:
     ```json
     {
       "passiveUserID": "recipient-user-id",
       "amount": 1000000,
       "transactionID": "unique-transaction-id"
     }
     ```
   - TransactionID ensures idempotency, so it is safe to retry
   - Using PUT instead of POST to let client know it is a idempotent operation
   - Creates paired TransferOut/TransferIn transactions at the same time for quick query
   - Error handling:
     - If passiveUserID is not found, return error
     - If passiveUserID is the same as userID, return error
     - If balance is insufficient, return error
7. Get Transaction History
   - GET /api/v1/users/{userID}/wallet/transactions
   - Lists wallet transactions with pagination(createdBefore or IDBefore) and limit
   - Use pagination for efficient large data retrieval
   - User can query with time
   - Query parameters:
     - createdBefore: RFC3339 timestamp
     - IDBefore: Transaction ID for pagination
     - limit: Max number of records (default 100)
   - Returns transactions sorted by creation time

## Postman Collection
[Postman Collection](./cryptocom.postman_collection.json)

# Project Structure
* cmd/ => any executable code, like server main.go or any other executable tools
* deploy/ => any deploy code, like db migration or any other deploy tools
* pkg/ => main package code, like domain, service, utl, etc.
* pkg/utl/ => utility package code, like postgres, zlog, etc.
* pkg/domain/ => interface code and data model code, used for interaction between service
* pkg/service/ => business logic code, like wallet service etc.
* pkg/service/wallet/ => wallet service code, wired database and repository to provide business logic.
* pkg/service/wallet/repository/ => database repository code, used for database operation.
* pkg/service/wallet/transport/ => http transport code, used for wrap http request and response for service.
* pkg/service/wallet/logging/ => logging code, used for wrap logger for http service.


# How to test
```
go test ./... -cover -p=1

```
Why we use -p=1?
Since we use embedded postgres, it will start a new instance for each test, so we need to use -p=1 to avoid file io issue.
It is better to use mock database for each test. But I don't have enough time for this demo.

I believe goleak is not working well with sqlx/db sql/db
since they maintain their own connection pool, and cannot be closed by our code
and goleak willdetect the connection pool, and report the error

## Result
most of the test are covered, but some of the test are not covered because of the mock interface.

```
❯ go test ./... -cover -p=1 -count=1
        github.com/sappy5678/cryptocom/cmd/api          coverage: 0.0% of statements
ok      github.com/sappy5678/cryptocom/pkg/domain       0.005s  coverage: 100.0% of statements
        github.com/sappy5678/cryptocom/pkg/service              coverage: 0.0% of statements
ok      github.com/sappy5678/cryptocom/pkg/service/wallet       0.011s  coverage: 80.0% of statements
ok      github.com/sappy5678/cryptocom/pkg/service/wallet/logging       0.010s  coverage: 100.0% of statements
ok      github.com/sappy5678/cryptocom/pkg/service/wallet/repository    8.947s  coverage: 74.8% of statements
ok      github.com/sappy5678/cryptocom/pkg/service/wallet/transport     0.017s  coverage: 75.2% of statements
ok      github.com/sappy5678/cryptocom/pkg/utl/config   0.011s  coverage: 100.0% of statements
ok      github.com/sappy5678/cryptocom/pkg/utl/postgres 3.624s  coverage: 83.3% of statements
ok      github.com/sappy5678/cryptocom/pkg/utl/server   0.007s  coverage: 27.8% of statements
ok      github.com/sappy5678/cryptocom/pkg/utl/zlog     0.006s  coverage: 100.0% of statements
```

# Total Time Used
```
25hr
```
