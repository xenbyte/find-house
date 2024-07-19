package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// checkURLStatus checks the HTTP status of a URL
func CheckURLStatus(url string) bool {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Println("ERROR")
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error checking URL %s: %v", url, err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently
}

func ExtractPriceNumber(input string) (int, error) {
	// Define a regular expression to match the price part
	re := regexp.MustCompile(`[\d,]+`)
	priceStr := re.FindString(input)

	// Remove any commas from the matched string
	cleanPriceStr := strings.Replace(priceStr, ",", "", -1)

	// Convert the cleaned string to an integer
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

	// Find the matching part of the URL
	matches := r.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("id not found in URL")
	}

	// Extract and return the id
	id := matches[1]
	return id, nil
}

func WriteListingsToCSV(filename string, listings []pararius.Listing, existingListings map[string]bool) error {
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
			fmt.Println(listing.Link + " (new)")
		}
	}

	return nil
}
