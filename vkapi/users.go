package vkapi

// User struct.
type User struct {
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	ID              UserID  `json:"id"`
	Deactivated     string  `json:"deactivated"`
	IsClosed        boolean `json:"is_closed"`
	CanAccessClosed bool    `json:"can_access_closed"`

	// opt fields
	Photo50      *string  `json:"photo_50"`
	Photo100     *string  `json:"photo_100"`
	Photo200     *string  `json:"photo_200"`
	PhotoMax     *string  `json:"photo_max"`
	Photo200Orig *string  `json:"photo_200_orig"`
	Photo400Orig *string  `json:"photo_400_orig"`
	PhotoMaxOrig *string  `json:"photo_max_orig"`
	Online       *boolean `json:"online"`
	ScreenName   *string  `json:"screen_name"`
}

// name cases
const (
	UserNameCaseNom = "nom"
	UserNameCaseGen = "gen"
	UserNameCaseDat = "dat"
	UserNameCaseAcc = "acc"
	UserNameCaseIns = "ins"
	UserNameCaseAbl = "abl"
)

// users opt fields
const (
	UserOptFieldPhoto50      = "photo_50"
	UserOptFieldPhoto100     = "photo_100"
	UserOptFieldPhoto200     = "photo_200"
	UserOptFieldPhotoMax     = "photo_max"
	UserOptFieldPhoto200Orig = "photo_200_orig"
	UserOptFieldPhoto400Orig = "photo_400_orig"
	UserOptFieldPhotoMaxOrig = "photo_max_orig"
	UserOptFieldOnline       = "online"
	UserOptFieldScreenName   = "screen_name"
)

// UsersGet returns detailed information on users.
func (c *Client) UsersGet(userIds []UserID, nameCase string, fields ...string) ([]*User, *Error) {
	var users []*User
	return users, c.method(&users, "users.get", ArgsMap{
		"user_ids":  userIds,
		"name_case": nameCase,
		"fields":    fields,
	})
}
