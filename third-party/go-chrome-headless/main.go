package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	// Create a context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set up Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36"),
	)

	// Create a Chrome instance
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	// Create a new tab
	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	// Navigate to the target website
	err := chromedp.Run(taskCtx, chromedp.Tasks{
		chromedp.Navigate("https://example.com"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the page to load completely
	err = chromedp.Run(taskCtx, chromedp.Tasks{
		chromedp.WaitReady("body"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Extract data from the page
	var pageTitle string
	err = chromedp.Run(taskCtx, chromedp.Tasks{
		chromedp.Title(&pageTitle),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Page Title:", pageTitle)

	// Get the text value of an element by its CSS selector
	var elementText string
	err = chromedp.Run(taskCtx, chromedp.Tasks{
		chromedp.Text("h1", &elementText, chromedp.NodeVisible, chromedp.ByQuery),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Element Text:", elementText)
}
