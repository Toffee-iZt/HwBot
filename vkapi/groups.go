package vkapi

// Group describes VK community.
type Group struct {
	Name        string  `json:"name"`
	ScreenName  string  `json:"screen_name"`
	Type        string  `json:"type"`
	ID          GroupID `json:"id"`
	IsClosed    BoolInt `json:"is_closed"`
	Deactivated string  `json:"deactivated"`
	Photo50     string  `json:"photo_50"`
	Photo100    string  `json:"photo_100"`
	Photo200    string  `json:"photo_200"`
}

// GroupsGetByID returns information about communities by their IDs.
func (c *Client) GroupsGetByID(groupIds ...GroupID) ([]*Group, error) {
	var groups []*Group
	return groups, c.method(&groups, "groups.getById", vkargs{
		"group_ids": groupIds,
	})
}

// LongPollServer struct.
type LongPollServer struct {
	Server string `json:"server"`
	Key    string `json:"key"`
	Ts     string `json:"ts"`
}

// GetLongPollServer returns the data needed to query a Long Poll server for events.
func (c *Client) GetLongPollServer(groupID GroupID) (*LongPollServer, error) {
	var lps LongPollServer
	return &lps, c.method(&lps, "groups.getLongPollServer", vkargs{
		"group_id": groupID,
	})
}
