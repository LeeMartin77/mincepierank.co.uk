package types

import (
	"time"

	"github.com/google/uuid"
)

type Maker struct {
	Id string `json:"id" col:"id,primary"`

	Name    string `json:"name,omitempty" col:"name"`
	Logo    string `json:"logo,omitempty" col:"logo"`
	Website string `json:"website,omitempty" col:"website"`
}

type Category struct {
	Id    uuid.UUID `json:"id" col:"id,primary"`
	Slug  string    `json:"slug" col:"slug"`
	Label string    `json:"label" col:"label"`
}

type MakerPieYearly struct {
	Year    int32  `json:"year" col:"year,primary"`
	MakerId string `json:"makerid" col:"makerid,primary"`
	Id      string `json:"id" col:"id,primary"`

	DisplayName      string   `json:"displayname,omitempty" col:"displayname"`
	Fresh            bool     `json:"fresh,omitempty" col:"fresh"`
	Labels           []string `json:"labels,omitempty" col:"labels"`
	ImageFile        string   `json:"image_file,omitempty" col:"image_file"`
	WebLink          string   `json:"web_link,omitempty" col:"web_link"`
	PackCount        int32    `json:"pack_count,omitempty" col:"pack_count"`
	PackPriceInPence int32    `json:"pack_price_in_pence,omitempty" col:"pack_price_in_pence"`
}

type MakerPieRankingYearly struct {
	Year    int32  `json:"year" col:"year,primary"`
	MakerId string `json:"makerid" col:"makerid,primary"`
	Id      string `json:"id" col:"id,primary"`
	UserId  string `json:"userid" col:"userid,primary"`

	Pastry  int8 `json:"pastry" col:"pastry"`
	Filling int8 `json:"filling" col:"filling"`
	Topping int8 `json:"topping" col:"topping"`
	Looks   int8 `json:"looks" col:"looks"`
	Value   int8 `json:"value" col:"value"`

	Notes       string    `json:"notes" col:"notes"`
	LastUpdated time.Time `json:"last_updated" col:"last_updated"`
}

type UserPieYearly struct {
	Year        int32  `json:"year" col:"year,primary"`
	Id          string `json:"id" col:"id,primary"`
	OwnerUserid string `json:"owner_userid" col:"owner_userid,primary"`

	Maker    string `json:"maker,omitempty" col:"maker"`
	Location string `json:"location,omitempty" col:"location"`

	DisplayName      string   `json:"displayname,omitempty" col:"displayname"`
	Fresh            bool     `json:"fresh,omitempty" col:"fresh"`
	Labels           []string `json:"labels,omitempty" col:"labels"`
	ImageFile        string   `json:"image_file,omitempty" col:"image_file"`
	WebLink          string   `json:"web_link,omitempty" col:"web_link"`
	PackCount        int32    `json:"pack_count,omitempty" col:"pack_count"`
	PackPriceInPence int32    `json:"pack_price_in_pence,omitempty" col:"pack_price_in_pence"`

	Clean bool `json:"clean,omitempty" col:"clean"`
}

type UserPieRankingYearly struct {
	Year   int32  `json:"year" col:"year,primary"`
	PieId  string `json:"pieid" col:"pieid,primary"`
	UserId string `json:"userid" col:"userid,primary"`

	Pastry  int8 `json:"pastry" col:"pastry"`
	Filling int8 `json:"filling" col:"filling"`
	Topping int8 `json:"topping" col:"topping"`
	Looks   int8 `json:"looks" col:"looks"`
	Value   int8 `json:"value" col:"value"`

	Notes       string    `json:"notes" col:"notes"`
	LastUpdated time.Time `json:"last_updated" col:"last_updated"`
}

type Config struct {
	Key   string `json:"key" col:"key,primary"`
	Value string `json:"value" col:"value"`
}

type Admins struct {
	Id string `json:"id" col:"id,primary"`
}
