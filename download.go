package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

func downloadImage(n string) {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout as a safety net to prevent any infinite wait loops
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// set up a channel so we can block later while we monitor the download progress
	downloadComplete := make(chan bool)

	// this will be used to capture the file name later
	var downloadGUID string

	// set up a listener to watch the download events and close the channel when complete
	// this could be expanded to handle multiple downloads through creating a guid map,
	// monitor download urls via EventDownloadWillBegin, etc
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			fmt.Printf("current download state: %s\n", ev.State.String())
			if ev.State == browser.DownloadProgressStateCompleted {
				downloadGUID = ev.GUID
				close(downloadComplete)
			}
		}
	})

	// download the zip of the chromedp/examples repo from github. we use a link click method here
	// but this could also be done with a chromedp.Navigate task which points directly at the file we
	// want to download, as long as you run browser.SetDownloadBehavior first
	if err := chromedp.Run(ctx,
		// navigate to the page
		chromedp.Navigate(`http://bl.ocks.org/Rokotyan/raw/0556f8facbaf344507cdc45dc3622177`),

		// configure headless browser downloads. note that SetDownloadBehaviorBehaviorAllowAndName is
		// preferred here over SetDownloadBehaviorBehaviorAllow so that the file will be named as the GUID.
		// please note that it only works with 92.0.4498.0 or later due to issue 1204880,
		// see https://bugs.chromium.org/p/chromium/issues/detail?id=1204880
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(".").
			WithEventsEnabled(true),

		chromedp.Click(`//button[contains(@id, "saveButton")]`, chromedp.NodeVisible),
	); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		// Note: Ignoring the net::ERR_ABORTED page error is essential here since downloads
		// will cause this error to be emitted, although the download will still succeed.
		log.Fatal(err)
	}

	// This will block until the chromedp listener closes the channel
	<-downloadComplete

	os.Rename(fmt.Sprintf("./%s", downloadGUID), fmt.Sprintf("./%s", n))

	// We can predict the exact file location and name here because of how we configured
	// SetDownloadBehavior and WithDownloadPath
	log.Printf("Download Complete: image.png")
}