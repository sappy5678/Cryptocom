BEGIN;
CREATE TABLE IF NOT EXISTS UserWallet (
    ID BIGSERIAL PRIMARY KEY,
    userID CHAR(36) UNIQUE NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    balance BIGINT NOT NULL DEFAULT 0 -- 1 dollor will be 10*6
    constraint balanceNonnegative check (balance >= 0)
);

CREATE TABLE IF NOT EXISTS UserWalletTransaction (
    ID BIGSERIAL PRIMARY KEY,
    userID CHAR(36) NOT NULL,
    operationType INT NOT NULL,
    passiveUserID CHAR(36),
    amount BIGINT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idxUserWalletUserID ON UserWallet(userID);
CREATE INDEX idxUserWalletTransactionUserID ON UserWalletTransaction(userID);
COMMIT;
