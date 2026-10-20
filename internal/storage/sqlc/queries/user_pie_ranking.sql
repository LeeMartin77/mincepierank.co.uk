-- name: UpsertUserPieRanking :exec
INSERT INTO user_pie_ranking_yearly (year, pieid, userid, pastry, filling, topping, looks, value, notes, last_updated)
VALUES (
    sqlc.arg(year)::int,
    sqlc.arg(pieid)::text,
    sqlc.arg(userid)::text,
    sqlc.arg(pastry)::int,
    sqlc.arg(filling)::int,
    sqlc.arg(topping)::int,
    sqlc.arg(looks)::int,
    sqlc.arg(value)::int,
    sqlc.arg(notes)::text,
    NOW()
)
ON CONFLICT (year, pieid, userid)
DO UPDATE SET pastry = EXCLUDED.pastry,
              filling = EXCLUDED.filling,
              topping = EXCLUDED.topping,
              looks = EXCLUDED.looks,
              value = EXCLUDED.value,
              notes = EXCLUDED.notes,
              last_updated = NOW();

-- name: GetUserPieRanking :one
SELECT year, pieid, userid, pastry, filling, topping, looks, value, notes, last_updated
FROM user_pie_ranking_yearly
WHERE year = sqlc.arg(year)::int AND pieid = sqlc.arg(pieid)::text AND userid = sqlc.arg(userid)::text;

-- name: DeleteUserPieRankings :exec
DELETE FROM user_pie_ranking_yearly
WHERE year = sqlc.arg(year)::int AND pieid = sqlc.arg(pieid)::text;

-- name: DeleteUserPieRankingsByUser :exec
DELETE FROM user_pie_ranking_yearly
WHERE year = sqlc.arg(year)::int AND userid = sqlc.arg(userid)::text;
