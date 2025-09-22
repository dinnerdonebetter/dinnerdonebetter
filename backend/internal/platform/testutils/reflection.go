package testutils

import (
	"reflect"
	"runtime"
	"strings"
)

// GetMethodName is meant to fetch the name of a given method passed in as an argument.
func GetMethodName(method any) string {
	v := reflect.ValueOf(method)
	if v.Kind() != reflect.Func {
		return ""
	}

	if pc := v.Pointer(); pc != 0 {
		if f := runtime.FuncForPC(pc); f != nil {
			fullName := f.Name()
			parts := strings.Split(fullName, ".")
			if len(parts) > 0 {
				return strings.TrimSuffix(parts[len(parts)-1], "-fm")
			}
		}
	}
	return ""
}
