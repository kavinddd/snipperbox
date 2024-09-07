package main

import (
	"html/template"
	"path/filepath"
	"time"
)

var funcMap template.FuncMap = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t *time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		ts, err := template.ParseFiles("./ui/html/base.html")
		ts.Funcs(funcMap)

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")

		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)

		if err != nil {
			return nil, err
		}

		fileName := filepath.Base(page)

		cache[fileName] = ts

	}

	return cache, nil

}
