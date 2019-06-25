package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape(website string) {
	// Request the HTML page.
	res, err := http.Get(website)
	if err != nil {
		log.Fatalf("[error] [http.Get()]: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("[error] [res.StatusCode] code:%d status:%s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Cannot load the HTML document: %s", err)
	}

	// Find the review items
	doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}

func main() {
	ExampleScrape("http://www.metalsucks.net/")
}
