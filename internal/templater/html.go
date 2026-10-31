package templater

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/leemartin77/mincepierank.co.uk/assets"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
	"github.com/leemartin77/mincepierank.co.uk/internal/templatefuncs"
)

func getTemplatesFromFS() (*template.Template, error) {
	funcs := templatefuncs.Common()
	funcs["containscat"] = func(cat interface{}, catlist interface{}) bool {
		category, ok := cat.(types.Category)
		if !ok {
			return false
		}
		category_list, ok := catlist.(*[]types.Category)
		if !ok {
			return false
		}
		return slices.ContainsFunc(*category_list, func(cl types.Category) bool {
			return cl.Id == category.Id
		})
	}

	t := template.New("").Funcs(funcs)

	err := fs.WalkDir(assets.TemplatesFS, "templates/html", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".gohtml") {
			return nil
		}
		content, err := assets.TemplatesFS.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = t.Parse(string(content))
		return err
	})
	if err != nil {
		return nil, err
	}
	return t, nil
}

func NewHtmlTemplater() Templater {

	embedded, err := getTemplatesFromFS()
	if err != nil {
		panic(err)
	}
	fmt.Println(embedded.DefinedTemplates())
	tmpltr := HtmlTemplaterWrapper{
		Templates: embedded,
	}
	return &tmpltr
}

type HtmlTemplaterWrapper struct {
	Templates *template.Template
}

func (wrp *HtmlTemplaterWrapper) GeneratePage(template string, data PageData) (io.Reader, error) {
	var buf bytes.Buffer
	err := wrp.Templates.ExecuteTemplate(&buf, "page:"+template, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return &buf, nil
}

// WritePage implements Templater: it renders page via GeneratePage and
// writes it to w with the given status and an html content type.
func (wrp *HtmlTemplaterWrapper) WritePage(w http.ResponseWriter, status int, page string, data PageData) {
	rdr, err := wrp.GeneratePage(page, data)
	if err != nil {
		http.Error(w, "error rendering page", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if _, err := io.Copy(w, rdr); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing page response: %v\n", err)
	}
}

func (wrp *HtmlTemplaterWrapper) GenerateTemplate(template string, data interface{}) (io.Reader, error) {
	var buf bytes.Buffer
	err := wrp.Templates.ExecuteTemplate(&buf, template, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return &buf, nil
}

func (wrp *HtmlTemplaterWrapper) GetTemplater() *template.Template {
	return wrp.Templates
}
