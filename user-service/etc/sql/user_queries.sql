-- name: RegisterUser
INSERT INTO users (auth_id, email, first_name, last_name, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id;

-- name: RegisterUserProfile
INSERT INTO user_profiles (user_id, created_at, updated_at)
VALUES ($1, NOW(), NOW());

-- name: GetUserByAuthID
SELECT id, email, first_name, last_name
FROM users
WHERE auth_id = $1;

-- name: GetUserByID
SELECT id, email, first_name, last_name
FROM users
WHERE id = $1;

-- name: GetUserAddresses
SELECT id, address_type, street_address, city, state, postal_code, country, is_default
FROM user_addresses
WHERE user_id = $1;

-- name: CreateUserAddress
INSERT INTO user_addresses (user_id, address_type, street_address, city, state, postal_code, country, is_default)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;
