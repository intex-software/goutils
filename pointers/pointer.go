package pointers

import "reflect"

func Ptr[T any](t T) *T {
	return &t
}

func Deref[T any](t T) any {
	if val := reflect.ValueOf(t); val.Kind() == reflect.Ptr {
		return val.Elem()
	}

	return t
}
