package optional

import (
	"errors"
	"reflect"
)

var ErrNil = errors.New("value required")

func IsNil[T any](val T) bool {
	v := reflect.ValueOf(val)
	return !v.IsValid() || v.IsNil()
}

func RequiredNonNull[T any](val T) T {
	if IsNil(val) {
		panic(ErrNil)
	}
	return val
}
