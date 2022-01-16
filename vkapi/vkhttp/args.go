package vkhttp

import (
	"reflect"
	"strconv"
	"strings"
)

// ArgsToMap converts args to map.
func ArgsToMap(args interface{}) map[string]string {
	q := marshalArgs(args)
	m := make(map[string]string, len(q.args))
	for _, v := range q.args {
		m[string(v.key)] = string(v.val)
	}
	releaseQuery(q)
	return m
}

func marshalArgs(args interface{}) *query {
	if args == nil {
		return nil
	}
	val := reflect.ValueOf(args)
	switch val.Kind() {
	case reflect.Struct:
		q := acquireQuery()
		marshalArgsStruct(q, val)
		return q
	case reflect.Map:
		q := acquireQuery()
		for iter := val.MapRange(); iter.Next(); {
			q.Set(iter.Key().String(), valof(iter.Value().Elem())...)
		}
		return q
	default:
		panic("vkhttp: invalid args interface kind " + val.Kind().String())
	}
}

func marshalArgsStruct(q *query, val reflect.Value) {
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

// Args is allias for string map.
type Args map[string]interface{}

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
		panic("BUG: query does not support type kind of field " + rval.Type().Name())
	}
}
