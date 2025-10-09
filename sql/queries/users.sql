-- name: CreateUser :one
-- Açıklama: Email ve şifre ile yeni kullanıcı kayıt etme (Sadece JWT sistemi)
-- API Key kullanmıyoruz artık - sadece modern JWT authentication
INSERT INTO users (id, created_at, updated_at, name, email, password_hash)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByEmail :one
-- Açıklama: Email ile kullanıcı bulma (login için gerekli)
-- Login yaparken kullanıcının email'ini alıp şifresini kontrol edeceğiz
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
-- Açıklama: ID ile kullanıcı bulma (JWT token'dan user_id alınca gerekli)
SELECT * FROM users WHERE id = $1;