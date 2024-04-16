// Code generated by generator, DO NOT EDIT.
package generated

import (
	"context"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	types "github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

// Insert new 'Maker' into table 'maker' - duplicate ids is a no-op
func MakerCreate(ctx context.Context, pg *pgxpool.Pool, c types.Maker) (*types.Maker, error) {
	args := []interface{}{
		c.Id,
		c.Name,
		c.Logo,
		c.Website,
	}
	sql := "INSERT INTO maker (id,name,logo,website) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING"
	_, err := pg.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Update 'Maker' in table 'maker' based on id columns
func MakerUpdate(ctx context.Context, pg *pgxpool.Pool, u types.Maker) (*types.Maker, error) {
	parameters := []interface{}{}
	parameters = append(parameters, u.Id)
	parameters = append(parameters, u.Name)
	parameters = append(parameters, u.Logo)
	parameters = append(parameters, u.Website)
	sql := "UPDATE maker SET name=$2,logo=$3,website=$4 WHERE id=$1"
	_, err := pg.Exec(ctx, sql, parameters...)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Read 'Maker' in table 'maker' based on id columns - nil but no error if not found
func MakerRead(ctx context.Context, pg *pgxpool.Pool, Id string) (*types.Maker, error) {
	r := types.Maker{}
	identifiers := []interface{}{}
	identifiers = append(identifiers, Id)
	sql := "SELECT id, name, logo, website FROM maker  WHERE id = $1"
	res := pg.QueryRow(ctx, sql, identifiers...)
	err := res.Scan(&r.Id, &r.Name, &r.Logo, &r.Website)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Deletes 'Maker' in table 'maker' based on id columns - no error if not found
func MakerDelete(ctx context.Context, pg *pgxpool.Pool, Id string) error {
	identifiers := []interface{}{}
	identifiers = append(identifiers, Id)
	sql := "DELETE FROM maker WHERE id = $1"
	_, err := pg.Exec(ctx, sql, identifiers...)
	if err != nil {
		return err
	}
	return nil
}
