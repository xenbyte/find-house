package main

import (
	"fmt"

	"github.com/xenbyte/find-house/services/pararius"
)

func main() {
	city := "almere"

	listings := pararius.ScrapeListings(city)
	for _, listing := range listings {
		fmt.Printf("ID: %v\nLink: %s\nTitle: %s\nSubtitle: %s\nPrice: %v\n\n", listing.ID, listing.Link, listing.Title, listing.Subtitle, listing.Price)
	}
}
