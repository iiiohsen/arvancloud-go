package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/S4eedb/arvancloud-go"
)

func main() {

	apiKey, ok := os.LookupEnv("ARVANCLOUD_TOKEN")
	if !ok {
		log.Fatal("Could not find ARVANCLOUD_TOKEN, please assert it is set.")
	}
	api, err := arvancloud.New(apiKey)
	if err != nil {
		log.Print(err)
	}
	api.SetDebug(true)
	// account, err := api.GetAccount(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// out, err := json.Marshal(account)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%s", out)
	domains, err := api.ListDomains(context.Background(), nil)
	if err != nil {
		log.Print("Error listing domains, expected struct, got error %v", err)
	}
	fmt.Print(domains)

}
