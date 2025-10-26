-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, description, logo_url, priority)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsByPriority :many
SELECT * FROM feeds ORDER BY priority DESC, updated_at ASC;