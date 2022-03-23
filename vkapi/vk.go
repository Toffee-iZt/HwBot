package vkapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Toffee-iZt/HwBot/vkapi/vkargs"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// Version is a vk api version.
const Version = "5.120"

// Auth returns vk authorized client.
func Auth(accessToken string) (*Client, *Error) {
	if accessToken == "" {
		panic("no vk access token provided")
	}

	c := Client{
		http:  &fasthttp.Client{},
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
	http       *fasthttp.Client
	bufferPool bytebufferpool.Pool

	token string
	self  *Group
	rndID int32
}

// Self returns self group info.
func (c *Client) Self() Group {
	return *c.self
}

func (c *Client) do(req *fasthttp.Request, obj interface{}) {
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := c.http.Do(req, resp)
	if err != nil {
		panic(fmt.Sprintf("vk request %s\n%s\n%s", err.Error(), string(req.RequestURI()), string(req.Body())))
	}
	body := resp.SwapBody(nil)
	if status := resp.StatusCode(); status != fasthttp.StatusOK {
		panic(fmt.Sprintf("vk http status %d\n%s", status, string(body)))
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		panic(fmt.Sprintf("vk invalid body json format(%s)", err.Error()))
	}
}

func (c *Client) doContext(ctx context.Context, req *fasthttp.Request, obj interface{}) error {
	reqCopy := fasthttp.AcquireRequest()
	req.Header.CopyTo(&reqCopy.Header)
	req.PostArgs().CopyTo(reqCopy.PostArgs())
	reqCopy.SwapBody(req.SwapBody(nil))

	ch := make(chan struct{})
	mu := sync.Mutex{}

	go func() {
		c.do(reqCopy, obj)
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

func (c *Client) buildRequest(url string, q *vkargs.Query) *fasthttp.Request {
	uribuf := c.bufferPool.Get()
	uribuf.WriteString(url)

	if q != nil {
		uribuf.WriteByte('?')
		uribuf.B = q.AppendBytes(uribuf.B)
		vkargs.ReleaseQuery(q)
	}

	req := fasthttp.AcquireRequest()
	req.Header.SetRequestURIBytes(uribuf.Bytes())

	c.bufferPool.Put(uribuf)
	return req
}

// GET ...
func (c *Client) GET(ctx context.Context, endpoint string, args Args, res interface{}) error {
	req := c.buildRequest(endpoint, vkargs.Marshal(args))
	req.Header.SetMethod(fasthttp.MethodGet)
	return c.doContext(ctx, req, res)
}

func (c *Client) method(obj interface{}, method string, args Args) *Error {
	const VkAPI string = "https://api.vk.com/method/"

	q := vkargs.Marshal(args)
	q.Set("access_token", c.token)
	q.Set("v", Version)

	req := c.buildRequest(VkAPI+method, q)
	req.Header.SetMethod(fasthttp.MethodPost)

	var res struct {
		Error *struct {
			Message string `json:"error_msg"`
			Code    int    `json:"error_code"`
		} `json:"error"`
		Response json.RawMessage `json:"response"`
	}

	c.do(req, &res)

	if res.Error != nil {
		return &Error{
			Method:  method,
			Args:    vkargs.ToMap(args),
			Code:    res.Error.Code,
			Message: res.Error.Message,
		}
	}

	if obj != nil {
		jerr := json.Unmarshal(res.Response, obj)
		if jerr != nil {
			panic("vk method response invalid format")
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

func (e *Error) Error() string {
	return fmt.Sprintf("vk.%s(%s) error %d %s", e.Method, e.Args, e.Code, e.Message)
}

// Args for vk.
type Args interface{}

// ArgsMap is allias for string map.
type ArgsMap map[string]interface{}
