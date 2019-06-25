package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Movie struct {
	title string // e.g. Venom 4K UHD
	link  string // e.g. https://blu-ray-rezensionen.net/venom-4k-uhd/
}

var Collection = []Movie{}

func MoviesTitleScrape(website string) {
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

	// Find UHD links and titles
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(string(s.Text()), "UHD") {
			Collection = append(Collection, Movie{s.Text(), href})
		}
	})
}

func Real4K(website string) {
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

	// Find UHD links and titles
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		p, _ := s.Attr("Real 4k")
		if strings.Contains(string(s.Text()), "Real 4K:") {
			//Collection = append(Collection, Movie{s.Text(), href})
			fmt.Printf("%s%s\n", s.Text(), p)
		}
	})
}

func main() {
	//MoviesTitleScrape("https://blu-ray-rezensionen.net/ultra-hd-blu-ray/")
	//fmt.Println(cap(Collection), len(Collection), Collection)

	/*
		for _, Collection := range Collection {
			fmt.Println(Collection.title, Collection.link)

		}
	*/
	Real4K("https://blu-ray-rezensionen.net/wiedersehen-in-howards-end-4k-uhd/")
}
