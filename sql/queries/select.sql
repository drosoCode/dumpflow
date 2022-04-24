----------- POSTS

-- name: GetPost :one
SELECT id, post_type_id, parent_id, accepted_answer_id, creation_date, closed_date, score, view_count, body, tags, answer_count, comment_count, favorite_count, content_license
FROM posts WHERE id = $1;

-- name: ListAnswers :many
SELECT id, post_type_id, parent_id, accepted_answer_id, creation_date, closed_date, score, view_count, body, tags, answer_count, comment_count, favorite_count, content_license
FROM posts WHERE parent_id = $1;

-- name: ListHistoryFromPost :many
SELECT id, post_history_type_id, post_id, revision_guid, creation_date, user_id,comment, text, content_license
FROM post_history WHERE post_id = $1;

-- name: ListRelatedPosts :many
SELECT * FROM post_links WHERE post_id = $1 LIMIT $2;

-- name: GetVotesFromPost :many
SELECT COUNT(*) AS votes, vote_type_id FROM votes WHERE post_id = $1 GROUP BY vote_type_id;

-- name: ListCommentsFromPost :many
SELECT id, post_id, score, text, creation_date, user_id, content_license
FROM comments WHERE post_id = $1;

-- name: ListUsersFromPost :many
SELECT * FROM users WHERE id IN (SELECT user_id FROM post_history WHERE post_id = $1 ORDER BY creation_date);


----------- COMMENTS

-- name: GetComment :one
SELECT id, post_id, score, text, creation_date, user_id, content_license
FROM comments WHERE id = $1;


----------- USERS

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

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