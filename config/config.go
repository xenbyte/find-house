package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	ls "github.com/xenbyte/find-house/listings"
	"github.com/xenbyte/find-house/notify"
	"github.com/xenbyte/find-house/services/pararius"
	"github.com/xenbyte/find-house/utils"
)

type Config struct {
	MaxPrice     int
	City         string
	CSVFile      string
	APIToken     string
	TelegramUser string
	ChannelName  string
}

func LoadConfig() (*Config, error) {
	maxPrice, err := strconv.Atoi(getEnv("MAX_PRICE", "0"))
	if err != nil {
		return nil, fmt.Errorf("could not parse MAX_PRICE: %w", err)
	}

	return &Config{
		MaxPrice:     maxPrice,
		City:         getEnv("CITY", ""),
		CSVFile:      getEnv("CSV_FILE", ""),
		APIToken:     getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramUser: getEnv("TELEGRAM_USER", ""),
		ChannelName:  getEnv("CHANNEL_NAME", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func ValidateConfig(cfg *Config) error {
	if cfg.City == "" {
		return fmt.Errorf("CITY environment variable is not set")
	}
	if cfg.CSVFile == "" {
		return fmt.Errorf("CSV_FILE environment variable is not set")
	}
	if cfg.APIToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is not set")
	}
	if cfg.TelegramUser == "" {
		return fmt.Errorf("TELEGRAM_USER environment variable is not set")
	}
	return nil
}

func ProcessListings(cfg *Config) error {
	existingListings, fileExisted, err := utils.ReadExistingListings(cfg.CSVFile)
	if err != nil {
		return fmt.Errorf("error reading existing listings: %w", err)
	}

	listings := pararius.ScrapeListings(cfg.City)
	var newListings []ls.Listing

	for _, listing := range listings {
		if listing.Price < cfg.MaxPrice && !existingListings[listing.ID] {
			newListings = append(newListings, listing)

			if fileExisted {
				// chatID, err := notify.GetChatID(cfg.APIToken, cfg.TelegramUser)
				// if err != nil {
				// 	log.Println("Error getting chat ID:", err)
				// 	continue
				// }

				chatID, err := notify.GetChannelID(cfg.APIToken, cfg.ChannelName)
				if err != nil {
					log.Println("Error getting chat ID:", err)
					continue
				}

				message := fmt.Sprintf(
					"New listing found!\n\nCity: %v\nName: %s\nLocation: %s\nLink: %s\nPrice: %v\n",
					cfg.City, listing.Title, listing.Subtitle, listing.Link, listing.Price,
				)

				if err := notify.SendTelegramMessage(cfg.APIToken, chatID, message); err != nil {
					log.Println("Error sending Telegram message:", err)
				}
			}
		}
	}

	if err := utils.WriteListingsToCSV(cfg.CSVFile, newListings, existingListings); err != nil {
		return fmt.Errorf("error writing listings to CSV: %w", err)
	}

	return nil
}
