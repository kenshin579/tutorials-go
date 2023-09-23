package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func Test_Nodes(t *testing.T) {
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

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

	var nodes []*cdp.Node
	selector := "#main ul li a"
	// selector := "#main > ul > li > h2 > a" //이렇게 표현을 동일하게 동작을 함
	pageURL := "https://notepad-plus-plus.org/downloads/"

	err := chromedp.Run(taskCtx, chromedp.Tasks{
		chromedp.Navigate(pageURL),
		chromedp.WaitReady(selector),
		chromedp.Nodes(selector, &nodes),
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, n := range nodes {
		u := n.AttributeValue("href")
		fmt.Printf("node: %s | href = %s\n", n.LocalName, u)
	}

}
