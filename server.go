package main

import (
  "flag"
  "log"

  "github.com/valyala/fasthttp"
)

var (
  addr  = flag.String("h", "localhost:8080", "TCP address to listen to")
  root  = flag.String("d", "/usr/share/httpd", "Directory to serve static files from")
  strip = flag.Int("s", 0, "Number of mount point to skip")
  // compress  = flag.Bool("compress", false, "Enables transparent response compression if set to true")
)

func main() {
  flag.Parse()

  fsHandler := fasthttp.FSHandler(*root, *strip)
  requestHandler := func(ctx *fasthttp.RequestCtx) {
    fsHandler(ctx)
  }

  // Start HTTP server.
  if len(*addr) > 0 {
    log.Printf("Starting HTTP server on %q", *addr)
    log.Printf("Serving files from directory %q", *root)
    if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
      log.Fatalf("error in ListenAndServe: %s", err)
    }
  }
}
