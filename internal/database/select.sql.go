// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: select.sql

package database

import (
	"context"
)

const getBadge = `-- name: GetBadge :one
SELECT id, user_id, name, date, class, tag_based FROM badges WHERE id = $1
`

func (q *Queries) GetBadge(ctx context.Context, id int64) (Badge, error) {
	row := q.db.QueryRowContext(ctx, getBadge, id)
	var i Badge
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Date,
		&i.Class,
		&i.TagBased,
	)
	return i, err
}

const getComment = `-- name: GetComment :one

SELECT (id, post_id, score, text, creation_date, user_id, content_license)
FROM comments WHERE id = $1
`

//--------- COMMENTS
func (q *Queries) GetComment(ctx context.Context, id int64) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getComment, id)
	var column_1 interface{}
	err := row.Scan(&column_1)
	return column_1, err
}

const getPost = `-- name: GetPost :one

SELECT (id, post_type_id, parent_id, accepted_answer_id, creation_date, closed_date, score, view_count, body, tags, answer_count, comment_count, favorite_count, content_license)
FROM posts WHERE id = $1
`

//--------- POSTS
func (q *Queries) GetPost(ctx context.Context, id int64) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var column_1 interface{}
	err := row.Scan(&column_1)
	return column_1, err
}

const getTag = `-- name: GetTag :many
SELECT id, tag_name, count, is_required, is_moderator_only, wiki_post_id, excerpt_post_id FROM tags WHERE id = $1
`

func (q *Queries) GetTag(ctx context.Context, id int64) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, getTag, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(
			&i.ID,
			&i.TagName,
			&i.Count,
			&i.IsRequired,
			&i.IsModeratorOnly,
			&i.WikiPostID,
			&i.ExcerptPostID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTagFromName = `-- name: GetTagFromName :many
SELECT id, tag_name, count, is_required, is_moderator_only, wiki_post_id, excerpt_post_id FROM tags WHERE tag_name = $1
`

func (q *Queries) GetTagFromName(ctx context.Context, tagName string) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, getTagFromName, tagName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(
			&i.ID,
			&i.TagName,
			&i.Count,
			&i.IsRequired,
			&i.IsModeratorOnly,
			&i.WikiPostID,
			&i.ExcerptPostID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one

SELECT id, reputation, creation_date, display_name, last_access_date, location, website_url, about_me, views, upvotes, downvotes, account_id, profile_image_url FROM users WHERE account_id = $1
`

//--------- USERS
func (q *Queries) GetUser(ctx context.Context, accountID int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, accountID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Reputation,
		&i.CreationDate,
		&i.DisplayName,
		&i.LastAccessDate,
		&i.Location,
		&i.WebsiteUrl,
		&i.AboutMe,
		&i.Views,
		&i.Upvotes,
		&i.Downvotes,
		&i.AccountID,
		&i.ProfileImageUrl,
	)
	return i, err
}

const listAnswers = `-- name: ListAnswers :many
SELECT (id, post_type_id, parent_id, accepted_answer_id, creation_date, closed_date, score, view_count, body, tags, answer_count, comment_count, favorite_count, content_license)
FROM posts WHERE parent_id = $1
`

func (q *Queries) ListAnswers(ctx context.Context, parentID int64) ([]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, listAnswers, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var column_1 interface{}
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBadges = `-- name: ListBadges :many

SELECT id, user_id, name, date, class, tag_based FROM badges
`

//--------- BADGES
func (q *Queries) ListBadges(ctx context.Context) ([]Badge, error) {
	rows, err := q.db.QueryContext(ctx, listBadges)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Badge
	for rows.Next() {
		var i Badge
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Date,
			&i.Class,
			&i.TagBased,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBadgesFromUser = `-- name: ListBadgesFromUser :many
SELECT id, user_id, name, date, class, tag_based FROM badges WHERE user_id = $1
`

func (q *Queries) ListBadgesFromUser(ctx context.Context, userID int64) ([]Badge, error) {
	rows, err := q.db.QueryContext(ctx, listBadgesFromUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Badge
	for rows.Next() {
		var i Badge
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Date,
			&i.Class,
			&i.TagBased,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCommentsFromPost = `-- name: ListCommentsFromPost :many
SELECT (id, post_id, score, text, creation_date, user_id, content_license)
FROM comments WHERE post_id = $1
`

func (q *Queries) ListCommentsFromPost(ctx context.Context, postID int64) ([]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, listCommentsFromPost, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var column_1 interface{}
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listHistoryFromPost = `-- name: ListHistoryFromPost :many
SELECT (id, post_history_type_id, post_id, revision_guid, creation_date, user_id,comment, text, content_license) 
FROM post_history WHERE post_id = $1
`

func (q *Queries) ListHistoryFromPost(ctx context.Context, postID int64) ([]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, listHistoryFromPost, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var column_1 interface{}
		if err := rows.Scan(&column_1); err != nil {
			return nil, err
		}
		items = append(items, column_1)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRelatedPosts = `-- name: ListRelatedPosts :many
SELECT id, creation_date, post_id, related_post_id, link_type_id FROM post_links WHERE post_id = $1 LIMIT $2
`

type ListRelatedPostsParams struct {
	PostID int64 `json:"postID"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) ListRelatedPosts(ctx context.Context, arg ListRelatedPostsParams) ([]PostLink, error) {
	rows, err := q.db.QueryContext(ctx, listRelatedPosts, arg.PostID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PostLink
	for rows.Next() {
		var i PostLink
		if err := rows.Scan(
			&i.ID,
			&i.CreationDate,
			&i.PostID,
			&i.RelatedPostID,
			&i.LinkTypeID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTags = `-- name: ListTags :many

SELECT id, tag_name, count, is_required, is_moderator_only, wiki_post_id, excerpt_post_id FROM tags ORDER BY tag_name
`

//--------- TAGS
func (q *Queries) ListTags(ctx context.Context) ([]Tag, error) {
	rows, err := q.db.QueryContext(ctx, listTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Tag
	for rows.Next() {
		var i Tag
		if err := rows.Scan(
			&i.ID,
			&i.TagName,
			&i.Count,
			&i.IsRequired,
			&i.IsModeratorOnly,
			&i.WikiPostID,
			&i.ExcerptPostID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsersFromPost = `-- name: ListUsersFromPost :many
SELECT id, reputation, creation_date, display_name, last_access_date, location, website_url, about_me, views, upvotes, downvotes, account_id, profile_image_url FROM users WHERE id IN (SELECT user_id FROM post_history WHERE post_id = $1 ORDER BY creation_date)
`

func (q *Queries) ListUsersFromPost(ctx context.Context, postID int64) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsersFromPost, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Reputation,
			&i.CreationDate,
			&i.DisplayName,
			&i.LastAccessDate,
			&i.Location,
			&i.WebsiteUrl,
			&i.AboutMe,
			&i.Views,
			&i.Upvotes,
			&i.Downvotes,
			&i.AccountID,
			&i.ProfileImageUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listVotesFromPost = `-- name: ListVotesFromPost :many
SELECT id, post_id, vote_type_id, creation_date FROM votes WHERE post_id = $1
`

func (q *Queries) ListVotesFromPost(ctx context.Context, postID int64) ([]Vote, error) {
	rows, err := q.db.QueryContext(ctx, listVotesFromPost, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Vote
	for rows.Next() {
		var i Vote
		if err := rows.Scan(
			&i.ID,
			&i.PostID,
			&i.VoteTypeID,
			&i.CreationDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
