package vkapi

import (
	"encoding/json"
	"fmt"

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
	Client vkhttp.Client

	token string
	self  *Group
	rndID int32
}

// Self returns self group info.
func (c *Client) Self() Group {
	return *c.self
}

func (c *Client) method(dst interface{}, method string, args interface{}) error {
	var res struct {
		Error *struct {
			Message string `json:"error_msg"`
			Code    int    `json:"error_code"`
		} `json:"error"`
		Response json.RawMessage `json:"response"`
	}

	c.Client.Method(method, args, c.token, Version, &res)

	if res.Error != nil {
		return &Error{
			Method:  method,
			Args:    vkhttp.ArgsToMap(args),
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

func (e *Error) Error() string {
	return fmt.Sprintf("vk.%s(%s) error %d %s", e.Method, e.Args, e.Code, e.Message)
}

type argmap map[string]interface{}
