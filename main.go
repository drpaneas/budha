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
	err := errors.New("Tag not found")
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
	err := errors.New("Tag not found")
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
	err := errors.New("Tag not found")
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
	err := errors.New("Tag not found")
	return "?", err
}

func getTagsFromList(tags []string, phrases []string) (string, error) {
	for _, phrase := range phrases {
		for index := range tags {
			if strings.Contains(tags[index], phrase) {
				// Phrase matched the value of the tag
				return strings.TrimSpace(strings.Split(tags[index], ":")[1]), nil
			}
		}
	}
	err := errors.New("Tag not found")
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
	fmt.Println(link)
	html := getHTMLDocument(link)
	fmt.Println(getTitle(html))

	// Test tags
	printSlice(getTags(html, "p"))
	pTags := getTags(html, "p")

	// Real 4K
	real4K, err := getTagWithPhrase(pTags, "Real 4K:")
	if err != nil {
		fmt.Printf("Couldn't find 'Real 4K:': %s", err)
	}

	// High Dynamic Range
	hdr, err := getTagWithPhrase(pTags, "High Dynamic Range:")
	if err != nil {
		fmt.Printf("Couldn't find 'High Dynamic Range:': %s", err)
	}

	// Director
	director, err := getTagWithPhrase(pTags, "Regie:")
	if err != nil {
		fmt.Printf("Couldn't find 'Regie:': %s", err)
	}

	// Studio
	studio, err := getTagWithPhrase(pTags, "Anbieter:")
	if err != nil {
		fmt.Printf("Couldn't find 'Anbieter:': %s", err)
	}

	// Runtime
	runtime, err := getTagWithPhrase(pTags, "Laufzeit:")
	if err != nil {
		fmt.Printf("Couldn't find 'Laufzeit:': %s", err)
	}

	// Screen
	screen, err := getScreen(pTags, "Bildformat:")
	if err != nil {
		fmt.Printf("Couldn't find 'Bildformat:': %s", err)
	}

	// Country
	country, err := getCountry(pTags, "Land/")
	if err != nil {
		fmt.Printf("Couldn't find 'Land/': %s", err)
	}

	// Year
	year, err := getYear(pTags, "/Jahr")
	if err != nil {
		fmt.Printf("Couldn't find '/Jahr': %s", err)
	}

	// Actors
	actors, err := getTagsFromList(pTags, []string{"Darsteller:", "Sprecher:"})
	if err != nil {
		fmt.Printf("Couldn't find actors : %s", err)
	}

	// Video Codec
	codec, err := getTagsFromList(pTags, []string{"Codec UHD:", "Code UHD:", "Codec (UHD):", "Codec:"})
	if err != nil {
		fmt.Printf("Couldn't find codec : %s", err)
	}

	// Audio
	audio, err := getTagsFromList(pTags, []string{"Tonformate UHD:", "Tonformate BD/UHD:", "Tonformate Blu-ray/UHD:", "Tonformate (UHD):", "UHD-Fassung", "Tonformate:", "Tonformate BD:"})
	if err != nil {
		fmt.Printf("Couldn't find audio : %s", err)
	}

	// Picture Quality
	pq, err := getTagsFromList(getTags(html, "p strong"), []string{"Bildqualität UHD (HDR10)", "Bildqualität UHD (DV)", "Bildqualität UHD:"})
	if err != nil {
		fmt.Printf("Couldn't find pq : %s", err)
	}

	// Title
	title := getTitle(html)

	fmt.Println("Link:", link)
	fmt.Println("Title:", title)
	fmt.Println("Real 4K:", real4K)
	fmt.Println("Director:", director)
	fmt.Println("Production Studio:", studio)
	fmt.Println("Runtime:", runtime)
	fmt.Println("Screen:", screen)
	fmt.Println("Country:", country)
	fmt.Println("Year:", year)
	fmt.Println("Actors:", actors)
	fmt.Println("HDR:", hdr)
	fmt.Println("Video Code:", codec)
	fmt.Println("Audio Format:", audio)
	fmt.Println("Picture Quality:", pq)
}

func main() {
	var err error
	url := "https://blu-ray-rezensionen.net/ultra-hd-blu-ray"
	html := getHTMLDocument(url)
	// See the HTML code
	//html := HtmlToStr(htmlDocument)

	links := getLinksWithPhrase(html, "UHD")
	numberLinks := len(links)
	test("https://blu-ray-rezensionen.net/india-4k-uhd/") // test it

	// Common mistake: You can't instantiate an array like that with a value calculated at runtime
	// movies := [numberLinks]Movie{} // non-constant array bound numberLinksgo

	// Initialize a slice with the desired length
	movies := make([]Movie, numberLinks)

	// Risk-free and safe loop to process each element of an array
	for index, link := range links {
		fmt.Println(link)
		movies[index].url = link
		html = getHTMLDocument(link)
		pTags := getTags(html, "p")

		// Title
		movies[index].title = getTitle(html)

		// 4K
		movies[index].real4K, err = getTagWithPhrase(pTags, "Real 4K:")
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// High Dynamic Range
		movies[index].hdr, err = getTagWithPhrase(pTags, "High Dynamic Range:")
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Director
		movies[index].director, err = getTagWithPhrase(pTags, "Regie:")
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Studio
		movies[index].studio, err = getTagWithPhrase(pTags, "Anbieter:")
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Runtime
		movies[index].runtime, err = getTagWithPhrase(pTags, "Laufzeit:")
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Screen
		movies[index].screen, err = getScreen(pTags, "Bildformat:")
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Country
		movies[index].country, err = getCountry(pTags, "Land/")
		if err != nil {
			fmt.Println("Couldn't find")
		}
		// Year
		movies[index].year, err = getYear(pTags, "/Jahr")
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Actors
		movies[index].actors, err = getTagsFromList(pTags, []string{"Darsteller:", "Sprecher:"})
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Video Codec
		movies[index].codec, err = getTagsFromList(pTags, []string{"Codec UHD:", "Code UHD:", "Codec (UHD):", "Codec:"})
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Audio
		movies[index].audioFormat, err = getTagsFromList(pTags, []string{"Tonformate UHD:", "Tonformate BD/UHD:", "Tonformate Blu-ray/UHD:", "Tonformate (UHD):", "UHD-Fassung", "Tonformate:", "Tonformate BD:"})
		if err != nil {
			fmt.Println("Couldn't find")
		}

		// Picture Quality
		movies[index].pq, err = getTagsFromList(getTags(html, "p strong"), []string{"Bildqualität UHD (HDR10)", "Bildqualität UHD (DV)", "Bildqualität UHD:"})
		if err != nil {
			fmt.Println("Couldn't find")
		}
	}

	printMovieSlice(movies)

}
