-- name: CreateUserPie :one
INSERT INTO user_pie_yearly
    (year, id, owner_userid, maker, location, displayname, fresh, labels,
     image_file, web_link, pack_count, pack_price_in_pence, clean,
     status, submitted_by, submitted_at, duplicate_of)
VALUES
    (sqlc.arg(year)::int, sqlc.arg(id)::text, sqlc.arg(owner_userid)::text,
     sqlc.arg(maker)::text, sqlc.arg(location)::text,
     sqlc.arg(displayname)::text, sqlc.arg(fresh)::boolean,
     sqlc.arg(labels)::text[], sqlc.arg(image_file)::text,
     sqlc.arg(web_link)::text, sqlc.arg(pack_count)::int,
     sqlc.arg(pack_price_in_pence)::int, sqlc.arg(clean)::boolean,
     'pending', sqlc.arg(submitted_by)::text, NOW(), sqlc.arg(duplicate_of)::text)
RETURNING year, id, owner_userid, maker, location, displayname, fresh, labels,
          image_file, web_link, pack_count, pack_price_in_pence, clean,
          status, moderation_note, duplicate_of, promoted_maker_pie_oid,
          submitted_by, submitted_at, moderated_by, moderated_at;

-- name: GetUserPie :one
SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
       image_file, web_link, pack_count, pack_price_in_pence, clean,
       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
       submitted_by, submitted_at, moderated_by, moderated_at
FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text;

-- name: GetUserPieByOwner :one
SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
       image_file, web_link, pack_count, pack_price_in_pence, clean,
       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
       submitted_by, submitted_at, moderated_by, moderated_at
FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text
  AND owner_userid = sqlc.arg(owner_userid)::text;

-- name: UpdateUserPie :exec
UPDATE user_pie_yearly
SET maker = sqlc.arg(maker)::text,
    location = sqlc.arg(location)::text,
    displayname = sqlc.arg(displayname)::text,
    fresh = sqlc.arg(fresh)::boolean,
    labels = sqlc.arg(labels)::text[],
    image_file = sqlc.arg(image_file)::text,
    web_link = sqlc.arg(web_link)::text,
    pack_count = sqlc.arg(pack_count)::int,
    pack_price_in_pence = sqlc.arg(pack_price_in_pence)::int
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text;

-- name: UpdateUserPieImage :exec
UPDATE user_pie_yearly
SET image_file = sqlc.arg(image_file)::text
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text;

-- name: DeleteUserPie :exec
DELETE FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text;

-- name: DeleteUserPieByOwner :exec
DELETE FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text
  AND owner_userid = sqlc.arg(owner_userid)::text;

-- name: GetUserPieYearIDs :many
SELECT id FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int;

-- name: GetUserPieCountByYear :one
SELECT COUNT(*) FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int;

-- name: GetUserPieCountByOwnerAndYear :one
SELECT COUNT(*) FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int AND owner_userid = sqlc.arg(owner_userid)::text;
