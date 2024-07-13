CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Wallets (
    wallet_id SERIAL PRIMARY KEY,
    user_id INT,
    balance DECIMAL(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Transactions (
    transaction_id SERIAL PRIMARY KEY,
    wallet_id INT REFERENCES Wallets(wallet_id),
    amount DECIMAL(10, 2),
    transaction_type VARCHAR(20),  -- e.g., 'credit', 'debit'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE GatewayLogs (
    log_id SERIAL PRIMARY KEY,
    request_path VARCHAR(255),
    request_method VARCHAR(10), -- e.g., 'GET', 'POST'
    status_code INT,
    response_time_ms INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
