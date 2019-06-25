package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yosssi/gohtml"
)

func getTitleHTML(data string) (pageTitle string) {
	// Find a substring
	titleStartIndex := strings.Index(data, "<title>")
	if titleStartIndex == -1 {
		fmt.Println("No title element found")
		os.Exit(0)
	}

	// The start index of the title is the index of the first character, the < symbol.
	// We don't want to include <title> as part of the final value, so let's offset the index
	// by the number of characters in <title> that is 7 chararacters
	titleStartIndex += 7

	// Find the index of the closing tag
	titleEndIndex := strings.Index(data, "</title>")
	if titleEndIndex == -1 {
		fmt.Println("No closing tag for title found.")
		os.Exit(0)
	}

	// Copy the substring into a separate variable
	// so the variables with the full document data can be garbage collected
	pageTitle = string([]byte(data[titleStartIndex:titleEndIndex]))
	return pageTitle
}

func getLinks(website string, anchor string) (string, error) {

	// Create a goquery document from the HTTP response
	doc, err := goquery.NewDocument(website)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
		os.Exit(1)
	}

	fmt.Println(doc)

	// Find all links and process them
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		// fmt.Printf("link: %s - anchor text: %s\n", href, item.Text())
		if strings.Contains(string(item.Text()), anchor) {
			fmt.Printf("link: %s - anchor text: %s\n", href, item.Text())
			return href
		}
	})
	return "", nil
}

func main() {
	// Call http.Get with the URL we want to retrieve
	res, err := http.Get("https://blu-ray-rezensionen.net/ultra-hd-blu-ray/")
	if err != nil {
		log.Fatal(err)
	}
	// Release the network connection once the "main" function exists
	defer res.Body.Close()

	// Read all the data in the response
	dataInBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	dataInString := string(dataInBytes)
	dataInFormattedString := gohtml.Format(dataInString)

	// Convert the data to a string and print it
	fmt.Println(dataInFormattedString)

	// Print out the result
	fmt.Println(getTitleHTML(dataInFormattedString))

}
