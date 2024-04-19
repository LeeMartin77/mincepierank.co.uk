package templater

import "io"

type PageDataHead struct {
	Title            string
	ShowCookieNotice bool
	ShowNotice       bool
	ReadOnly         bool
	HasNotice        bool
	Notice           string
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
