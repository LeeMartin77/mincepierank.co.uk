
-- name: GetUserMakerPieRanking :one
SELECT mpry.* 
FROM maker_pie_ranking_yearly mpry
WHERE
    userid = sqlc.arg(userid)::text
    AND year = sqlc.arg(year)::int
    AND makerid = sqlc.arg(makerid)::text
    AND pieid = sqlc.arg(pieid)::text;

-- name: UpsertUserMakerPieRanking :exec
INSERT INTO maker_pie_ranking_yearly (year, makerid, pieid, userid, pastry, filling, topping, looks, value, last_updated)
VALUES (
        sqlc.arg(year)::int,
        sqlc.arg(makerid)::text,
        sqlc.arg(pieid)::text,
        sqlc.arg(userid)::text,
        sqlc.arg(pastry)::int,
        sqlc.arg(filling)::int,
        sqlc.arg(topping)::int,
        sqlc.arg(looks)::int,
        sqlc.arg(value)::int,
        NOW()
)
ON CONFLICT (year, makerid, pieid, userid)
DO UPDATE SET pastry = EXCLUDED.pastry, filling = EXCLUDED.filling, topping = EXCLUDED.topping, looks = EXCLUDED.looks, value = EXCLUDED.value, last_updated = NOW();
-- name: ListMakerPieRankingSummaries :many
SELECT mpy.year,
       mpy.makerid,
       mpy.id AS pieid,
       mpy.displayname,
       m.name AS makername,
       count(mpry.userid)::bigint AS ranking_count
FROM maker_pie_yearly mpy
JOIN maker m ON m.id = mpy.makerid
LEFT JOIN maker_pie_ranking_yearly mpry
    ON mpry.year = mpy.year AND mpry.makerid = mpy.makerid AND mpry.pieid = mpy.id
WHERE mpy.year = sqlc.arg(year)::int
GROUP BY mpy.year, mpy.makerid, mpy.id, mpy.displayname, m.name
ORDER BY m.name, mpy.displayname;

-- name: ListMakerPieRankingsForPie :many
SELECT *
FROM maker_pie_ranking_yearly
WHERE year = sqlc.arg(year)::int
    AND makerid = sqlc.arg(makerid)::text
    AND pieid = sqlc.arg(pieid)::text
ORDER BY userid;

-- name: DeleteMakerPieRanking :exec
DELETE FROM maker_pie_ranking_yearly
WHERE year = sqlc.arg(year)::int
    AND makerid = sqlc.arg(makerid)::text
    AND pieid = sqlc.arg(pieid)::text
    AND userid = sqlc.arg(userid)::text;
