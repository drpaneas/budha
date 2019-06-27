package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/yosssi/gohtml"

	"github.com/PuerkitoBio/goquery"
)

type Movie struct {
	title       string // e.g. Venom 4K UHD
	url         string // e.g. https://blu-ray-rezensionen.net/venom-4k-uhd/
	real4K      string // e.g. Ja (4K DI)
	director    string // e.g. James Ivory
	studio      string // e.g. Disney
	audioFormat string // e.g. Dolby Atmos
	codec       string // e.g. HEVC
	hdr         string // e.g. HDR10
	runtime     string // e.g. 130
	screen      string // e.g. 2,35:1
	country     string // e.g. USA
	year        string // e.g. 2017
	actors      string // e.g. Chris Hemsworth, Tom Hiddleston, Cate Blanchett
	pq          string // 85%
}

var phrase string
var phraseList []string
var err error

// Initialize slice

func getHTMLDocument(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal("Cannot load the HTML document:", err)
	}
	return doc
}

func htmlToStr(htmlDocument *goquery.Document) string {
	htmlCode, err := htmlDocument.Html()
	if err != nil {
		log.Fatalf("Error while converting HTML Document into string: %s", err)
	}
	htmlPrettyCode := gohtml.Format(htmlCode)
	return htmlPrettyCode
}

func printSlice(mySlice []string) {
	fmt.Printf("len=%d cap=%d\n", len(mySlice), cap(mySlice))
	for index, element := range mySlice {
		fmt.Printf("%4d : %v\n", index, element)
	}
}

func printMovieSlice(mySlice []Movie) {
	fmt.Printf("len=%d cap=%d\n", len(mySlice), cap(mySlice))
	for index, element := range mySlice {
		fmt.Printf("%4d : %v\n", index, element)
	}
	fmt.Printf("\n\n\n")
}

func getLinksWithPhrase(htmlDocument *goquery.Document, phrase string) []string {
	linkSlice := []string{}
	htmlDocument.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(string(s.Text()), phrase) {
			linkSlice = append(linkSlice, href)
		}
	})
	return linkSlice
}

func getTitle(htmlDocument *goquery.Document) string {
	var title string
	htmlDocument.Find("title").Each(func(i int, s *goquery.Selection) {
		title = string(s.Text())
	})
	return title
}

func getTags(htmlDocument *goquery.Document, tag string) []string {
	tags := []string{}
	htmlDocument.Find(tag).Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err != nil {
			log.Fatalf("Cannot parse the HTML: %s", err)
		}
		// Replace the HTML breaks with newline escape characters, creating a multiline string
		str := strings.Replace(html, "<br/>", "\n", -1)
		// Iterate over this multiline string
		scanner := bufio.NewScanner(strings.NewReader(str))
		for scanner.Scan() {
			tags = append(tags, string(scanner.Text()))
		}
	})
	return tags
}

func getTagWithPhrase(tags []string, phrase string) (string, error) {
	for index := range tags {
		if strings.Contains(tags[index], phrase) {
			// Phrase matched the value of the tag
			return strings.TrimSpace(strings.Split(tags[index], ":")[1]), nil
		}
	}
	err := errors.New("not found")
	return "?", err
}

func getScreen(tags []string, phrase string) (string, error) {
	for index := range tags {
		if strings.Contains(tags[index], phrase) {
			screen := strings.TrimSpace(strings.Split(tags[index], ":")[1])
			// in case there are more than one ':'
			// only the first one is considered as delimeter
			number := strings.Count(tags[index], ":")
			if number > 1 {
				for i := 2; i <= number; i++ {
					screen += ":" + strings.TrimSpace(strings.Split(tags[index], ":")[i])
				}
			}
			return screen, nil
		}
	}
	err := errors.New("not found")
	return "?", err
}

func getCountry(tags []string, phrase string) (string, error) {
	for index := range tags {
		if strings.Contains(tags[index], phrase) {
			parts := strings.Split(tags[index], ":")      // e.g. Land/Jahr: USA 2017
			trimmedPart := strings.TrimSpace(parts[1])    // e.g. "USA 2017"
			country := strings.Split(trimmedPart, " ")[0] //e.g. USA
			return country, nil
		}
	}
	err := errors.New("not found")
	return "?", err
}

func getYear(tags []string, phrase string) (string, error) {
	for index := range tags {
		if strings.Contains(tags[index], phrase) {
			parts := strings.Split(tags[index], ":")   // e.g. Land/Jahr: USA 2017
			trimmedPart := strings.TrimSpace(parts[1]) // e.g. "USA 2017"
			year := strings.Split(trimmedPart, " ")[1] //e.g. 2017
			return year, nil
		}
	}
	err := errors.New("not found")
	return "?", err
}

func getTagsFromList(tags []string, phrases []string) (string, error) {
	for _, phrase := range phrases {
		for index := range tags {
			if strings.Contains(tags[index], phrase) {
				return strings.TrimSpace(strings.Split(tags[index], ":")[1]), nil
			}
		}
	}
	err := errors.New("not found")
	return "?", err
}

func containsTag(tags []string, phrase string) bool {
	for index := range tags {
		if strings.Contains(tags[index], phrase) {
			// Phrase matched the value of the tag
			return true
		}
	}
	return false
}

func test(link string) {
	testInfo := parse4KUHD(link)

	fmt.Println("Link:", testInfo.url)
	fmt.Println("Title:", testInfo.title)
	fmt.Println("Real 4K:", testInfo.real4K)
	fmt.Println("Director:", testInfo.director)
	fmt.Println("Production Studio:", testInfo.studio)
	fmt.Println("Runtime:", testInfo.runtime)
	fmt.Println("Screen:", testInfo.screen)
	fmt.Println("Country:", testInfo.country)
	fmt.Println("Year:", testInfo.year)
	fmt.Println("Actors:", testInfo.actors)
	fmt.Println("HDR:", testInfo.hdr)
	fmt.Println("Video Code:", testInfo.codec)
	fmt.Println("Audio Format:", testInfo.audioFormat)
	fmt.Println("Picture Quality:", testInfo.pq)
}

func parse4KUHD(link string) Movie {
	// Create a Movie variable to store the parsed data
	var info Movie

	// Load the Go Query HTML Document
	html := getHTMLDocument(link)

	// Grab all the values of <p></p> tags
	pTags := getTags(html, "p")

	// Title
	info.title = getTitle(html)

	// URL
	info.url = link

	// 4K
	phrase = "Real 4K:"
	info.real4K, err = getTagWithPhrase(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// High Dynamic Range
	phrase = "High Dynamic Range:"
	info.hdr, err = getTagWithPhrase(pTags, phrase)
	if err != nil {

		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Director
	phrase = "Regie:"
	info.director, err = getTagWithPhrase(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Runtime
	phrase = "Laufzeit:"
	info.runtime, err = getTagWithPhrase(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Screen
	phrase = "Bildformat:"
	info.screen, err = getScreen(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Country
	phrase = "Land/"
	info.country, err = getCountry(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Year
	phrase = "/Jahr"
	info.year, err = getYear(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Actors
	phraseList = append(phraseList, "Darsteller:", "Sprecher:")
	info.actors, err = getTagsFromList(pTags, phraseList)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phraseList, err)
	}
	phraseList = nil // Go to garbage collector

	// Video Codec
	phraseList = append(phraseList, "Codec UHD:", "Codec UHD:", "Codec (UHD):", "Codec:")
	info.codec, err = getTagsFromList(pTags, phraseList)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phraseList, err)
	}
	phraseList = nil // Go to garbage collector

	// Audio
	phraseList = append(phraseList, "Tonformate UHD:", "Tonformate BD/UHD:", "Tonformate Blu-ray/UHD:", "Tonformate (UHD):", "UHD-Fassung:", "Tonformate:", "Tonformate BD:")
	info.audioFormat, err = getTagsFromList(pTags, phraseList)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phraseList, err)
	}
	phraseList = nil // Go to garbage collector

	// Picture Quality
	phraseList = append(phraseList, "Bildqualität UHD (HDR10):", "Bildqualität UHD (DV):", "Bildqualität UHD:")
	info.pq, err = getTagsFromList(getTags(html, "p strong"), phraseList)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phraseList, err)
	}
	phraseList = nil // Go to garbage collector

	// Studio
	phrase = "Anbieter:"
	info.studio, err = getTagWithPhrase(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	return info

}

func main() {
	url := "https://blu-ray-rezensionen.net/ultra-hd-blu-ray"
	html := getHTMLDocument(url)
	links := getLinksWithPhrase(html, "UHD")

	// Initialize a Movie slice with the desired length
	movies := make([]Movie, len(links))

	// Parse all the 4k UHD reviews
	for index, link := range links {
		movies[index] = parse4KUHD(link)
	}

	printMovieSlice(movies)

}
