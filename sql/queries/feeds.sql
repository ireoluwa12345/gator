-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetAllFeeds :many
SELECT feeds.id, feeds.name, feeds.url, feeds.user_id, feeds.created_at, feeds.updated_at, users.name AS user_name
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByName :one
SELECT feeds.id, feeds.name, feeds.url, feeds.user_id, feeds.created_at, feeds.updated_at, users.name AS user_name
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id
WHERE feeds.name = $1;

-- name: GetFeedByUrl :one
SELECT feeds.id, feeds.name, feeds.url, feeds.user_id, feeds.created_at, feeds.updated_at, users.name AS user_name
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id
WHERE feeds.url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT feeds.id, feeds.name, feeds.url, feeds.user_id, feeds.created_at, feeds.updated_at, users.name AS user_name
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id
ORDER BY feeds.last_fetched_at NULLS FIRST, feeds.created_at
LIMIT 1;