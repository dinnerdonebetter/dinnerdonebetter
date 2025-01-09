package pages

import (
	"context"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
)

func mustValidateTextProps(props components.TextInputsProps) components.ValidatedTextInput {
	validatedProps, err := components.BuildValidatedTextInputPrompt(context.Background(), &props)
	if err != nil {
		panic(err)
	}

	return validatedProps
}
