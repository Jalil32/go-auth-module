-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,							-- Auto-incremented unique user ID
    email VARCHAR(255) UNIQUE NOT NULL,				-- Email, unique and not null
    password_hash TEXT,					-- Hashed password
    first_name VARCHAR(100) NOT NULL,						-- Optional first name
    last_name VARCHAR(100) NOT NULL,							-- Optional last name
    is_active BOOLEAN DEFAULT TRUE,					-- Account status flag (default is true)
	verified BOOLEAN DEFAULT FALSE,
	provider TEXT,
 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Auto-generated timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP	-- Auto-generated timestamp
);

-- Create an index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);

-- Create function to automatically update udpated at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Set the updated_at column to the current timestamp
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
