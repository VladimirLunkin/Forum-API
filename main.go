package main

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
)

func Index(ctx *fasthttp.RequestCtx) {
	a := "Просто привет"
	b, _ := json.Marshal(a)
	ctx.SetBody(b)
}

func main() {
	r := router.New()
	r.GET("/", Index)

	log.Fatal(fasthttp.ListenAndServe(":8000", r.Handler))
}
