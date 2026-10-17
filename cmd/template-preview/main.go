// Command template-preview renders the site's HTML templates with fixed sample
// data so UI changes can be eyeballed (and screenshotted) without Postgres,
// Redis, Dex or imgprssr.
//
// It deliberately uses only the standard library so it runs with nothing but a
// Go toolchain. The template FuncMap below mirrors internal/templater/html.go;
// keep the two in sync when adding template functions.
//
// See README.md in this directory for usage.
package main

import (
	"errors"
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
)

// ---- Minimal mirrors of the shapes the templates read ----------------------
// Field names must match what the templates access; types are kept close to
// the real ones (int32 ratings, float64 summaries) so `eq` / `nicefloat`
// behave the same as in production.

type link struct{ Label, URL string }

type menuSettings struct {
	ActiveYear int64
	SignedIn   bool
}

type pageHead struct {
	Title, Description, Keywords                      string
	ShowCookieNotice, ShowNotice, ReadOnly, HasNotice bool
	Notice                                            string
	ShowBreadcrumb                                    bool
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
			return singlePie(userRanking{}, false, []string{"pastry is required", "value is required"}, activeYear)
		}},
	{"singlepie-readonly", "Site in read-only mode - ranking shown, no form", "singlepie",
		func() pageData { return singlePie(ranked, true, nil, activeYear) }},
	{"singlepie-pastyear", "Pie from a previous year - ranking shown, no form", "singlepie",
		func() pageData { r := ranked; r.Year = 2025; return singlePie(r, false, nil, 2025) }},
}

// ---- Templates -------------------------------------------------------------

// funcMap mirrors internal/templater/html.go. Keep in sync.
var funcMap = template.FuncMap{
	"dict": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid dict call")
		}
		d := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			k, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}
			d[k] = values[i+1]
		}
		return d, nil
	},
	"jsonserial":  func(v interface{}) (string, error) { return fmt.Sprintf("%+v", v), nil },
	"containscat": func(cat, catlist interface{}) bool { return false },
	"ppencediv": func(v interface{}) (float32, error) {
		switch n := v.(type) {
		case float32:
			return n / 100, nil
		case int32:
			return float32(n) / 100, nil
		}
		return 0, fmt.Errorf("first value not floatable")
	},
	"timetohuman": func(v interface{}) (string, error) {
		tm, ok := v.(time.Time)
		if !ok {
			return "", fmt.Errorf("is not a time")
		}
		return tm.Format(time.RFC1123), nil
	},
	"nicefloat": func(v interface{}) (string, error) {
		var f float32
		switch n := v.(type) {
		case float32:
			f = n
		case float64:
			f = float32(n)
		default:
			return "", fmt.Errorf("first value not floatable")
		}
		s := fmt.Sprintf("%.2f", f)
		return strings.TrimSuffix(s, ".00"), nil
	},
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
	return template.New("").Funcs(funcMap).ParseFiles(files...)
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
