package templater

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

func NewHtmlTemplater() Templater {

	tmplts, err := getTemplates("assets/templates/html")
	if err != nil {
		panic(err)
	}
	fmt.Println(tmplts.DefinedTemplates())
	tmpltr := HtmlTemplaterWrapper{
		Templates: tmplts,
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

func getFilesInDirAndChildren(dir string) []string {
	var allFiles []string
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		filename := file.Name()
		if file.IsDir() {
			filePath := filepath.Join(dir, filename)
			allFiles = append(allFiles, getFilesInDirAndChildren(filePath)...)
		} else if strings.HasSuffix(filename, ".gohtml") {
			filePath := filepath.Join(dir, filename)
			allFiles = append(allFiles, filePath)
		}
	}
	return allFiles
}

func getTemplates(tmpltdir string) (templates *template.Template, err error) {
	allFiles := getFilesInDirAndChildren(tmpltdir)
	return template.New("").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"jsonserial": func(value interface{}) (string, error) {
			b, err := json.MarshalIndent(value, "", "  ")
			if err != nil {
				fmt.Println(value)
				return "Error serialising", err
			}
			return string(b), nil
		},
		"ppencediv": func(value interface{}) (float32, error) {
			first, ok := value.(float32)
			if !ok {
				fint, ok := value.(int32)
				if !ok {

					return 0, fmt.Errorf("first value not floatable")
				}
				first = float32(fint)
			}
			return first / 100, nil
		},
		"containscat": func(cat interface{}, catlist interface{}) bool {
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
		},
	}).ParseFiles(allFiles...)
}
