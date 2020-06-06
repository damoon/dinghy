package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	//	opts := append(
	//		chromedp.DefaultExecAllocatorOptions[:],
	//		chromedp.Flag("headless", false),
	//	)
	//	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	//	defer cancel()

	ctx := context.Background()

	ctx, cancel := chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(`http://dinghy-backend:8080`),
		//		chromedp.Click(`.//*[contains(text(),'subfolder')]`, chromedp.NodeVisible),
		//		chromedp.WaitVisible(`.//*[contains(text(),'favicon.png')]`),
	)
	if err != nil {
		log.Fatal(err)
	}
}
