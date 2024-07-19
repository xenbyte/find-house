package main

import (
	"flag"
	"fmt"

	"github.com/xenbyte/find-house/services/pararius"
)

type Listing struct {
	ID       string
	Link     string
	Title    string
	Subtitle string
	Price    int
}

func main() {
	// Define the city flag
	city := flag.String("city", "lelystad", "The city to search for listings")

	// Define the max price flag
	maxPrice := flag.Int("maxPrice", 1600, "The maximum price for listings")

	// Parse the flags
	flag.Parse()

	// Use the city and max price flag values
	listings := pararius.ScrapeListings(*city)
	for _, listing := range listings {
		if listing.Price < *maxPrice {
			fmt.Println(listing.Link)
		}
	}
}
