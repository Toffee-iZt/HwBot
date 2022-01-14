package vkapi

import (
	"encoding/json"

	"github.com/Toffee-iZt/HwBot/common/strbytes"
)

// ID is a general id that can point to anything.
type ID int

// ToUser converts ID to UserID.
func (id ID) ToUser() UserID {
	if id > 0 && id < 2e9 {
		return UserID(id)
	}
	return 0
}

// ToGroup converts ID to GroupID.
func (id ID) ToGroup() GroupID {
	if id < 0 {
		return GroupID(-id)
	}
	return 0
}

// ToChat converts ID to ChatID.
func (id ID) ToChat() ChatID {
	if id > 2e9 {
		return ChatID(id - 2e9)
	}
	return 0
}

// UserID is equal to id but points only to users.
type UserID uint

// ToID ...
func (u UserID) ToID() ID {
	return ID(u)
}

// GroupID points to group.
type GroupID uint

// ToID ...
func (g GroupID) ToID() ID {
	return -ID(g)
}

// ChatID points to chat.
type ChatID uint

// ToID ...
func (c ChatID) ToID() ID {
	return ID(c + 2e9)
}

// JSONData represents json as string.
type JSONData string

// NewJSONData creates new JSONData from string.
func NewJSONData(s string) (JSONData, bool) {
	ok := json.Valid(strbytes.S2b(s))
	if !ok {
		return "{}", false
	}
	return JSONData(s), true
}

// NewJSONDataBytes creates new JSONData from bytes slice.
func NewJSONDataBytes(b []byte) (JSONData, bool) {
	ok := json.Valid(b)
	if !ok {
		return "{}", false
	}
	return JSONData(strbytes.B2s(b)), true
}

// NewJSONDataMarshal creates new JSONData from object.
func NewJSONDataMarshal(v interface{}) (JSONData, bool) {
	d, err := json.Marshal(v)
	if err != nil {
		return "{}", false
	}
	return JSONData(strbytes.B2s(d)), true
}

// BoolInt is bool as 0 or 1.
type BoolInt bool

// MarshalJSON implements json.Marshaler.
func (b BoolInt) MarshalJSON() ([]byte, error) {
	if b {
		return []byte{'1'}, nil
	}
	return []byte{'0'}, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BoolInt) UnmarshalJSON(data []byte) error {
	*b = string(data) != "0"
	return nil
}
