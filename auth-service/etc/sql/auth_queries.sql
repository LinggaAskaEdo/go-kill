-- name: CheckUser
SELECT EXISTS(SELECT 1 FROM users_auth WHERE email = $1);

-- name: SaveUSer
INSERT INTO users_auth (email, password_hash, is_active, created_at, updated_at)
VALUES ($1, $2, true, NOW(), NOW()) 
RETURNING id;

-- name: GetUserByEmail
SELECT id, email, password_hash, is_active 
FROM users_auth 
WHERE email = $1;

-- name: StoreRefreshToken
INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) 
VALUES ($1, $2, $3, NOW());

-- name: GetUserWithID
SELECT email, is_active 
FROM users_auth 
WHERE id = $1;

-- name: DeleteRefreshToken
DELETE FROM refresh_tokens 
WHERE user_id = $1;