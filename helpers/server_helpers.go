package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	browser "webproxy/browser"
)

func Render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

func GetResolution(request *http.Request) (int64, int64) {
	resolution, _ := request.Cookie("resolution")

	res := strings.Split(resolution.Value, "x")

	var resInt = []int64{}

	for _, i := range res {
		j, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		resInt = append(resInt, j)
	}

	return resInt[0], resInt[1]
}

func GetTemplate(filename string) string {
	content, _ := ioutil.ReadFile(filename)
	return string(content)
}

func Redirect(w http.ResponseWriter, r *http.Request, pageId string) {
	count := browser.PageCount(pageId)
	http.Redirect(w, r, "/"+pageId+"/"+strconv.FormatInt(count, 10), http.StatusMovedPermanently)
}
