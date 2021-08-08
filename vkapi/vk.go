package vkapi

import (
	"HwBot/shttp"
	"encoding/json"
	"fmt"
)

// Version is a vk api version.
const Version = "5.120"

var apibuilder = shttp.NewURIBuilder("https://api.vk.com", "method", "")

// Auth ...
func Auth(accessToken string) (*Client, *Error) {
	if accessToken == "" {
		return nil, nil
	}

	var c Client

	c = Client{
		token: accessToken,

		Groups:   ProvideGroups(&c),
		Messages: ProvideMessages(&c),
		Photos:   ProvidePhotos(&c),
		Users:    ProvideUsers(&c),
	}

	g, err := c.Groups.GetByID(nil)
	if err != nil {
		return nil, err
	}
	c.group = g[0].ID
	return &c, nil
}

// Client ...
type Client struct {
	client shttp.Client

	token string
	group int

	Groups   *GroupsProvider
	Messages *MessagesProvider
	Photos   *PhotosProvider
	Users    *UsersProvider
}

// HTTP returns http client.
func (c *Client) HTTP() *shttp.Client {
	return &c.client
}

// Group ...
func (c *Client) Group() int {
	return c.group
}

// Method ...
func (c *Client) Method(method string, args Args, dst interface{}) *Error {
	args.
		Set("access_token", c.token).
		Set("v", Version)
	query := args.q
	defer shttp.ReleaseQuery(query)

	uri := apibuilder.Build(query, method)
	req, resp := shttp.New(shttp.POSTStr, uri)
	defer shttp.Release(req, resp)

	err := c.client.Do(req, resp)
	if err != nil {
		e := Error{Message: err.Error()}
		return e.fill(method, query)
	}

	var res struct {
		Error    *Error          `json:"error"`
		Response json.RawMessage `json:"response"`
	}

	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		e := Error{Message: err.Error()}
		return e.fill(method, query)
	}

	if res.Error != nil {
		return res.Error.fill(method, query)
	}

	err = json.Unmarshal(res.Response, dst)
	if err != nil {
		e := Error{Message: err.Error()}
		return e.fill(method, query)
	}

	return nil
}

// APIProvider struct.
type APIProvider struct {
	client *Client
}

// Args ...
type Args struct {
	q *shttp.Query
}

// Set sets key=val arg.
func (a Args) Set(key string, val ...string) Args {
	a.q.Set(key, val...)
	return a
}

// Add adds key=value,val arg.
func (a Args) Add(key string, val string) Args {
	a.q.Add(key, val)
	return a
}

// NewArgs ...
func NewArgs() Args {
	return Args{shttp.AcquireQuery()}
}

// Error is vk api error.
type Error struct {
	Args    map[string]string
	Method  string
	Message string `json:"error_msg"`
	Code    int    `json:"error_code"`
}

func (e *Error) fill(m string, q *shttp.Query) *Error {
	e.Method = m
	e.Args = make(map[string]string)
	q.VisitAll(func(k, v []byte) {
		key := string(k)
		if !(key == "access_token") {
			e.Args[key] = string(v)
		}
	})
	return e
}

// IsVK ...
func (e *Error) IsVK() bool {
	return e.Message != "" && e.Code == 0
}

func (e *Error) Error() string {
	str := fmt.Sprintf("vk.%s: (%d) %s", e.Method, e.Code, e.Message)
	if e.Args != nil {
		str += "\n" + fmt.Sprint(e.Args)
	}
	return str
}
