package shttp

import (
	"github.com/valyala/bytebufferpool"
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
func New(method string, uri []byte) *Request {
	r := fasthttp.AcquireRequest()
	r.Header.SetMethod(method)
	r.Header.SetRequestURIBytes(uri)
	return r
}

var bpool bytebufferpool.Pool

// NewRequestsBuilder creates new requests builder.
func NewRequestsBuilder(dst string, path ...string) *RequestsBuilder {
	b := RequestsBuilder{
		dst:  dst,
		path: path,
	}
	for _, p := range path {
		if p == "" || p == "*" {
			b.epath++
		}
	}
	return &b
}

// RequestsBuilder struct.
type RequestsBuilder struct {
	dst   string
	path  []string
	epath int
}

// Build builds request with given params.
func (b *RequestsBuilder) Build(method string, q *Query, params ...string) *Request {
	if len(params) < b.epath {
		return nil
	}

	uribuf := bpool.Get()
	uribuf.WriteString(b.dst)

	for i, ip := 0, 0; i < len(b.path); i++ {
		uribuf.WriteByte('/')
		p := b.path[i]
		if p == "" || p == "*" {
			p = params[ip]
			ip++
		}
		uribuf.WriteString(p)
	}

	if q != nil {
		uribuf.WriteByte('?')
		uribuf.B = q.AppendBytes(uribuf.B)
	}

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.Header.SetRequestURIBytes(uribuf.Bytes())

	bpool.Put(uribuf)

	return req
}
