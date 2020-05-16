package fasthttp-middleware
//CORSMiddleWare ...
type CORSMiddleWare struct {
	next MiddleWare
}

//Apply ...
func (m CORSMiddleWare) Apply(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	corsAllowHeaders := "authorization"
	corsAllowMethods := "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin := "*"
	corsAllowCredentials := "true"
	newHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
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
