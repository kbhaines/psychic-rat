package tmpl

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var (
	path      string
	templates = map[string]*template.Template{}
)

func Init(templatePath string) {
	path = templatePath
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
}

func RenderTemplate(writer http.ResponseWriter, templateName string, variables interface{}) {
	template, ok := templates[templateName]
	if !ok {
		template = loadTemplate(templateName)
		templates[templateName] = template
	}
	if err := template.Execute(writer, variables); err != nil {
		panic(fmt.Sprintf("template error: %v", err))
	}
}

func loadTemplate(name string) *template.Template {
	tFiles := []string{path + name, path + "header.html.tmpl", path + "footer.html.tmpl", path + "navi.html.tmpl"}
	return template.Must(template.New(name).ParseFiles(tFiles...))
}
