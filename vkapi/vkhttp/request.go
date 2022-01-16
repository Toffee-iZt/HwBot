package vkhttp

import (
	"strings"

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

var bpool bytebufferpool.Pool

// NewRequestsBuilder creates new requests builder.
func NewRequestsBuilder(dst string) RequestsBuilder {
	return RequestsBuilder(strings.TrimSuffix(dst, "/"))
}

// RequestsBuilder struct.
type RequestsBuilder string

// Build builds request with given args.
func (b *RequestsBuilder) Build(args Args, cred ...string) *Request {
	return b.BuildMethod("", args, cred...)
}

// BuildMethod builds request with vkmethod and given args.
func (b *RequestsBuilder) BuildMethod(vkmethod string, args interface{}, cred ...string) *Request {
	uribuf := bpool.Get()
	uribuf.WriteString(string(*b))

	if vkmethod != "" {
		uribuf.WriteByte('/')
		uribuf.WriteString(vkmethod)
	}

	if args != nil {
		uribuf.WriteByte('?')
		q := marshalArgs(args)
		for i := 0; len(cred)%2 == 0 && i < len(cred); i += 2 {
			q.Set(cred[i], cred[i+1])
		}
		uribuf.B = q.AppendBytes(uribuf.B)
		releaseQuery(q)
	}

	req := fasthttp.AcquireRequest()
	meth := GETStr
	if len(vkmethod) != 0 {
		meth = POSTStr
	}
	req.Header.SetMethod(meth)
	req.Header.SetRequestURIBytes(uribuf.Bytes())

	bpool.Put(uribuf)

	return req
}
