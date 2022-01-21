package vkhttp

// Short version of github.com/valyala/fasthttp.Args

import (
	"sync"
	"unsafe"
)

func acquireQuery() *query {
	return queryPool.Get().(*query)
}

func releaseQuery(q *query) {
	q.args = q.args[:0]
	queryPool.Put(q)
}

var queryPool = sync.Pool{New: func() interface{} {
	return new(query)
}}

type query struct {
	args []kv
	buf  []byte
}

type kv struct {
	key []byte
	val []byte
}

// String returns string representation of query.
func (q *query) String() string {
	q.buf = q.AppendBytes(q.buf[:0])
	return *(*string)(unsafe.Pointer(&q.buf))
}

// AppendBytes appends query string to dst and returns the extended dst.
func (q *query) AppendBytes(dst []byte) []byte {
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

func (q *query) set(key, val string) {
	n := len(q.args)
	for i := 0; i < n; i++ {
		kv := &q.args[i]
		if key == string(kv.key) {
			kv.val = append(kv.val[:0], val...)
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
	kv.val = append(kv.val[:0], val...)
}

// Set sets 'key=value1,value2...' argument.
func (q *query) Set(key string, value ...string) {
	switch len(value) {
	case 0:
		q.set(key, "")
	case 1:
		q.set(key, value[0])
	default:
		q.buf = q.buf[:0]
		for i, v := range value {
			if i > 0 {
				q.buf = append(q.buf, ',')
			}
			q.buf = append(q.buf, v...)
		}
		q.set(key, string(q.buf))
	}
}

const hex = "0123456789ABCDEF"
const escapeTable = "\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
	"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x01\x01" +
	"\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x01\x00" +
	"\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x00\x01" +
	"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
	"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
	"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01" +
	"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01\x01"
