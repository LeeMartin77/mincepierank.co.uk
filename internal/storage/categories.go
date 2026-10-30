package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
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

// categoryFromIdstr builds a types.Category from the text id the sqlc
// queries return (ids are cast to text in the queries so generated rows
// hold plain Go types).
func categoryFromIdstr(idstr, slug, label string) (types.Category, error) {
	id, err := uuid.Parse(idstr)
	if err != nil {
		return types.Category{}, err
	}
	return types.Category{Id: id, Slug: slug, Label: label}, nil
}

// GetAllCategoriesAdmin returns every category with a total pie count across
// all years, ordered by label.
func (o *OperationWrapper) GetAllCategoriesAdmin(c context.Context) (*[]AdminCategoryRow, error) {
	rows, err := o.querier.ListCategoriesAdmin(c)
	if err != nil {
		return nil, err
	}

	r := []AdminCategoryRow{}
	for _, row := range rows {
		cat, err := categoryFromIdstr(row.Idstr, row.Slug, row.Label)
		if err != nil {
			return nil, err
		}
		r = append(r, AdminCategoryRow{Category: cat, PieCount: row.Piecount})
	}
	return &r, nil
}

// GetCategoryById returns the category or nil (no error) when not found.
func (o *OperationWrapper) GetCategoryById(c context.Context, id uuid.UUID) (*types.Category, error) {
	row, err := o.querier.GetCategoryById(c, id.String())
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	cat, err := categoryFromIdstr(row.Idstr, row.Slug, row.Label)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

// GetCategoryPieCount counts brand pies (all years) using the category, used
// to guard category deletes against orphaning pie references.
func (o *OperationWrapper) GetCategoryPieCount(c context.Context, id uuid.UUID) (int64, error) {
	return o.querier.GetCategoryPieCount(c, id.String())
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

	idstr, err := o.querier.CreateCategory(c, sqlcgen.CreateCategoryParams{
		Slug:  cat.Slug,
		Label: cat.Label,
	})
	if err != nil {
		return nil, err
	}

	created, err := categoryFromIdstr(idstr, cat.Slug, cat.Label)
	if err != nil {
		return nil, err
	}
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

	tag, err := o.querier.UpdateCategory(c, sqlcgen.UpdateCategoryParams{
		Slug:  cat.Slug,
		Label: cat.Label,
		ID:    cat.Id.String(),
	})
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

	qtx := sqlcgen.New(tx)
	if err := qtx.DeletePieCategoriesByCategory(c, id.String()); err != nil {
		return err
	}
	tag, err := qtx.DeleteCategory(c, id.String())
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
	return o.querier.CategorySlugTaken(c, sqlcgen.CategorySlugTakenParams{
		Slug:     slug,
		ExceptID: exceptId.String(),
	})
}
