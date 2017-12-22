package tmpl

import (
	"html/template"
	"net/http"
)

var templates = map[string]*template.Template{}

// TODO: handle template errors
func RenderTemplate(writer http.ResponseWriter, templateName string, variables interface{}) {
	template, ok := templates[templateName]
	if !ok {
		template = loadTemplate(templateName)
		templates[templateName] = template
	}
	template.Execute(writer, variables)
}

func loadTemplate(name string) *template.Template {
	tFiles := []string{"res/" + name, "res/header.html.tmpl", "res/footer.html.tmpl", "res/navi.html.tmpl"}
	return template.Must(template.New(name).ParseFiles(tFiles...))
}
