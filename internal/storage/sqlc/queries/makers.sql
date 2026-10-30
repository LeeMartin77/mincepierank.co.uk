-- name: ListMakersAdmin :many
SELECT mkr.id,
       coalesce(mkr.name, '') as name,
       coalesce(mkr.logo, '') as logo,
       coalesce(mkr.website, '') as website,
       coalesce(pies.cnt, 0)::bigint as piecount
FROM maker mkr
LEFT JOIN (
    SELECT makerid, count(1) as cnt
    FROM maker_pie_yearly
    GROUP BY makerid
) pies on pies.makerid = mkr.id
ORDER BY lower(mkr.name), mkr.id;

-- name: GetMakerPieCount :one
SELECT count(1)::bigint
FROM maker_pie_yearly
WHERE makerid = sqlc.arg(makerid)::text;

-- name: GetMakerByID :one
SELECT id, coalesce(name, '') as name, coalesce(logo, '') as logo, coalesce(website, '') as website
FROM maker
WHERE id = sqlc.arg(id)::text;

-- name: CreateMaker :exec
INSERT INTO maker (id, name, logo, website)
VALUES (
    sqlc.arg(id)::text,
    sqlc.arg(name)::text,
    sqlc.arg(logo)::text,
    sqlc.arg(website)::text
)
ON CONFLICT DO NOTHING;

-- name: UpdateMakerByID :execresult
UPDATE maker
SET name = sqlc.arg(name)::text, logo = sqlc.arg(logo)::text, website = sqlc.arg(website)::text
WHERE id = sqlc.arg(id)::text;

-- name: DeleteMakerByID :exec
DELETE FROM maker
WHERE id = sqlc.arg(id)::text;

-- name: GetMakersForYear :many
WITH ordered_brand_ids as (
    SELECT makerid, count(1) as cnt
    FROM maker_pie_yearly
    WHERE year = sqlc.arg(year)::int
    GROUP BY makerid
)
SELECT mkr.id,
       coalesce(mkr.name, '') as name,
       coalesce(mkr.logo, '') as logo,
       coalesce(mkr.website, '') as website,
       mpy.cnt::bigint as piecount
FROM maker mkr
INNER JOIN ordered_brand_ids mpy
        on mpy.makerid = mkr.id
ORDER BY mpy.cnt DESC;
