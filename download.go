package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func downloadImage(n string) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	downloadComplete := make(chan bool)
	var downloadGUID string

	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			fmt.Printf("current download state: %s\n", ev.State.String())
			if ev.State == browser.DownloadProgressStateCompleted {
				downloadGUID = ev.GUID
				close(downloadComplete)
			}
		}
	})

	if err := chromedp.Run(ctx,
		emulation.SetDeviceMetricsOverride(950, 600, 2.0, false),
		chromedp.Navigate(`http://localhost:8080`),
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(".").
			WithEventsEnabled(true),

		chromedp.Click(`//button[contains(@id, "download")]`, chromedp.NodeVisible),
	); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		// Note: Ignoring the net::ERR_ABORTED page error is essential here since downloads
		// will cause this error to be emitted, although the download will still succeed.
		log.Fatal(err)
	}

	<-downloadComplete

	os.Rename(fmt.Sprintf("./%s", downloadGUID), fmt.Sprintf("./%s", n))

	log.Printf(fmt.Sprintf("Download Complete: %s", n))
}
