package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen"
)

const (
	timeType            = "time.Time"
	mapStringToBoolType = "map[string]bool"
	stringType          = "string"
	boolType            = "bool"
)

var (
	// Times I've tried to optimize this regex before realizing it already accounts for
	// every edge case and there is no value (in either performance or readability terms)
	// in making it smaller: 1.
	numberMatcherRegex = regexp.MustCompile(`((u)?int(8|16|32|64)?|float(32|64))`)
)

func isCustomType(x string) bool {
	switch x {
	case "int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"float32",
		"float64",
		boolType,
		mapStringToBoolType,
		timeType,
		stringType:
		return false
	default:
		return true
	}
}

func generateProtoForType[T any](x T) (out string, err error) {
	typ := reflect.TypeOf(x)
	typeName := typ.Name()
	fieldsForType := reflect.VisibleFields(typ)

	output := fmt.Sprintf("message %s {\n", typeName)

	importedTypes := []string{}
	fieldIndex := 1
	for i := range fieldsForType {
		field := fieldsForType[i]
		if field.Name == "_" {
			continue
		}

		fieldName := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		if fieldName == "-" {
			continue
		}

		fieldType := strings.Replace(strings.TrimPrefix(strings.Replace(field.Type.String(), "[]", "", 1), "*"), "types.", "", 1)
		isPointer := field.Type.Kind() == reflect.Ptr
		isSlice := field.Type.Kind() == reflect.Slice

		if fieldType == "UserLoginInput" {
			continue
		}

		if isCustomType(fieldType) {
			importedTypes = append(importedTypes, fieldType)
		}

		switch fieldType {
		case timeType:
			fieldType = "google.protobuf.Timestamp"
		case mapStringToBoolType:
			fieldType = "map<string, bool>"
		case boolType:
			fieldType = "bool"
		case "float64":
			fieldType = "double"
		case "float32":
			fieldType = "double"
		case "int8":
			fieldType = "int32"
		case "int16":
			fieldType = "int32"
		case "uint8":
			fieldType = "uint32"
		case "uint16":
			fieldType = "uint32"
		case "ErrorCode":
			fieldType = "uint64"
		}

		var optional string
		if isPointer {
			optional = "optional "
		}

		var repeated string
		if isSlice {
			repeated = "repeated "
		}

		output += fmt.Sprintf("\t%s%s%s %s = %d;\n", optional, repeated, fieldType, field.Name, fieldIndex)
		fieldIndex++

		if t, ok := codegen.CustomTypeMap[fmt.Sprintf("%s.%s", typ.Name(), field.Name)]; ok {
			fieldType = t
			importedTypes = append(importedTypes, t)
		}
	}

	output += "}\n"

	return output, nil
}

func main() {
	if err := os.MkdirAll("artifacts", os.ModeDir); err != nil {
		log.Println("artifacts directory already exists")
	}

	finalOutput := `syntax = "proto3";
package dinnerdonebetter;

import "google/protobuf/timestamp.proto";

option go_package = "github.com/dinnerdonebetter/backend/internal/proto/types;types";

`
	for _, typesToGenerateFor := range codegen.TypeDefinitionFilesToGenerate {
		for _, typ := range typesToGenerateFor {
			output, err := generateProtoForType(typ)
			if err != nil {
				log.Fatal(err)
			}

			finalOutput += output + "\n"
		}
	}

	if err := os.WriteFile("../proto/types.proto", []byte(finalOutput), 0o644); err != nil {
		log.Fatal(err)
	}
}
