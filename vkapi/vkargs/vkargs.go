package vkargs

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

func ToMap(args interface{}) map[string]string {
	q := Marshal(args)
	m := make(map[string]string, len(q.args))
	for _, v := range q.args {
		m[string(v.key)] = string(v.val)
	}
	ReleaseQuery(q)
	return m
}

func Marshal(args interface{}) *Query {
	q := acquireQuery()
	if args == nil {
		return q
	}
	val := reflect.ValueOf(args)
	switch val.Kind() {
	case reflect.Struct:
		marshalArgsStruct(q, val)
		return q
	case reflect.Map:
		for iter := val.MapRange(); iter.Next(); {
			q.Set(iter.Key().String(), valof(iter.Value().Elem())...)
		}
		return q
	default:
		panic("invalid vkhttp args kind " + val.Kind().String())
	}
}

func marshalArgsStruct(q *Query, val reflect.Value) {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag, ok := field.Tag.Lookup("vkargs")
		if !ok {
			continue
		}

		tagVals := strings.Split(tag, ",")
		key := tagVals[0]
		if key == "" || key == "omitempty" {
			continue
		}

		vf := val.Field(i)
		if key == "embed" {
			marshalArgsStruct(q, vf)
			continue
		}

		if vf.IsZero() && len(tagVals) > 1 && tagVals[1] == "omitempty" {
			continue
		}
		q.Set(key, valof(vf)...)
	}
}

func valof(rval reflect.Value) []string {
	if k := rval.Kind(); k == reflect.Array || k == reflect.Slice {
		val := make([]string, rval.Len())
		for i := 0; i < len(val); i++ {
			val[i] = valof1(rval.Index(i))
		}
		return val
	}
	return []string{valof1(rval)}
}

func valof1(rval reflect.Value) string {
	switch rval.Kind() {
	case reflect.Bool:
		if rval.Bool() {
			return "1"
		}
		return "0"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rval.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rval.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rval.Float(), 'f', 10, 64)
	case reflect.String:
		return rval.String()
	default:
		panic("vkargs does not support type kind of field " + rval.Type().Name())
	}
}

// Modified version of github.com/valyala/fasthttp.Args

func acquireQuery() *Query {
	return queryPool.Get().(*Query)
}

func ReleaseQuery(q *Query) {
	q.args = q.args[:0]
	queryPool.Put(q)
}

var queryPool = sync.Pool{New: func() interface{} {
	return new(Query)
}}

type Query struct {
	args []kv
	buf  []byte
}

type kv struct {
	key []byte
	val []byte
}

// String returns string representation of query.
func (q *Query) String() string {
	q.buf = q.AppendBytes(q.buf[:0])
	return *(*string)(unsafe.Pointer(&q.buf))
}

// AppendBytes appends query string to dst and returns the extended dst.
func (q *Query) AppendBytes(dst []byte) []byte {
	const hex = "0123456789ABCDEF"
	const escapeTable = "\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x01\x01" +
		"\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x00" +
		"\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x00\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
		"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01"

	for i, n := 0, len(q.args); i < n; i++ {
		kv := &q.args[i]
		if i > 0 {
			dst = append(dst, '&')
		}
		dst = append(dst, kv.key...)
		dst = append(dst, '=')
		for _, c := range kv.val {
			switch {
			case c == ' ':
				dst = append(dst, '+')
			case escapeTable[int(c)] != 0:
				dst = append(dst, '%', hex[c>>4], hex[c&0xf])
			default:
				dst = append(dst, c)
			}
		}
	}
	return dst
}

// Set sets 'key=value1,value2...' argument.
func (q *Query) Set(key string, value ...string) {
	q.buf = q.buf[:0]
	if len(value) == 1 {
		q.buf = append(q.buf, value[0]...)
	} else {
		for i, v := range value {
			if i > 0 {
				q.buf = append(q.buf, ',')
			}
			q.buf = append(q.buf, v...)
		}
	}

	n := len(q.args)
	for i := 0; i < n; i++ {
		kv := &q.args[i]
		if key == string(kv.key) {
			kv.val = append(kv.val[:0], q.buf...)
			return
		}
	}

	if cap(q.args) > n {
		q.args = q.args[:n+1]
	} else {
		q.args = append(q.args, kv{})
	}
	kv := &q.args[n]
	kv.key = append(kv.key[:0], key...)
	kv.val = append(kv.val[:0], q.buf...)
}
