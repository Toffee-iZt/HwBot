package vkapi

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// Version is a vk api version.
const Version = "5.120"

// Auth returns vk authorized client.
func Auth(accessToken string) (*Client, *Error) {
	if accessToken == "" {
		panic("vk: no access token provided")
	}

	c := Client{
		token: accessToken,
		rndID: -(1 << 31),
	}

	g, err := c.GroupsGetByID()
	if err != nil {
		return nil, err
	}
	c.self = g[0]

	return &c, nil
}

// Client struct.
type Client struct {
	HTTP       *fasthttp.Client
	bufferPool bytebufferpool.Pool

	token string
	self  *Group
	rndID int32
}

// Self returns self group info.
func (c *Client) Self() Group {
	return *c.self
}

func (c *Client) buildURI(q *query, endpoint ...string) *fasthttp.Request {
	uribuf := c.bufferPool.Get()
	for i := range endpoint {
		uribuf.WriteString(endpoint[i])
	}

	if q != nil {
		uribuf.WriteByte('?')
		uribuf.B = q.AppendBytes(uribuf.B)
	}
	releaseQuery(q)

	req := fasthttp.AcquireRequest()
	req.Header.SetRequestURIBytes(uribuf.Bytes())

	c.bufferPool.Put(uribuf)
	return req
}

// GET ...
func (c *Client) GET(ctx context.Context, endpoint string, args interface{}, res interface{}) error {
	req := c.buildURI(marshalArgs(args), endpoint)
	req.Header.SetMethod("GET")
	return c.DoContext(ctx, req, res)
}

// Do performs the given http request and parses response body.
func (c *Client) Do(req *fasthttp.Request, obj interface{}) {
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := c.HTTP.Do(req, resp)
	if err != nil {
		panic(fmt.Sprintf("vk: request %s\n%s\n%s", err.Error(), string(req.RequestURI()), string(req.Body())))
	}
	body := resp.SwapBody(nil)
	if status := resp.StatusCode(); status != fasthttp.StatusOK {
		panic(fmt.Sprintf("vk: http status %d\n%s", status, string(body)))
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		panic(fmt.Sprintf("vk: invalid body json format(%s)", err.Error()))
	}
}

// DoContext performs the given http request with context and parses response body.
func (c *Client) DoContext(ctx context.Context, req *fasthttp.Request, obj interface{}) error {
	reqCopy := fasthttp.AcquireRequest()
	req.Header.CopyTo(&reqCopy.Header)
	req.PostArgs().CopyTo(reqCopy.PostArgs())
	reqCopy.SwapBody(req.SwapBody(nil))

	ch := make(chan struct{})
	mu := sync.Mutex{}

	go func() {
		c.Do(reqCopy, obj)
		mu.Lock()
		select {
		case <-ch: // closed when ctx is done
		default:
			req.SwapBody(reqCopy.SwapBody(nil))
			ch <- struct{}{}
		}
		mu.Unlock()
	}()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		mu.Lock()
		close(ch)
		mu.Unlock()
	}

	return ctx.Err()
}

func (c *Client) method(dst interface{}, method string, args interface{}) *Error {
	const VkAPI string = "https://api.vk.com/method/"

	q := marshalArgs(args)
	q.Set("access_token", c.token)
	q.Set("v", Version)

	req := c.buildURI(q, VkAPI, method)
	req.Header.SetMethod("POST")

	var res struct {
		Error *struct {
			Message string `json:"error_msg"`
			Code    int    `json:"error_code"`
		} `json:"error"`
		Response json.RawMessage `json:"response"`
	}

	c.Do(req, &res)

	if res.Error != nil {
		return &Error{
			Method:  method,
			Args:    ArgsToMap(args),
			Code:    res.Error.Code,
			Message: res.Error.Message,
		}
	}

	if dst != nil {
		jerr := json.Unmarshal(res.Response, dst)
		if jerr != nil {
			panic("vk: method response invalid format")
		}
	}

	return nil
}

// Error struct.
type Error struct {
	Method  string
	Args    map[string]string
	Code    int
	Message string
}

func (e *Error) ErrorString() string {
	return fmt.Sprintf("vk.%s(%s) error %d %s", e.Method, e.Args, e.Code, e.Message)
}

// NewRequest returns new request instance.
func NewRequest(uri string, method string) *fasthttp.Request {
	r := fasthttp.AcquireRequest()
	r.SetRequestURI(uri)
	r.Header.SetMethod(method)
	return r
}

// ArgsMap is allias for string map.
type ArgsMap map[string]interface{}

// ArgsToMap converts args to map.
func ArgsToMap(args interface{}) map[string]string {
	q := marshalArgs(args)
	m := make(map[string]string, len(q.args))
	for _, v := range q.args {
		m[string(v.key)] = string(v.val)
	}
	releaseQuery(q)
	return m
}

func marshalArgs(args interface{}) *query {
	q := acquireQuery()
	if args == nil {
		return q
	}
	val := reflect.ValueOf(args)
	switch val.Kind() {
	case reflect.Struct:
		marshalArgsStruct(q, val)
		return q
	case reflect.Map:
		for iter := val.MapRange(); iter.Next(); {
			q.Set(iter.Key().String(), valof(iter.Value().Elem())...)
		}
		return q
	default:
		panic("vkhttp: invalid args interface kind " + val.Kind().String())
	}
}

func marshalArgsStruct(q *query, val reflect.Value) {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag, ok := field.Tag.Lookup("vkargs")
		if !ok {
			continue
		}

		tagVals := strings.Split(tag, ",")
		key := tagVals[0]
		if key == "" || key == "omitempty" {
			continue
		}

		vf := val.Field(i)
		if key == "embed" {
			marshalArgsStruct(q, vf)
			continue
		}

		if vf.IsZero() && len(tagVals) > 1 && tagVals[1] == "omitempty" {
			continue
		}
		q.Set(key, valof(vf)...)
	}
}

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

// Modified version of github.com/valyala/fasthttp.Args

func acquireQuery() *query {
	return queryPool.Get().(*query)
}

func releaseQuery(q *query) {
	q.args = q.args[:0]
	queryPool.Put(q)
}

var queryPool = sync.Pool{New: func() interface{} {
	return new(query)
}}

type query struct {
	args []kv
	buf  []byte
}

type kv struct {
	key []byte
	val []byte
}

// String returns string representation of query.
func (q *query) String() string {
	q.buf = q.AppendBytes(q.buf[:0])
	return *(*string)(unsafe.Pointer(&q.buf))
}

// AppendBytes appends query string to dst and returns the extended dst.
func (q *query) AppendBytes(dst []byte) []byte {
	const hex = "0123456789ABCDEF"
	const escapeTable = "\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x01\x01" +
		"\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x00" +
		"\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x00\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01"

	for i, n := 0, len(q.args); i < n; i++ {
		kv := &q.args[i]
		if i > 0 {
			dst = append(dst, '&')
		}
		dst = append(dst, kv.key...)
		dst = append(dst, '=')
		for _, c := range kv.val {
			switch {
			case c == ' ':
				dst = append(dst, '+')
			case escapeTable[int(c)] != 0:
				dst = append(dst, '%', hex[c>>4], hex[c&0xf])
			default:
				dst = append(dst, c)
			}
		}
	}
	return dst
}

// Set sets 'key=value1,value2...' argument.
func (q *query) Set(key string, value ...string) {
	q.buf = q.buf[:0]
	switch len(value) {
	case 0:
	case 1:
		q.buf = append(q.buf, value[0]...)
	default:
		for i, v := range value {
			if i > 0 {
				q.buf = append(q.buf, ',')
			}
			q.buf = append(q.buf, v...)
		}
	}

	n := len(q.args)
	for i := 0; i < n; i++ {
		kv := &q.args[i]
		if key == string(kv.key) {
			kv.val = append(kv.val[:0], q.buf...)
			return
		}
	}

	if cap(q.args) > n {
		q.args = q.args[:n+1]
	} else {
		q.args = append(q.args, kv{})
	}
	kv := &q.args[n]
	kv.key = append(kv.key[:0], key...)
	kv.val = append(kv.val[:0], q.buf...)
}
