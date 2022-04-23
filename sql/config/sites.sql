-- name: ListSites :many
SELECT * FROM sites;

-- name: GetSite :one
SELECT * FROM sites WHERE id = $1;

-- name: GetSiteFromDB :one
SELECT * FROM sites WHERE db_name = $1;

-- name: AddSite :exec
INSERT INTO sites (db_name, link, update_date, auto_update, enabled) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO UPDATE SET update_date = EXCLUDED.update_date;

-- name: RemoveSite :exec
DELETE FROM sites WHERE id = $1;

-- name: RemoveSiteFromDB :exec
DELETE FROM sites WHERE db_name = $1;