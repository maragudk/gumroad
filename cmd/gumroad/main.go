package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/maragudk/gumroad"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) == 1 {
		printUsage()
	}

	client := gumroad.NewClientWithOptions(gumroad.NewClientOptions{AccessToken: os.Getenv("GUMROAD_ACCESS_TOKEN")})
	ctx := context.Background()

	var err error
	var res interface{}
	switch os.Args[1] {
	case "get-products":
		res, err = client.GetProducts(ctx)
	case "get-resource-subscriptions":
		if len(os.Args) == 2 {
			log.Fatalln("Usage: gumroad get-resource-subscriptions", gumroad.ResourceSubscriptions)
		}
		res, err = client.GetResourceSubscriptions(ctx, gumroad.ResourceSubscription(os.Args[2]))
	default:
		printUsage()
	}

	if err != nil {
		log.Fatalln(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(res); err != nil {
		log.Fatalln(err)
	}
}

func printUsage() {
	log.Fatalln(`Usage: gumroad [cmd]
where [cmd] is one of
	get-products
	get-resource-subscriptions`)
}
