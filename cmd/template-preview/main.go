// Command template-preview renders the site's HTML templates with fixed sample
// data so UI changes can be eyeballed (and screenshotted) without Postgres,
// Redis, Dex or imgprssr.
//
// It depends only on the standard library and internal/templatefuncs (the same
// template functions the real templater uses), so it runs with nothing but a
// Go toolchain.
//
// See README.md in this directory for usage.
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/leemartin77/mincepierank.co.uk/internal/templatefuncs"
)

// ---- Minimal mirrors of the shapes the templates read ----------------------
// Field names must match what the templates access; types are kept close to
// the real ones (int32 ratings, float64 summaries) so `eq` / `nicefloat`
// behave the same as in production.

type link struct{ Label, URL string }

type menuSettings struct {
	ActiveYear      int64
	SignedIn        bool
	UserPiesEnabled bool
}

type pageHead struct {
	Title, Description, Keywords                      string
	ShowCookieNotice, ShowNotice, ReadOnly, HasNotice bool
	Notice                                            string
	ShowBreadcrumb                                    bool
	IsAdminPage                                       bool
	MenuSettings                                      menuSettings
	Breadcrumbs                                       []link
}

type pageFoot struct{ Title string }

type pageData struct {
	Head     pageHead
	PageData interface{}
	Foot     pageFoot
}

type pie struct {
	Year                        int32
	MakerId, Id, DisplayName    string
	ImageFile, WebLink          string
	PackCount, PackPriceInPence int32
}

type maker struct{ Id, Name string }

type rankingSummary struct {
	Pastry, Filling, Topping, Looks, Value, Average float64
	Count                                           int64
}

type timestamp struct{ Time time.Time }

type userRanking struct {
	Year                                   int32
	Pastry, Filling, Topping, Looks, Value int32
	LastUpdated                            timestamp
}

type userPie struct {
	Year             int32
	Id               string
	Maker            string
	DisplayName      string
	ImageFile        string
	WebLink          string
	PackCount        int32
	PackPriceInPence int32
	Fresh            bool
	Location         string
	Labels           []string
	Status           string
}

type userPieRanking struct {
	Year        int32
	Pastry, Filling, Topping, Looks, Value int32
	LastUpdated timestamp
}

type userPieCard struct {
	PieLink        string
	ImgprssrPrefix string
	Pie            userPie
	RankingSummary rankingSummary
	Status         string
	ModerationNote string
}

type makerOpt struct {
	Name string
}

// ---- Scenarios -------------------------------------------------------------

type scenario struct {
	Name        string
	Description string
	Page        string // template name without the "page:" prefix
	Build       func() pageData
}

const activeYear = 2026

func singlePie(ranking userRanking, readonly bool, errs []string, pieYear int32) pageData {
	p := &pie{
		Year: pieYear, MakerId: "greggs", Id: "classic-mince-pie",
		DisplayName: "Classic Mince Pie", ImageFile: "greggs/logo.png",
		WebLink: "#", PackCount: 4, PackPriceInPence: 200,
	}
	return pageData{
		Head: pageHead{
			Title:        "Classic Mince Pie from Greggs :: 3.21/5",
			MenuSettings: menuSettings{ActiveYear: activeYear, SignedIn: true},
		},
		PageData: map[string]interface{}{
			"CategoryLinks":  []link{{"Traditional", "#"}},
			"ImgprssrPrefix": "/img/",
			"HasMaker":       true,
			"Pie":            p,
			"RankingSummary": &rankingSummary{Pastry: 3.4, Filling: 2.75, Topping: 4.1, Looks: 3.6, Value: 2.2, Average: 3.21, Count: 12},
			"Maker":          &maker{Id: "greggs", Name: "Greggs"},
			"Breadcrumb":     []link{{"Home", "/"}, {fmt.Sprint(pieYear), "#"}, {"Greggs", "#"}, {"Classic Mince Pie", "#"}},
			"Readonly":       readonly,
			"UserRanking":    ranking,
			"Errors":         errs,
		},
	}
}

var ranked = userRanking{
	Year: activeYear, Pastry: 2, Filling: 4, Topping: 3, Looks: 5, Value: 4,
	LastUpdated: timestamp{time.Date(activeYear, 10, 16, 19, 30, 0, 0, time.UTC)},
}

var scenarios = []scenario{
	{"singlepie-ranked", "Signed in, active year, existing ranking (P2 F4 T3 L5 V4) - editable form", "singlepie",
		func() pageData { return singlePie(ranked, false, nil, activeYear) }},
	{"singlepie-unranked", "Signed in, active year, not yet ranked - empty form", "singlepie",
		func() pageData { return singlePie(userRanking{}, false, nil, activeYear) }},
	{"singlepie-errors", "Submission rejected with validation errors", "singlepie",
		func() pageData {
			return singlePie(userRanking{}, false, []string{"Pastry is required", "Value for Money is required"}, activeYear)
		}},
	{"singlepie-readonly", "Site in read-only mode - ranking shown, no form", "singlepie",
		func() pageData { return singlePie(ranked, true, nil, activeYear) }},
	{"singlepie-pastyear", "Pie from a previous year - ranking shown, no form", "singlepie",
		func() pageData { r := ranked; r.Year = 2025; return singlePie(r, false, nil, 2025) }},

	// User pie scenarios
	{"submit", "Community pie submission form with brand dropdown and 'Other' option", "submit",
		func() pageData {
			return pageData{
				Head: pageHead{
					Title:        "Submit a pie for 2026",
					MenuSettings: menuSettings{ActiveYear: activeYear, SignedIn: true},
				},
				PageData: map[string]interface{}{
					"Year":       activeYear,
					"Breadcrumb": []link{{"Home", "/"}, {fmt.Sprint(activeYear), "#"}, {"Submit", "#"}},
					"Readonly":   false,
					"Makers":     []makerOpt{{"Greggs"}, {"M&S"}, {"Co-op"}, {"Tesco"}, {"Sainsbury's"}},
				},
			}
		}},
	{"userpie-ranked", "Community pie detail with existing ranking - signed in", "userpie",
		func() pageData {
			p := userPie{
				Year: activeYear, Id: "abc123", Maker: "Tesco Value", DisplayName: "Value Mince Pie",
				ImageFile: "tesco/value-mince.jpg", WebLink: "#", PackCount: 4, PackPriceInPence: 150,
				Fresh: false, Location: "London, UK", Labels: []string{"Value"}, Status: "approved",
			}
			return pageData{
				Head: pageHead{
					Title:        "Value Mince Pie (Community) :: 3.50/5",
					MenuSettings: menuSettings{ActiveYear: activeYear, SignedIn: true},
				},
				PageData: map[string]interface{}{
					"ImgprssrPrefix": "/img/",
					"Pie":            p,
					"RankingSummary": rankingSummary{Pastry: 3.2, Filling: 3.8, Topping: 3.5, Looks: 3.0, Value: 4.0, Average: 3.50, Count: 8},
					"Breadcrumb":     []link{{"Home", "/"}, {fmt.Sprint(activeYear), "#"}, {"Community Pies", "#"}, {"Value Mince Pie", "#"}},
					"Readonly":       false,
					"UserRanking":    userPieRanking{Year: activeYear, Pastry: 3, Filling: 4, Topping: 3, Looks: 4, Value: 5, LastUpdated: timestamp{time.Date(activeYear, 10, 17, 14, 0, 0, 0, time.UTC)}},
					"IsOwner":        false,
					"Status":         "approved",
					"Errors":         nil,
					"RankingFormURL": fmt.Sprintf("/years/%d/user-pies/abc123/rating", activeYear),
				},
			}
		}},
	{"userpie-unranked", "Community pie detail without ranking - signed in", "userpie",
		func() pageData {
			p := userPie{
				Year: activeYear, Id: "def456", Maker: "M&S", DisplayName: "Premium Mince Pie",
				ImageFile: "ms/premium-mince.jpg", WebLink: "#", PackCount: 6, PackPriceInPence: 360,
				Fresh: true, Location: "Manchester, UK", Labels: []string{"Premium", "Free-range"}, Status: "approved",
			}
			return pageData{
				Head: pageHead{
					Title:        "Premium Mince Pie (Community) :: 4.20/5",
					MenuSettings: menuSettings{ActiveYear: activeYear, SignedIn: true},
				},
				PageData: map[string]interface{}{
					"ImgprssrPrefix": "/img/",
					"Pie":            p,
					"RankingSummary": rankingSummary{Pastry: 4.5, Filling: 4.2, Topping: 4.0, Looks: 4.3, Value: 3.8, Average: 4.20, Count: 15},
					"Breadcrumb":     []link{{"Home", "/"}, {fmt.Sprint(activeYear), "#"}, {"Community Pies", "#"}, {"Premium Mince Pie", "#"}},
					"Readonly":       false,
					"UserRanking":    userPieRanking{},
					"IsOwner":        false,
					"Status":         "approved",
					"Errors":         nil,
					"RankingFormURL": fmt.Sprintf("/years/%d/user-pies/def456/rating", activeYear),
				},
			}
		}},
	{"userpielist", "Community pies list with multiple pies", "userpielist",
		func() pageData {
			pies := []userPieCard{
				{PieLink: "/years/2026/user-pies/abc123", ImgprssrPrefix: "/img/", Pie: userPie{Year: activeYear, Id: "abc123", Maker: "Tesco Value", DisplayName: "Value Mince Pie", ImageFile: "tesco/value-mince.jpg", PackPriceInPence: 150}, RankingSummary: rankingSummary{Average: 3.50, Count: 8}},
				{PieLink: "/years/2026/user-pies/def456", ImgprssrPrefix: "/img/", Pie: userPie{Year: activeYear, Id: "def456", Maker: "M&S", DisplayName: "Premium Mince Pie", ImageFile: "ms/premium-mince.jpg", PackPriceInPence: 360}, RankingSummary: rankingSummary{Average: 4.20, Count: 15}},
				{PieLink: "/years/2026/user-pies/ghi789", ImgprssrPrefix: "/img/", Pie: userPie{Year: activeYear, Id: "ghi789", Maker: "Co-op", DisplayName: "Free Range Mince Pie", ImageFile: "coop/freerange-mince.jpg", PackPriceInPence: 280}, RankingSummary: rankingSummary{Average: 3.80, Count: 12}},
			}
			return pageData{
				Head: pageHead{
					Title:        "Community Pies for 2026",
					MenuSettings: menuSettings{ActiveYear: activeYear, SignedIn: true},
				},
				PageData: map[string]interface{}{
					"Heading":            fmt.Sprintf("Community Pies for %d", activeYear),
					"Breadcrumb":         []link{{"Home", "/"}, {fmt.Sprint(activeYear), "#"}, {"Community Pies", "#"}},
					"PieCards":           pies,
					"NextPage":           2,
					"TotalPages":         3,
					"ShouldLoadNextPage": true,
					"URL":                fmt.Sprintf("/years/%d/user-pies", activeYear),
					"Query":              "",
				},
			}
		}},
	{"profilepies", "My pies profile page with mixed statuses", "profilepies",
		func() pageData {
			pies := []userPieCard{
				{PieLink: "/years/2026/user-pies/abc123", ImgprssrPrefix: "/img/", Pie: userPie{Year: activeYear, Id: "abc123", Maker: "Tesco Value", DisplayName: "Value Mince Pie", ImageFile: "tesco/value-mince.jpg", PackPriceInPence: 150}, RankingSummary: rankingSummary{Average: 3.50, Count: 8}, Status: "approved"},
				{PieLink: "/profile/pies/2026/def456", ImgprssrPrefix: "/img/", Pie: userPie{Year: activeYear, Id: "def456", Maker: "M&S", DisplayName: "Premium Mince Pie", ImageFile: "ms/premium-mince.jpg", PackPriceInPence: 360}, RankingSummary: rankingSummary{Average: 0, Count: 0}, Status: "pending"},
				{PieLink: "/profile/pies/2026/ghi789", ImgprssrPrefix: "/img/", Pie: userPie{Year: activeYear, Id: "ghi789", Maker: "Co-op", DisplayName: "Budget Mince Pie", ImageFile: "coop/budget-mince.jpg", PackPriceInPence: 100}, RankingSummary: rankingSummary{Average: 0, Count: 0}, Status: "rejected", ModerationNote: "Duplicate submission"},
			}
			return pageData{
				Head: pageHead{
					Title:        "My Community Pies for 2026",
					MenuSettings: menuSettings{ActiveYear: activeYear, SignedIn: true},
				},
				PageData: map[string]interface{}{
					"Heading":    fmt.Sprintf("My Community Pies for %d", activeYear),
					"Breadcrumb": []link{{"Home", "/"}, {"Profile", "#"}, {"My Pies", "#"}},
					"PieCards":   pies,
					"Year":       activeYear,
				},
			}
		}},
	{"editpie", "Edit pie form for a pending submission", "editpie",
		func() pageData {
			p := userPie{Year: activeYear, Id: "def456", Maker: "M&S", DisplayName: "Premium Mince Pie", ImageFile: "ms/premium-mince.jpg", WebLink: "#", PackCount: 6, PackPriceInPence: 360, Fresh: true, Location: "Manchester, UK", Labels: []string{"Premium", "Free-range"}, Status: "pending"}
			return pageData{
				Head: pageHead{
					Title:        "Edit pie",
					MenuSettings: menuSettings{ActiveYear: activeYear, SignedIn: true},
				},
				PageData: map[string]interface{}{
					"Year":        activeYear,
					"Pie":         p,
					"Breadcrumb":  []link{{"Home", "/"}, {"Profile", "#"}, {"My Pies", "#"}, {"Edit", "#"}},
					"Makers":      []makerOpt{{"Greggs"}, {"M&S"}, {"Co-op"}, {"Tesco"}, {"Sainsbury's"}},
					"LabelsStr":   "Premium, Free-range",
					"DisplayName": "Premium Mince Pie",
					"Maker":       "M&S",
					"Location":    "Manchester, UK",
					"Fresh":       true,
					"WebLink":     "#",
					"PackCount":   "6",
					"PackPrice":   "3.60",
					"Readonly":    false,
					"Errors":      nil,
				},
			}
		}},
}

// ---- Templates -------------------------------------------------------------

// funcMap is the shared standard-library template functions plus a stub for
// containscat, which needs storage types (it only affects category pages).
func funcMap() template.FuncMap {
	f := templatefuncs.Common()
	f["containscat"] = func(cat, catlist interface{}) bool { return false }
	return f
}

func parseTemplates(dir string) (*template.Template, error) {
	var files []string
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(p, ".gohtml") {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no .gohtml files under %s", dir)
	}
	return template.New("").Funcs(funcMap()).ParseFiles(files...)
}

func render(w io.Writer, tmplDir string, s scenario) error {
	t, err := parseTemplates(tmplDir)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, "page:"+s.Page, s.Build())
}

// ---- Entry point -----------------------------------------------------------

func main() {
	tmplDir := flag.String("templates", "assets/templates/html", "template directory")
	staticDir := flag.String("static", "assets/static", "directory served at /static/")
	imgDir := flag.String("images", "data/images", "directory served at /img/ (stands in for imgprssr)")
	addr := flag.String("addr", "127.0.0.1:8099", "listen address for serve mode")
	out := flag.String("out", "", "if set, write each scenario to <out>/<name>.html and exit instead of serving")
	flag.Parse()

	if *out != "" {
		if err := os.MkdirAll(*out, 0o755); err != nil {
			log.Fatal(err)
		}
		failed := false
		for _, s := range scenarios {
			path := filepath.Join(*out, s.Name+".html")
			f, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			if err := render(f, *tmplDir, s); err != nil {
				log.Printf("FAIL %s: %v", s.Name, err)
				failed = true
			} else {
				log.Printf("wrote %s", path)
			}
			f.Close()
		}
		if failed {
			os.Exit(1)
		}
		return
	}

	byName := map[string]scenario{}
	var names []string
	for _, s := range scenarios {
		byName[s.Name] = s
		names = append(names, s.Name)
	}
	sort.Strings(names)

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*staticDir))))
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(*imgDir))))
	mux.HandleFunc("/preview/", func(w http.ResponseWriter, r *http.Request) {
		s, ok := byName[strings.TrimPrefix(r.URL.Path, "/preview/")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		// Templates are re-parsed per request, so edits show on refresh.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := render(w, *tmplDir, s); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, `<!doctype html><meta charset="utf-8"><title>Template preview</title>`+
			`<link href="/static/bootstrap@5.3.8.min.css" rel="stylesheet">`+
			`<div class="container py-4"><h1 class="h3">Template preview</h1><ul class="list-group">`)
		for _, n := range names {
			fmt.Fprintf(w, `<a class="list-group-item list-group-item-action" href="/preview/%s"><strong>%s</strong><br><small class="text-muted">%s</small></a>`,
				template.HTMLEscapeString(n), template.HTMLEscapeString(n), template.HTMLEscapeString(byName[n].Description))
		}
		fmt.Fprint(w, `</ul></div>`)
	})

	log.Printf("template preview on http://%s/ (templates: %s)", *addr, *tmplDir)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
