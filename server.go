package main

import (
	"flag"
	"log"

	"github.com/valyala/fasthttp"
)

var (
	addr      = flag.String("h", "0.0.0.0:8080", "TCP address to listen to")
	root      = flag.String("d", "/usr/share/httpd", "Directory to serve static files from")
	strip     = flag.Int("s", 0, "Number of mount point to skip")
	compress  = flag.Bool("c", false, "Enables transparent response compression if set to true")
	fsHandler fasthttp.RequestHandler
)

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.Logger().Printf("File %s not found, defaulting to index.html", ctx.Path())
	ctx.Request.SetRequestURI("/index.html")
	fsHandler(ctx)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fsHandler(ctx)
}

func createFsHandler(stripSlashes int) fasthttp.RequestHandler {
	fs := &fasthttp.FS{
		Root:               *root,
		Compress:           *compress,
		IndexNames:         []string{"index.html"},
		PathNotFound:       notFoundHandler,
		GenerateIndexPages: false,
		AcceptByteRange:    true,
	}
	if stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewPathSlashesStripper(stripSlashes)
	}
	return fs.NewRequestHandler()
}

func main() {
	flag.Parse()
	fsHandler = createFsHandler(*strip)

	// Start HTTP server.
	if len(*addr) > 0 {
		log.Printf("Starting HTTP server on %q", *addr)
		log.Printf("Serving files from directory %q", *root)
		if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}
}
