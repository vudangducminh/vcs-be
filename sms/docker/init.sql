-- Initialize SMS application database
CREATE DATABASE IF NOT EXISTS postgres;

-- Connect to the database
\c postgres;

-- Create account table
CREATE TABLE IF NOT EXISTS account (
    "ID" SERIAL PRIMARY KEY,
    "fullname" VARCHAR(255),
    "email" VARCHAR(255),
    "username" VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "role" VARCHAR(50) NOT NULL
);

-- Insert default admin user
INSERT INTO account ("fullname", "email", "username", "password", "role") 
VALUES ('Administrator', 'admin@example.com', 'admin1', '123', 'admin')
ON CONFLICT ("username") DO NOTHING;

-- Insert test user
INSERT INTO account ("fullname", "email", "username", "password", "role") 
VALUES ('Test User', 'test@example.com', 'test1', '123', 'user')
ON CONFLICT ("username") DO NOTHING;

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO vudangducminh;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO vudangducminh;
