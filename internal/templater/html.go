package templater

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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
	fmt.Println("Got here")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return &buf, nil
}

func (wrp *HtmlTemplaterWrapper) GenerateTemplate(template string, data interface{}) (io.Reader, error) {
	var buf bytes.Buffer
	err := wrp.Templates.ExecuteTemplate(&buf, template, data)
	fmt.Println("Got here")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return &buf, nil
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
	return template.New("").ParseFiles(allFiles...)
}
