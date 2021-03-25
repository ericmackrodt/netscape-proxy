package browser

import (
	"bytes"
	"context"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

var (
	localKeys      = GetArrangedKeys()
	browserContext context.Context
	pages          = make(map[string]*context.Context)
)

func Initialize() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("hide-scrollbars", false),
	)

	browserContext, _ = chromedp.NewExecAllocator(context.Background(), opts...)
}

func NewPage(url string, width int64, height int64) string {
	id := uuid.New().String()

	ctx, _ := chromedp.NewContext(browserContext) //, chromedp.WithDebugf(log.Printf))

	pages[id] = &ctx

	chromedp.Run(ctx,
		emulation.SetDeviceMetricsOverride(width, height, 1, false),
		chromedp.Navigate(url),
	)

	return id
}

func GetPage(id string) *context.Context {
	return pages[id]
}

func Screenshot(id string, width int64, height int64) []byte {
	var captureBytes []byte

	page := GetPage(id)

	chromedp.Run(*page, emulation.SetDeviceMetricsOverride(width, height, 1, false))

	chromedp.Run(*page, chromedp.CaptureScreenshot(&captureBytes))

	raw, _ := png.Decode(bytes.NewReader(captureBytes))

	var jpegImage bytes.Buffer

	jpeg.Encode(&jpegImage, raw, &jpeg.Options{Quality: 10})

	return jpegImage.Bytes()
}

func Click(id string, x float64, y float64) {
	page := GetPage(id)

	chromedp.Run(*page, chromedp.MouseClickXY(x, y))
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

	page := GetPage(id)

	chromedp.Run(*page, chromedp.KeyEvent(result))
}
