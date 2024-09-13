-- name: GetAllMakerPies :many
SELECT mpy.oid::text as oidstr , mpy.* FROM maker_pie_yearly mpy;

-- name: GetMakerPieByOid :one
SELECT * FROM maker_pie_yearly where oid = uuid(sqlc.arg(oid)::text);

-- name: DeleteMakerPieByOid :exec
DELETE FROM maker_pie_yearly where oid = uuid(sqlc.arg(oid)::text);

-- name: CreateMakerPie :one
INSERT INTO maker_pie_yearly
    (
        year,
        makerid,
        id,
        displayname,
        fresh,
        image_file,
        web_link,
        pack_count,
        pack_price_in_pence,
        validated
    ) 
    VALUES
    (
        sqlc.arg(year)::int,
        sqlc.arg(makerid)::text,
        sqlc.arg(id)::text,
        sqlc.arg(displayname)::text,
        sqlc.arg(fresh)::boolean,
        sqlc.arg(image_file)::text,
        sqlc.arg(web_link)::text,
        sqlc.arg(pack_count)::int,
        sqlc.arg(pack_price_in_pence)::int,
        sqlc.arg(validated)::boolean
    )
RETURNING oid::text as oid;

-- name: UpdateMakerPieByOid :exec
UPDATE maker_pie_yearly
SET
    -- year=sqlc.arg(year)::int,
    -- makerid=sqlc.arg(makerid)::text,
    -- id=sqlc.arg(id)::text,
    displayname=sqlc.arg(displayname)::text,
    fresh=sqlc.arg(fresh)::boolean,
    -- image_file=sqlc.arg(image_file)::text,
    web_link=sqlc.arg(web_link)::text,
    pack_count=sqlc.arg(pack_count)::int,
    pack_price_in_pence=sqlc.arg(pack_price_in_pence)::int,
    validated=sqlc.arg(validated)::boolean
WHERE oid = uuid(sqlc.arg(oid)::text);