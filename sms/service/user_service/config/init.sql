-- Initialize user service database

-- Create account table (singular to match Go code)
CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user'
);

-- Insert default admin user (password: admin123)
INSERT INTO account (fullname, email, username, password, role) 
VALUES ('Admin User', 'admin@example.com', 'admin', '123', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Also insert a simple test user
INSERT INTO account (fullname, email, username, password, role)
VALUES ('Test User', 'test@example.com', 'test', '123', 'user')
ON CONFLICT (username) DO NOTHING;