package main

import (
	"log"

	"github.com/xenbyte/find-house/config"
)

type Listing struct {
	ID       string
	Link     string
	Title    string
	Subtitle string
	Price    int
}

var (
	maxPrice     int
	city         string
	csvFile      string
	apiToken     string
	telegramUser string
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Configuration error:", err)
	}

	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatal("Invalid configuration:", err)
	}

	if err := config.ProcessListings(cfg); err != nil {
		log.Fatal(err)
	}
}
