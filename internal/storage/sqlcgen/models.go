// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlcgen

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Admin struct {
	ID string
}

type Category struct {
	ID    pgtype.UUID
	Slug  pgtype.Text
	Label pgtype.Text
}

type Config struct {
	Key   string
	Value pgtype.Text
}

type Maker struct {
	ID      string
	Name    pgtype.Text
	Logo    pgtype.Text
	Website pgtype.Text
}

type MakerPieCategory struct {
	MakerPieOid pgtype.UUID
	CategoryID  pgtype.UUID
}

type MakerPieRankingYearly struct {
	Year        int32
	Makerid     string
	Pieid       string
	Userid      string
	Pastry      pgtype.Int4
	Filling     pgtype.Int4
	Topping     pgtype.Int4
	Looks       pgtype.Int4
	Value       pgtype.Int4
	Notes       pgtype.Text
	LastUpdated pgtype.Timestamp
}

type MakerPieYearly struct {
	Year             int32
	Makerid          string
	ID               string
	Displayname      pgtype.Text
	Fresh            pgtype.Bool
	Labels           []string
	ImageFile        pgtype.Text
	WebLink          pgtype.Text
	PackCount        pgtype.Int4
	PackPriceInPence pgtype.Int4
	Validated        pgtype.Bool
	Oid              pgtype.UUID
}

type MakerPieYearlyCategory struct {
	Oid           pgtype.UUID
	Categories    []byte
	CategorySlugs interface{}
}

type UserPieRankingYearly struct {
	Year        int32
	Pieid       string
	Userid      string
	Pastry      pgtype.Int4
	Filling     pgtype.Int4
	Topping     pgtype.Int4
	Looks       pgtype.Int4
	Value       pgtype.Int4
	Notes       pgtype.Text
	LastUpdated pgtype.Timestamp
}

type UserPieYearly struct {
	Year             int32
	ID               string
	OwnerUserid      string
	Maker            pgtype.Text
	Location         pgtype.Text
	Displayname      pgtype.Text
	Fresh            pgtype.Bool
	Labels           []string
	ImageFile        pgtype.Text
	WebLink          pgtype.Text
	PackCount        pgtype.Int4
	PackPriceInPence pgtype.Int4
	Clean            pgtype.Bool
}