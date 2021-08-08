package shttp

import (
	"github.com/valyala/fasthttp"
)

//
const (
	GETStr  = "GET"
	POSTStr = "POST"
)

// Request as fasthttp.Request allias.
type Request = fasthttp.Request

// NewRequest returns an empty request instance from request pool.
func NewRequest(method string, b *URIBuffer) *Request {
	r := fasthttp.AcquireRequest()
	r.Header.SetMethod(method)
	r.Header.SetRequestURIBytes(b.B)
	uripool.Put(b)
	return r
}

// ReleaseRequest returns req to request pool.
func ReleaseRequest(req *Request) {
	fasthttp.ReleaseRequest(req)
}

// Response as fasthttp.Response allias.
type Response = fasthttp.Response

// NewResponse returns an empty Response instance from response pool.
func NewResponse() *Response {
	return fasthttp.AcquireResponse()
}

// ReleaseResponse return resp to response pool.
func ReleaseResponse(r *Response) {
	fasthttp.ReleaseResponse(r)
}

// New returns new (acquired) request and response instances.
func New(method string, b *URIBuffer) (*Request, *Response) {
	return NewRequest(method, b), NewResponse()
}

// Release releases request and response instances.
func Release(req *Request, resp *Response) {
	ReleaseRequest(req)
	ReleaseResponse(resp)
}
