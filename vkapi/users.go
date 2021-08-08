package vkapi

import "strconv"

// ProvideUsers makes users provider.
func ProvideUsers(c *Client) *UsersProvider {
	return &UsersProvider{
		APIProvider: APIProvider{c},
	}
}

// UsersProvider provides users api.
type UsersProvider struct {
	APIProvider
}

// User struct.
type User struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ID              int    `json:"id"`
	Deactivated     string `json:"deactivated"`
	IsClosed        bool   `json:"is_closed"`
	CanAccessClosed bool   `json:"can_access_closed"`
}

// Get returns detailed information on users.
// TODO: support fields and name_case
func (u *UsersProvider) Get(userIds ...int) ([]User, *Error) {
	if len(userIds) == 0 {
		return nil, nil
	}

	args := NewArgs()

	var n int
	for _, id := range userIds {
		if id != 0 {
			args.Add("user_ids", strconv.Itoa(id))
			n++
		}
	}
	if n == 0 {
		return nil, nil
	}

	var users []User
	err := u.client.Method("users.get", args, &users)
	return users, err
}
