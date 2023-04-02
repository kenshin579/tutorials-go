package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type PokemonProduct struct {
	url, image, name, price string
}

func main() {
	// scraping logic...
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		// turning on the asynchronous request mode in Colly
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 4 request at a time
		Parallelism: 4,
	})

	var pokemonProducts []PokemonProduct

	pagesToScrape := []string{
		"https://scrapeme.live/shop/page/1/",
		"https://scrapeme.live/shop/page/2/",
		"https://scrapeme.live/shop/page/47/",
		"https://scrapeme.live/shop/page/48/",
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
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
	})

	for _, pageToScrape := range pagesToScrape {
		c.Visit(pageToScrape)
	}

	c.Wait()

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
