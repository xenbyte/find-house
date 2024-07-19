package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xenbyte/find-house/listings"
)

func ExtractPriceNumber(input string) (int, error) {
	re := regexp.MustCompile(`[\d,]+`)
	priceStr := re.FindString(input)
	cleanPriceStr := strings.Replace(priceStr, ",", "", -1)
	price, err := strconv.Atoi(cleanPriceStr)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func ExtractIDFromURL(url string) (string, error) {
	pattern := `/([0-9a-f]{8})/`
	r, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	matches := r.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("id not found in URL")
	}

	id := matches[1]
	return id, nil
}

func WriteListingsToCSV(filename string, listings []listings.Listing, existingListings map[string]bool) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, listing := range listings {
		if !existingListings[listing.ID] {
			err := writer.Write([]string{
				listing.ID,
				listing.Link,
				listing.Title,
				listing.Subtitle,
				fmt.Sprintf("%d", listing.Price),
			})
			if err != nil {
				return err
			}
			existingListings[listing.ID] = true
		}
	}

	return nil

}

func ReadExistingListings(filename string) (map[string]bool, bool, error) {
	existingListings := make(map[string]bool)
	fileExisted := true

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, return empty map and print init message
			fmt.Println("CSV file not found, initializing...")
			// Create the file for future use
			if createErr := createCSVFile(filename); createErr != nil {
				return nil, false, createErr
			}
			fileExisted = false
			return existingListings, fileExisted, nil
		}
		return nil, false, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, false, err
	}

	for _, record := range records {
		if len(record) > 0 {
			existingListings[record[0]] = true
		}
	}

	return existingListings, fileExisted, nil
}

func createCSVFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header to CSV file
	headers := []string{"ID", "Link", "Title", "Subtitle", "Price"}
	return writer.Write(headers)
}
