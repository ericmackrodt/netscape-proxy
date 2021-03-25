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

type PageDefinition struct {
	PageCtx *context.Context
	Width   int64
	Height  int64
	Count   int64
}

var (
	localKeys      = GetArrangedKeys()
	browserContext context.Context
	pages          = make(map[string]PageDefinition)
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

	pages[id] = PageDefinition{PageCtx: &ctx, Width: width, Height: height, Count: 0}

	chromedp.Run(ctx,
		emulation.SetDeviceMetricsOverride(width, height, 1, false),
		chromedp.Navigate(url),
	)

	return id
}

func GetPage(id string) (*context.Context, int64, int64, int64) {
	page := pages[id]
	return page.PageCtx, page.Width, page.Height, page.Count
}

func Screenshot(id string) []byte {
	var captureBytes []byte

	page, width, height, _ := GetPage(id)

	chromedp.Run(*page, emulation.SetDeviceMetricsOverride(width, height, 1, false))

	chromedp.Run(*page, chromedp.CaptureScreenshot(&captureBytes))

	raw, _ := png.Decode(bytes.NewReader(captureBytes))

	var jpegImage bytes.Buffer

	jpeg.Encode(&jpegImage, raw, &jpeg.Options{Quality: 50})

	return jpegImage.Bytes()
}

func Click(id string, x float64, y float64) {
	page, _, _, _ := GetPage(id)

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

	page, _, _, _ := GetPage(id)

	chromedp.Run(*page, chromedp.KeyEvent(result))
}

func ResizeBrowser(id string, newWidth int64, newHeight int64) {
	page, width, height, Count := GetPage(id)

	if height != newHeight && width != newWidth {
		println("resolution changed")

		chromedp.Run(*page,
			emulation.SetDeviceMetricsOverride(newWidth, newHeight, 1, false),
		)

		pages[id] = PageDefinition{Width: newWidth, Height: newHeight, PageCtx: page, Count: Count}
	}
}

func PageCount(id string) int64 {
	page, width, height, Count := GetPage(id)
	count := Count + 1
	pages[id] = PageDefinition{Width: width, Height: height, PageCtx: page, Count: count}
	return count
}

func GoBack(id string) {
	page, _, _, _ := GetPage(id)
	chromedp.Run(*page,
		chromedp.NavigateBack(),
	)
}
