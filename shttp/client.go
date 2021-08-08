package shttp

import (
	"context"
	"sync"

	"github.com/valyala/fasthttp"
)

// Client is fasthttp client.
type Client struct {
	client fasthttp.Client
}

// Do performs the given http request and fills the given http response.
func (c *Client) Do(req *Request, resp *Response) error {
	return c.client.Do(req, resp)
}

// DoContext performs the given http request with context and fills the given http response.
func (c *Client) DoContext(ctx context.Context, req *Request, resp *Response) error {
	reqCopy := fasthttp.AcquireRequest()
	respCopy := fasthttp.AcquireResponse()

	req.Header.CopyTo(&reqCopy.Header)
	req.URI().CopyTo(reqCopy.URI())
	req.PostArgs().CopyTo(reqCopy.PostArgs())
	reqCopy.SwapBody(req.SwapBody(nil))

	ch := make(chan error)
	mu := sync.Mutex{}

	go func() {
		err := c.client.Do(reqCopy, respCopy)
		mu.Lock()
		select {
		case <-ch: // closed when ctx is done
		default:
			req.SwapBody(reqCopy.SwapBody(nil))
			respCopy.CopyTo(resp)

			ch <- err
		}
		fasthttp.ReleaseRequest(reqCopy)
		fasthttp.ReleaseResponse(respCopy)
		mu.Unlock()
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		mu.Lock()
		close(ch)
		mu.Unlock()
	}

	return ctx.Err()
}
