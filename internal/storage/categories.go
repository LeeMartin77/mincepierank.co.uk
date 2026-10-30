package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/rs/zerolog/log"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("a category with that slug already exists")
	// ErrUnknownCategories is returned by SetPieCategories when a requested
	// slug does not match any category, so typos are visible instead of being
	// silently dropped.
	ErrUnknownCategories = errors.New("unknown category slugs")
)

// AdminCategoryRow is a category plus how many brand pies use it, for the
// admin categories list and the delete guard.
type AdminCategoryRow struct {
	types.Category
	PieCount int64
}

// GetAllCategoriesAdmin returns every category with a total pie count across
// all years, ordered by label.
func (o *OperationWrapper) GetAllCategoriesAdmin(c context.Context) (*[]AdminCategoryRow, error) {
	r := []AdminCategoryRow{}
	sql := `
	select c.id, coalesce(c.slug, ''), coalesce(c.label, ''),
	       coalesce(usage.cnt, 0) as piecount
	from categories c
	left join (
		select category_id, count(1) as cnt
		from maker_pie_categories
		group by category_id
	) usage on usage.category_id = c.id
	order by lower(coalesce(c.label, c.slug)), c.id
	`
	rows, err := o.db.Query(c, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itm AdminCategoryRow
		if err := rows.Scan(&itm.Id, &itm.Slug, &itm.Label, &itm.PieCount); err != nil {
			log.Error().Err(err).Msg("Scan failed in get all categories")
			continue
		}
		r = append(r, itm)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during all categories iteration")
	}
	return &r, nil
}

// GetCategoryById returns the category or nil (no error) when not found.
func (o *OperationWrapper) GetCategoryById(c context.Context, id uuid.UUID) (*types.Category, error) {
	cat := types.Category{}
	err := o.db.QueryRow(c,
		`select id, coalesce(slug, ''), coalesce(label, '') from categories where id = $1`, id).
		Scan(&cat.Id, &cat.Slug, &cat.Label)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// GetCategoryPieCount counts brand pies (all years) using the category, used
// to guard category deletes against orphaning pie references.
func (o *OperationWrapper) GetCategoryPieCount(c context.Context, id uuid.UUID) (int64, error) {
	var cnt int64
	err := o.db.QueryRow(c,
		`select count(1) from maker_pie_categories where category_id = $1`, id).Scan(&cnt)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// CreateCategory inserts a new category. Returns ErrCategoryAlreadyExists
// when another category already has the slug (there is no unique constraint
// on slug, so the check happens here).
func (o *OperationWrapper) CreateCategory(c context.Context, cat types.Category) (*types.Category, error) {
	exists, err := o.categorySlugTaken(c, cat.Slug, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryAlreadyExists
	}

	created := types.Category{}
	err = o.db.QueryRow(c,
		`insert into categories (slug, label) values ($1, $2) returning id`,
		cat.Slug, cat.Label).Scan(&created.Id)
	if err != nil {
		return nil, err
	}
	created.Slug = cat.Slug
	created.Label = cat.Label
	return &created, nil
}

// UpdateCategory updates the slug and label of an existing category. Returns
// ErrCategoryNotFound when the id does not exist and ErrCategoryAlreadyExists
// when another category already has the slug.
func (o *OperationWrapper) UpdateCategory(c context.Context, cat types.Category) (*types.Category, error) {
	exists, err := o.categorySlugTaken(c, cat.Slug, cat.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryAlreadyExists
	}

	tag, err := o.db.Exec(c,
		`update categories set slug = $2, label = $3 where id = $1`,
		cat.Id, cat.Slug, cat.Label)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrCategoryNotFound
	}
	return &cat, nil
}

// DeleteCategory removes a category and any pie references to it. Callers
// guard deletes with GetCategoryPieCount; the reference cleanup runs in the
// same transaction as a safety net for the race between the two.
func (o *OperationWrapper) DeleteCategory(c context.Context, id uuid.UUID) error {
	tx, err := o.db.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback(c)
	}()

	if _, err := tx.Exec(c, `delete from maker_pie_categories where category_id = $1`, id); err != nil {
		return err
	}
	tag, err := tx.Exec(c, `delete from categories where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}
	return tx.Commit(c)
}

// categorySlugTaken reports whether any category other than exceptId already
// uses the slug.
func (o *OperationWrapper) categorySlugTaken(c context.Context, slug string, exceptId uuid.UUID) (bool, error) {
	var taken bool
	err := o.db.QueryRow(c,
		`select exists(select 1 from categories where slug = $1 and id <> $2)`, slug, exceptId).Scan(&taken)
	if err != nil {
		return false, err
	}
	return taken, nil
}
