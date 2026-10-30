package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

type MakerPieYearlyWithRankings struct {
	Year    int32  `json:"year"`
	MakerId string `json:"makerid"`
	Id      string `json:"id"`

	DisplayName      string           `json:"displayname"`
	Fresh            bool             `json:"fresh"`
	Categories       types.Categories `json:"categories"`
	ImageFile        string           `json:"image_file"`
	WebLink          string           `json:"web_link"`
	PackCount        int32            `json:"pack_count"`
	PackPriceInPence int32            `json:"pack_price_in_pence" `

	Pastry  float64 `json:"pastry"`
	Filling float64 `json:"filling"`
	Topping float64 `json:"topping"`
	Looks   float64 `json:"looks"`
	Value   float64 `json:"value"`

	Average float64 `json:"average"`

	Count int64 `json:"count"`
}

type RankingSummary struct {
	Pastry  float64 `json:"pastry"`
	Filling float64 `json:"filling"`
	Topping float64 `json:"topping"`
	Looks   float64 `json:"looks"`
	Value   float64 `json:"value"`

	Average float64 `json:"average"`

	Count int64 `json:"count"`
}

func (o *OperationWrapper) GetFilterableMakerPies(c context.Context, year int64, pageSize int64, zeroIdxPage int64, filters PieFilters) (*[]MakerPieYearlyWithRankings, error) {
	rows, err := o.querier.GetFilterableMakerPies(c, sqlcgen.GetFilterableMakerPiesParams{
		Year:          int32(year),
		BrandIds:      strings.Join(filters.BrandIds, ","),
		CategorySlugs: strings.Join(filters.CategorySlugs, ","),
		PageOffset:    int32(zeroIdxPage * pageSize),
		PageLimit:     int32(pageSize),
	})
	if err != nil {
		return nil, err
	}

	r := []MakerPieYearlyWithRankings{}
	for _, row := range rows {
		r = append(r, MakerPieYearlyWithRankings{
			Year:             row.Year,
			MakerId:          row.Makerid,
			Id:               row.ID,
			DisplayName:      row.Displayname,
			Fresh:            row.Fresh,
			Categories:       row.Categories,
			ImageFile:        row.ImageFile,
			WebLink:          row.WebLink,
			PackCount:        row.PackCount,
			PackPriceInPence: row.PackPriceInPence,
			Pastry:           row.Pastry,
			Filling:          row.Filling,
			Topping:          row.Topping,
			Looks:            row.Looks,
			Value:            row.Value,
			Average:          row.Avg,
			Count:            row.Count,
		})
	}
	return &r, nil
}

func (o *OperationWrapper) GetFilterableMakerPiesForUser(c context.Context, year int64, userid string, pageSize int64, zeroIdxPage int64, filters PieFilters) (*[]MakerPieYearlyWithRankings, error) {
	rows, err := o.querier.GetFilterableMakerPiesForUser(c, sqlcgen.GetFilterableMakerPiesForUserParams{
		Year:          int32(year),
		BrandIds:      strings.Join(filters.BrandIds, ","),
		CategorySlugs: strings.Join(filters.CategorySlugs, ","),
		PageOffset:    int32(zeroIdxPage * pageSize),
		PageLimit:     int32(pageSize),
		Userid:        userid,
	})
	if err != nil {
		return nil, err
	}

	r := []MakerPieYearlyWithRankings{}
	for _, row := range rows {
		r = append(r, MakerPieYearlyWithRankings{
			Year:             row.Year,
			MakerId:          row.Makerid,
			Id:               row.ID,
			DisplayName:      row.Displayname,
			Fresh:            row.Fresh,
			Categories:       row.Categories,
			ImageFile:        row.ImageFile,
			WebLink:          row.WebLink,
			PackCount:        row.PackCount,
			PackPriceInPence: row.PackPriceInPence,
			Pastry:           row.Pastry,
			Filling:          row.Filling,
			Topping:          row.Topping,
			Looks:            row.Looks,
			Value:            row.Value,
			Average:          row.Avg,
			// one user ranking per pie
			Count: 1,
		})
	}
	return &r, nil
}

func (o *OperationWrapper) GetTopMakerPie(c context.Context, activeYear int64) (*MakerPieYearlyWithRankings, error) {
	row, err := o.querier.GetTopMakerPie(c, int32(activeYear))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	r := MakerPieYearlyWithRankings{
		Year:             row.Year,
		MakerId:          row.Makerid,
		Id:               row.ID,
		DisplayName:      row.Displayname,
		Fresh:            row.Fresh,
		Categories:       row.Categories,
		ImageFile:        row.ImageFile,
		WebLink:          row.WebLink,
		PackCount:        row.PackCount,
		PackPriceInPence: row.PackPriceInPence,
		Pastry:           row.Pastry,
		Filling:          row.Filling,
		Topping:          row.Topping,
		Looks:            row.Looks,
		Value:            row.Value,
		Average:          row.Avg,
		Count:            row.Count,
	}
	return &r, nil
}

// GetLatestPieRanking implements Operations.
func (o *OperationWrapper) GetLatestPieRanking(c context.Context, activeYear int64) (*MakerPieYearlyWithRankings, *time.Time, error) {
	row, err := o.querier.GetLatestPieRanking(c, int32(activeYear))
	if err == pgx.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	r := MakerPieYearlyWithRankings{
		Year:             row.Year,
		MakerId:          row.Makerid,
		Id:               row.ID,
		DisplayName:      row.Displayname,
		Fresh:            row.Fresh,
		Categories:       row.Categories,
		ImageFile:        row.ImageFile,
		WebLink:          row.WebLink,
		PackCount:        row.PackCount,
		PackPriceInPence: row.PackPriceInPence,
		Pastry:           row.Pastry,
		Filling:          row.Filling,
		Topping:          row.Topping,
		Looks:            row.Looks,
		Value:            row.Value,
		Average:          row.Avg,
		Count:            row.Count,
	}
	tm := row.LastUpdated.Time
	return &r, &tm, nil
}

func (o *OperationWrapper) GetMakerPieCategoriesForYear(c context.Context, year int64, requiredSlugs []string, requiredBrand *string) (*[]types.Category, error) {
	brand := ""
	if requiredBrand != nil {
		brand = *requiredBrand
	}

	rows, err := o.querier.GetMakerPieCategoriesForYear(c, sqlcgen.GetMakerPieCategoriesForYearParams{
		Year:          int32(year),
		Brand:         brand,
		RequiredSlugs: strings.Join(requiredSlugs, ","),
	})
	if err != nil {
		return nil, err
	}

	r := []types.Category{}
	for _, row := range rows {
		cat, err := categoryFromIdstr(row.Idstr, row.Slug, row.Label)
		if err != nil {
			return nil, err
		}
		r = append(r, cat)
	}
	return &r, nil
}

func (o *OperationWrapper) GetMakerPieCategoriesForMakerPieOid(ctx context.Context, oid uuid.UUID) (*[]types.Category, error) {
	rows, err := o.querier.GetMakerPieCategoriesForMakerPieOid(ctx, oid.String())
	if err != nil {
		return nil, err
	}

	r := []types.Category{}
	for _, row := range rows {
		cat, err := categoryFromIdstr(row.Idstr, row.Slug, row.Label)
		if err != nil {
			return nil, err
		}
		r = append(r, cat)
	}
	return &r, nil
}

func (o *OperationWrapper) GetMakerPieCategories(ctx context.Context) (*[]types.Category, error) {
	rows, err := o.querier.GetMakerPieCategories(ctx)
	if err != nil {
		return nil, err
	}

	r := []types.Category{}
	for _, row := range rows {
		cat, err := categoryFromIdstr(row.Idstr, row.Slug, row.Label)
		if err != nil {
			return nil, err
		}
		r = append(r, cat)
	}
	return &r, nil
}

func (o *OperationWrapper) SetPieCategories(ctx context.Context, oid uuid.UUID, category_keys []string) error {
	// start a transaction
	tx, err := o.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback(ctx)
	}()

	err = o.SetPieCategoriesTx(ctx, tx, oid, category_keys)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

func (o *OperationWrapper) SetPieCategoriesTx(ctx context.Context, tx pgx.Tx, oid uuid.UUID, category_keys []string) error {
	qtx := sqlcgen.New(tx)
	if err := qtx.DeletePieCategoriesByPieOid(ctx, oid.String()); err != nil {
		return err
	}

	// Blank keys are meaningless (bulk CSV splits can produce them) and never
	// match a real slug, so drop them before anything else.
	keys := make([]string, 0, len(category_keys))
	for _, k := range category_keys {
		if k = strings.TrimSpace(k); k != "" {
			keys = append(keys, k)
		}
	}
	if len(keys) == 0 {
		return nil
	}

	// Slugs are looked up on tx, not the pool, so the function is correct
	// inside a transaction that has just inserted the categories or the pie.
	// Unknown slugs are reported instead of being silently dropped.
	rows, err := qtx.GetCategoriesBySlugs(ctx, keys)
	if err != nil {
		return err
	}

	catids := []string{}
	found := map[string]bool{}
	for _, row := range rows {
		catids = append(catids, row.Idstr)
		found[row.Slug.String] = true
	}

	var unknown []string
	for _, key := range keys {
		if !found[key] {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) > 0 {
		return fmt.Errorf("%w: %s", ErrUnknownCategories, strings.Join(unknown, ", "))
	}

	return qtx.InsertPieCategories(ctx, sqlcgen.InsertPieCategoriesParams{
		Oid:         oid.String(),
		CategoryIds: catids,
	})
}

// UpdateMakerPieLabels updates the labels array for a maker pie by oid.
func (o *OperationWrapper) UpdateMakerPieLabels(ctx context.Context, oid uuid.UUID, labels []string) error {
	return o.querier.UpdateMakerPieLabels(ctx, sqlcgen.UpdateMakerPieLabelsParams{
		Oid:    oid.String(),
		Labels: labels,
	})
}
