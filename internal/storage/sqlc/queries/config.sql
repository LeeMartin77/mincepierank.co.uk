-- name: GetConfig :one
SELECT key, value FROM config WHERE key = sqlc.arg(key)::text;

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