package vkhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/valyala/fasthttp"
)

// StatusOK 200(OK).
const StatusOK = fasthttp.StatusOK

// Client is fasthttp client.
type Client struct {
	fhttp fasthttp.Client
}

// Do performs the given http request and parses response body.
func (c *Client) Do(req *Request, obj interface{}) {
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
func (c *Client) DoContext(ctx context.Context, req *Request, obj interface{}) error {
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
