package main

import (
  "log"

  "github.com/valyala/fasthttp"
)

var (
  filesHandler = fasthttp.FSHandler("build", 1)
)

// Main request handler
func requestHandler(ctx *fasthttp.RequestCtx) {
  switch {

  default:
    filesHandler(ctx)
  }
}

func main() {
  if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
    log.Fatalf("Error in server: %s", err)
  }
}
