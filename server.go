package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	browser "webproxy/browser"
	serverHelpers "webproxy/helpers"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func homeRoute(response http.ResponseWriter, request *http.Request) {
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

		width, height := serverHelpers.GetResolution(request)

		pageId := browser.NewPage(url, width, height)

		http.Redirect(response, request, "/"+pageId+"/"+uuid.New().String(), http.StatusMovedPermanently)
		return
	}

	serverHelpers.Render(response, request, homeTpl, "homepage_view", nil)
}

func screenshotRoute(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	pageId := chi.URLParam(request, "pageId")
	keypresses := query.Get("keypresses")
	if keypresses != "" {
		browser.TypeText(pageId, keypresses)
	}

	width, height := serverHelpers.GetResolution(request)
	data := browser.Screenshot(pageId, width, height)
	response.Write(data)
}

func browserRoute(response http.ResponseWriter, request *http.Request) {
	pageId := chi.URLParam(request, "pageId")
	data := BrowserPage{PageId: pageId, ImageId: "0", Title: "Something"}

	serverHelpers.Render(response, request, browserTpl, "browser_view", data)
}

func clickRoute(response http.ResponseWriter, request *http.Request) {
	pageId := chi.URLParam(request, "pageId")
	var x int64
	var y int64
	fmt.Sscanf(request.URL.RawQuery, "%d,%d", &x, &y)

	browser.Click(pageId, float64(x), float64(y))

	http.Redirect(response, request, "/"+pageId+"/"+uuid.New().String(), http.StatusMovedPermanently)
}

func main() {
	browser.Initialize()

	homeHTML := serverHelpers.GetTemplate("templates/home.html")
	homeTpl = template.Must(template.New("homepage_view").Parse(homeHTML))

	browserHTML := serverHelpers.GetTemplate("templates/browser.html")
	browserTpl = template.Must(template.New("browser_view").Parse(browserHTML))

	router := chi.NewRouter()

	router.Get("/click/{pageId}", clickRoute)
	router.Get("/{pageId}/{num}", browserRoute)
	router.Get("/ss/{pageId}/{imageId}", screenshotRoute)
	router.HandleFunc("/", homeRoute)

	fmt.Println("Starting Proxy Server")

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+*flagPort, router))
}
