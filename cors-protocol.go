package fastmiddle

import (
	"strings"

	"github.com/valyala/fasthttp"
)

//CorsProtocol defines CORS protocol
type CorsProtocol struct {
	isAllowCredentials string
	maxAge             string
	allowedOrigins     map[string]bool
	allowedHeaders     map[string]bool
	allowedMethods     map[string]bool
}

func (cors CorsProtocol) handle(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	//based on the OPTIONS method to divide the requests to 2 categories
	if cors.isPreflight(ctx) {
		//if the request is preflight then processing preflight step
		cors.handlePreflight(ctx)
		return ctx
	}
	//if not, handling the request as actual request
	cors.handleActualRequest(ctx)
	return ctx
}
func (cors CorsProtocol) isPreflight(ctx *fasthttp.RequestCtx) bool {
	return ctx.Request.Header.IsOptions()
}
func (cors CorsProtocol) isAllowedOrigins(origin string) bool {
	if _, isOK := cors.allowedOrigins["*"]; isOK {
		return true
	}
	if _, isOK := cors.allowedOrigins[origin]; isOK {
		return true
	}

	return false
}
func (cors CorsProtocol) isAllowedRequestsMethods(methods []string) bool {
	for i := 0; i < len(methods); i++ {
		if _, isOk := cors.allowedMethods[methods[i]]; !isOk {
			return false
		}
	}
	return true
}
func (cors CorsProtocol) isAllowedRequestsHeaders(headers []string) bool {
	for i := 0; i < len(headers); i++ {
		if _, isOk := cors.allowedHeaders[headers[i]]; !isOk {
			return false
		}
	}
	return true
}
func (cors CorsProtocol) handlePreflight(ctx *fasthttp.RequestCtx) *fasthttp.RequestCtx {
	//get request methods
	requestMethodStrs := string(ctx.Request.Header.Peek("Access-Control-Request-Method"))
	requestMethods := strings.Split(requestMethodStrs, ",")
	//check request methods are in the allowed method set
	if !cors.isAllowedRequestsMethods(requestMethods) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return ctx
	}
	//get request headers
	requestHeadersStrs := string(ctx.Request.Header.Peek("Access-Control-Request-Headers"))
	requestHeaders := strings.Split(requestHeadersStrs, ",")
	//check request headers are in the allowed headers set
	//if requests are not allowed by our server, sending a response with 403 status in HTTP response
	//and dont include any header in the below session
	if !cors.isAllowedRequestsHeaders(requestHeaders) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return ctx
	}

	//if request methods and request headers which will be sent to our server are allowed,
	//sending a response to requestor
	//the response has headers

	//1.Access-Control-Allow-Origin
	requestOrigin := string(ctx.Request.Header.Peek("Origin"))
	if !cors.isAllowedOrigins(requestOrigin) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		return ctx
	}

	ctx.Response.Header.Set("Access-Control-Allow-Origin", requestOrigin)
	//2.Access-Control-Allow-Credentials
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", cors.isAllowCredentials)
	//3.`Access-Control-Allow-Methods`
	ctx.Response.Header.Set("Access-Control-Allow-Methods", requestMethodStrs)
	//4.`Access-Control-Allow-Headers`
	ctx.Response.Header.Set("Access-Control-Allow-Headers", requestHeadersStrs)
	//5.`Access-Control-Max-Age`
	ctx.Response.Header.Set("Access-Control-Max-Age", cors.maxAge)
	// //6.`Access-Control-Expose-Headers`
	// ctx.Response.Header.Set("Access-Control-Expose-Headers", requestHeadersStrs)
	return ctx
}
func (cors CorsProtocol) handleActualRequest(ctx *fasthttp.RequestCtx) {
	requestOrigin := string(ctx.Request.Header.Peek("Origin"))
	ctx.Response.Header.Set("Access-Control-Allow-Origin", requestOrigin)
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", cors.isAllowCredentials)
}
func NewDefaultCorsProtocol() CorsProtocol {
	return CorsProtocol{
		isAllowCredentials: "true",
		maxAge:             "20",
		allowedHeaders:     map[string]bool{"accept": true, "content-type": true, "content-length": true, "accept-encoding": true, " X-CSRF-token": true, "authorization": true},
		allowedOrigins:     map[string]bool{"*": true},
		allowedMethods:     map[string]bool{"HEAD": true, "GET": true, "POST": true, "PUT": true, "DELETE": true, "OPTIONS": true},
	}
}
