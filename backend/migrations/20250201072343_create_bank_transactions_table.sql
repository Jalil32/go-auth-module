-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bank_transactions (
    transaction_id SERIAL PRIMARY KEY,							-- Auto incrementing transaction ID
    user_id INT NOT NULL,						-- Foreign key to the user who made the transaction
    date TIMESTAMP NOT NULL,				-- Date of transaction
    amount_cents INT NOT NULL,					-- Transaction amount in cents
    description TEXT NOT NULL,						-- Transaction description: Default is "No description"
    CONSTRAINT fk_user_id 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE -- Foreign key constraint: When a user is deleted, delete all their transactions
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bank_transactions;
-- +goose StatementEnd
