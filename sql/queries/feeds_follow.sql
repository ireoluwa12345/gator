-- name: CreateFeedFollow :one

WITH inserted_feed_follow AS (
    INSERT INTO feed_follow (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5) 
    RETURNING *
)

SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedsFollowForUser :many
SELECT 
    feed_follow.id,
    feed_follow.created_at,
    feed_follow.updated_at,
    feed_follow.user_id,
    feed_follow.feed_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follow
INNER JOIN feeds ON feed_follow.feed_id = feeds.id
INNER JOIN users ON feed_follow.user_id = users.id
WHERE feed_follow.user_id = $1;

-- name: GetFollowedFeedByID :one
SELECT feeds.id, feeds.name, feeds.url, feeds.user_id, feeds.created_at, feeds.updated_at, users.name AS user_name
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id
WHERE feeds.id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follow
WHERE feed_id = $1 AND user_id = $2;