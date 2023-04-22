package converters

import (
	"fmt"
	"reflect"
)

func Convert[X, Y any](x X, y Y) {
	//nolint:gocritic // I don't control this function
	for _, field := range reflect.VisibleFields(reflect.TypeOf(x)) {
		if !field.IsExported() {
			continue
		}

		fmt.Println(field.Name)

		newValue := reflect.ValueOf(y).FieldByName(field.Name)
		f := reflect.ValueOf(x).FieldByName(field.Name)
		f.Set(newValue)
	}
}
