package vkapi

import (
	"encoding/json"

	"github.com/Toffee-iZt/HwBot/common/strbytes"
)

// ID is a general id that can point to anything.
type ID int

// IsUser returns true if it is UserID.
func (id ID) IsUser() bool {
	return id > 0 && id < 2e9
}

// ToUser converts ID to UserID.
func (id ID) ToUser() UserID {
	if !id.IsUser() {
		return 0
	}
	return UserID(id)
}

// IsGroup returns true if it is GroupID.
func (id ID) IsGroup() bool {
	return id < 0
}

// ToGroup converts ID to GroupID.
func (id ID) ToGroup() GroupID {
	if !id.IsGroup() {
		return 0
	}
	return GroupID(-id)
}

// IsChat returns true if it is ChatID.
func (id ID) IsChat() bool {
	return id > 2e9
}

// ToChat converts ID to ChatID.
func (id ID) ToChat() ChatID {
	if !id.IsChat() {
		return 0
	}
	return ChatID(id - 2e9)
}

// UserID is equal to id but points only to users.
type UserID uint

// ToID converts UserID to ID.
func (id UserID) ToID() ID {
	return ID(id)
}

// GroupID points to group.
type GroupID uint

// ToID converts GroupID to ID.
func (id GroupID) ToID() ID {
	return -ID(id)
}

// ChatID points to chat.
type ChatID uint

// ToID converts ChatID to ID.
func (id ChatID) ToID() ID {
	return ID(id + 2e9)
}

// NewJSONData creates new json_data from string.
func NewJSONData(s string) (JSONData, bool) {
	ok := json.Valid(strbytes.S2b(s))
	if !ok {
		return "{}", false
	}
	return JSONData(s), true
}

// NewJSONDataBytes creates new json_data from bytes slice.
func NewJSONDataBytes(b []byte) (JSONData, bool) {
	ok := json.Valid(b)
	if !ok {
		return "{}", false
	}
	return JSONData(strbytes.B2s(b)), true
}

// JSONData represents json as string.
type JSONData string
