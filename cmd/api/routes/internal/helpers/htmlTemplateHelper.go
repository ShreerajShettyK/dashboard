package helpers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, data interface{}, templateName string) error {
	// Use the relative path
	tmplPath, err := filepath.Abs("internal/templates/" + templateName)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	// Execute the template with the provided data
	if err := tmpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}
