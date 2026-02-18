-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_profiles (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID UNIQUE NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    bio TEXT,
    avatar_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles;
-- +goose StatementEnd
