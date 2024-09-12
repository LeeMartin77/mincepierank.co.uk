-- name: GetAllMakerPies :many
SELECT mpy.oid::text as oidstr , mpy.* FROM maker_pie_yearly mpy;

-- name: GetMakerPieByOid :one
SELECT * FROM maker_pie_yearly where oid = uuid(sqlc.arg(oid)::text);