package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/generated"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

func TestGeneratedFunctions(t *testing.T) {
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
	assert.NoError(t, err)

	ops, err := storage.NewOperations(connStr)
	assert.NoError(t, err)

	err = ops.Migrate(ctx)
	assert.NoError(t, err)

	insert := types.MakerPieYearly{
		Year:        2014,
		MakerId:     "Bramwell",
		Id:          "super-delicious",
		DisplayName: "Super Delicious Pies",
	}

	_, err = generated.MakerPieYearlyCreate(ctx, ops.GetDatabase(), insert)
	assert.NoError(t, err)

	res, err := generated.MakerPieYearlyRead(ctx, ops.GetDatabase(), 2014, "Bramwell", "super-delicious")
	assert.NoError(t, err)
	assert.Equal(t, "Super Delicious Pies", res.DisplayName)

	res.DisplayName = "Not so delicious really"

	_, err = generated.MakerPieYearlyUpdate(ctx, ops.GetDatabase(), *res)
	assert.NoError(t, err)

	res2, err := generated.MakerPieYearlyRead(ctx, ops.GetDatabase(), 2014, "Bramwell", "super-delicious")
	assert.NoError(t, err)
	assert.Equal(t, "Not so delicious really", res2.DisplayName)

	err = generated.MakerPieYearlyDelete(ctx, ops.GetDatabase(), 2014, "Bramwell", "super-delicious")
	assert.NoError(t, err)

	res3, err := generated.MakerPieYearlyRead(ctx, ops.GetDatabase(), 2014, "Bramwell", "super-delicious")
	assert.NoError(t, err)
	assert.Nil(t, res3)
}
