package vkapi

// User struct.
type User struct {
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	ID              UserID  `json:"id"`
	Deactivated     string  `json:"deactivated"`
	IsClosed        BoolInt `json:"is_closed"`
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

// UsersGet returns detailed information on users.
func (c *Client) UsersGet(userIds []UserID, nameCase string, fields ...string) ([]*User, error) {
	args := struct {
		UserIDs  []UserID `vkargs:"user_ids"`
		Fields   []string `vkargs:"fields,omitempty"`
		NameCase string   `vkargs:"name_case,omitempty"`
	}{
		userIds,
		fields,
		nameCase,
	}

	var users []*User
	return users, c.method(&users, "users.get", args)
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
	Online       *BoolInt `json:"online"`
	ScreenName   *string  `json:"screen_name"`
}
