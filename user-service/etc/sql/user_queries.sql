-- name: CreateUser
INSERT INTO users (name, email, age)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at;

-- name: FindUserByID
SELECT id, name, email, age, created_at, updated_at 
FROM users 
WHERE id = $1;