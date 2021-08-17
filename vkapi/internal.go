package vkapi

import (
	"encoding/json"
	"strconv"
)

type boolInt bool

func (b boolInt) MarshalJSON() ([]byte, error) {
	if b {
		return []byte{'1'}, nil
	}
	return []byte{'0'}, nil
}

func (b boolInt) String() string {
	if b {
		return "1"
	}
	return "0"
}

func (b *boolInt) UnmarshalJSON(data []byte) error {
	if string(data) != "0" {
		*b = true
	}
	return nil
}

func itoa(a int) string {
	return strconv.Itoa(a)
}

func ftoa(a float64) string {
	return strconv.FormatFloat(a, 'f', 7, 64)
}

func marshal(dst interface{}) []byte {
	b, _ := json.Marshal(dst)
	return b
}

func unmarshal(data []byte, dst interface{}) error {
	return json.Unmarshal(data, dst)
}
