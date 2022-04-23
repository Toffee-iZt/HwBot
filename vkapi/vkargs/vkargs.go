package vkargs

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Marshal args.
func Marshal(args interface{}) url.Values {
	q := url.Values{}
	if args == nil {
		return q
	}
	val := reflect.ValueOf(args)
	switch val.Kind() {
	case reflect.Struct:
		marshalArgsStruct(q, val)
	case reflect.Map:
		for iter := val.MapRange(); iter.Next(); {
			q[iter.Key().String()] = valof(iter.Value().Elem())
		}
	default:
		panic("invalid vkhttp args kind " + val.Kind().String())
	}
	return q
}

func marshalArgsStruct(q url.Values, val reflect.Value) {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fValue := val.Field(i)
		if field.Anonymous {
			marshalArgsStruct(q, fValue)
			continue
		}

		tag, ok := field.Tag.Lookup("vkargs")
		if !ok || fValue.IsZero() {
			continue
		}

		tagVals := strings.Split(tag, ",")
		key := tagVals[0]
		if key == "" {
			continue
		}

		q[key] = valof(fValue)
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
	case reflect.Struct:
		s, ok := rval.Interface().(interface{ String() string })
		if ok {
			return s.String()
		}
	default:
	}
	panic("invalid kind of vk args field " + rval.Type().Name())
}
