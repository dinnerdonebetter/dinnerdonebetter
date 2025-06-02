package components

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"maragu.dev/gomponents"
	"maragu.dev/gomponents/components"
	ghtml "maragu.dev/gomponents/html"
)

func buildInputAndLabel(input ValidatedTextInput) []gomponents.Node {
	props := input.data()

	return []gomponents.Node{
		ghtml.Label(
			ghtml.For(fmt.Sprintf("%s_input", props.ID)),
			components.Classes{
				"block":         true,
				"mb-2":          true,
				"text-sm":       true,
				"font-medium":   true,
				"text-gray-700": true,
			},
			gomponents.Text(props.LabelText),
		),
		ghtml.Input(
			ghtml.Type(props.Type),
			ghtml.Name(props.Name),
			ghtml.ID(fmt.Sprintf("%s_input", props.ID)),
			components.Classes{
				"w-full":                true,
				"px-2":                  true,
				"py-1":                  true,
				"text-gray-700":         true,
				"transition-all":        true,
				"duration-300":          true,
				"bg-white":              true,
				"border":                true,
				"border-gray-300":       true,
				"rounded-lg":            true,
				"shadow-sm":             true,
				"focus:ring-2":          true,
				"focus:ring-blue-400":   true,
				"focus:border-blue-400": true,
				"focus:outline-none":    true,
			},
			ghtml.Placeholder(props.Placeholder),
		),
	}
}

type TextInputsProps struct {
	ID          string
	Name        string
	LabelText   string
	Type        string
	Placeholder string
}

func (t *TextInputsProps) Validate(ctx context.Context) (ValidatedTextInput, error) {
	if err := t.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return &validatedInputProps{
		ID:          t.ID,
		Name:        t.Name,
		LabelText:   t.LabelText,
		Type:        t.Type,
		Placeholder: t.Placeholder,
	}, nil
}

func BuildValidatedTextInputPrompt(ctx context.Context, props *TextInputsProps) (ValidatedTextInput, error) {
	return props.Validate(ctx)
}

func MustBuildValidatedTextInputPrompt(ctx context.Context, props *TextInputsProps) ValidatedTextInput {
	x, err := BuildValidatedTextInputPrompt(ctx, props)
	if err != nil {
		panic(err)
	}

	return x
}

type validatedInputProps struct {
	ID          string
	Name        string
	LabelText   string
	Type        string
	Placeholder string
}

func (v *validatedInputProps) validated() {}
func (v *validatedInputProps) textInput() {}
func (v *validatedInputProps) data() *validatedInputProps {
	return v
}

func (t *TextInputsProps) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		t,
		validation.Field(&t.ID, validation.Required),
		validation.Field(&t.Name, validation.Required),
		validation.Field(&t.LabelText, validation.Required),
		validation.Field(&t.Type, validation.Required),
		validation.Field(&t.Placeholder, validation.Required),
	)
}

type Validated interface {
	validated()
}

type TextInput interface {
	textInput()
}

type ValidatedTextInput interface {
	Validated
	TextInput
	data() *validatedInputProps
}

func FormTextInput(ctx context.Context, props ValidatedTextInput) gomponents.Node {
	return ghtml.Div(
		buildInputAndLabel(props)...,
	)
}

func TextInputs(ctx context.Context, props ...ValidatedTextInput) gomponents.Node {
	inputs := []gomponents.Node{}
	for _, prop := range props {
		inputs = append(inputs,
			ghtml.Div(
				append(
					[]gomponents.Node{ghtml.Class("w-full")},
					buildInputAndLabel(prop)...,
				)...,
			),
		)
	}

	return ghtml.Div(
		ghtml.Div(
			ghtml.Class("flex"),
			ghtml.Div(
				ghtml.Class("w-full max-w-2xl"),
				ghtml.Div(append(
					[]gomponents.Node{ghtml.Class("flex space-x-4")},
					inputs...,
				)...),
			),
		),
	)
}
