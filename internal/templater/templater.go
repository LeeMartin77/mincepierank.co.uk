package templater

import (
	"io"
	"strings"
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
	Breadcrumbs      []BreadcrumbLink
}

type MenuSettings struct {
	ActiveYear int64
	LoggedIn   bool
}

type BreadcrumbLink struct {
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

func BreadcrumbsFromUrl(url string) []BreadcrumbLink {
	retval := []BreadcrumbLink{}
	compoundUrl := ""
	for _, prt := range strings.Split(url, "/") {
		if prt == "" {
			continue
		}
		compoundUrl += "/" + prt
		retval = append(retval, BreadcrumbLink{URL: compoundUrl, Label: prt})
	}
	return retval
}
