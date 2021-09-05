package shttp

import (
	"context"
	"sync"

	"github.com/valyala/fasthttp"
)

// StatusOK 200(OK).
const StatusOK = fasthttp.StatusOK

// Client is fasthttp client.
type Client struct {
	fhttp fasthttp.Client
}

// Do performs the given http request and returns status and response body.
func (c *Client) Do(req *Request) (int, []byte, error) {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := c.fhttp.Do(req, resp)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode(), resp.SwapBody(nil), nil
}

// DoContext performs the given http request with context and returns status and response body.
func (c *Client) DoContext(ctx context.Context, req *Request) (int, []byte, error) {
	reqCopy := fasthttp.AcquireRequest()
	req.Header.CopyTo(&reqCopy.Header)
	req.PostArgs().CopyTo(reqCopy.PostArgs())
	reqCopy.SwapBody(req.SwapBody(nil))

	var status int
	var body []byte
	ch := make(chan error)
	mu := sync.Mutex{}

	go func() {
		var err error
		status, body, err = c.Do(reqCopy)
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
		return status, body, err
	case <-ctx.Done():
		mu.Lock()
		close(ch)
		mu.Unlock()
	}

	return 0, nil, ctx.Err()
}
