package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/akhyar02/bookings/pkg/config"
	"github.com/akhyar02/bookings/pkg/models"
	"github.com/justinas/nosurf"
)

var app config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = *a
}

func addDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData, r *http.Request) {
	var tc map[string]*template.Template
	var err error
	if app.UseCache {
		log.Println("using cache")
		tc = app.TemplateCache
	} else {
		log.Println("creating new tc")
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Println("Error creating template cache")
		}
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Println("Error getting template from cache")
		fmt.Fprintf(w, "404 Page not found")
		return
	}

	buf := new(bytes.Buffer)
	td = addDefaultData(td, r)
	err = t.Execute(buf, td)
	if err != nil {
		log.Println("Error executing template:", err)
		fmt.Fprintf(w, "404 Page not found")
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing buffer to response:", err)
		fmt.Fprintf(w, "404 Page not found")
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	tc := make(map[string]*template.Template)

	templateFiles, err := filepath.Glob("./templates/*.page.gtpl")
	if err != nil {
		log.Println("Error reading page files:", err)
		return nil, err
	}

	for _, page := range templateFiles {
		fileName := filepath.Base(page)
		pTmpl, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			log.Println("Error parsing template:", err)
			return nil, err
		}

		layoutFiles, err := filepath.Glob("./templates/*.layout.gtpl")
		if err != nil {
			log.Println("Error reading layout files:", err)
			return nil, err
		}
		if len(layoutFiles) > 0 {
			pTmpl, err = pTmpl.ParseGlob("./templates/*.layout.gtpl")
			if err != nil {
				log.Println("Error parsing layout glob:", err)
				return nil, err
			}
		}
		tc[fileName] = pTmpl
	}
	return tc, nil
}
