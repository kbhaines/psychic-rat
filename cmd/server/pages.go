package main

import (
	"html/template"
	"log"
	"net/http"
	"psychic-rat/api"
	"psychic-rat/mdl"
)

type variables struct {
	Username string
	Items    []api.ItemElement
}

type renderFunction func(writer http.ResponseWriter, templateName string, variables interface{})

var renderPage renderFunction = renderPageUsingTemplate

func HomePageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "home.html.tmpl", vars)
}

func renderPageUsingTemplate(writer http.ResponseWriter, templateName string, variables interface{}) {
	tpt := template.Must(template.New(templateName).ParseFiles(templateName, "header.html.tmpl", "footer.html.tmpl"))
	tpt.Execute(writer, variables)
}

func SignInPageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "signin.html.tmpl", vars)
}

func PledgePageHandler(writer http.ResponseWriter, request *http.Request) {
	report, err := apis.Item.ListItems()
	if err != nil {
		log.Fatal(err)
	}

	if request.Method == "GET" {
		vars := variables{Username: "Kevin", Items: report.Items}
		renderPage(writer, "pledge.html.tmpl", vars)
		return
	}

	if request.Method == "POST" {
		if err := request.ParseForm(); err != nil {
			log.Print(err)
			return
		}
		log.Printf("request = %+v\n", request)
		log.Printf("request.Form= %+v\n", request.Form)
		itemId := mdl.Id(request.FormValue("item"))
		if itemId == "" {
			log.Print("item field not passed")
			return
		}
		if _, err := apis.Item.GetById(itemId); err != nil {
			log.Printf("error looking up item %v : %v", itemId, err)
			return
		}
		log.Printf("pledge item %v from user", itemId)
		userId := mdl.Id("1234")
		plId, err := apis.Pledge.NewPledge(itemId, userId)
		if err != nil {
			log.Print("unable to pledge : ", err)
			return
		}
		log.Printf("pledge %v created", plId)
		return
	}
}

func ThanksPageHandler(writer http.ResponseWriter, request *http.Request) {
	vars := variables{Username: "Kevin"}
	renderPage(writer, "thanks.html.tmpl", vars)
}
