package fastmiddle

import "github.com/valyala/fasthttp"

//CORSMiddleWare wrap CORS protocol inside
type CORSMiddleWare struct {
	cors CorsProtocol
	next MiddleWare
}

func (m CORSMiddleWare) Apply(next fasthttp.RequestHandler) fasthttp.RequestHandler {

	newHandler := func(ctx *fasthttp.RequestCtx) {
		ctx = m.cors.handle(ctx)
		if ctx.Response.StatusCode() == fasthttp.StatusForbidden {
			return
		}
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

func NewDefaultCORSMiddleWare() CORSMiddleWare {
	return CORSMiddleWare{
		cors: NewDefaultCorsProtocol(),
	}
}
