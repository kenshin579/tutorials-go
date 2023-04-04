package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/thoas/go-funk"
)

func main() {
	// scraping logic...
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	var pagesToScrape []string

	pageToScrape := "https://scrapeme.live/shop/"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})
	// iterating over the list of pagination links to implement the crawling logic
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.Attr("href")

		// if the page discovered is new
		if !funk.Contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !funk.Contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	// visiting the first page
	c.Visit(pageToScrape)

	fmt.Println(pagesDiscovered)
}
