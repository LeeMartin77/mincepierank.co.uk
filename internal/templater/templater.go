package templater

import (
	"io"
	"strings"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

type PageDataHead struct {
	Title            string
	Description      string
	Keywords         string
	ShowCookieNotice bool
	ShowNotice       bool
	ReadOnly         bool
	HasNotice        bool
	Notice           string
	ShowBreadcrumb   bool
	MenuSettings     MenuSettings
	Breadcrumbs      []Link
}

type MenuSettings struct {
	ActiveYear int64
	LoggedIn   bool
}

type Link struct {
	Label string
	URL   string
}

type PageDataFoot struct {
	Title string
}

type PageData struct {
	Head     PageDataHead
	PageData interface{}
	Foot     PageDataFoot
}

type Templater interface {
	GenerateTemplate(template string, data interface{}) (io.Reader, error)
	// will prepend "page:" when finding template
	GeneratePage(page string, data PageData) (io.Reader, error)
}

func BreadcrumbsFromUrl(url string) []Link {
	retval := []Link{}
	compoundUrl := ""
	for _, prt := range strings.Split(url, "/") {
		if prt == "" {
			continue
		}
		compoundUrl += "/" + prt
		retval = append(retval, Link{URL: compoundUrl, Label: prt})
	}
	return retval
}

type PieCardData struct {
	ImgprssrPrefix string
	HasDate        bool
	PieLink        string
	Pie            interface{}
}

type MakerCardData struct {
	ImgprssrPrefix string
	Year           int64
	Maker          types.Maker
}
