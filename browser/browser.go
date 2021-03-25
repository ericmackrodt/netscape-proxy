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
	"github.com/chromedp/chromedp/kb"
)

func arrangeMap(oldMap map[rune]*kb.Key) map[string]string {
	newMap := make(map[string]string)
	for key, value := range oldMap {
		newMap[value.Key] = string(key)
	}
	return newMap
}

var (
	ctx       context.Context
	localKeys = arrangeMap(kb.Keys)
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
	var result string
	for _, i := range codes {
		if i == "@@nbsp@@" {
			result = result + localKeys[" "]
		} else if i == "@@comma@@" {
			result = result + localKeys[","]
		} else if i == "@@enter@@" {
			result = result + localKeys["Enter"]
		} else if i == "@@backspace@@" {
			result = result + localKeys["Backspace"]
		} else {
			result = result + localKeys[i]
		}
	}

	println(result)

	chromedp.Run(ctx, chromedp.KeyEvent(result))
}
