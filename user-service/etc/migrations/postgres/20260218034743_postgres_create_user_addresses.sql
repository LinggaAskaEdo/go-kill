-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    address_type VARCHAR(20) CHECK (address_type IN ('shipping', 'billing', 'both')),
    street_address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_addresses_user_id ON user_addresses(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_addresses;
-- +goose StatementEnd
