package vkapi

// ProvideGroups makes groups provider.
func ProvideGroups(c *Client) *GroupsProvider {
	return &GroupsProvider{
		APIProvider: APIProvider{c},
	}
}

// GroupsProvider provides groups api.
type GroupsProvider struct {
	APIProvider
}

// Group describes VK community.
type Group struct {
	Name        string  `json:"name"`
	ScreenName  string  `json:"screen_name"`
	Type        string  `json:"type"`
	ID          GroupID `json:"id"`
	IsClosed    boolInt `json:"is_closed"`
	Deactivated string  `json:"deactivated"`
	Photo50     string  `json:"photo_50"`
	Photo100    string  `json:"photo_100"`
	Photo200    string  `json:"photo_200"`
}

// GetByID returns information about communities by their IDs.
func (g *GroupsProvider) GetByID(groupIds []GroupID) ([]Group, error) {
	args := NewArgs()
	for _, id := range groupIds {
		args.Add("group_ids", itoa(int(id)))
	}

	var groups []Group
	err := g.client.Method("groups.getById", args, &groups)
	return groups, err
}

// LongPollServer struct.
type LongPollServer struct {
	Server string `json:"server"`
	Key    string `json:"key"`
	Ts     string `json:"ts"`
}

// GetLongPollServer returns the data needed to query a Long Poll server for events.
func (g *GroupsProvider) GetLongPollServer(groupID GroupID) (*LongPollServer, error) {
	args := NewArgs().Set("group_id", itoa(int(groupID)))
	var lps LongPollServer
	err := g.client.Method("groups.getLongPollServer", args, &lps)
	return &lps, err
}
