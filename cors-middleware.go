package fastmiddle

import "github.com/valyala/fasthttp"

//CORSMiddleWare ...
type CORSMiddleWare struct {
	next MiddleWare
}

//Apply ...
func (m CORSMiddleWare) Apply(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	corsAllowOrigin := "*"
	corsAllowCredentials := "true"
	newHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", string(ctx.Request.Header.Peek("Origin")))
		ctx.Response.Header.Set("Access-Control-Allow-Methods", string(ctx.Request.Header.Peek("Access-Control-Request-Method")))
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		next(ctx)
	}
	if m.next != nil {
		return m.next.Apply(newHandler)
	}
	return newHandler
}

//SetNext ...
func (m CORSMiddleWare) SetNext(next MiddleWare) {
	m.next = next
}
