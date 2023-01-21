package server

import (
	"errors"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

var (
	errUnknownError = errors.New("unknown error")
)

func readTemplates(folder string) (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}
	if err := filepath.Walk(folder, func(name string, info os.FileInfo, _ error) error {
		if info == nil {
			return os.ErrNotExist
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(name, ".html") {
			return nil
		}

		noHTML := name[len(folder) : len(name)-5]
		tmpl, err := newTemplate(name, noHTML)
		if err != nil {
			return err
		}

		templates[noHTML] = tmpl

		return nil

	}); err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		return nil, errors.New("found no templates in folder")
	}

	return templates, nil
}

func newTemplate(path, name string) (*template.Template, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(b))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

// notify ...
func notify(e error, meta ...interface{}) {
	// log.Println(e, meta)
	if err := bugsnag.Notify(e, meta...); err != nil {
		log.Println(err)
		log.Println(e)
		log.Println(meta...)
	}
}
