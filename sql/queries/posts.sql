-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, feed_id, published_at, description)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: FetchUserPosts :many
SELECT posts.title, posts.url FROM posts
JOIN feed_follow ON posts.feed_id = feed_follow.feed_id
WHERE feed_follow.user_id = $1;