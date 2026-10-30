package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

func newMigratedOps(t *testing.T) storage.Operations {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:13.13"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	ops, err := storage.NewOperations(connStr)
	require.NoError(t, err)
	t.Cleanup(ops.Shutdown)

	require.NoError(t, ops.Migrate(ctx))
	return ops
}

func TestBrandCrud(t *testing.T) {
	ctx := context.Background()
	ops := newMigratedOps(t)

	// create
	created, err := ops.CreateMaker(ctx, types.Maker{
		Id: "bramwell", Name: "Bramwell's", Website: "https://bramwells.example",
	})
	require.NoError(t, err)
	assert.Equal(t, "bramwell", created.Id)

	// duplicate create is refused
	_, err = ops.CreateMaker(ctx, types.Maker{Id: "bramwell", Name: "Other"})
	assert.ErrorIs(t, err, storage.ErrBrandAlreadyExists)

	// read
	got, err := ops.GetMaker(ctx, "bramwell")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Bramwell's", got.Name)

	// update
	updated, err := ops.UpdateMaker(ctx, types.Maker{
		Id: "bramwell", Name: "Bramwell Pies", Logo: "/brands/bramwell/logo.png", Website: "https://bramwells.example",
	})
	require.NoError(t, err)
	assert.Equal(t, "Bramwell Pies", updated.Name)

	got, err = ops.GetMaker(ctx, "bramwell")
	require.NoError(t, err)
	assert.Equal(t, "/brands/bramwell/logo.png", got.Logo)

	// update of a missing brand is refused
	_, err = ops.UpdateMaker(ctx, types.Maker{Id: "nosuchbrand", Name: "x"})
	assert.ErrorIs(t, err, storage.ErrBrandNotFound)

	// a second brand, out of name order, to check list ordering
	_, err = ops.CreateMaker(ctx, types.Maker{Id: "aaa-pies", Name: "AAA Pies"})
	require.NoError(t, err)

	// brand with a pie: counts show up on the list and in the delete guard
	q := ops.GetQuerier()
	_, err = q.CreateMakerPie(ctx, sqlcgen.CreateMakerPieParams{
		Year: 2026, Makerid: "aaa-pies", ID: "flagship", Displayname: "AAA Flagship",
	})
	require.NoError(t, err)

	count, err := ops.GetMakerPieCount(ctx, "aaa-pies")
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	count, err = ops.GetMakerPieCount(ctx, "bramwell")
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	brands, err := ops.GetAllMakers(ctx)
	require.NoError(t, err)
	require.Len(t, *brands, 2)
	assert.Equal(t, "aaa-pies", (*brands)[0].Id, "list is ordered by name")
	assert.Equal(t, "bramwell", (*brands)[1].Id)
	assert.Equal(t, int64(1), (*brands)[0].PieCount)
	assert.Equal(t, int64(0), (*brands)[1].PieCount)

	// delete unreferenced brand works; the row is gone
	require.NoError(t, ops.DeleteMaker(ctx, "bramwell"))
	got, err = ops.GetMaker(ctx, "bramwell")
	require.NoError(t, err)
	assert.Nil(t, got, "deleted brand should read back as nil")
}

func TestGetAllMakersEmpty(t *testing.T) {
	ctx := context.Background()
	ops := newMigratedOps(t)

	brands, err := ops.GetAllMakers(ctx)
	require.NoError(t, err)
	assert.Empty(t, *brands)
}
