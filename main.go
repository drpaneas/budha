package main

import (
	"github.com/drpaneas/budha/bdrezensionen"
	"github.com/drpaneas/budha/goquerywrapper"
)

/*
func test(link string) {
	testInfo := bdrezensionen.Parse4KUHD(link)

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
*/

func main() {
	url := "https://blu-ray-rezensionen.net/ultra-hd-blu-ray"
	html := goquerywrapper.GetHTMLDocument(url)
	links := goquerywrapper.GetLinksWithPhrase(html, "UHD")

	// Initialize a Movie slice with the desired length
	movies := make([]bdrezensionen.Movie, len(links))

	// Parse all the 4k UHD reviews
	for index, link := range links {
		movies[index] = bdrezensionen.Parse4KUHD(link)
	}

	bdrezensionen.PrintMovieSlice(movies)

}
