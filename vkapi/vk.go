package vkapi

import (
	"encoding/json"
	"fmt"

	"github.com/Toffee-iZt/HwBot/shttp"
)

// Version is a vk api version.
const Version = "5.120"

// Auth returns vk authorized client.
func Auth(accessToken string) (*Client, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	c := Client{
		api:   shttp.NewRequestsBuilder("https://api.vk.com", "method", ""),
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
	client shttp.Client

	api *shttp.RequestsBuilder

	token string
	self  *Group
	rndID int32
}

// Self returns self group info.
func (c *Client) Self() Group {
	return *c.self
}

// HTTP returns http client.
func (c *Client) HTTP() *shttp.Client {
	return &c.client
}

func (c *Client) method(method string, args args, dst interface{}) error {
	call := MethodCall{
		Method: method,
		Args:   make(map[string]string),
	}
	args.VisitAll(func(k, v []byte) {
		call.Args[string(k)] = string(v)
	})

	args.Set("access_token", c.token)
	args.Set("v", Version)

	req := c.api.Build(shttp.POSTStr, args.Query, method)
	releaseArgs(args)

	var res struct {
		Error *struct {
			Message string `json:"error_msg"`
			Code    int    `json:"error_code"`
		} `json:"error"`
		Response json.RawMessage `json:"response"`
	}

	status, body, err := c.client.Do(req)
	if err != nil || status != shttp.StatusOK {
		goto ret
	}

	err = unmarshal(body, &res)
	if err != nil {
		goto ret
	}

	if res.Error != nil {
		return &Error{
			Call:    &call,
			Code:    res.Error.Code,
			Message: res.Error.Message,
		}
	}

	if dst != nil {
		err = unmarshal(res.Response, dst)
	}

	if err == nil {
		return nil
	}

ret:
	return &Error{
		Call:       &call,
		Message:    err.Error(),
		HTTPStatus: status,
		Body:       body,
	}
}

// MethodCall struct.
type MethodCall struct {
	Method string
	Args   map[string]string
}

// Error struct.
type Error struct {
	Call       *MethodCall
	Code       int
	Message    string
	HTTPStatus int
	Body       []byte
}

func (e *Error) Error() string {
	var errbody string
	if e.Code != 0 {
		errbody = fmt.Sprintf("vk: (%d) %s", e.Code, e.Message)
	} else {
		errbody = fmt.Sprintf("Status: %d\n%s\n\n%s", e.HTTPStatus, e.Message, string(e.Body))
	}
	return fmt.Sprintf("vk.%s(%s)\n%s", e.Call.Method, e.Call.Args, errbody)
}
