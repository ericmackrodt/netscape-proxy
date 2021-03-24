package browser

import (
	"bytes"
	"context"
	"strings"

	// 	"fmt"
	// 	"io/ioutil"

	// 	"math"
	"image/jpeg"
	"image/png"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

var (
	ctx context.Context
	// cancel context.CancelFunc
	// acancel context.CancelFunc
	typeQueue []string
)

func Initialize() context.Context {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("hide-scrollbars", false),
	)

	actx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// defer acancel()
	ctx, _ = chromedp.NewContext(actx) //, chromedp.WithDebugf(log.Printf))
	// defer cancel()
	return ctx
}

func NewPage(url string, width int64, height int64) {
	chromedp.Run(ctx,
		emulation.SetDeviceMetricsOverride(width, height, 1, false),
		chromedp.Navigate(url),
	)
}

func Screenshot(width int64, height int64) []byte {
	var captureBytes []byte

	chromedp.Run(ctx, emulation.SetDeviceMetricsOverride(width, height, 1, false))

	chromedp.Run(ctx, chromedp.CaptureScreenshot(&captureBytes))

	raw, _ := png.Decode(bytes.NewReader(captureBytes))

	var jpegImage bytes.Buffer

	jpeg.Encode(&jpegImage, raw, &jpeg.Options{Quality: 10})

	return jpegImage.Bytes()
}

func Click(id string, x float64, y float64) {
	chromedp.Run(ctx, chromedp.MouseClickXY(x, y))
}

func TypeText(id string, keys string) {
	codes := strings.Split(keys, ",")

	for _, i := range codes {
		chromedp.Run(ctx, chromedp.KeyEvent(i))
	}
}
