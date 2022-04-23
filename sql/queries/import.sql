-- name: AddBadge :exec
INSERT INTO badges (id, user_id, name, date, class, tag_based) 
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT DO NOTHING;

-- name: AddComment :exec
INSERT INTO comments (id, post_id, score, text, creation_date, user_id, content_license)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT DO NOTHING;

-- name: AddPostHistory :exec
INSERT INTO post_history (id, post_history_type_id, post_id, revision_guid, creation_date, user_id, comment, text, content_license)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    ON CONFLICT DO NOTHING;
    
-- name: AddPostLink :exec
INSERT INTO post_links (id, creation_date, post_id, related_post_id, link_type_id)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT DO NOTHING;

-- name: AddPost :exec
INSERT INTO posts (id, post_type_id, parent_id, accepted_answer_id, creation_date, score, view_count, body, answer_count, comment_count, favorite_count, content_license, closed_date, tags)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    ON CONFLICT DO NOTHING;

-- name: AddTag :exec
INSERT INTO tags (id, tag_name, count, is_required, is_moderator_only, wiki_post_id, excerpt_post_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT DO NOTHING;

-- name: AddUser :exec
INSERT INTO users (id, reputation, creation_date, display_name, last_access_date, location, website_url, about_me, views, upvotes, downvotes, account_id, profile_image_url)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    ON CONFLICT DO NOTHING;

-- name: AddVote :exec
INSERT INTO votes (id, post_id, vote_type_id, creation_date)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT DO NOTHING;
