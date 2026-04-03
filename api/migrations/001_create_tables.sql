CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id VARCHAR(255) PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    description VARCHAR(255) NOT NULL,
    category VARCHAR(255),
    deposit DECIMAL(10, 2) NOT NULL,
    withdrawal DECIMAL(10, 2) NOT NULL,
    bank_account_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (account_id) REFERENCES users(account_id) ON DELETE CASCADE
);
