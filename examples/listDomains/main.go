package main

import (
	"context"
	"encoding/json"
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
	domains, err := api.ListDomains(context.Background(), nil)
	if err != nil {
		log.Printf("Error listing domains, expected struct, got error %v", err)
	}
	out, err := json.Marshal(domains)
	if err != nil {
		log.Print(err)
	}
	fmt.Print(string(out))

}
