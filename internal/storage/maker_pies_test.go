package storage_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

// UpsertMakerPie is what the admin bulk import uses, so uploading the same
// year/brand/pie twice must update the row in place rather than fail on the
// (year, makerid, id) primary key.
func TestUpsertMakerPieReplacesExistingRow(t *testing.T) {
	ctx := context.Background()
	ops := newMigratedOps(t)

	_, err := ops.CreateMaker(ctx, types.Maker{Id: "greggs", Name: "Greggs"})
	require.NoError(t, err)
	_, err = ops.CreateCategory(ctx, types.Category{Slug: "sweet", Label: "Sweet"})
	require.NoError(t, err)
	_, err = ops.CreateCategory(ctx, types.Category{Slug: "savory", Label: "Savoury"})
	require.NoError(t, err)

	q := ops.GetQuerier()

	oidStr, err := q.UpsertMakerPie(ctx, sqlcgen.UpsertMakerPieParams{
		Year: 2026, Makerid: "greggs", ID: "steak-slice",
		Displayname: "Steak Slice", ImageFile: "/2026/greggs/steak-slice/a.jpg",
		WebLink: "https://example.com/a", PackCount: 1, PackPriceInPence: 299,
		Validated: true,
	})
	require.NoError(t, err)
	oid, err := uuid.Parse(oidStr)
	require.NoError(t, err)
	require.NoError(t, ops.SetPieCategories(ctx, oid, []string{"sweet"}))

	// re-uploading the same pie updates every column in place and keeps the
	// oid, so category links (and anything else keyed on oid) survive
	reupload, err := q.UpsertMakerPie(ctx, sqlcgen.UpsertMakerPieParams{
		Year: 2026, Makerid: "greggs", ID: "steak-slice",
		Displayname: "The Steak Slice", Fresh: true, ImageFile: "/2026/greggs/steak-slice/b.jpg",
		WebLink: "https://example.com/b", PackCount: 2, PackPriceInPence: 329,
		Validated: true,
	})
	require.NoError(t, err)
	assert.Equal(t, oidStr, reupload)

	pie, err := ops.GetMakerPieYearly(ctx, 2026, "greggs", "steak-slice")
	require.NoError(t, err)
	require.NotNil(t, pie)
	assert.Equal(t, "The Steak Slice", pie.DisplayName)
	assert.True(t, pie.Fresh)
	assert.Equal(t, "/2026/greggs/steak-slice/b.jpg", pie.ImageFile)
	assert.Equal(t, "https://example.com/b", pie.WebLink)
	assert.Equal(t, int32(2), pie.PackCount)
	assert.Equal(t, int32(329), pie.PackPriceInPence)

	// the bulk import re-sets categories per row, replacing the old set
	require.NoError(t, ops.SetPieCategories(ctx, oid, []string{"savory"}))
	pieCats, err := ops.GetMakerPieCategoriesForMakerPieOid(ctx, oid)
	require.NoError(t, err)
	require.Len(t, *pieCats, 1)
	assert.Equal(t, "savory", (*pieCats)[0].Slug)

	// a different pie id under the same brand still inserts a new row
	otherOid, err := q.UpsertMakerPie(ctx, sqlcgen.UpsertMakerPieParams{
		Year: 2026, Makerid: "greggs", ID: "sausage-roll",
		Displayname: "Sausage Roll", ImageFile: "/2026/greggs/sausage-roll/c.jpg",
		PackCount: 1, PackPriceInPence: 199, Validated: true,
	})
	require.NoError(t, err)
	assert.NotEqual(t, oidStr, otherOid)
}
