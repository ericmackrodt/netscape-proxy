package browser

import (
	"github.com/chromedp/chromedp/kb"
)

func GetArrangedKeys() map[string]string {
	newMap := make(map[string]string)
	for key, value := range kb.Keys {
		newMap[value.Key] = string(key)
	}
	return newMap
}
