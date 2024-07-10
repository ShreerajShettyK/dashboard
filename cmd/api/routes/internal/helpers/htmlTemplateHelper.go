package helpers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	// Making these variables to allow mocking in tests
	filepathAbs     = filepath.Abs
	templateParse   = template.ParseFiles
	templateExecute = func(t *template.Template, wr http.ResponseWriter, data interface{}) error {
		return t.Execute(wr, data)
	}
)

func RenderTemplate(w http.ResponseWriter, data interface{}, templateName string) error {
	tmplPath, err := filepathAbs("internal/templates/" + templateName)
	if err != nil {
		return err
	}

	tmpl, err := templateParse(tmplPath)
	if err != nil {
		return err
	}

	if err := templateExecute(tmpl, w, data); err != nil {
		return err
	}

	return nil
}
