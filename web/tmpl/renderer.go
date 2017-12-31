package tmpl

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type renderer struct {
	path      string
	templates map[string]*template.Template
}

func NewRenderer(templatePath string) *renderer {
	path := templatePath
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return &renderer{path: path, templates: map[string]*template.Template{}}
}

func (r *renderer) Render(w http.ResponseWriter, templateName string, variables interface{}) error {
	template, ok := r.templates[templateName]
	if !ok {
		template = r.loadTemplate(templateName)
		r.templates[templateName] = template
	}
	if err := template.Execute(w, variables); err != nil {
		panic(fmt.Sprintf("template error: %v", err))
	}
	return nil
}

func (r *renderer) loadTemplate(name string) *template.Template {
	tFiles := []string{r.path + name, r.path + "header.html.tmpl", r.path + "footer.html.tmpl", r.path + "navi.html.tmpl"}
	return template.Must(template.New(name).ParseFiles(tFiles...))
}
