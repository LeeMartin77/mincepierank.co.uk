package storage_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

func TestCategoryCrud(t *testing.T) {
	ctx := context.Background()
	ops := newMigratedOps(t)

	// create
	created, err := ops.CreateCategory(ctx, types.Category{Slug: "vegetarian", Label: "Vegetarian"})
	require.NoError(t, err)
	assert.Equal(t, "vegetarian", created.Slug)
	assert.NotEqual(t, uuid.Nil, created.Id)

	// duplicate slug is refused (no unique constraint, the check is in code)
	_, err = ops.CreateCategory(ctx, types.Category{Slug: "vegetarian", Label: "Other"})
	assert.ErrorIs(t, err, storage.ErrCategoryAlreadyExists)

	// read
	got, err := ops.GetCategoryById(ctx, created.Id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Vegetarian", got.Label)

	// missing read is nil, not an error
	missing, err := ops.GetCategoryById(ctx, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, missing)

	// update
	updated, err := ops.UpdateCategory(ctx, types.Category{Id: created.Id, Slug: "vegetarian", Label: "Veggie"})
	require.NoError(t, err)
	assert.Equal(t, "Veggie", updated.Label)

	got, err = ops.GetCategoryById(ctx, created.Id)
	require.NoError(t, err)
	assert.Equal(t, "Veggie", got.Label)

	// updating a missing category errors
	_, err = ops.UpdateCategory(ctx, types.Category{Id: uuid.New(), Slug: "nope", Label: "Nope"})
	assert.ErrorIs(t, err, storage.ErrCategoryNotFound)

	// list with usage counts
	_, err = ops.CreateCategory(ctx, types.Category{Slug: "sweet", Label: "Sweet"})
	require.NoError(t, err)
	list, err := ops.GetAllCategoriesAdmin(ctx)
	require.NoError(t, err)
	require.Len(t, *list, 2)
}

func TestCategoryDeleteGuardAndPieUsage(t *testing.T) {
	ctx := context.Background()
	ops := newMigratedOps(t)

	cat, err := ops.CreateCategory(ctx, types.Category{Slug: "sweet", Label: "Sweet"})
	require.NoError(t, err)

	// a brand pie using the category
	_, err = ops.CreateMaker(ctx, types.Maker{Id: "greggs", Name: "Greggs"})
	require.NoError(t, err)
	oidStr, err := ops.GetQuerier().CreateMakerPie(ctx, sqlcgen.CreateMakerPieParams{
		Year: 2026, Makerid: "greggs", ID: "steak-slice",
		Displayname: "Steak Slice", ImageFile: "/2026/greggs/steak-slice/x.jpg",
	})
	require.NoError(t, err)
	oid, err := uuid.Parse(oidStr)
	require.NoError(t, err)

	require.NoError(t, ops.SetPieCategories(ctx, oid, []string{"sweet"}))

	count, err := ops.GetCategoryPieCount(ctx, cat.Id)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	list, err := ops.GetAllCategoriesAdmin(ctx)
	require.NoError(t, err)
	require.Len(t, *list, 1)
	assert.Equal(t, int64(1), (*list)[0].PieCount)

	// delete removes the category and its references
	require.NoError(t, ops.DeleteCategory(ctx, cat.Id))
	count, err = ops.GetCategoryPieCount(ctx, cat.Id)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	pieCats, err := ops.GetMakerPieCategoriesForMakerPieOid(ctx, oid)
	require.NoError(t, err)
	assert.Empty(t, *pieCats)

	// deleting again is a not-found
	assert.ErrorIs(t, ops.DeleteCategory(ctx, cat.Id), storage.ErrCategoryNotFound)
}

func TestSetPieCategoriesRejectsUnknownSlugs(t *testing.T) {
	ctx := context.Background()
	ops := newMigratedOps(t)

	_, err := ops.CreateCategory(ctx, types.Category{Slug: "sweet", Label: "Sweet"})
	require.NoError(t, err)
	_, err = ops.CreateMaker(ctx, types.Maker{Id: "greggs", Name: "Greggs"})
	require.NoError(t, err)
	oidStr, err := ops.GetQuerier().CreateMakerPie(ctx, sqlcgen.CreateMakerPieParams{
		Year: 2026, Makerid: "greggs", ID: "steak-slice",
		Displayname: "Steak Slice", ImageFile: "/2026/greggs/steak-slice/x.jpg",
	})
	require.NoError(t, err)
	oid, err := uuid.Parse(oidStr)
	require.NoError(t, err)

	// a typo is an error naming the slug, not a silent drop
	err = ops.SetPieCategories(ctx, oid, []string{"sweet", "sweat"})
	assert.ErrorIs(t, err, storage.ErrUnknownCategories)
	assert.Contains(t, err.Error(), "sweat")

	// the failed set left the pie untouched (still no categories)
	pieCats, err := ops.GetMakerPieCategoriesForMakerPieOid(ctx, oid)
	require.NoError(t, err)
	assert.Empty(t, *pieCats)

	// a good set works, and an empty set clears without error
	require.NoError(t, ops.SetPieCategories(ctx, oid, []string{"sweet", " ", ""}))
	pieCats, err = ops.GetMakerPieCategoriesForMakerPieOid(ctx, oid)
	require.NoError(t, err)
	require.Len(t, *pieCats, 1)

	require.NoError(t, ops.SetPieCategories(ctx, oid, []string{}))
	pieCats, err = ops.GetMakerPieCategoriesForMakerPieOid(ctx, oid)
	require.NoError(t, err)
	assert.Empty(t, *pieCats)
}
