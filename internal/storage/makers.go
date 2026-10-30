package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

var (
	ErrBrandAlreadyExists = errors.New("a brand with that id already exists")
	ErrBrandNotFound      = errors.New("brand not found")
)

// AdminMakerRow is a brand plus how many brand pies reference it, for the
// admin brands list and the delete guard.
type AdminMakerRow struct {
	types.Maker
	PieCount int64
}

// GetAllMakers returns every brand (not just those with pies this year),
// ordered by name, with a total pie count across all years.
func (o *OperationWrapper) GetAllMakers(c context.Context) (*[]AdminMakerRow, error) {
	rows, err := o.querier.ListMakersAdmin(c)
	if err != nil {
		return nil, err
	}

	r := []AdminMakerRow{}
	for _, row := range rows {
		r = append(r, AdminMakerRow{
			Maker: types.Maker{
				Id:      row.ID,
				Name:    row.Name,
				Logo:    row.Logo,
				Website: row.Website,
			},
			PieCount: row.Piecount,
		})
	}
	return &r, nil
}

// GetMakerPieCount counts brand pies (all years) referencing the maker id.
// Used to guard brand deletes against orphaning pies.
func (o *OperationWrapper) GetMakerPieCount(c context.Context, makerid string) (int64, error) {
	return o.querier.GetMakerPieCount(c, makerid)
}

// CreateMaker inserts a new brand. Returns ErrBrandAlreadyExists when the id
// is taken (the generated insert is ON CONFLICT DO NOTHING, so the existence
// check happens here first).
func (o *OperationWrapper) CreateMaker(c context.Context, m types.Maker) (*types.Maker, error) {
	_, err := o.querier.GetMakerByID(c, m.Id)
	if err == nil {
		return nil, ErrBrandAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	err = o.querier.CreateMaker(c, sqlcgen.CreateMakerParams{
		ID:      m.Id,
		Name:    m.Name,
		Logo:    m.Logo,
		Website: m.Website,
	})
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// UpdateMaker updates name/logo/website for an existing brand.
// Returns ErrBrandNotFound when the brand does not exist.
func (o *OperationWrapper) UpdateMaker(c context.Context, m types.Maker) (*types.Maker, error) {
	_, err := o.querier.GetMakerByID(c, m.Id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrBrandNotFound
	}
	if err != nil {
		return nil, err
	}

	_, err = o.querier.UpdateMakerByID(c, sqlcgen.UpdateMakerByIDParams{
		Name:    m.Name,
		Logo:    m.Logo,
		Website: m.Website,
		ID:      m.Id,
	})
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// DeleteMaker removes a brand. Callers should check GetMakerPieCount first so
// pies are not silently orphaned (brand pages inner join maker).
func (o *OperationWrapper) DeleteMaker(c context.Context, makerid string) error {
	return o.querier.DeleteMakerByID(c, makerid)
}

// GetMakersForYear returns the brands with pies in the year, ordered by how
// many pies they have. The pie count is only used for ordering.
func (o *OperationWrapper) GetMakersForYear(c context.Context, activeYear int64) (*[]types.Maker, error) {
	rows, err := o.querier.GetMakersForYear(c, int32(activeYear))
	if err != nil {
		return nil, err
	}

	r := []types.Maker{}
	for _, row := range rows {
		r = append(r, types.Maker{
			Id:      row.ID,
			Name:    row.Name,
			Logo:    row.Logo,
			Website: row.Website,
		})
	}
	return &r, nil
}
