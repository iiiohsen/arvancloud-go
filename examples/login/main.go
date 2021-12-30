package main

import (
	"fmt"
	"log"
	"os"

	"github.com/S4eedb/arvancloud-go"
)

func main() {
	api, err := arvancloud.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
	if err != nil {
		log.Print(err)
	}
	fmt.Println(api)
}
