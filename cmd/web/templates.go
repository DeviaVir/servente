package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/DeviaVir/servente/pkg/forms"
	"github.com/DeviaVir/servente/pkg/models"
)

type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	Form            *forms.Form
	IsAuthenticated bool
	Service         *models.Service
	Services        []*models.Service
	User            *models.User
	Organization    *models.Organization
	Organizations   []*models.Organization
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			pages, err := filepath.Glob(filepath.Join(path, "*.page.tmpl"))
			if err != nil {
				return err
			}

			dirName := filepath.Base(path)
			for _, page := range pages {
				name := fmt.Sprintf("%s/%s", dirName, filepath.Base(page))

				ts, err := template.New(filepath.Base(page)).Funcs(functions).ParseFiles(page)
				if err != nil {
					return err
				}

				ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
				if err != nil {
					return err
				}

				ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
				if err != nil {
					return err
				}

				cache[name] = ts
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return cache, err
}
