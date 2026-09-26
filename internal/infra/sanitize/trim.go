package sanitize

import (
	"reflect"
	"strings"
)

// TrimStrings removes leading and trailing whitespace from every string
// reachable from v, which must be a pointer. It walks structs, pointers,
// slices, arrays, map values and interfaces. Map keys, unexported fields and
// fields tagged `trim:"-"` are left untouched.
func TrimStrings(v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return
	}

	trim(rv)
}

func trim(v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		if v.CanSet() {
			v.SetString(strings.TrimSpace(v.String()))
		}
	case reflect.Pointer:
		if !v.IsNil() {
			trim(v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() || !v.CanSet() {
			return
		}
		v.Set(trimmedCopy(v.Elem()))
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() || field.Tag.Get("trim") == "-" {
				continue
			}
			trim(v.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			trim(v.Index(i))
		}
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			v.SetMapIndex(iter.Key(), trimmedCopy(iter.Value()))
		}
	}
}

// trimmedCopy trims an addressable copy of v, for values that cannot be set in
// place (map values and the dynamic value of an interface).
func trimmedCopy(v reflect.Value) reflect.Value {
	cp := reflect.New(v.Type()).Elem()
	cp.Set(v)
	trim(cp)

	return cp
}
