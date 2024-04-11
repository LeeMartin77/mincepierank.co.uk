package types

// `CREATE TABLE IF NOT EXISTS maker_pie_yearly
// (
// 	year int,
// 	makerid text,
// 	id text,
// 	displayname text,
// 	fresh boolean,
// 	labels text[],
// 	image_file text,
// 	web_link text,
// 	pack_count int,
// 	pack_price_in_pence int,
// 	PRIMARY KEY (year, makerId, id)
// );`,

// maker_pie_yearly
type MakerPieYearly struct {
	Year             int32    `json:"year" col:"year"`
	MakerId          string   `json:"makerid" col:"makerid"`
	Id               string   `json:"id" col:"id"`
	DisplayName      string   `json:"displayname,omitempty" col:"displayname"`
	Fresh            bool     `json:"fresh,omitempty" col:"fresh"`
	Labels           []string `json:"labels,omitempty" col:"labels"`
	ImageFile        string   `json:"image_file,omitempty" col:"image_file"`
	WebLink          string   `json:"web_link,omitempty" col:"web_link"`
	PackCount        int32    `json:"pack_count,omitempty" col:"pack_count"`
	PackPriceInPence int32    `json:"pack_price_in_pence,omitempty" col:"pack_price_in_pence"`
}
