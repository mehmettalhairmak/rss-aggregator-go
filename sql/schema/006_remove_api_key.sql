-- +goose Up
ALTER TABLE users 
DROP COLUMN api_key;

-- +goose Down
ALTER TABLE users 
ADD COLUMN api_key TEXT UNIQUE NOT NULL DEFAULT encode(sha256(random()::text::bytea), 'hex');
