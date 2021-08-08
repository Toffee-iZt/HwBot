package vkapi

import (
	"fmt"
	"strconv"
)

type baseBoolInt bool

func (b baseBoolInt) MarshalJSON() ([]byte, error) {
	if b {
		return []byte{'1'}, nil
	}
	return []byte{'0'}, nil
}

func (b baseBoolInt) String() string {
	if b {
		return "1"
	}
	return "0"
}

func argsSetAny(a Args, k string, v interface{}) {
	var s string

	switch v.(type) {
	case string:
		s = v.(string)
	case int:
		i := v.(int)
		if i == 0 {
			return
		}
		s = strconv.Itoa(i)
	case []string:
		ar := v.([]string)
		if len(ar) == 0 {
			return
		}
		a.Set(k, ar...)
		return
	case []int:
		ar := v.([]int)
		if len(ar) == 0 {
			return
		}
		for i := range ar {
			a.Add(k, strconv.Itoa(ar[i]))
		}
		return
	default:
		s = fmt.Sprint(v)
	}

	if s == "" {
		return
	}

	a.Set(k, s)
}
