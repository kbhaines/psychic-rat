package main

import (
	"html/template"
	"net/http"
)

type variables struct {
	Username string
}

func renderPage(writer http.ResponseWriter, templateName string, variables interface{}) {
	tpt := template.Must(template.New(templateName).ParseFiles(templateName, "header.html", "footer.html"))
	tpt.Execute(writer, variables)
}

func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "home.html", vars)
}

func SignInPageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "signin.html", vars)
}

func PledgePageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "pledge.html", vars)
}

func ThanksPageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "thanks.html", vars)
}
