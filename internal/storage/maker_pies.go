package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/rs/zerolog/log"
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
	r := []MakerPieYearlyWithRankings{}
	sql := `
	WITH pie_rankings AS (
	    SELECT
	  		year,
	  		makerid,
	  		pieid,
	  		SUM(pastry::float )/COUNT(1) as pastry,
	       	SUM(filling::float )/COUNT(1)as filling,
	        SUM(topping::float )/COUNT(1) as topping,
	        SUM(looks::float )/COUNT(1) as looks,
	        SUM(value::float )/COUNT(1) as value,
	  		(
		        SUM(pastry::float )/COUNT(1) +
		       	SUM(filling::float )/COUNT(1) +
		        SUM(topping::float )/COUNT(1) +
		        SUM(looks::float )/COUNT(1) +
		        SUM(value::float )/COUNT(1)
	        )/5 as avg,
	  		COUNT(1) as count
	    FROM maker_pie_ranking_yearly
	  	WHERE year = $1
	    GROUP BY year, makerid, pieid
	)
	SELECT
		mpy.year,
		mpy.makerid,
		mpy.id,
		mpy.displayname,
		mpy.fresh,
		cts.categories,
		mpy.image_file,
		mpy.web_link,
		mpy.pack_count,
		mpy.pack_price_in_pence,
	    COALESCE(tp.pastry, 0::float),
	    COALESCE(tp.filling, 0::float),
	    COALESCE(tp.topping, 0::float),
	    COALESCE(tp.looks, 0::float),
	    COALESCE(tp.value, 0::float),
	    COALESCE(tp.avg, 0::float),
	    COALESCE(tp.count, 0)
		FROM maker_pie_yearly mpy
		LEFT JOIN pie_rankings tp ON
	  	tp.year = mpy.year AND
	  	tp.makerid = mpy.makerid AND
	    tp.pieid = mpy.id
		INNER JOIN maker_pie_yearly_categories cts ON
		cts.oid = mpy.oid
		WHERE mpy.year = $1
		AND mpy.validated = true
		AND ($4 = '' OR mpy.makerid=ANY(string_to_array($4, ',')))
		AND ($5 = '' OR cts.category_slugs @> string_to_array($5, ','))
		ORDER BY tp.avg DESC NULLS LAST
		LIMIT $2
		OFFSET $3
	`
	rows, err := o.db.Query(c, sql, year, pageSize, zeroIdxPage*pageSize, strings.Join(filters.BrandIds, ","), strings.Join(filters.CategorySlugs, ","))
	if err == pgx.ErrNoRows {
		return &r, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sigl := MakerPieYearlyWithRankings{}

		err := rows.Scan(
			&sigl.Year,
			&sigl.MakerId,
			&sigl.Id,
			&sigl.DisplayName,
			&sigl.Fresh,
			&sigl.Categories,
			&sigl.ImageFile,
			&sigl.WebLink,
			&sigl.PackCount,
			&sigl.PackPriceInPence,

			&sigl.Pastry,
			&sigl.Filling,
			&sigl.Topping,
			&sigl.Looks,
			&sigl.Value,
			&sigl.Average,
			&sigl.Count,
		)
		if err != nil {
			log.Error().Err(err).Msg("Scan failed in get categories for year")
			continue
		}
		r = append(r, sigl)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during categories for year iteration")
	}
	return &r, nil
}

func (o *OperationWrapper) GetFilterableMakerPiesForUser(c context.Context, year int64, userid string, pageSize int64, zeroIdxPage int64, filters PieFilters) (*[]MakerPieYearlyWithRankings, error) {
	r := []MakerPieYearlyWithRankings{}
	sql := `
	SELECT
		mpy.year,
		mpy.makerid,
		mpy.id,
		mpy.displayname,
		mpy.fresh,
		cts.categories,
		mpy.image_file,
		mpy.web_link,
		mpy.pack_count,
		mpy.pack_price_in_pence,
	    tp.pastry,
	    tp.filling,
	    tp.topping,
	    tp.looks,
	    tp.value,
		(
	        tp.pastry +
	       	tp.filling +
	        tp.topping +
	        tp.looks +
	        tp.value
        )/5 as avg
		FROM maker_pie_yearly mpy
		INNER JOIN maker_pie_ranking_yearly tp ON
		  	tp.year = mpy.year AND
		  	tp.makerid = mpy.makerid AND
		    tp.pieid = mpy.id AND
			tp.userid = $2
		INNER JOIN maker_pie_yearly_categories cts ON
		cts.oid = mpy.oid
		WHERE mpy.year = $1
		AND mpy.validated = true
		AND ($5 = '' OR mpy.makerid=ANY(string_to_array($5, ',')))
		AND ($6 = '' OR cts.category_slugs @> string_to_array($6, ','))
		ORDER BY (
	        tp.pastry +
	       	tp.filling +
	        tp.topping +
	        tp.looks +
	        tp.value
        ) DESC NULLS LAST
		LIMIT $3
		OFFSET $4
	`
	rows, err := o.db.Query(c, sql, year, userid, pageSize, zeroIdxPage*pageSize, strings.Join(filters.BrandIds, ","), strings.Join(filters.CategorySlugs, ","))
	if err == pgx.ErrNoRows {
		return &r, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sigl := MakerPieYearlyWithRankings{}

		err := rows.Scan(
			&sigl.Year,
			&sigl.MakerId,
			&sigl.Id,
			&sigl.DisplayName,
			&sigl.Fresh,
			&sigl.Categories,
			&sigl.ImageFile,
			&sigl.WebLink,
			&sigl.PackCount,
			&sigl.PackPriceInPence,

			&sigl.Pastry,
			&sigl.Filling,
			&sigl.Topping,
			&sigl.Looks,
			&sigl.Value,
			&sigl.Average,
		)
		if err != nil {
			log.Error().Err(err).Msg("Scan failed in getting ranked pies for user")
			continue
		}
		sigl.Count = 1
		r = append(r, sigl)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during categories for year iteration")
	}
	return &r, nil
}

func (o *OperationWrapper) GetTopMakerPie(c context.Context, activeYear int64) (*MakerPieYearlyWithRankings, error) {
	r := MakerPieYearlyWithRankings{}
	sql := `
	WITH aggregated_pies AS (
	    SELECT
	  		year,
	  		makerid,
	  		pieid,
	  			SUM(pastry::float )/COUNT(1) as pastry,
	       	SUM(filling::float )/COUNT(1)as filling,
	        SUM(topping::float )/COUNT(1) as topping,
	        SUM(looks::float )/COUNT(1) as looks,
	        SUM(value::float )/COUNT(1) as value,
	  		(
	        SUM(pastry::float )/COUNT(1) +
	       	SUM(filling::float )/COUNT(1) +
	        SUM(topping::float )/COUNT(1) +
	        SUM(looks::float )/COUNT(1) +
	        SUM(value::float )/COUNT(1)
	      )/5 as avg,
	  		COUNT(1) as count
	    FROM maker_pie_ranking_yearly
	  	WHERE year = $1
	    GROUP BY year, makerid, pieid
	), top_pie AS (
	    SELECT *
	    FROM aggregated_pies
	  	ORDER BY avg DESC
	  	LIMIT 1
	)
	SELECT
		mpy.year,
		mpy.makerid,
		mpy.id,
		mpy.displayname,
		mpy.fresh,
		cts.categories,
		mpy.image_file,
		mpy.web_link,
		mpy.pack_count,
		mpy.pack_price_in_pence,
	    tp.pastry,
	    tp.filling,
	    tp.topping,
	    tp.looks,
	    tp.value,
	    tp.avg,
	    tp.count
		FROM maker_pie_yearly mpy
		INNER JOIN top_pie tp ON
	  	tp.year = mpy.year AND
	  	tp.makerid = mpy.makerid AND
	    tp.pieid = mpy.id
		INNER JOIN maker_pie_yearly_categories cts ON
		cts.oid = mpy.oid
	`
	res := o.db.QueryRow(c, sql, activeYear)
	err := res.Scan(
		&r.Year,
		&r.MakerId,
		&r.Id,
		&r.DisplayName,
		&r.Fresh,
		&r.Categories,
		&r.ImageFile,
		&r.WebLink,
		&r.PackCount,
		&r.PackPriceInPence,

		&r.Pastry,
		&r.Filling,
		&r.Topping,
		&r.Looks,
		&r.Value,
		&r.Average,
		&r.Count,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (o *OperationWrapper) GetMakerPieCategoriesForYear(c context.Context, year int64, requiredSlugs []string, requiredBrand *string) (*[]types.Category, error) {
	brndprm := ""
	if requiredBrand != nil {
		brndprm = *requiredBrand
	}
	r := []types.Category{}
	sql := `
	SELECT distinct c.id, c.slug, c.label
	FROM categories c
	inner join maker_pie_categories mpc on
		mpc.category_id = c.id
	inner join maker_pie_yearly mpy on
		mpc.maker_pie_oid = mpy.oid and
		mpy.year = $1
	INNER JOIN maker_pie_yearly_categories cts ON
		cts.oid = mpy.oid
	WHERE
		($2 = '' OR mpy.makerid=$2)
		AND ($3 = '' OR cts.category_slugs @> string_to_array($3, ','))
	`
	rows, err := o.db.Query(c, sql, year, brndprm, strings.Join(requiredSlugs, ","))
	if err == pgx.ErrNoRows {
		return &r, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ctg types.Category
		if err := rows.Scan(&ctg.Id, &ctg.Slug, &ctg.Label); err != nil {
			log.Error().Err(err).Msg("Scan failed in get categories for year")
			continue
		}
		r = append(r, ctg)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during categories for year iteration")
	}
	return &r, nil
}
func (o *OperationWrapper) GetMakerPieCategoriesForMakerPieOid(ctx context.Context, oid uuid.UUID) (*[]types.Category, error) {
	r := []types.Category{}
	sql := `
	SELECT c.id, c.slug, c.label
	FROM categories c
	inner join maker_pie_categories mpc on
		mpc.category_id = c.id
	WHERE
		mpc.maker_pie_oid = $1
	`
	rows, err := o.db.Query(ctx, sql, oid)
	if err == pgx.ErrNoRows {
		return &r, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ctg types.Category
		if err := rows.Scan(&ctg.Id, &ctg.Slug, &ctg.Label); err != nil {
			log.Error().Err(err).Msg("Scan failed in get categories for year")
			continue
		}
		r = append(r, ctg)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during categories for year iteration")
	}
	return &r, nil
}

func (o *OperationWrapper) GetMakerPieCategories(ctx context.Context) (*[]types.Category, error) {
	r := []types.Category{}
	sql := `
	SELECT c.id, c.slug, c.label
	FROM categories c;
	`
	rows, err := o.db.Query(ctx, sql)
	if err == pgx.ErrNoRows {
		return &r, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ctg types.Category
		if err := rows.Scan(&ctg.Id, &ctg.Slug, &ctg.Label); err != nil {
			log.Error().Err(err).Msg("Scan failed in get categories")
			continue
		}
		r = append(r, ctg)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during categories")
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
	_, err := tx.Exec(ctx, `DELETE from maker_pie_categories WHERE maker_pie_oid = $1`, oid)
	if err != nil {
		return err
	}

	catids := []uuid.UUID{}
	sql := `
	SELECT c.id, c.slug, c.label
	FROM categories c
	WHERE c.slug = ANY($1::text[])
	`
	rows, err := o.db.Query(ctx, sql, category_keys)
	if err == pgx.ErrNoRows || len(category_keys) == 0 {
		return nil
	}
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var ctg types.Category
		if err := rows.Scan(&ctg.Id, &ctg.Slug, &ctg.Label); err != nil {
			log.Error().Err(err).Msg("Scan failed")
			continue
		}
		catids = append(catids, ctg.Id)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error during cat id retreival")
	}

	insql := `INSERT INTO maker_pie_categories (maker_pie_oid, category_id) VALUES `
	argidx := 1
	args := []any{}
	for _, cid := range catids {
		insql += fmt.Sprintf(" ($%d, $%d),", argidx, argidx+1)
		argidx += 2
		args = append(args, oid)
		args = append(args, cid)
	}
	insql = insql[:len(insql)-1]
	insql += ";"
	_, err = tx.Exec(ctx, insql, args...)
	if err != nil {
		return err
	}
	return nil
}
