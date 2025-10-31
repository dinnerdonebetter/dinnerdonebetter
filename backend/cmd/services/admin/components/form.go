package components

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/design"
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// FieldConfig holds configuration for a single form field
type FieldConfig struct {
	// Name is the struct field name
	Name string
	// DisplayName is the label to show (if empty, uses camelCaseToTitleCase)
	DisplayName string
	// Enabled determines if the field is editable (default is disabled)
	Enabled bool
	// Validation holds HTMX-compatible validation rules
	Validation *FieldValidation
	// InputType specifies the HTML input type (default: text)
	InputType string
	// Placeholder text for the input
	Placeholder string
	// CustomRenderer allows custom rendering of the entire field
	CustomRenderer func(fieldName string, value any, config FieldConfig, palette *design.Palette) g.Node
}

// FieldValidation holds HTMX-compatible validation rules
type FieldValidation struct {
	// Required makes the field required
	Required bool
	// MinLength specifies minimum length for text inputs
	MinLength int
	// MaxLength specifies maximum length for text inputs
	MaxLength int
	// Pattern is a regex pattern for validation
	Pattern string
	// Min is minimum value for number inputs
	Min *float64
	// Max is maximum value for number inputs
	Max *float64
	// HTMXValidate enables HTMX validation endpoint
	HTMXValidate string
	// CustomMessage is a custom validation message
	CustomMessage string
}

// FormRow represents a row of fields in the form
type FormRow struct {
	// Fields are the field names to include in this row
	Fields []string
	// Columns specifies the number of columns (defaults to len(Fields))
	Columns int
}

// FormOptions holds configuration options for the form
type FormOptions[T any] struct {
	// FormID sets the HTML ID attribute for the form
	FormID string

	// Action is the form submission URL
	Action string

	// Method is the HTTP method (default: POST)
	Method string

	// FieldConfigs maps field names to their configuration
	FieldConfigs map[string]*FieldConfig

	// EnabledFields is a shorthand to enable multiple fields by name
	EnabledFields []string

	// FormRows defines the layout of fields in rows
	// If nil or empty, each field gets its own row
	FormRows []FormRow

	// Palette allows customization of colors
	Palette *design.Palette

	// HTMXSwap specifies the HTMX swap strategy (default: innerHTML)
	HTMXSwap string

	// HTMXTarget specifies the HTMX target element
	HTMXTarget string

	// HTMXPushURL determines if URL should be pushed to history
	HTMXPushURL bool

	// SubmitButtonText customizes the submit button text
	SubmitButtonText string

	// ShowCancelButton determines if cancel button is shown
	ShowCancelButton bool

	// CancelButtonText customizes the cancel button text
	CancelButtonText string

	// CancelURL is where the cancel button navigates to
	CancelURL string

	// AdditionalButtons allows adding custom buttons
	AdditionalButtons []g.Node

	// CSSClasses allows adding custom CSS classes to the form
	CSSClasses string
}

// Form creates a generic HTML form from a struct with HTMX support
func Form[T any](data T, options *FormOptions[T]) (g.Node, error) {
	if options == nil {
		options = &FormOptions[T]{}
	}

	// Get palette (use standard if not provided)
	palette := options.Palette
	if palette == nil {
		palette = &design.StandardPalette
	}

	// Extract field information from the data
	fields, values, err := extractFieldsAndValues(data)
	if err != nil {
		return nil, fmt.Errorf("failed to extract fields: %w", err)
	}

	// Apply field configurations
	fields = applyFieldConfigs(fields, options)

	// Build form attributes
	formAttrs := buildFormAttributes(options)

	// Generate form rows
	formRows := generateFormRows(fields, values, options, palette)

	// Generate buttons section
	buttonsSection := generateFormButtons(options, palette)

	return ghtml.Form(
		append(formAttrs,
			ghtml.Class(buildFormClasses(options.CSSClasses)),
			ghtml.Div(
				ghtml.Class("space-y-6"),
				g.Group(formRows),
			),
			buttonsSection,
		)...,
	), nil
}

// extractFieldsAndValues uses reflection to get field information and values from a struct
func extractFieldsAndValues(item any) ([]fieldInfo, map[string]any, error) {
	var fields []fieldInfo
	values := make(map[string]any)

	v := reflect.ValueOf(item)
	if !v.IsValid() {
		return nil, nil, fmt.Errorf("invalid value")
	}

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil, fmt.Errorf("nil pointer")
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("expected struct or pointer to struct, got %s", v.Kind())
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Create a display name
		displayName := camelCaseToTitleCase(field.Name)

		fields = append(fields, fieldInfo{
			Name:        field.Name,
			DisplayName: displayName,
			Type:        field.Type,
		})

		// Extract value
		fieldValue := v.Field(i)
		if fieldValue.CanInterface() {
			values[field.Name] = fieldValue.Interface()
		}
	}

	return fields, values, nil
}

// applyFieldConfigs applies field configurations and enabled fields list
func applyFieldConfigs[T any](fields []fieldInfo, options *FormOptions[T]) []fieldInfo {
	// Create a set of enabled fields for quick lookup
	enabledSet := make(map[string]bool)
	for _, name := range options.EnabledFields {
		enabledSet[name] = true
	}

	// Apply custom display names from configs
	for i := range fields {
		if config, exists := options.FieldConfigs[fields[i].Name]; exists {
			if config.DisplayName != "" {
				fields[i].DisplayName = config.DisplayName
			}
		}
	}

	return fields
}

// buildFormAttributes constructs the form HTML attributes
func buildFormAttributes[T any](options *FormOptions[T]) []g.Node {
	attrs := []g.Node{}

	if options.FormID != "" {
		attrs = append(attrs, ghtml.ID(options.FormID))
	}

	if options.Action != "" {
		attrs = append(attrs, ghtml.Action(options.Action))
	}

	method := options.Method
	if method == "" {
		method = "POST"
	}
	attrs = append(attrs, ghtml.Method(method))

	// HTMX attributes
	if options.Action != "" {
		attrs = append(attrs, g.Attr("hx-post", options.Action))
	}

	if options.HTMXTarget != "" {
		attrs = append(attrs, g.Attr("hx-target", options.HTMXTarget))
	}

	swap := options.HTMXSwap
	if swap == "" {
		swap = "innerHTML"
	}
	attrs = append(attrs, g.Attr("hx-swap", swap))

	if options.HTMXPushURL {
		attrs = append(attrs, g.Attr("hx-push-url", "true"))
	}

	return attrs
}

// generateFormRows generates the form field rows based on layout configuration
func generateFormRows[T any](fields []fieldInfo, values map[string]any, options *FormOptions[T], palette *design.Palette) []g.Node {
	var rows []g.Node

	// If no custom layout specified, create one field per row
	if len(options.FormRows) == 0 {
		for _, field := range fields {
			config := getFieldConfig(field.Name, options)
			fieldNode := generateFormField(field, values[field.Name], config, palette)
			rows = append(rows, ghtml.Div(
				ghtml.Class("grid grid-cols-1 gap-4"),
				fieldNode,
			))
		}
		return rows
	}

	// Create a map for quick field lookup
	fieldMap := make(map[string]fieldInfo)
	for _, field := range fields {
		fieldMap[field.Name] = field
	}

	// Generate rows based on custom layout
	for _, formRow := range options.FormRows {
		var rowFields []g.Node
		columns := formRow.Columns
		if columns == 0 {
			columns = len(formRow.Fields)
		}

		for _, fieldName := range formRow.Fields {
			if field, exists := fieldMap[fieldName]; exists {
				config := getFieldConfig(fieldName, options)
				fieldNode := generateFormField(field, values[fieldName], config, palette)
				rowFields = append(rowFields, fieldNode)
			}
		}

		if len(rowFields) > 0 {
			gridClass := fmt.Sprintf("grid grid-cols-1 md:grid-cols-%d gap-4", columns)
			rows = append(rows, ghtml.Div(
				ghtml.Class(gridClass),
				g.Group(rowFields),
			))
		}
	}

	return rows
}

// getFieldConfig retrieves the field configuration or returns defaults
func getFieldConfig[T any](fieldName string, options *FormOptions[T]) FieldConfig {
	// Check if field is in enabled list
	enabled := false
	for _, name := range options.EnabledFields {
		if name == fieldName {
			enabled = true
			break
		}
	}

	// Get or create config
	var config FieldConfig
	if options.FieldConfigs != nil {
		if c, exists := options.FieldConfigs[fieldName]; exists {
			config = *c
		}
	}

	// Apply enabled status
	if enabled {
		config.Enabled = true
	}

	config.Name = fieldName

	return config
}

// generateFormField generates a single form field
func generateFormField(field fieldInfo, value any, config FieldConfig, palette *design.Palette) g.Node {
	// Use custom renderer if provided
	if config.CustomRenderer != nil {
		return config.CustomRenderer(field.Name, value, config, palette)
	}

	displayName := field.DisplayName
	if config.DisplayName != "" {
		displayName = config.DisplayName
	}

	inputType := config.InputType
	if inputType == "" {
		inputType = inferInputType(field.Type)
	}

	// Build input attributes
	inputAttrs := []g.Node{
		ghtml.Type(inputType),
		ghtml.ID(field.Name),
		ghtml.Name(field.Name),
		ghtml.Class(buildInputClasses(palette)),
	}

	// Add value
	if value != nil {
		inputAttrs = append(inputAttrs, ghtml.Value(formatValue(value)))
	}

	// Add disabled attribute if not enabled
	if !config.Enabled {
		inputAttrs = append(inputAttrs, ghtml.Disabled())
	}

	// Add placeholder
	if config.Placeholder != "" {
		inputAttrs = append(inputAttrs, ghtml.Placeholder(config.Placeholder))
	}

	// Add validation attributes
	if config.Validation != nil {
		inputAttrs = append(inputAttrs, buildValidationAttributes(config.Validation)...)
	}

	return ghtml.Div(
		ghtml.Class("flex flex-col"),
		ghtml.Label(
			ghtml.For(field.Name),
			ghtml.Class(fmt.Sprintf("block text-sm font-medium %s mb-2", design.TextColor(palette.Text))),
			g.Text(displayName),
			g.If(config.Validation != nil && config.Validation.Required,
				ghtml.Span(
					ghtml.Class("text-red-500 ml-1"),
					g.Text("*"),
				),
			),
		),
		ghtml.Input(inputAttrs...),
		g.If(config.Validation != nil && config.Validation.CustomMessage != "",
			ghtml.P(
				ghtml.Class("mt-1 text-sm text-gray-500"),
				g.Text(config.Validation.CustomMessage),
			),
		),
	)
}

// inferInputType infers the HTML input type from the Go type
func inferInputType(t reflect.Type) string {
	// Handle pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "checkbox"
	default:
		// Check if it's a time.Time
		if t.String() == "time.Time" {
			return "datetime-local"
		}
		return "text"
	}
}

// formatValue formats a value for display in an input field
func formatValue(value any) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case time.Time:
		if v.IsZero() {
			return ""
		}
		// Format for datetime-local input
		return v.Format("2006-01-02T15:04")
	case *time.Time:
		if v == nil || v.IsZero() {
			return ""
		}
		return v.Format("2006-01-02T15:04")
	case string:
		return v
	case *string:
		if v == nil {
			return ""
		}
		return *v
	case bool:
		if v {
			return "checked"
		}
		return ""
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// buildValidationAttributes constructs HTML5 and HTMX validation attributes
func buildValidationAttributes(validation *FieldValidation) []g.Node {
	var attrs []g.Node

	if validation.Required {
		attrs = append(attrs, ghtml.Required())
	}

	if validation.MinLength > 0 {
		attrs = append(attrs, ghtml.MinLength(fmt.Sprintf("%d", validation.MinLength)))
	}

	if validation.MaxLength > 0 {
		attrs = append(attrs, ghtml.MaxLength(fmt.Sprintf("%d", validation.MaxLength)))
	}

	if validation.Pattern != "" {
		attrs = append(attrs, ghtml.Pattern(validation.Pattern))
	}

	if validation.Min != nil {
		attrs = append(attrs, ghtml.Min(fmt.Sprintf("%f", *validation.Min)))
	}

	if validation.Max != nil {
		attrs = append(attrs, ghtml.Max(fmt.Sprintf("%f", *validation.Max)))
	}

	if validation.HTMXValidate != "" {
		attrs = append(attrs,
			g.Attr("hx-post", validation.HTMXValidate),
			g.Attr("hx-trigger", "blur"),
			g.Attr("hx-target", "next .validation-message"),
		)
	}

	return attrs
}

// buildInputClasses constructs the CSS classes for input fields
func buildInputClasses(palette *design.Palette) string {
	return fmt.Sprintf("block w-full px-3 py-2 border %s rounded-md shadow-sm focus:outline-none focus:ring-1 focus:ring-%s focus:border-%s sm:text-sm disabled:bg-gray-100 disabled:cursor-not-allowed",
		design.BorderColor(palette.Background),
		palette.Primary.Value,
		palette.Primary.Value,
	)
}

// generateFormButtons generates the form buttons section
func generateFormButtons[T any](options *FormOptions[T], palette *design.Palette) g.Node {
	var buttons []g.Node

	// Submit button
	submitText := options.SubmitButtonText
	if submitText == "" {
		submitText = "Submit"
	}

	buttons = append(buttons, ghtml.Button(
		ghtml.Type("submit"),
		ghtml.Class(fmt.Sprintf("inline-flex justify-center px-4 py-2 text-sm font-medium %s %s border border-transparent rounded-md shadow-sm hover:%s focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-%s",
			design.TextColor(design.Color{Value: "white"}),
			design.Background(palette.Primary),
			design.Background(design.Color{Value: palette.Primary.Value + "-700"}),
			palette.Primary.Value,
		)),
		g.Text(submitText),
	))

	// Cancel button
	if options.ShowCancelButton {
		cancelText := options.CancelButtonText
		if cancelText == "" {
			cancelText = "Cancel"
		}

		cancelAttrs := []g.Node{
			ghtml.Type("button"),
			ghtml.Class(fmt.Sprintf("inline-flex justify-center px-4 py-2 text-sm font-medium %s %s border %s rounded-md shadow-sm hover:%s focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-%s",
				design.TextColor(palette.Text),
				design.Background(design.Color{Value: "white"}),
				design.BorderColor(palette.Text),
				design.Background(palette.Background),
				palette.Primary.Value,
			)),
			g.Text(cancelText),
		}

		if options.CancelURL != "" {
			cancelAttrs = append(cancelAttrs,
				g.Attr("hx-get", options.CancelURL),
				g.Attr("hx-target", "body"),
				g.Attr("hx-swap", "innerHTML"),
				g.Attr("hx-push-url", "true"),
			)
		}

		buttons = append(buttons, ghtml.Button(cancelAttrs...))
	}

	// Additional buttons
	buttons = append(buttons, options.AdditionalButtons...)

	return ghtml.Div(
		ghtml.Class("flex justify-end space-x-3 mt-6 pt-6 border-t border-gray-200"),
		g.Group(buttons),
	)
}

// buildFormClasses constructs the CSS classes for the form
func buildFormClasses(customClasses string) string {
	baseClasses := "bg-white shadow rounded-lg p-6"
	if customClasses != "" {
		return fmt.Sprintf("%s %s", baseClasses, customClasses)
	}
	return baseClasses
}
