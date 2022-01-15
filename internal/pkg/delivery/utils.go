package delivery

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

func Send(ctx *fasthttp.RequestCtx, status int, data interface{}) {
	ctx.SetContentType("application/json")

	ctx.SetStatusCode(status)
	body, err := json.Marshal(data)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
	ctx.SetBody(body)
}

func SendOK(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusOK)
}

func SendError(ctx *fasthttp.RequestCtx, status int, err string) {
	Send(ctx, status, Error{Message: err})
}
