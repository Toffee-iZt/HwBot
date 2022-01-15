package vkhttp

import (
	"reflect"
	"strconv"
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
func (b *RequestsBuilder) BuildMethod(vkmethod string, args Args, cred ...string) *Request {
	uribuf := bpool.Get()
	uribuf.WriteString(string(*b))

	if vkmethod != "" {
		uribuf.WriteByte('/')
		uribuf.WriteString(vkmethod)
	}

	if args != nil {
		uribuf.WriteByte('?')
		q := acquireQuery()
		for k, v := range args {
			q.Set(k, valof(reflect.ValueOf(v))...)
		}
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

// Args is allias for string map.
type Args map[string]interface{}

func valof(rval reflect.Value) []string {
	if k := rval.Kind(); k == reflect.Array || k == reflect.Slice {
		val := make([]string, rval.Len())
		for i := 0; i < len(val); i++ {
			val[i] = valof1(rval.Index(i))
		}
		return val
	}
	return []string{valof1(rval)}
}

func valof1(rval reflect.Value) string {
	switch rval.Kind() {
	case reflect.Bool:
		if rval.Bool() {
			return "1"
		}
		return "0"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rval.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rval.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rval.Float(), 'f', 10, 64)
	case reflect.String:
		return rval.String()
	default:
		panic("BUG: query does not support type kind of field " + rval.Type().Name())
	}
}
