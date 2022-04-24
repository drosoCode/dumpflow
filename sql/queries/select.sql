----------- POSTS

-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: ListAnswers :many
SELECT * FROM posts WHERE parent_id = $1;

-- name: ListHistoryFromPost :many
SELECT * FROM post_history WHERE post_id = $1;

-- name: ListRelatedPosts :many
SELECT * FROM post_links WHERE post_id = $1 LIMIT $2;

-- name: ListVotesFromPost :many
SELECT * FROM votes WHERE post_id = $1;

-- name: ListCommentsFromPost :many
SELECT * FROM comments WHERE post_id = $1;

-- name: ListUsersFromPost :many
SELECT * FROM users WHERE id IN (SELECT user_id FROM post_history WHERE post_id = $1 ORDER BY creation_date);


----------- COMMENTS

-- name: GetComment :one
SELECT * FROM comments WHERE id = $1;


----------- USERS

-- name: GetUser :one
SELECT * FROM users WHERE account_id = $1;

-- name: ListBadgesFromUser :many
SELECT * FROM badges WHERE user_id = $1;

----------- TAGS

-- name: ListTags :many
SELECT * FROM tags ORDER BY tag_name;

-- name: GetTag :many
SELECT * FROM tags WHERE id = $1;

-- name: GetTagFromName :many
SELECT * FROM tags WHERE tag_name = $1;

----------- BADGES

-- name: ListBadges :many
SELECT * FROM badges;

-- name: GetBadge :one
SELECT * FROM badges WHERE id = $1;