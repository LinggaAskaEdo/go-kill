-- name: RegisterUser
INSERT INTO users (auth_id, email, first_name, last_name, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW()) 
RETURNING id;

-- name: RegisterUserProfile
INSERT INTO user_profiles (user_id, created_at, updated_at)
VALUES ($1, NOW(), NOW());