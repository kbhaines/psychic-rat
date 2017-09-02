package main

import (
	"fmt"
	"os"
	"psychic-rat/api"
)

func exitIfErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	localhost := "http://localhost:8080"
	company := api.GetRestfulCompanyApi(localhost)
	companies, err := company.GetCompanies()
	exitIfErr(err)
	for _, co := range companies.Companies {
		fmt.Printf("%s %s\n", co.Id, co.Name)
	}

	itemsApi := api.GetRestfulItemApi(localhost)
	itemReport, err := api.GetItems(company.Id("1"))
	exitIfErr(err)
	for _, i := range itemReport.Items {
		fmt.Printf("%s %s %s\n", i.Id, i.Make, i.Model)
	}

	//pledges, err := api.NewPledge(itemReport.Items[0].Id)
	//exitIfErr(err)
	//for _, p := range pledges.Pledges {
	//	fmt.Printf("%s %s %s\n", p.PledgeId, p.Item.Make, p.Item.Model)
	//}
}
