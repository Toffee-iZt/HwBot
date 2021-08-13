package shttp

import (
	"context"
	"sync"

	"github.com/valyala/fasthttp"
)

// Client is fasthttp client.
type Client struct {
	fasthttp.Client
}

// Do performs the given http request and returns response body.
func (c *Client) Do(req *Request) ([]byte, error) {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := c.Client.Do(req, resp)
	if err != nil {
		return nil, err
	}
	b := resp.SwapBody(nil)
	return b, nil
}

// DoContext performs the given http request with context and returns http response.
func (c *Client) DoContext(ctx context.Context, req *Request) ([]byte, error) {
	reqCopy := fasthttp.AcquireRequest()
	req.Header.CopyTo(&reqCopy.Header)
	req.PostArgs().CopyTo(reqCopy.PostArgs())
	reqCopy.SwapBody(req.SwapBody(nil))

	var body []byte
	ch := make(chan error)
	mu := sync.Mutex{}

	go func() {
		var err error
		body, err = c.Do(reqCopy)
		mu.Lock()
		select {
		case <-ch: // closed when ctx is done
		default:
			req.SwapBody(reqCopy.SwapBody(nil))
			ch <- err
		}
		fasthttp.ReleaseRequest(reqCopy)
		mu.Unlock()
	}()

	select {
	case err := <-ch:
		return body, err
	case <-ctx.Done():
		mu.Lock()
		close(ch)
		mu.Unlock()
	}

	return nil, ctx.Err()
}
