// Code generated by generator, DO NOT EDIT.
package generated

import (
	"context"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	types "github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

// Insert new 'UserPieYearly' into table 'user_pie_yearly' - duplicate ids is a no-op
func UserPieYearlyCreate(ctx context.Context, pg *pgxpool.Pool, c types.UserPieYearly) (*types.UserPieYearly, error) {
	args := []interface{}{
		c.Year,
		c.Id,
		c.OwnerUserid,
		c.Maker,
		c.Location,
		c.DisplayName,
		c.Fresh,
		c.Labels,
		c.ImageFile,
		c.WebLink,
		c.PackCount,
		c.PackPriceInPence,
		c.Clean,
	}
	sql := "INSERT INTO user_pie_yearly (year,id,owner_userid,maker,location,displayname,fresh,labels,image_file,web_link,pack_count,pack_price_in_pence,clean) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) ON CONFLICT DO NOTHING"
	_, err := pg.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Update 'UserPieYearly' in table 'user_pie_yearly' based on id columns
func UserPieYearlyUpdate(ctx context.Context, pg *pgxpool.Pool, u types.UserPieYearly) (*types.UserPieYearly, error) {
	parameters := []interface{}{}
	parameters = append(parameters, u.Year)
	parameters = append(parameters, u.Id)
	parameters = append(parameters, u.OwnerUserid)
	parameters = append(parameters, u.Maker)
	parameters = append(parameters, u.Location)
	parameters = append(parameters, u.DisplayName)
	parameters = append(parameters, u.Fresh)
	parameters = append(parameters, u.Labels)
	parameters = append(parameters, u.ImageFile)
	parameters = append(parameters, u.WebLink)
	parameters = append(parameters, u.PackCount)
	parameters = append(parameters, u.PackPriceInPence)
	parameters = append(parameters, u.Clean)
	sql := "UPDATE user_pie_yearly SET maker=$4,location=$5,displayname=$6,fresh=$7,labels=$8,image_file=$9,web_link=$10,pack_count=$11,pack_price_in_pence=$12,clean=$13 WHERE year=$1 AND id=$2 AND owner_userid=$3"
	_, err := pg.Exec(ctx, sql, parameters...)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Read 'UserPieYearly' in table 'user_pie_yearly' based on id columns - nil but no error if not found
func UserPieYearlyRead(ctx context.Context, pg *pgxpool.Pool, Year int32, Id string, OwnerUserid string) (*types.UserPieYearly, error) {
	r := types.UserPieYearly{}
	identifiers := []interface{}{}
	identifiers = append(identifiers, Year)
	identifiers = append(identifiers, Id)
	identifiers = append(identifiers, OwnerUserid)
	sql := "SELECT year, id, owner_userid, maker, location, displayname, fresh, labels, image_file, web_link, pack_count, pack_price_in_pence, clean FROM user_pie_yearly  WHERE year = $1 AND id = $2 AND owner_userid = $3"
	res := pg.QueryRow(ctx, sql, identifiers...)
	err := res.Scan(&r.Year, &r.Id, &r.OwnerUserid, &r.Maker, &r.Location, &r.DisplayName, &r.Fresh, &r.Labels, &r.ImageFile, &r.WebLink, &r.PackCount, &r.PackPriceInPence, &r.Clean)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Deletes 'UserPieYearly' in table 'user_pie_yearly' based on id columns - no error if not found
func UserPieYearlyDelete(ctx context.Context, pg *pgxpool.Pool, Year int32, Id string, OwnerUserid string) error {
	identifiers := []interface{}{}
	identifiers = append(identifiers, Year)
	identifiers = append(identifiers, Id)
	identifiers = append(identifiers, OwnerUserid)
	sql := "DELETE FROM user_pie_yearly WHERE year = $1 AND id = $2 AND owner_userid = $3"
	_, err := pg.Exec(ctx, sql, identifiers...)
	if err != nil {
		return err
	}
	return nil
}
