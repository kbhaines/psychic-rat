package main

import (
	"html/template"
	"log"
	"net/http"
	"psychic-rat/api"
)

type variables struct {
	Username string
	Items    []api.ItemElement
}

func renderPage(writer http.ResponseWriter, templateName string, variables interface{}) {
	tpt := template.Must(template.New(templateName).ParseFiles(templateName, "header.html.tmpl", "footer.html.tmpl"))
	//tpt.ExecuteTemplate(writer, "header.html", variables)
	tpt.Execute(writer, variables)
	//tpt.ExecuteTemplate(writer, "footer.html", variables)
}

func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "home.html.tmpl", vars)
}

func SignInPageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "signin.html.tmpl", vars)
}

func PledgePageHandler(writer http.ResponseWriter, request *http.Request) {
	itemApi := api.GetRepoItemApi()
	report, err := itemApi.ListItems()
	if err != nil {
		log.Fatal(err)
	}

	vars := variables{Username: "Kevin", Items: report.Items}
	renderPage(writer, "pledge.html.tmpl", vars)
}

func ThanksPageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "thanks.html.tmpl", vars)
}
