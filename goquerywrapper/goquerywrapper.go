package goquerywrapper

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yosssi/gohtml"
)

func GetHTMLDocument(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal("Cannot load the HTML document:", err)
	}
	return doc
}

func HTMLToStr(htmlDocument *goquery.Document) string {
	htmlCode, err := htmlDocument.Html()
	if err != nil {
		log.Fatalf("Error while converting HTML Document into string: %s", err)
	}
	htmlPrettyCode := gohtml.Format(htmlCode)
	return htmlPrettyCode
}

func PrintSlice(mySlice []string) {
	fmt.Printf("len=%d cap=%d\n", len(mySlice), cap(mySlice))
	for index, element := range mySlice {
		fmt.Printf("%4d : %v\n", index, element)
	}
}

func GetLinksWithPhrase(htmlDocument *goquery.Document, phrase string) []string {
	linkSlice := []string{}
	htmlDocument.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(string(s.Text()), phrase) {
			linkSlice = append(linkSlice, href)
		}
	})
	return linkSlice
}

func GetTitle(htmlDocument *goquery.Document) string {
	var title string
	htmlDocument.Find("title").Each(func(i int, s *goquery.Selection) {
		title = string(s.Text())
	})
	return title
}

func GetTags(htmlDocument *goquery.Document, tag string) []string {
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

func GetTagWithPhrase(tags []string, phrase string) (string, error) {
	for index := range tags {
		if strings.Contains(tags[index], phrase) {
			// Phrase matched the value of the tag
			return strings.TrimSpace(strings.Split(tags[index], ":")[1]), nil
		}
	}
	err := errors.New("not found")
	return "?", err
}

func GetTagsFromList(tags []string, phrases []string) (string, error) {
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

func ContainsTag(tags []string, phrase string) bool {
	for index := range tags {
		if strings.Contains(tags[index], phrase) {
			// Phrase matched the value of the tag
			return true
		}
	}
	return false
}
