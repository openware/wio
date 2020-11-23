package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("h", "0.0.0.0:8080", "TCP address to listen to")
	root     = flag.String("d", "/usr/share/httpd", "Directory to serve static files from")
	prefix   = flag.String("p", "/", "URL prefix used to access the server")
	compress = flag.Bool("c", false, "Enables transparent response compression if set to true")
)

var fsHandler fasthttp.RequestHandler

func reponseNotFound(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(http.StatusNotFound)
	ctx.Response.SetBody([]byte("Not Found\n"))
}

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	if strings.HasSuffix(string(ctx.Request.RequestURI()), ".map") {
		reponseNotFound(ctx)
		return
	}
	ctx.Logger().Printf("File %s not found, defaulting to index.html", ctx.Path())
	ctx.Request.SetRequestURI("/index.html")
	fsHandler(ctx)
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
	fsHandler = createFsHandler()

	handler := func(ctx *fasthttp.RequestCtx) {
		uri := string(ctx.RequestURI())
		if strings.HasPrefix(uri, *prefix) {
			uri = uri[len(*prefix):]
			ctx.Request.SetRequestURI(uri)
			fsHandler(ctx)
			return
		}
		reponseNotFound(ctx)
	}

	// Start HTTP server.
	if len(*addr) > 0 {
		log.Printf("Starting HTTP server on %q", *addr)
		log.Printf("Serving files from directory %q", *root)
		if err := fasthttp.ListenAndServe(*addr, handler); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}
}
