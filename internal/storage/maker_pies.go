package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)


type MakerPieYearlyWithRankings struct {
	Year    int32  `json:"year"`
	MakerId string `json:"makerid"`
	Id      string `json:"id"`

	DisplayName      string   `json:"displayname"`
	Fresh            bool     `json:"fresh"`
	Labels           []string `json:"labels"`
	ImageFile        string   `json:"image_file"`
	WebLink          string   `json:"web_link"`
	PackCount        int32    `json:"pack_count"`
	PackPriceInPence int32    `json:"pack_price_in_pence" `

	Pastry  int8 `json:"pastry"`
	Filling int8 `json:"filling"`
	Topping int8 `json:"topping"`
	Looks   int8 `json:"looks"`
	Value   int8 `json:"value"`
	Average int8 `json:"average"`

	Count int64 `json:"count"`
}


func (o *OperationWrapper) GetTopMakerPie(c context.Context, activeYear int64) (*MakerPieYearlyWithRankings, error) {
	r := MakerPieYearlyWithRankings{}
	sql := `
	WITH aggregated_pies AS (
	    SELECT
	  		year,
	  		makerid,
	  		pieid,
	  			SUM(pastry)/COUNT(1) as pastry,
	       	SUM(filling)/COUNT(1)as filling,
	        SUM(topping)/COUNT(1) as topping,
	        SUM(looks)/COUNT(1) as looks,
	        SUM(value)/COUNT(1) as value,
	  		(
	        SUM(pastry)/COUNT(1) +
	       	SUM(filling)/COUNT(1) +
	        SUM(topping)/COUNT(1) +
	        SUM(looks)/COUNT(1) +
	        SUM(value)/COUNT(1)
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
		mpy.labels,
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
	`
	res := o.db.QueryRow(c, sql, activeYear)
	err := res.Scan(
		&r.Year,
		&r.MakerId,
		&r.Id,
		&r.DisplayName,
		&r.Fresh,
		&r.Labels,
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