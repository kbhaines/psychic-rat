package main

import (
	"fmt"
	"os"
	"psychic-rat/api/client"
	"psychic-rat/mdl/company"
)

func exitIfErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	api := client.New("http://localhost:8080")
	companies, err := api.GetCompanies()
	exitIfErr(err)
	for _, co := range companies {
		fmt.Printf("%s %s\n", co.Id(), co.Name())
	}

	items, err := api.GetItems(company.Id("1"))
	exitIfErr(err)
	for _, i := range items {
		fmt.Printf("%s %s %s\n", i.Id(), i.Make(), i.Model())
	}
}
