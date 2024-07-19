package pararius

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/xenbyte/find-house/listings"
	"github.com/xenbyte/find-house/utils"
)

func ScrapeListings(city string) []listings.Listing {
	c := colly.NewCollector()
	var ls []listings.Listing
	page := 1

	for {
		url := fmt.Sprintf("https://www.pararius.com/apartments/%s/page-%d", city, page)
		if !CheckParariusURLStatus(url) {
			break
		}

		c.OnHTML(`ul.search-list[data-controller="search-list"]`, func(e *colly.HTMLElement) {
			e.ForEach("li.search-list__item--listing", func(_ int, el *colly.HTMLElement) {
				link := el.ChildAttr("div.listing-search-item__depiction a.listing-search-item__link--depiction", "href")
				title := el.ChildText("h2.listing-search-item__title")
				subtitle := el.ChildText("div.listing-search-item__sub-title")
				priceStr := el.ChildText("div.listing-search-item__price")
				price, err := utils.ExtractPriceNumber(priceStr)
				if err != nil {
					log.Println("couldn't parse ID: ", err.Error())
				}

				id, err := utils.ExtractIDFromURL(link)
				if err != nil {
					log.Println("couldn't parse ID: ", err.Error())
				}

				listing := listings.Listing{
					ID:       id,
					Link:     fmt.Sprintf("https://pararius.com%v", el.ChildAttr("div.listing-search-item__depiction a.listing-search-item__link--depiction", "href")),
					Title:    title,
					Subtitle: subtitle,
					Price:    price,
				}
				ls = append(ls, listing)
			})
		})

		err := c.Visit(url)
		if err != nil {
			log.Printf("Error visiting URL %s: %v", url, err)
		}

		page++
	}

	c.Wait()
	return ls
}

func CheckParariusURLStatus(url string) bool {
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
