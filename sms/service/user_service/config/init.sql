-- Initialize user service database

-- Create account table (singular to match Go code)
CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(100) DEFAULT ''
);

CREATE TABLE IF NOT EXISTS email_manager (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL
);
