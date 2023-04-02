package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
	"github.com/thoas/go-funk"
)

type PokemonProduct struct {
	url, image, name, price string
}

func main() {
	// scraping logic...
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	var pokemonProducts []PokemonProduct
	var pagesToScrape []string

	pageToScrape := "https://scrapeme.live/shop/page/1/"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	i := 1
	limit := 5

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

	// scraping logic
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}

		pokemonProduct.url = e.ChildAttr("a", "href")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")

		// until there is still a page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			// getting the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// incrementing the iteration counter
			i++

			// visiting a new page
			c.Visit(pageToScrape)
		}
	})

	// visiting the first page
	c.Visit(pageToScrape)

	// convert the data to CSV...
	createCSV("products.csv", pokemonProducts)
}

func createCSV(filename string, products []PokemonProduct) {
	// opening the CSV file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// writing the CSV headers
	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}
	writer.Write(headers)

	// writing each Pokemon product as a CSV row
	for _, pokemonProduct := range products {
		// converting a PokemonProduct to an array of strings
		record := []string{
			pokemonProduct.url,
			pokemonProduct.image,
			pokemonProduct.name,
			pokemonProduct.price,
		}

		// adding a CSV record to the output file
		writer.Write(record)
	}
	defer writer.Flush()
}
