package vkapi

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
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	ID              UserID  `json:"id"`
	Deactivated     string  `json:"deactivated"`
	IsClosed        boolInt `json:"is_closed"`
	CanAccessClosed bool    `json:"can_access_closed"`

	UserOptFields
}

// name cases
const (
	NameCaseNom = "nom"
	NameCaseGen = "gen"
	NameCaseDat = "dat"
	NameCaseAcc = "acc"
	NameCaseIns = "ins"
	NameCaseAbl = "abl"
)

// Get returns detailed information on users.
func (u *UsersProvider) Get(userIds []UserID, nameCase string, fields ...string) ([]User, error) {
	if len(userIds) == 0 {
		return nil, nil
	}

	args := NewArgs()

	var n int
	for _, id := range userIds {
		if id != 0 {
			args.Add("user_ids", itoa(int(id)))
			n++
		}
	}
	if n == 0 {
		return nil, nil
	}
	args.Set("fields", fields...)
	if nameCase != "" {
		args.Set("name_case")
	}

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
	Photo50      *string  `json:"photo_50"`
	Photo100     *string  `json:"photo_100"`
	Photo200     *string  `json:"photo_200"`
	PhotoMax     *string  `json:"photo_max"`
	Photo200Orig *string  `json:"photo_200_orig"`
	Photo400Orig *string  `json:"photo_400_orig"`
	PhotoMaxOrig *string  `json:"photo_max_orig"`
	Online       *boolInt `json:"online"`
	ScreenName   *string  `json:"screen_name"`
}
