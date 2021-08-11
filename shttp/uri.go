package shttp

import "github.com/valyala/bytebufferpool"

var uripool bytebufferpool.Pool

// NewURIBuilder creates new uri builder.
func NewURIBuilder(host string, path ...string) *URIBuilder {
	u := URIBuilder{
		host: host,
		path: path,
	}
	for _, p := range path {
		if p == "" || p == "*" {
			u.empty++
		}
	}
	return &u
}

// URIBuilder struct.
type URIBuilder struct {
	host  string
	path  []string
	empty int
}

// AppendTo builds uri and appends it to dst.
func (u *URIBuilder) AppendTo(dst []byte, query *Query, params ...string) []byte {
	if len(params) < u.empty {
		return nil
	}

	ip := 0
	dst = append(dst, u.host...)

	for _, p := range u.path {
		dst = append(dst, '/')
		if p == "" || p == "*" {
			dst = append(dst, params[ip]...)
			ip++
		} else {
			dst = append(dst, p...)
		}
	}

	if query != nil {
		dst = append(dst, '?')
		dst = query.AppendBytes(dst)
	}

	return dst
}

// Build builds uri and writes it to the request.
func (u *URIBuilder) Build(query *Query, params ...string) *URIBuffer {
	buf := uripool.Get()
	buf.B = u.AppendTo(buf.B, query, params...)
	return buf
}

// URIBuffer struct.
type URIBuffer = bytebufferpool.ByteBuffer

// URIFromString copies string to the uri buffer.
func URIFromString(s string) *URIBuffer {
	buf := uripool.Get()
	buf.SetString(s)
	return buf
}
