package vkhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// StatusOK 200(OK).
const StatusOK = fasthttp.StatusOK

var bpool bytebufferpool.Pool

// VkAPI is default vk method uri.
const VkAPI string = "https://api.vk.com/method"

// NewRequest returns empty request.
func NewRequest() *fasthttp.Request {
	return fasthttp.AcquireRequest()
}

// Client is fasthttp client.
type Client struct {
	fhttp fasthttp.Client
}

// LongPoll ...
func (c *Client) LongPoll(ctx context.Context, server string, args interface{}, res interface{}) error {
	uribuf := bpool.Get()
	uribuf.WriteString(server)

	if args != nil {
		uribuf.WriteByte('?')
		uribuf.B = appendArgs(uribuf.B, args)
	}

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.SetRequestURIBytes(uribuf.Bytes())

	bpool.Put(uribuf)

	return c.DoContext(ctx, req, res)
}

// Method ...
func (c *Client) Method(method string, args interface{}, token, v string, res interface{}) {
	uribuf := bpool.Get()
	uribuf.WriteString(VkAPI)

	uribuf.WriteByte('/')
	uribuf.WriteString(method)

	uribuf.WriteByte('?')
	q := marshalArgs(args)
	q.Set("access_token", token)
	q.set("v", v)

	uribuf.B = q.AppendBytes(uribuf.B)
	releaseQuery(q)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.SetRequestURIBytes(uribuf.Bytes())

	bpool.Put(uribuf)

	c.Do(req, res)
}

// Do performs the given http request and parses response body.
func (c *Client) Do(req *fasthttp.Request, obj interface{}) {
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := c.fhttp.Do(req, resp)
	if err != nil {
		panic(fmt.Sprintf("vkhttp: bug %s\n%s\n%s", err.Error(), string(req.RequestURI()), string(req.Body())))
	}
	body := resp.SwapBody(nil)
	if status := resp.StatusCode(); status != StatusOK {
		panic(fmt.Sprintf("vkhttp: http status %d\n%s", status, string(body)))
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		panic(fmt.Sprintf("vkhttp: invalid body json format(%s)", err.Error()))
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
