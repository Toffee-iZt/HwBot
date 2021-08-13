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

// New returns an empty request instance from request pool.
func New(method string, b *URIBuffer) *Request {
	r := fasthttp.AcquireRequest()
	r.Header.SetMethod(method)
	r.Header.SetRequestURIBytes(b.B)
	uripool.Put(b)
	return r
}
