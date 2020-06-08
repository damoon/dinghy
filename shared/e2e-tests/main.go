package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	ctx := context.Background()
	// enable to show browser window
	//	opts := append(
	//		chromedp.DefaultExecAllocatorOptions[:],
	//		chromedp.Flag("headless", false),
	//	)
	//	ctx, cancel2 := chromedp.NewExecAllocator(ctx, opts...)
	//	defer cancel2()

	ctx, cancel := chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(`http://backend:8080`),
		//		chromedp.Sleep(2*time.Second),
		chromedp.Click(`.//*[contains(text(),'subfolder')]`, chromedp.NodeVisible),
		chromedp.Click(`.//*[contains(text(),'hello-world.txt')]`, chromedp.NodeVisible),
		chromedp.WaitVisible(`.//*[contains(text(),'Hello, world.')]`),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}
