package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/suite"
)

type chromedpTestSuite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
}

func TestChromedpTestSuite(t *testing.T) {
	suite.Run(t, new(chromedpTestSuite))
}

func (suite *chromedpTestSuite) SetupSuite() {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36"),
	)

	// Create a Chrome instance
	suite.ctx, suite.cancel = chromedp.NewExecAllocator(ctx, opts...)

}

func (suite *chromedpTestSuite) Test_Chromedp() {
	suite.Run("Nodes", func() {
		// Create a new tab
		taskCtx, taskCancel := chromedp.NewContext(suite.ctx)
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
			suite.T().Fatal(err)
		}

		for _, n := range nodes {
			u := n.AttributeValue("href")
			fmt.Printf("node: %s | href = %s\n", n.LocalName, u)
		}
	})

	suite.Run("Tag Exists", func() {
		// Create a new tab
		taskCtx, taskCancel := chromedp.NewContext(suite.ctx)
		defer taskCancel()

		selector := `"#js-category-content > div.js-symbol-page-tab-overview-root > div > section > div:nth-child(3) > div.content-gdSWdaJr > div.container-GRoarMHL > div:nth-child(1) > div.wrapper-GgmpMpKr > a > div"`
		pageURL := "https://www.tradingview.com/symbols/KRX-069620/"

		const expr = `(function(d, sel) {
		var element = d.querySelector(sel);
        return !!element;
	})(document, %s);`

		err := chromedp.Run(taskCtx, chromedp.Tasks{
			chromedp.Navigate(pageURL),
			chromedp.WaitReady("body"),
			chromedp.ActionFunc(func(ctx context.Context) error {
				s := fmt.Sprintf(expr, selector)
				result, exp, err := runtime.Evaluate(s).Do(ctx)
				fmt.Printf("result:%s, exp:%v, err:%v\n", result.Value, exp, err)
				if err != nil {
					return err
				}

				if exp != nil {
					return exp
				}
				return nil
			}),
		})
		if err != nil {
			suite.T().Fatal(err)
		}

	})
}
