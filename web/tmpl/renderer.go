package tmpl

import (
	"html/template"
	"net/http"
	"strings"
)

type renderer struct {
	path     string
	cache    bool
	template *template.Template
}

func NewRenderer(templatePath string, cache bool) *renderer {
	path := templatePath
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return &renderer{path: path, cache: cache, template: nil}
}

func (r *renderer) Render(w http.ResponseWriter, templateName string, variables interface{}) error {
	if r.template == nil || !r.cache {
		r.template = template.Must(template.New(templateName).ParseGlob(r.path + "/*.tmpl"))
	}
	if err := r.template.Execute(w, variables); err != nil {
		return err
	}
	return nil
}
