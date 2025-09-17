-- Initialize SMS database schema

-- Create account table
CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(255),
    email VARCHAR(255),
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user'
);

-- Insert some sample data
INSERT INTO account (fullname, email, username, password, role) 
VALUES 
    ('Admin User', 'admin@example.com', 'admin1', '123', 'admin'),
    ('Test User', 'user@example.com', 'user1', 'password', 'user')
ON CONFLICT (username) DO NOTHING;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_account_username ON account(username);
CREATE INDEX IF NOT EXISTS idx_account_email ON account(email);
