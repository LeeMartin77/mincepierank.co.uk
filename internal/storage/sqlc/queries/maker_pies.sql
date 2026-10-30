-- name: GetAllMakerPies :many
SELECT mpy.oid::text as oidstr , mpy.* FROM maker_pie_yearly mpy;

-- name: GetMakerPieByOid :one
SELECT mpy.oid::text as oidstr , mpy.* FROM maker_pie_yearly mpy where mpy.oid = uuid(sqlc.arg(oid)::text);

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

-- name: UpdateMakerPieImageByOid :exec
UPDATE maker_pie_yearly
SET image_file=sqlc.arg(image_file)::text
WHERE oid = uuid(sqlc.arg(oid)::text);

-- name: UpdateMakerPieLabels :exec
UPDATE maker_pie_yearly
SET labels = sqlc.arg(labels)::text[]
WHERE oid = uuid(sqlc.arg(oid)::text);

-- name: GetMakerPieYears :many
SELECT DISTINCT year
FROM maker_pie_yearly
WHERE validated = true;

-- name: GetFilterableMakerPies :many
WITH pie_rankings AS (
    SELECT
        year,
        makerid,
        pieid,
        SUM(pastry::float )/COUNT(1) as pastry,
        SUM(filling::float )/COUNT(1) as filling,
        SUM(topping::float )/COUNT(1) as topping,
        SUM(looks::float )/COUNT(1) as looks,
        SUM(value::float )/COUNT(1) as value,
        (
            SUM(pastry::float )/COUNT(1) +
            SUM(filling::float )/COUNT(1) +
            SUM(topping::float )/COUNT(1) +
            SUM(looks::float )/COUNT(1) +
            SUM(value::float )/COUNT(1)
        )/5 as avg,
        COUNT(1)::bigint as count
    FROM maker_pie_ranking_yearly
    WHERE year = sqlc.arg(year)::int
    GROUP BY year, makerid, pieid
)
SELECT
    mpy.year,
    mpy.makerid,
    mpy.id,
    coalesce(mpy.displayname, '') as displayname,
    coalesce(mpy.fresh, false) as fresh,
    cts.categories::jsonb as categories,
    coalesce(mpy.image_file, '') as image_file,
    coalesce(mpy.web_link, '') as web_link,
    coalesce(mpy.pack_count, 0)::int as pack_count,
    coalesce(mpy.pack_price_in_pence, 0)::int as pack_price_in_pence,
    COALESCE(tp.pastry, 0::float)::float as pastry,
    COALESCE(tp.filling, 0::float)::float as filling,
    COALESCE(tp.topping, 0::float)::float as topping,
    COALESCE(tp.looks, 0::float)::float as looks,
    COALESCE(tp.value, 0::float)::float as value,
    COALESCE(tp.avg, 0::float)::float as avg,
    COALESCE(tp.count, 0)::bigint as count
FROM maker_pie_yearly mpy
LEFT JOIN pie_rankings tp ON
    tp.year = mpy.year AND
    tp.makerid = mpy.makerid AND
    tp.pieid = mpy.id
INNER JOIN maker_pie_yearly_categories cts ON
    cts.oid = mpy.oid
WHERE mpy.year = sqlc.arg(year)::int
    AND mpy.validated = true
    AND (sqlc.arg(brand_ids)::text = '' OR mpy.makerid = ANY(string_to_array(sqlc.arg(brand_ids)::text, ',')))
    AND (sqlc.arg(category_slugs)::text = '' OR cts.category_slugs @> string_to_array(sqlc.arg(category_slugs)::text, ','))
ORDER BY tp.avg DESC NULLS LAST
LIMIT sqlc.arg(page_limit)::int
OFFSET sqlc.arg(page_offset)::int;

-- name: GetFilterableMakerPiesForUser :many
WITH user_rankings AS (
    SELECT
        year,
        makerid,
        pieid,
        pastry,
        filling,
        topping,
        looks,
        value,
        (
            pastry::float +
            filling::float +
            topping::float +
            looks::float +
            value::float
        )/5 as avg
    FROM maker_pie_ranking_yearly
    WHERE userid = sqlc.arg(userid)::text
)
SELECT
    mpy.year,
    mpy.makerid,
    mpy.id,
    coalesce(mpy.displayname, '') as displayname,
    coalesce(mpy.fresh, false) as fresh,
    cts.categories::jsonb as categories,
    coalesce(mpy.image_file, '') as image_file,
    coalesce(mpy.web_link, '') as web_link,
    coalesce(mpy.pack_count, 0)::int as pack_count,
    coalesce(mpy.pack_price_in_pence, 0)::int as pack_price_in_pence,
    COALESCE(tp.pastry, 0::float)::float as pastry,
    COALESCE(tp.filling, 0::float)::float as filling,
    COALESCE(tp.topping, 0::float)::float as topping,
    COALESCE(tp.looks, 0::float)::float as looks,
    COALESCE(tp.value, 0::float)::float as value,
    COALESCE(tp.avg, 0::float)::float as avg
FROM maker_pie_yearly mpy
INNER JOIN user_rankings tp ON
    tp.year = mpy.year AND
    tp.makerid = mpy.makerid AND
    tp.pieid = mpy.id
INNER JOIN maker_pie_yearly_categories cts ON
    cts.oid = mpy.oid
WHERE mpy.year = sqlc.arg(year)::int
    AND mpy.validated = true
    AND (sqlc.arg(brand_ids)::text = '' OR mpy.makerid = ANY(string_to_array(sqlc.arg(brand_ids)::text, ',')))
    AND (sqlc.arg(category_slugs)::text = '' OR cts.category_slugs @> string_to_array(sqlc.arg(category_slugs)::text, ','))
ORDER BY tp.avg DESC NULLS LAST
LIMIT sqlc.arg(page_limit)::int
OFFSET sqlc.arg(page_offset)::int;

-- name: GetTopMakerPie :one
WITH aggregated_pies AS (
    SELECT
        year,
        makerid,
        pieid,
        SUM(pastry::float )/COUNT(1) as pastry,
        SUM(filling::float )/COUNT(1) as filling,
        SUM(topping::float )/COUNT(1) as topping,
        SUM(looks::float )/COUNT(1) as looks,
        SUM(value::float )/COUNT(1) as value,
        (
            SUM(pastry::float )/COUNT(1) +
            SUM(filling::float )/COUNT(1) +
            SUM(topping::float )/COUNT(1) +
            SUM(looks::float )/COUNT(1) +
            SUM(value::float )/COUNT(1)
        )/5 as avg,
        COUNT(1)::bigint as count
    FROM maker_pie_ranking_yearly
    WHERE year = sqlc.arg(year)::int
    GROUP BY year, makerid, pieid
), top_pie AS (
    SELECT *
    FROM aggregated_pies
    ORDER BY avg DESC
    LIMIT 1
)
SELECT
    mpy.year,
    mpy.makerid,
    mpy.id,
    coalesce(mpy.displayname, '') as displayname,
    coalesce(mpy.fresh, false) as fresh,
    cts.categories::jsonb as categories,
    coalesce(mpy.image_file, '') as image_file,
    coalesce(mpy.web_link, '') as web_link,
    coalesce(mpy.pack_count, 0)::int as pack_count,
    coalesce(mpy.pack_price_in_pence, 0)::int as pack_price_in_pence,
    COALESCE(tp.pastry, 0::float)::float as pastry,
    COALESCE(tp.filling, 0::float)::float as filling,
    COALESCE(tp.topping, 0::float)::float as topping,
    COALESCE(tp.looks, 0::float)::float as looks,
    COALESCE(tp.value, 0::float)::float as value,
    COALESCE(tp.avg, 0::float)::float as avg,
    COALESCE(tp.count, 0)::bigint as count
FROM maker_pie_yearly mpy
INNER JOIN top_pie tp ON
    tp.year = mpy.year AND
    tp.makerid = mpy.makerid AND
    tp.pieid = mpy.id
INNER JOIN maker_pie_yearly_categories cts ON
    cts.oid = mpy.oid;

-- name: GetLatestPieRanking :one
WITH top_pie AS (
    SELECT
        year,
        makerid,
        pieid,
        pastry,
        filling,
        topping,
        looks,
        value,
        (
            pastry::float +
            filling::float +
            topping::float +
            looks::float +
            value::float
        )/5 as avg,
        1::bigint as count,
        last_updated
    FROM maker_pie_ranking_yearly
    WHERE year = sqlc.arg(year)::int
    ORDER BY last_updated DESC
    LIMIT 1
)
SELECT
    mpy.year,
    mpy.makerid,
    mpy.id,
    coalesce(mpy.displayname, '') as displayname,
    coalesce(mpy.fresh, false) as fresh,
    cts.categories::jsonb as categories,
    coalesce(mpy.image_file, '') as image_file,
    coalesce(mpy.web_link, '') as web_link,
    coalesce(mpy.pack_count, 0)::int as pack_count,
    coalesce(mpy.pack_price_in_pence, 0)::int as pack_price_in_pence,
    COALESCE(tp.pastry, 0::float)::float as pastry,
    COALESCE(tp.filling, 0::float)::float as filling,
    COALESCE(tp.topping, 0::float)::float as topping,
    COALESCE(tp.looks, 0::float)::float as looks,
    COALESCE(tp.value, 0::float)::float as value,
    COALESCE(tp.avg, 0::float)::float as avg,
    COALESCE(tp.count, 0)::bigint as count,
    COALESCE(tp.last_updated, 'epoch'::timestamp)::timestamptz as last_updated
FROM maker_pie_yearly mpy
INNER JOIN top_pie tp ON
    tp.year = mpy.year AND
    tp.makerid = mpy.makerid AND
    tp.pieid = mpy.id
INNER JOIN maker_pie_yearly_categories cts ON
    cts.oid = mpy.oid;

-- name: GetMakerPieCategoriesForYear :many
SELECT distinct c.id::text as idstr, coalesce(c.slug, '') as slug, coalesce(c.label, '') as label
FROM categories c
INNER JOIN maker_pie_categories mpc ON
    mpc.category_id = c.id
INNER JOIN maker_pie_yearly mpy ON
    mpc.maker_pie_oid = mpy.oid AND
    mpy.year = sqlc.arg(year)::int
INNER JOIN maker_pie_yearly_categories cts ON
    cts.oid = mpy.oid
WHERE
    (sqlc.arg(brand)::text = '' OR mpy.makerid = sqlc.arg(brand)::text)
    AND (sqlc.arg(required_slugs)::text = '' OR cts.category_slugs @> string_to_array(sqlc.arg(required_slugs)::text, ','));

-- name: GetMakerPieCategoriesForMakerPieOid :many
SELECT c.id::text as idstr, coalesce(c.slug, '') as slug, coalesce(c.label, '') as label
FROM categories c
INNER JOIN maker_pie_categories mpc ON
    mpc.category_id = c.id
WHERE
    mpc.maker_pie_oid = uuid(sqlc.arg(oid)::text);

-- name: GetMakerPieCategories :many
SELECT id::text as idstr, coalesce(slug, '') as slug, coalesce(label, '') as label
FROM categories;

-- name: DeletePieCategoriesByPieOid :exec
DELETE FROM maker_pie_categories
WHERE maker_pie_oid = uuid(sqlc.arg(oid)::text);

-- name: InsertPieCategories :exec
INSERT INTO maker_pie_categories (maker_pie_oid, category_id)
SELECT uuid(sqlc.arg(oid)::text), u.cid::uuid
FROM unnest(sqlc.arg(category_ids)::text[]) as u(cid);
