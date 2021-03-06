package middlewares

import "github.com/valyala/fasthttp"

func JsonRequestHandler(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		handler(ctx)
	}
}