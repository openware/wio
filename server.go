package main

import (
	"os"
	"flag"
	"log"

	"github.com/valyala/fasthttp"
)

var (
	addr      = flag.String("h", "0.0.0.0:8080", "TCP address to listen to")
	root      = flag.String("d", "/usr/share/httpd", "Directory to serve static files from")
	file      = flag.String("f", "index.html", "File to serve")
	strip     = flag.Int("s", 0, "Number of mount point to skip")
	compress  = flag.Bool("c", false, "Enables transparent response compression if set to true")
	fsHandler fasthttp.RequestHandler
)

func notFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.Logger().Printf("File %s not found, defaulting to index.html", ctx.Path())
	if isIndexFileExists() {
		ctx.Request.SetRequestURI("/index.html")
		fsHandler(ctx)
	} else {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
	}
}

func isIndexFileExists() bool {
	path := *root + "/" + *file
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Printf("File %s do not exists", *file)
			return false
		}
	}
	return true
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fsHandler(ctx)
}

func createFsHandler(stripSlashes int) fasthttp.RequestHandler {
	fs := &fasthttp.FS {
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
	if len(*addr) > 0 && isIndexFileExists() {
		log.Printf("Starting HTTP server on %q", *addr)
		log.Printf("Serving files from directory %q", *root)
		if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}
}
