package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	browser "webproxy/browser"

	"github.com/go-chi/chi/v5"
)

var (
	// flagPort is the open port the application listens on
	flagPort   = flag.String("port", "8080", "Port to listen on")
	homeTpl    *template.Template
	browserTpl *template.Template
)

type BrowserPage struct {
	Title   string
	PageId  string //exported field since it begins with a capital letter
	ImageId string
}

func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

func getResolution(request *http.Request) (int64, int64) {
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

func home(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	// always modified right now
	// response.Header().Set("Last-Modified",   . gmdate("D, d M Y H:i:s") . " GMT");
	// HTTP/1.1
	response.Header().Set(
		"Cache-Control",
		"private, no-store, max-age=0, no-cache, must-revalidate, post-check=0, pre-check=0",
	)
	// HTTP/1.0
	response.Header().Set("Pragma", "no-cache")

	if request.Method == "POST" {
		request.ParseForm()
		url := request.Form.Get("url")

		width, height := getResolution(request)

		browser.NewPage(url, width, height)

		http.Redirect(response, request, "/abc/0", http.StatusMovedPermanently)
		return
	}

	render(response, request, homeTpl, "homepage_view", nil)
}

func screenshot(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	pageId := chi.URLParam(request, "pageId")
	keypresses := query.Get("keypresses")
	if keypresses != "" {
		println(keypresses)
		browser.TypeText(pageId, keypresses)
	}

	width, height := getResolution(request)
	data := browser.Screenshot(width, height)
	response.Write(data)
}

func browsePage(response http.ResponseWriter, request *http.Request) {
	pageId := chi.URLParam(request, "pageId")
	data := BrowserPage{PageId: pageId, ImageId: "0", Title: "Something"}

	render(response, request, browserTpl, "browser_view", data)
}

func click(response http.ResponseWriter, request *http.Request) {
	pageId := chi.URLParam(request, "pageId")
	var x int64
	var y int64
	fmt.Sscanf(request.URL.RawQuery, "%d,%d", &x, &y)

	println(x, y)
	browser.Click(pageId, float64(x), float64(y))

	http.Redirect(response, request, "/abc/0", http.StatusMovedPermanently)
}

func getTemplate(filename string) string {
	content, _ := ioutil.ReadFile(filename)
	return string(content)
}

func main() {
	browser.Initialize()

	homeHTML := getTemplate("templates/home.html")
	homeTpl = template.Must(template.New("homepage_view").Parse(homeHTML))

	browserHTML := getTemplate("templates/browser.html")
	browserTpl = template.Must(template.New("browser_view").Parse(browserHTML))

	router := chi.NewRouter()

	router.Get("/click/{pageId}", click)
	router.Get("/{pageId}/{num}", browsePage)
	router.Get("/ss/{pageId}/{imageId}", screenshot)
	router.HandleFunc("/", home)

	fmt.Println("Starting Proxy Server")

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+*flagPort, router))
}
