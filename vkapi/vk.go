package vkapi

import (
	"encoding/json"
	"fmt"

	"github.com/Toffee-iZt/HwBot/shttp"
)

// Version is a vk api version.
const Version = "5.120"

var apibuilder = shttp.NewURIBuilder("https://api.vk.com", "method", "")

// Auth ...
func Auth(accessToken string) (*Client, error) {
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
func (c *Client) Method(method string, args Args, dst interface{}) error {
	args.
		Set("access_token", c.token).
		Set("v", Version)
	query := args.q
	defer shttp.ReleaseQuery(query)

	uri := apibuilder.Build(query, method)
	req := shttp.New(shttp.POSTStr, uri)

	body, err := c.client.Do(req)
	if err != nil {
		return &HTTPError{
			Args:    argMap(args),
			Method:  method,
			Message: err.Error(),
		}
	}

	var res struct {
		Error    *Error          `json:"error"`
		Response json.RawMessage `json:"response"`
	}

	

	err = json.Unmarshal(body, &res)
	if err != nil {
		return &JSONError{
			Args:    argMap(args),
			Method:  method,
			Message: err.Error(),
			Data:    body,
		}
	}

	if res.Error != nil {
		res.Error.Args = argMap(args)
		res.Error.Method = method
		return res.Error
	}

	err = json.Unmarshal(res.Response, dst)
	if err != nil {
		return &JSONError{
			Args:    argMap(args),
			Method:  method,
			Message: err.Error(),
			Data:    res.Response,
		}
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

func argMap(a Args) map[string]string {
	m := make(map[string]string)
	a.q.VisitAll(func(k, v []byte) {
		key := string(k)
		if key != "access_token" && key != "v" {
			m[key] = string(v)
		}
	})
	return m
}

// HTTPError represents vk method http error.
type HTTPError struct {
	Args    map[string]string
	Method  string
	Message string
}

func (e *HTTPError) Error() string {
	str := fmt.Sprintf("vk.%s: http %s", e.Method, e.Message)
	if e.Args != nil {
		str += "\n" + fmt.Sprint(e.Args)
	}
	return str
}

// JSONError represents vk method json error.
type JSONError struct {
	Args    map[string]string
	Method  string
	Message string
	Data    []byte
}

func (e *JSONError) Error() string {
	str := fmt.Sprintf("vk.%s: json %s", e.Method, e.Message)
	if e.Args != nil {
		str += "\n" + fmt.Sprint(e.Args)
	}
	str += "\n" + string(e.Data)
	return str
}

// Error is vk api error.
type Error struct {
	Args    map[string]string
	Method  string
	Message string `json:"error_msg"`
	Code    int    `json:"error_code"`
}

func (e *Error) Error() string {
	str := fmt.Sprintf("vk.%s: (%d) %s", e.Method, e.Code, e.Message)
	if e.Args != nil {
		str += "\n" + fmt.Sprint(e.Args)
	}
	return str
}
