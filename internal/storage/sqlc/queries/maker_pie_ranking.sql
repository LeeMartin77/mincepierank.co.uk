
-- name: GetUserMakerPieRanking :one
SELECT mpry.* 
FROM maker_pie_ranking_yearly mpry
WHERE
    userid = sqlc.arg(userid)::text
    AND year = sqlc.arg(year)::int
    AND makerid = sqlc.arg(makerid)::text
    AND pieid = sqlc.arg(pieid)::text;
