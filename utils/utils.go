package utils

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

func ExtractPriceNumber(priceStr string) int {
	re := regexp.MustCompile(`[0-9]+`)
	matches := re.FindAllString(priceStr, -1)
	if len(matches) > 0 {
		price, err := strconv.Atoi(matches[0])
		if err != nil {
			log.Printf("Error converting price string to number: %v", err)
			return 0
		}
		return price
	}
	return 0
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
