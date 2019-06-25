package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Movie struct {
	title       string // e.g. Venom 4K UHD
	link        string // e.g. https://blu-ray-rezensionen.net/venom-4k-uhd/
	real4K      string // e.g. Ja (4K DI)
	director    string // e.g. James Ivory
	studio      string // e.g. Disney
	audioFormat string // e.g. Dolby Atmos
	codec       string // e.g. HEVC
	hdr         string // e.g. HDR10
	length      string // e.g. 130
	ratio       string // e.g. 2,35:1
	country     string // e.g. USA
	year        string // e.g. 2017
	actors      string // e.g. Chris Hemsworth, Tom Hiddleston, Cate Blanchett
	pq          string // 85%
}

var Collection = []Movie{
	title:       movieTitle,
	link:        movieLink,
	real4K:      movieReal4K,
	director:    movieDirector,
	studio:      movieStudio,
	audioFormat: movieAudioFormat,
	codec:       movieCodec,
	hdr:         movieHDR,
	length:      movieLength,
	ratio:       movieRatio,
	country:     movieCountry,
	year:        movieYear,
	actors:      movieActors,
	pq:          moviePQ,
})
}

func titleAndURL(website string, Collection []Movie) *goquery.Document {
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
	return doc
}

func main() {
	doc := titleAndURL("https://blu-ray-rezensionen.net/ultra-hd-blu-ray/", Collection)
	doc.Find("a[href]").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, _ := s.Attr("href")
		if strings.Contains(string(s.Text()), "UHD") {
			movieTitle := s.Text()
			movieLink := href
		}
		return true
	})

	for _, Collection := range Collection {
		fmt.Println("Movie:", Collection.title)
		fmt.Println("Link:", Collection.link)
	}

}
