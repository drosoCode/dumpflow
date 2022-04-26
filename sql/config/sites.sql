-- name: ListSites :many
SELECT * FROM sites ORDER BY link;

-- name: GetSite :one
SELECT * FROM sites WHERE db_name = $1;

-- name: AddSite :exec
INSERT INTO sites (db_name, link, update_date, auto_update, enabled) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING;

-- name: RemoveSite :exec
DELETE FROM sites WHERE db_name = $1;
