package qbr

import (
	"reflect"
	"time"
)

// isZero check on zero value
func isZero(value any) bool {
	// check is nil
	if value == nil {
		return true
	}

	// get value by reflect
	v := reflect.ValueOf(value)

	// for pointer types
	if v.Kind() == reflect.Ptr ||
		v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Map ||
		v.Kind() == reflect.Chan ||
		v.Kind() == reflect.Interface {
		return v.IsNil()
	}

	// is string check on empty
	if v.Kind() == reflect.String {
		return v.String() == ""
	}

	// is number check on zero value
	if (v.Kind() >= reflect.Int && v.Kind() <= reflect.Int64) ||
		(v.Kind() >= reflect.Float32 && v.Kind() <= reflect.Float64) {
		return v.IsZero()
	}

	// for Time types
	if v.Kind() == reflect.Struct && v.Type() == reflect.TypeOf(time.Time{}) {
		return v.Interface().(time.Time).IsZero()
	}

	// for other types return false
	return false
}
