package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("h", "0.0.0.0:8080", "TCP address to listen to")
	root     = flag.String("d", "/usr/share/httpd", "Directory to serve static files from")
	compress = flag.Bool("c", false, "Enables transparent response compression if set to true")
)

var indexBody []byte

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	if strings.HasSuffix(string(ctx.Request.RequestURI()), ".map") {
		ctx.Response.SetStatusCode(http.StatusNotFound)
		return
	}
	ctx.Logger().Printf("File %s not found, defaulting to index.html", ctx.Path())
	ctx.Response.SetBody(indexBody)
	ctx.Response.SetStatusCode(http.StatusOK)
}

func createFsHandler() fasthttp.RequestHandler {
	fs := &fasthttp.FS{
		Root:               *root,
		Compress:           *compress,
		IndexNames:         []string{"index.html"},
		PathNotFound:       notFoundHandler,
		GenerateIndexPages: false,
		AcceptByteRange:    true,
	}
	return fs.NewRequestHandler()
}

func main() {
	flag.Parse()

	var err error
	indexBody, err = ioutil.ReadFile(*root + "/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	fsHandler := createFsHandler()

	// Start HTTP server.
	if len(*addr) > 0 {
		log.Printf("Starting HTTP server on %q", *addr)
		log.Printf("Serving files from directory %q", *root)
		if err := fasthttp.ListenAndServe(*addr, fsHandler); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}
}
