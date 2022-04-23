package vkapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Toffee-iZt/HwBot/vkapi/vkargs"
)

// Version is a vk api version.
const Version = "5.120"

// Auth returns vk authorized client.
func Auth(accessToken string) (*Client, *Error) {
	if accessToken == "" {
		panic("no vk access token provided")
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
	http http.Client

	token string
	self  *Group
	rndID int32
}

// Self returns self group info.
func (c *Client) Self() Group {
	return *c.self
}

// HTTP sends http request.
func (c *Client) HTTP(req *http.Request) (*http.Response, error) {
	return c.http.Do(req)
}

// Do sends http request and parses body as JSON.
func (c *Client) Do(req *http.Request, obj interface{}) (status int, canceled bool) {
	resp, err := c.http.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return 0, true
		}
		panic(fmt.Sprintf("vk request %s\n%s", err.Error(), string(req.URL.String())))
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(obj)
	if err != nil {
		panic(fmt.Sprintf("vk invalid body format(%s)", err.Error()))
	}

	return resp.StatusCode, false
}

func (c *Client) method(obj interface{}, method string, args Args) *Error {
	const VkAPI string = "https://api.vk.com/method/"

	q := vkargs.Marshal(args)
	q.Set("access_token", c.token)
	q.Set("v", Version)

	uri := VkAPI + method + "?" + q.Encode()
	req, _ := http.NewRequest(http.MethodPost, uri, nil)

	var res struct {
		Error *struct {
			Message string `json:"error_msg"`
			Code    int    `json:"error_code"`
		} `json:"error"`
		Response json.RawMessage `json:"response"`
	}

	c.Do(req, &res)

	if res.Error != nil {
		q.Del("access_token")
		return &Error{
			Method:  method,
			Args:    q,
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
	Args    url.Values
	Code    int
	Message string
}

func (e *Error) String() string {
	return fmt.Sprintf("vk.%s(%s)\nerror %d %s", e.Method, e.Args, e.Code, e.Message)
}

// Args for vk.
type Args interface{}

// ArgsMap is allias for string map.
type ArgsMap map[string]interface{}
