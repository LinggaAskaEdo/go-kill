-- name: CheckUser
SELECT EXISTS(SELECT 1 FROM users_auth WHERE email = $1);

-- name: SaveUSer
INSERT INTO users_auth (email, password_hash, is_active, created_at, updated_at)
VALUES ($1, $2, true, NOW(), NOW()) 
RETURNING id;