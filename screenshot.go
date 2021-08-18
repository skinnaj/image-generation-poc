package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
)

func screenshotImage(n string) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(`http://localhost:8080`, `//div[contains(@id, "bar-chart")]`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(n, buf, 0o644); err != nil {
		log.Fatal(err)
	}
}

func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
