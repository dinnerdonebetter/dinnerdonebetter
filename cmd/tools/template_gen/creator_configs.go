package main

import (
	"bytes"
	_ "embed"
)

//go:embed templates/creator.gotpl
var basicCreatorTemplateSrc string

func buildBasicCreatorTemplate(cfg *basicCreatorTemplateConfig) string {
	var b bytes.Buffer

	if err := parseTemplate("", basicCreatorTemplateSrc, nil).Execute(&b, cfg); err != nil {
		panic(err)
	}

	return b.String()
}

type basicCreatorTemplateConfig struct {
	Title         string
	SubmissionURL string
	Fields        []formField
}

var creatorConfigs = map[string]*basicCreatorTemplateConfig{
	"internal/services/frontend/templates/partials/generated/creators/account_creator.gotpl": {
		Title:         "New Account",
		SubmissionURL: "/accounts/new/submit",
		Fields: []formField{
			{
				LabelName:       "name",
				FormName:        "name",
				StructFieldName: "Name",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/api_client_creator.gotpl": {
		Title:         "New API Client",
		SubmissionURL: "/api_clients/new/submit",
		Fields: []formField{
			{
				LabelName:       "name",
				FormName:        "name",
				StructFieldName: "Name",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "client_id",
				FormName:        "client_id",
				StructFieldName: "ClientID",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "external ID",
				FormName:        "external_id",
				StructFieldName: "ExternalID",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/webhook_creator.gotpl": {
		Title:         "New Webhook",
		SubmissionURL: "/webhooks/new/submit",
		Fields: []formField{
			{
				LabelName:       "name",
				StructFieldName: "Name",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "Method",
				StructFieldName: "Method",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "ContentType",
				StructFieldName: "ContentType",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "URL",
				StructFieldName: "URL",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/valid_instrument_creator.gotpl": {
		Title:         "New Valid Instrument",
		SubmissionURL: "/valid_instruments/new/submit",
		Fields: []formField{
			{
				LabelName:       "name",
				FormName:        "name",
				StructFieldName: "Name",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "variant",
				FormName:        "variant",
				StructFieldName: "Variant",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "description",
				FormName:        "description",
				StructFieldName: "Description",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "iconPath",
				FormName:        "iconPath",
				StructFieldName: "IconPath",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/valid_preparation_creator.gotpl": {
		Title:         "New Valid Preparation",
		SubmissionURL: "/valid_preparations/new/submit",
		Fields: []formField{
			{
				LabelName:       "name",
				FormName:        "name",
				StructFieldName: "Name",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "description",
				FormName:        "description",
				StructFieldName: "Description",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "iconPath",
				FormName:        "iconPath",
				StructFieldName: "IconPath",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/valid_ingredient_creator.gotpl": {
		Title:         "New Valid Ingredient",
		SubmissionURL: "/valid_ingredients/new/submit",
		Fields: []formField{
			{
				LabelName:       "name",
				FormName:        "name",
				StructFieldName: "Name",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "variant",
				FormName:        "variant",
				StructFieldName: "Variant",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "description",
				FormName:        "description",
				StructFieldName: "Description",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "warning",
				FormName:        "warning",
				StructFieldName: "Warning",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsEgg",
				FormName:        "containsEgg",
				StructFieldName: "ContainsEgg",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsDairy",
				FormName:        "containsDairy",
				StructFieldName: "ContainsDairy",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsPeanut",
				FormName:        "containsPeanut",
				StructFieldName: "ContainsPeanut",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsTreeNut",
				FormName:        "containsTreeNut",
				StructFieldName: "ContainsTreeNut",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsSoy",
				FormName:        "containsSoy",
				StructFieldName: "ContainsSoy",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsWheat",
				FormName:        "containsWheat",
				StructFieldName: "ContainsWheat",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsShellfish",
				FormName:        "containsShellfish",
				StructFieldName: "ContainsShellfish",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsSesame",
				FormName:        "containsSesame",
				StructFieldName: "ContainsSesame",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsFish",
				FormName:        "containsFish",
				StructFieldName: "ContainsFish",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "containsGluten",
				FormName:        "containsGluten",
				StructFieldName: "ContainsGluten",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "animalFlesh",
				FormName:        "animalFlesh",
				StructFieldName: "AnimalFlesh",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "animalDerived",
				FormName:        "animalDerived",
				StructFieldName: "AnimalDerived",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "volumetric",
				FormName:        "volumetric",
				StructFieldName: "Volumetric",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "iconPath",
				FormName:        "iconPath",
				StructFieldName: "IconPath",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/valid_ingredient_preparation_creator.gotpl": {
		Title:         "New Valid Ingredient Preparation",
		SubmissionURL: "/valid_ingredient_preparations/new/submit",
		Fields: []formField{
			{
				LabelName:       "notes",
				FormName:        "notes",
				StructFieldName: "Notes",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "validIngredientID",
				FormName:        "validIngredientID",
				StructFieldName: "ValidIngredientID",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "validPreparationID",
				FormName:        "validPreparationID",
				StructFieldName: "ValidPreparationID",
				InputType:       "number",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/valid_preparation_instrument_creator.gotpl": {
		Title:         "New Valid Preparation Instrument",
		SubmissionURL: "/valid_preparation_instruments/new/submit",
		Fields: []formField{
			{
				LabelName:       "instrumentID",
				FormName:        "instrumentID",
				StructFieldName: "InstrumentID",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "preparationID",
				FormName:        "preparationID",
				StructFieldName: "PreparationID",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "notes",
				FormName:        "notes",
				StructFieldName: "Notes",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/recipe_step_creator.gotpl": {
		Title:         "New Recipe Step",
		SubmissionURL: "/recipe_steps/new/submit",
		Fields: []formField{
			{
				LabelName:       "index",
				FormName:        "index",
				StructFieldName: "Index",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "preparationID",
				FormName:        "preparationID",
				StructFieldName: "PreparationID",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "prerequisiteStep",
				FormName:        "prerequisiteStep",
				StructFieldName: "PrerequisiteStep",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "minEstimatedTimeInSeconds",
				FormName:        "minEstimatedTimeInSeconds",
				StructFieldName: "MinEstimatedTimeInSeconds",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "maxEstimatedTimeInSeconds",
				FormName:        "maxEstimatedTimeInSeconds",
				StructFieldName: "MaxEstimatedTimeInSeconds",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "temperatureInCelsius",
				FormName:        "temperatureInCelsius",
				StructFieldName: "TemperatureInCelsius",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "notes",
				FormName:        "notes",
				StructFieldName: "Notes",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "why",
				FormName:        "why",
				StructFieldName: "Why",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "recipeID",
				FormName:        "recipeID",
				StructFieldName: "RecipeID",
				InputType:       "number",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/recipe_step_ingredient_creator.gotpl": {
		Title:         "New Recipe Step Ingredient",
		SubmissionURL: "/recipe_step_ingredients/new/submit",
		Fields: []formField{
			{
				LabelName:       "ingredientID",
				FormName:        "ingredientID",
				StructFieldName: "IngredientID",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "name",
				FormName:        "name",
				StructFieldName: "Name",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "quantityType",
				FormName:        "quantityType",
				StructFieldName: "QuantityType",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "quantityValue",
				FormName:        "quantityValue",
				StructFieldName: "QuantityValue",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "quantityNotes",
				FormName:        "quantityNotes",
				StructFieldName: "QuantityNotes",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "productOfRecipeStep",
				FormName:        "productOfRecipeStep",
				StructFieldName: "ProductOfRecipeStep",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "ingredientNotes",
				FormName:        "ingredientNotes",
				StructFieldName: "IngredientNotes",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/recipe_step_product_creator.gotpl": {
		Title:         "New Recipe Step Product",
		SubmissionURL: "/recipe_step_products/new/submit",
		Fields: []formField{
			{
				LabelName:       "name",
				FormName:        "name",
				StructFieldName: "Name",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "quantityType",
				FormName:        "quantityType",
				StructFieldName: "QuantityType",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "quantityValue",
				FormName:        "quantityValue",
				StructFieldName: "QuantityValue",
				InputType:       "number",
				Required:        true,
			},
			{
				LabelName:       "quantityNotes",
				FormName:        "quantityNotes",
				StructFieldName: "QuantityNotes",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "recipeStepID",
				FormName:        "recipeStepID",
				StructFieldName: "RecipeStepID",
				InputType:       "number",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/invitation_creator.gotpl": {
		Title:         "New Invitation",
		SubmissionURL: "/invitations/new/submit",
		Fields: []formField{
			{
				LabelName:       "code",
				FormName:        "code",
				StructFieldName: "Code",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "consumed",
				FormName:        "consumed",
				StructFieldName: "Consumed",
				InputType:       "text",
				Required:        true,
			},
		},
	},
	"internal/services/frontend/templates/partials/generated/creators/report_creator.gotpl": {
		Title:         "New Report",
		SubmissionURL: "/reports/new/submit",
		Fields: []formField{
			{
				LabelName:       "reportType",
				FormName:        "reportType",
				StructFieldName: "ReportType",
				InputType:       "text",
				Required:        true,
			},
			{
				LabelName:       "concern",
				FormName:        "concern",
				StructFieldName: "Concern",
				InputType:       "text",
				Required:        true,
			},
		},
	},
}
