-- name: GetConfig :one
SELECT * FROM config WHERE key = sqlc.arg(key)::text;

-- name: GetAllConfig :many
SELECT key, coalesce(value, '') FROM config;

-- name: InsertConfig :exec
INSERT INTO config (key, value) VALUES (sqlc.arg(key)::text, sqlc.arg(value)::text);
-- name: UpdateConfig :exec
UPDATE config 
SET value=sqlc.arg(value)::text
WHERE key=sqlc.arg(key)::text;

-- name: DeleteConfig :exec
DELETE FROM config 
WHERE key=sqlc.arg(key)::text;


-- name: GetActiveYear :one
SELECT CAST(coalesce(value, '2023') AS bigint) FROM config WHERE key = 'activeYear';

-- name: IsAdminId :one
SELECT CASE 
    WHEN EXISTS (
        SELECT 1 
        FROM admins
        WHERE id = sqlc.arg(user_id)::text
    ) 
    THEN true 
    ELSE false 
END;