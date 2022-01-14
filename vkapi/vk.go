package vkapi

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Toffee-iZt/HwBot/vkapi/vkhttp"
)

// Version is a vk api version.
const Version = "5.120"

// Auth returns vk authorized client.
func Auth(accessToken string) (*Client, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("no access token provided")
	}

	c := Client{
		api:   vkhttp.NewRequestsBuilder("https://api.vk.com/method"),
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
	client vkhttp.Client

	api vkhttp.RequestsBuilder

	token string
	self  *Group
	rndID int32
}

// Self returns self group info.
func (c *Client) Self() Group {
	return *c.self
}

// HTTP returns http client.
func (c *Client) HTTP() *vkhttp.Client {
	return &c.client
}

func (c *Client) method(dst interface{}, method string, args vkargs) error {
	call := MethodCall{
		Method: method,
		Args:   args,
	}

	req := c.api.BuildMethod(method, args, "access_token", c.token, "v", Version)

	var res struct {
		Error *struct {
			Message string `json:"error_msg"`
			Code    int    `json:"error_code"`
		} `json:"error"`
		Response json.RawMessage `json:"response"`
	}

	status, body, err := c.client.Do(req)
	if err != nil || status != vkhttp.StatusOK {
		goto ret
	}

	err = json.Unmarshal(body, &res)
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
		err = json.Unmarshal(res.Response, dst)
		if err == nil {
			return nil
		}
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
	Args   vkargs
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

type vkargs = vkhttp.Args

func itoa(a int) string {
	return strconv.Itoa(a)
}

func ftoa(a float64) string {
	return strconv.FormatFloat(a, 'f', 7, 64)
}
