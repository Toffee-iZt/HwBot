package shttp

// Short version of github.com/valyala/fasthttp.Args

import (
	"sync"
	"unsafe"
)

// AcquireQuery returns an empty Query instance from query pool.
//
// The returned Query instance may be passed with Release when it is
// no longer needed. This allows Query recycling, reduces GC pressure
// and usually improves performance.
func AcquireQuery() *Query {
	return queryPool.Get().(*Query)
}

// ReleaseQuery returns Query acquired via AcquireQuery to query pool.
func ReleaseQuery(q *Query) {
	q.args = q.args[:0]
	queryPool.Put(q)
}

var queryPool = sync.Pool{New: func() interface{} {
	return new(Query)
}}

// Query represents HTTP query.
type Query struct {
	args []kv
	buf  []byte
}

type kv struct {
	key []byte
	val []byte
	has bool
}

// VisitAll calls f for each arg.
func (q *Query) VisitAll(f func([]byte, []byte)) {
	for _, kv := range q.args {
		f(kv.key, kv.val)
	}
}

// String returns string representation of query.
func (q *Query) String() string {
	q.buf = q.AppendBytes(q.buf[:0])
	return *(*string)(unsafe.Pointer(&q.buf))
}

// AppendBytes appends query string to dst and returns the extended dst.
func (q *Query) AppendBytes(dst []byte) []byte {
	for i, n := 0, len(q.args); i < n; i++ {
		kv := &q.args[i]
		if i > 0 {
			dst = append(dst, '&')
		}
		dst = append(dst, kv.key...)
		if kv.has {
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
	}
	return dst
}

func (q *Query) forceAppendArg(key, val string, has bool) {
	n := len(q.args)
	if cap(q.args) > n {
		q.args = q.args[:n+1]
	} else {
		q.args = append(q.args, kv{})
	}
	kv := &q.args[n]
	kv.key = append(kv.key[:0], key...)
	if has {
		kv.val = append(kv.val[:0], val...)
	}
	kv.has = has
}

func (q *Query) set(key, val string, has bool) {
	for i := 0; i < len(q.args); i++ {
		kv := &q.args[i]
		if key == string(kv.key) {
			if has {
				kv.val = append(kv.val[:0], val...)
			}
			kv.has = has
			return
		}
	}

	q.forceAppendArg(key, val, has)
}

// SetEmpty sets only 'key' as argument without the '='.
func (q *Query) SetEmpty(key string) *Query {
	q.set(key, "", false)
	return q
}

// Set sets 'key=value1,value2...' argument.
func (q *Query) Set(key string, value ...string) {
	switch len(value) {
	case 0:
		q.set(key, "", true)
	case 1:
		q.set(key, value[0], true)
	default:
		q.buf = q.buf[:0]
		for i, v := range value {
			if i > 0 {
				q.buf = append(q.buf, ',')
			}
			q.buf = append(q.buf, v...)
		}
		q.set(key, string(q.buf), true)
	}
}

// SetBytes sets 'key=value1,value2...' argument.
func (q *Query) SetBytes(key string, value ...[]byte) {
	switch len(value) {
	case 0:
		q.set(key, "", true)
	case 1:
		q.set(key, string(value[0]), true)
	default:
		q.buf = q.buf[:0]
		for i, v := range value {
			if i > 0 {
				q.buf = append(q.buf, ',')
			}
			q.buf = append(q.buf, v...)
		}
		q.set(key, string(q.buf), true)
	}
}

// Add adds 'key=value' argument.
//
// Multiple values for the same key may be added.
func (q *Query) Add(key string, value string) {
	for i := 0; i < len(q.args); i++ {
		kv := &q.args[i]
		if key == string(kv.key) {
			if kv.has {
				kv.val = append(kv.val, ',')
			} else {
				kv.has = true
			}
			kv.val = append(kv.val, value...)
			return
		}
	}

	q.forceAppendArg(key, value, true)
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
