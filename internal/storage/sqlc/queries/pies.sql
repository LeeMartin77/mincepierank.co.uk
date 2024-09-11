-- name: GetConfig :one
SELECT key, value FROM config WHERE key = sqlc.arg(key)::text;

-- name: GetActiveYear :one
SELECT CAST(coalesce(value, '2023') AS bigint) FROM config WHERE key = 'activeYear';