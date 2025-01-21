package components

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// GetFieldNames returns the field names of a struct in alphabetical order.
func GetFieldNames[T any](v T) []string {
	val := reflect.TypeOf(v)
	if val.Kind() != reflect.Struct {
		panic("GetFieldNames expects a struct")
	}

	var fields []string
	for i := 0; i < val.NumField(); i++ {
		stringValue := val.Field(i).Name
		if stringValue == "_" {
			continue
		}

		if val.Field(i).Tag.Get("json") == "-" {
			continue
		}

		fields = append(fields, stringValue)
	}

	sort.Strings(fields)
	return fields
}

// GetFieldValues returns the stringified values of a struct's fields in alphabetical order of field names.
func GetFieldValues[T any](v T) []string {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Struct {
		panic("GetFieldValues expects a struct")
	}

	fieldNames := GetFieldNames(v)
	values := make([]string, len(fieldNames))

	// Map field names to their values
	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < typ.NumField(); i++ {
		fieldMap[typ.Field(i).Name] = val.Field(i)
	}

	// Retrieve and stringify values in sorted field name order
	for i, name := range fieldNames {
		fieldValue := fieldMap[name]
		stringValue := stringifyValue(fieldValue.Interface())
		values[i] = stringValue
	}

	return values
}

// stringifyValue converts a value to a string. If the value is a struct, map, or slice, it returns its JSON-encoded string.
func stringifyValue[T any](value T) string {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		return string(jsonValue)
	default:
		return fmt.Sprintf("%v", value)
	}
}
