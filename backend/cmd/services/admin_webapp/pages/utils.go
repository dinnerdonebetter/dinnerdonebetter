package pages

import (
	"context"

	"github.com/dinnerdonebetter/backend/cmd/services/admin_webapp/components"
)

func mustValidateTextProps(props *components.TextInputsProps) components.ValidatedTextInput {
	if props == nil {
		panic("props cannot be nil")
	}

	validatedProps, err := components.BuildValidatedTextInputPrompt(context.Background(), props)
	if err != nil {
		panic(err)
	}

	return validatedProps
}
