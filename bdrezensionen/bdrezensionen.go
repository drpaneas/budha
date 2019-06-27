package bdrezensionen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/drpaneas/budha/goquerywrapper"
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

func GetScreen(tags []string, phrase string) (string, error) {
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

func GetCountry(tags []string, phrase string) (string, error) {
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

func GetYear(tags []string, phrase string) (string, error) {
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

func Parse4KUHD(link string) Movie {
	// Create a Movie variable to store the parsed data
	var info Movie

	// Load the Go Query HTML Document
	html := goquerywrapper.GetHTMLDocument(link)

	// Grab all the values of <p></p> tags
	pTags := goquerywrapper.GetTags(html, "p")

	// Title
	info.title = goquerywrapper.GetTitle(html)

	// URL
	info.url = link

	// 4K
	phrase = "Real 4K:"
	info.real4K, err = goquerywrapper.GetTagWithPhrase(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// High Dynamic Range
	phrase = "High Dynamic Range:"
	info.hdr, err = goquerywrapper.GetTagWithPhrase(pTags, phrase)
	if err != nil {

		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Director
	phrase = "Regie:"
	info.director, err = goquerywrapper.GetTagWithPhrase(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Runtime
	phrase = "Laufzeit:"
	info.runtime, err = goquerywrapper.GetTagWithPhrase(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Screen
	phrase = "Bildformat:"
	info.screen, err = GetScreen(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Country
	phrase = "Land/"
	info.country, err = GetCountry(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Year
	phrase = "/Jahr"
	info.year, err = GetYear(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	// Actors
	phraseList = append(phraseList, "Darsteller:", "Sprecher:")
	info.actors, err = goquerywrapper.GetTagsFromList(pTags, phraseList)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phraseList, err)
	}
	phraseList = nil // Go to garbage collector

	// Video Codec
	phraseList = append(phraseList, "Codec UHD:", "Codec UHD:", "Codec (UHD):", "Codec:")
	info.codec, err = goquerywrapper.GetTagsFromList(pTags, phraseList)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phraseList, err)
	}
	phraseList = nil // Go to garbage collector

	// Audio
	phraseList = append(phraseList, "Tonformate UHD:", "Tonformate BD/UHD:", "Tonformate Blu-ray/UHD:", "Tonformate (UHD):", "UHD-Fassung:", "Tonformate:", "Tonformate BD:")
	info.audioFormat, err = goquerywrapper.GetTagsFromList(pTags, phraseList)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phraseList, err)
	}
	phraseList = nil // Go to garbage collector

	// Picture Quality
	phraseList = append(phraseList, "Bildqualität UHD (HDR10):", "Bildqualität UHD (DV):", "Bildqualität UHD:")
	info.pq, err = goquerywrapper.GetTagsFromList(goquerywrapper.GetTags(html, "p strong"), phraseList)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phraseList, err)
	}
	phraseList = nil // Go to garbage collector

	// Studio
	phrase = "Anbieter:"
	info.studio, err = goquerywrapper.GetTagWithPhrase(pTags, phrase)
	if err != nil {
		fmt.Printf("[INFO] '%s' %s\n", phrase, err)
	}

	return info
}

func PrintMovieSlice(mySlice []Movie) {
	fmt.Printf("len=%d cap=%d\n", len(mySlice), cap(mySlice))
	for index, element := range mySlice {
		fmt.Printf("%4d : %v\n", index, element)
	}
	fmt.Printf("\n\n\n")
}
