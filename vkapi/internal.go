package vkapi

import (
	"encoding/json"
	"strconv"

	"github.com/Toffee-iZt/HwBot/shttp"
)

type args struct {
	*shttp.Query
}

func newArgs() args {
	return args{shttp.AcquireQuery()}
}

func releaseArgs(a args) {
	shttp.ReleaseQuery(a.Query)
}

func itoa(a int) string {
	return strconv.Itoa(a)
}

func ftoa(a float64) string {
	return strconv.FormatFloat(a, 'f', 7, 64)
}

type jsonRaw json.RawMessage

func marshal(dst interface{}) []byte {
	b, _ := json.Marshal(dst)
	return b
}

func unmarshal(data []byte, dst interface{}) error {
	return json.Unmarshal(data, dst)
}
