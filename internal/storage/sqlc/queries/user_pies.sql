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

-- name: GetUserPieByOid :one
SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
       image_file, web_link, pack_count, pack_price_in_pence, clean,
       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
       submitted_by, submitted_at, moderated_by, moderated_at
FROM user_pie_yearly
WHERE promoted_maker_pie_oid = sqlc.arg(promotedMakerPieOid)::text
   OR (year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text);

-- name: DeleteUserPieByOid :exec
DELETE FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text;

-- name: UpdateUserPieStatusWithNote :exec
UPDATE user_pie_yearly
SET status = sqlc.arg(status)::text,
    moderation_note = sqlc.arg(moderation_note)::text,
    moderated_by = sqlc.arg(moderated_by)::text,
    moderated_at = NOW()
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text;

-- name: UpdateUserPieDuplicateOf :exec
UPDATE user_pie_yearly
SET duplicate_of = sqlc.arg(duplicate_of)::text,
    status = sqlc.arg(status)::text,
    moderated_by = sqlc.arg(moderated_by)::text,
    moderated_at = NOW()
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text;

-- name: UpdateUserPiePromotedMakerPieOid :exec
UPDATE user_pie_yearly
SET promoted_maker_pie_oid = sqlc.arg(promotedMakerPieOid)::text,
    status = sqlc.arg(status)::text,
    moderated_by = sqlc.arg(moderated_by)::text,
    moderated_at = NOW()
WHERE year = sqlc.arg(year)::int AND id = sqlc.arg(id)::text;

-- name: GetUserPieAdminQueueCountByStatus :one
SELECT COUNT(*) FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int AND status = sqlc.arg(status)::text;

-- name: GetUserPieDuplicateGroups :many
SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
       image_file, web_link, pack_count, pack_price_in_pence, clean,
       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
       submitted_by, submitted_at, moderated_by, moderated_at
FROM user_pie_yearly
WHERE year = sqlc.arg(year)::int
  AND status IN ('pending', 'approved')
  AND id != sqlc.arg(excludeId)::text
  AND (
    -- Same maker and similar displayname (normalised)
    (maker IS NOT NULL AND maker = sqlc.arg(maker)::text
     AND lower(regexp_replace(displayname, '[^a-z0-9]+', '', 'g')) =
         lower(regexp_replace(sqlc.arg(displayname)::text, '[^a-z0-9]+', '', 'g')))
    OR
    -- Same non-empty web_link
    (web_link IS NOT NULL AND web_link = sqlc.arg(webLink)::text
     AND web_link != '')
  )
ORDER BY submitted_at ASC;

-- name: GetUserPieDuplicateGroupsAll :many
-- Find all user pies that have at least one potential duplicate
-- (same maker and normalised displayname, or same web_link).
-- Returns pies ordered by maker, then normalised name.
WITH normalised AS (
  SELECT *,
    lower(regexp_replace(displayname, '[^a-z0-9]+', '', 'g')) AS norm_name
  FROM user_pie_yearly
  WHERE year = sqlc.arg(year)::int
    AND status IN ('pending', 'approved')
)
SELECT a.year, a.id, a.owner_userid, a.maker, a.location, a.displayname,
       a.fresh, a.labels, a.image_file, a.web_link, a.pack_count,
       a.pack_price_in_pence, a.clean, a.status, a.moderation_note,
       a.duplicate_of, a.promoted_maker_pie_oid, a.submitted_by,
       a.submitted_at, a.moderated_by, a.moderated_at
FROM normalised a
JOIN normalised b
  ON a.maker = b.maker
  AND a.norm_name = b.norm_name
  AND a.id != b.id
  AND a.submitted_at <= b.submitted_at
GROUP BY a.year, a.id, a.owner_userid, a.maker, a.location, a.displayname,
         a.fresh, a.labels, a.image_file, a.web_link, a.pack_count,
         a.pack_price_in_pence, a.clean, a.status, a.moderation_note,
         a.duplicate_of, a.promoted_maker_pie_oid, a.submitted_by,
         a.submitted_at, a.moderated_by, a.moderated_at
ORDER BY a.maker, a.norm_name, a.submitted_at;

-- name: GetExistingBrandPiesForYear :many
SELECT year, makerid, id, displayname, fresh, labels, image_file, web_link,
       pack_count, pack_price_in_pence, validated, oid
FROM maker_pie_yearly
WHERE year = sqlc.arg(year)::int
  AND (makerid = sqlc.arg(makerid)::text OR web_link = sqlc.arg(webLink)::text)
ORDER BY displayname;
