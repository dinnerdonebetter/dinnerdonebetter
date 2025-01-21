package components

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"maragu.dev/gomponents"
)

// GenerateInputs generates a list of Gomponents inputs based on struct fields.
func GenerateInputs[T any](ctx context.Context, submitProps SubmissionFormProps, fieldColumnNames map[string]string) ([]gomponents.Node, error) {
	var s T
	// Get the type of the struct
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", t.Kind())
	}

	var inputs []gomponents.Node

	// Iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get field name and type
		fieldName := field.Name
		displayName := fieldName
		if x, ok := fieldColumnNames[fieldName]; ok {
			displayName = x
		}

		fieldType := field.Type.Kind()

		// Determine input type based on field type
		var inputType string
		switch fieldType {
		case reflect.String:
			inputType = "text"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			inputType = "number"
		case reflect.Bool:
			inputType = "checkbox"
		default:
			log.Printf("unknown field type %s", fieldType)
			continue
		}

		// Add the input to the list
		inputs = append(inputs, FormTextInput(ctx, MustBuildValidatedTextInputPrompt(ctx, &TextInputsProps{
			ID:          fmt.Sprintf("input%s", fieldName),
			Name:        fieldName,
			LabelText:   displayName,
			Type:        inputType,
			Placeholder: " ",
		})))
	}

	return []gomponents.Node{BuildHTMXPoweredSubmissionForm(submitProps, inputs...)}, nil
}
