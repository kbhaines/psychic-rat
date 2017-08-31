package main

import (
	"fmt"
	"psychic-rat/api/client"
	"os"
)

func main() {
	api := client.New("http://localhost:8080")
	companies, err := api.GetCompanies()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	for _, co := range companies {
		fmt.Printf("%s %s\n", co.Id(), co.Name())
	}
}
