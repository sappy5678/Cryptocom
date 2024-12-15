BEGIN;
CREATE TABLE IF NOT EXISTS UserWallet (
    ID BIGSERIAL PRIMARY KEY,
    userID VARCHAR(36) UNIQUE NOT NULL,
    -- 1 dollor will be 10^6, so if we want to store 10 dollor, we need to store balance to 10*10^6
    -- this should be enough for now and need to handle with client side
    balance BIGINT NOT NULL DEFAULT 0 
    constraint balanceNonnegative check (balance >= 0)
);

CREATE TABLE IF NOT EXISTS UserWalletTransaction (
    ID BIGSERIAL PRIMARY KEY,
    userID VARCHAR(36) NOT NULL,
    transactionID VARCHAR(60) UNIQUE NOT NULL,
    operationType INT NOT NULL,
    passiveUserID VARCHAR(36),
    amount BIGINT NOT NULL,
    createdAt TIMESTAMP NOT NULL
);

CREATE INDEX idxUserWalletUserID ON UserWallet(userID);
CREATE INDEX idxUserWalletTransactionUserIDCreatedAtID ON UserWalletTransaction(userID, createdAt, ID);
CREATE INDEX idxUserWalletTransactionTransactionID ON UserWalletTransaction(transactionID);
COMMIT;
