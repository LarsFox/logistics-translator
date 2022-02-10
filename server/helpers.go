package server

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
)

var (
	errNoParamsInBody  = errors.New("method has no params in body")
	errUnknownError    = errors.New("unknown error")
	errUnknownLanguage = errors.New("unknown language")
	errUnnamedMethod   = errors.New("unnamed method")
)

var tmplDefault = template.Must(template.New("go").Parse("blank"))

func readTemplates(folder string) map[string]*template.Template {
	templates := map[string]*template.Template{}
	if err := filepath.Walk(folder, func(name string, info os.FileInfo, err error) error {
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
		templates[noHTML] = newTemplate(name, noHTML)
		return nil
	}); err != nil {
		return nil
	}
	return templates
}

func newTemplate(path, name string) *template.Template {
	b, err := os.ReadFile(path)
	if err != nil {
		notify(err)
		return tmplDefault
	}
	tmpl, err := template.New(name).Parse(string(b))
	if err != nil {
		notify(err)
		return tmplDefault
	}
	return tmpl
}

func newLogger(local bool, path string) *log.Logger {
	if local {
		return log.Default()
	}

	outfile, err := os.Create(fmt.Sprintf("%s/log_%d.log", path, time.Now().Unix()))
	if err != nil {
		log.Printf("unable to touch log file: %v", err)
		return log.Default()
	}
	return log.New(outfile, "", 0)
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
