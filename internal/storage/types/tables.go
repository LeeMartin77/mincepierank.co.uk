package types

type MakerPieYearly struct {
	Year             int32    `json:"year" col:"year,primary"`
	MakerId          string   `json:"makerid" col:"makerid,primary"`
	Id               string   `json:"id" col:"id,primary"`
	DisplayName      string   `json:"displayname,omitempty" col:"displayname"`
	Fresh            bool     `json:"fresh,omitempty" col:"fresh"`
	Labels           []string `json:"labels,omitempty" col:"labels"`
	ImageFile        string   `json:"image_file,omitempty" col:"image_file"`
	WebLink          string   `json:"web_link,omitempty" col:"web_link"`
	PackCount        int32    `json:"pack_count,omitempty" col:"pack_count"`
	PackPriceInPence int32    `json:"pack_price_in_pence,omitempty" col:"pack_price_in_pence"`
}
