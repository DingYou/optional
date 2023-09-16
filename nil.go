package optional

import (
	"errors"
	"reflect"
)

var ErrNil = errors.New("value required")

func IsNil[T any](val T) bool {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return v.IsNil()
	}
	return false
}

func RequireNonNull[T any](val T) T {
	if IsNil(val) {
		panic(ErrNil)
	}
	return val
}
