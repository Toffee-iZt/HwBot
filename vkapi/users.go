package vkapi

import (
	"strconv"

	"github.com/Toffee-iZt/HwBot/vkapi/vktypes"
)

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

	UserOptFields
}

// Get returns detailed information on users.
// TODO: support fields and name_case
func (u *UsersProvider) Get(userIds []int, fields ...string) ([]User, error) {
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
	args.Set("fields", fields...)

	var users []User
	err := u.client.Method("users.get", args, &users)
	return users, err
}

// UserOptFields struct.
type UserOptFields struct {
	City *struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"city"`
	Country *struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"country"`
	Photo50      *string            `json:"photo_50"`
	Photo100     *string            `json:"photo_100"`
	Photo200     *string            `json:"photo_200"`
	Photo400     *string            `json:"photo_400"`
	PhotoMax     *string            `json:"photo_max"`
	PhotoMaxOrig *string            `json:"photo_max_orig"`
	Online       *baseBoolInt       `json:"online"`
	ScreenName   *string            `json:"screen_name"`
	CropPhoto    *vktypes.CropPhoto `json:"crop_photo"`
}
