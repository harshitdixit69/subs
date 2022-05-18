package utils

import "reflect"

func IsSlicePointer(v interface{}) (result bool) {
	defer func() {
		err := recover()
		if err != nil {
			result = false
		}
	}()

	return reflect.TypeOf(v).Elem().Kind() == reflect.Slice
}
