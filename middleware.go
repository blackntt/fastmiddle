package fashttp-middleware

import "github.com/valyala/fasthttp"

//MiddleWare ...
type MiddleWare interface {
	Apply(next fasthttp.RequestHandler) fasthttp.RequestHandler
	SetNext(MiddleWare)
}
