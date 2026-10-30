-- name: ListCategoriesAdmin :many
SELECT c.id::text as idstr,
       coalesce(c.slug, '') as slug,
       coalesce(c.label, '') as label,
       coalesce(usage_counts.cnt, 0)::bigint as piecount
FROM categories c
LEFT JOIN (
    SELECT category_id, count(1) as cnt
    FROM maker_pie_categories
    GROUP BY category_id
) usage_counts on usage_counts.category_id = c.id
ORDER BY lower(coalesce(c.label, c.slug)), c.id;

-- name: GetCategoryById :one
SELECT id::text as idstr, coalesce(slug, '') as slug, coalesce(label, '') as label
FROM categories
WHERE id = uuid(sqlc.arg(id)::text);

-- name: CategorySlugTaken :one
SELECT CASE
    WHEN EXISTS (
        SELECT 1
        FROM categories
        WHERE slug = sqlc.arg(slug)::text AND id <> uuid(sqlc.arg(except_id)::text)
    )
    THEN true
    ELSE false
END;

-- name: CreateCategory :one
INSERT INTO categories (slug, label)
VALUES (sqlc.arg(slug)::text, sqlc.arg(label)::text)
RETURNING id::text as idstr;

-- name: UpdateCategory :execresult
UPDATE categories
SET slug = sqlc.arg(slug)::text, label = sqlc.arg(label)::text
WHERE id = uuid(sqlc.arg(id)::text);

-- name: GetCategoryPieCount :one
SELECT count(1)::bigint
FROM maker_pie_categories
WHERE category_id = uuid(sqlc.arg(category_id)::text);

-- name: DeletePieCategoriesByCategory :exec
DELETE FROM maker_pie_categories
WHERE category_id = uuid(sqlc.arg(category_id)::text);

-- name: DeleteCategory :execresult
DELETE FROM categories
WHERE id = uuid(sqlc.arg(id)::text);

-- name: GetCategoriesBySlugs :many
SELECT id::text as idstr, slug
FROM categories
WHERE slug = ANY(sqlc.arg(slugs)::text[]);
