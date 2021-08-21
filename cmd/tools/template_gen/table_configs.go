package main

import (
	"bytes"
	_ "embed"
)

type basicTableTemplateConfig struct {
	SearchURL            string
	CreatorPageURL       string
	RowDataFieldName     string
	Title                string
	CreatorPagePushURL   string
	CellFields           []string
	Columns              []string
	EnableSearch         bool
	ExcludeIDRow         bool
	ExcludeLink          bool
	IncludeLastUpdatedOn bool
	IncludeCreatedOn     bool
	IncludeDeleteRow     bool
}

//go:embed templates/table.gotpl
var basicTableTemplateSrc string

func buildBasicTableTemplate(cfg *basicTableTemplateConfig) string {
	var b bytes.Buffer

	if err := parseTemplate("", basicTableTemplateSrc, nil).Execute(&b, cfg); err != nil {
		panic(err)
	}

	return b.String()
}

var tableConfigs = map[string]*basicTableTemplateConfig{
	"internal/services/frontend/templates/partials/generated/tables/api_clients_table.gotpl": {
		Title:              "API Clients",
		CreatorPagePushURL: "/api_clients/new",
		CreatorPageURL:     "/dashboard_pages/api_clients/new",
		Columns: []string{
			"ID",
			"Name",
			"External ID",
			"Client ID",
			"Belongs To User",
			"Created On",
		},
		CellFields: []string{
			"ID",
			"Name",
			"ExternalID",
			"ClientID",
			"BelongsToUser",
			"CreatedOn",
		},
		RowDataFieldName:     "Clients",
		IncludeLastUpdatedOn: false,
		IncludeCreatedOn:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/households_table.gotpl": {
		Title:              "Households",
		CreatorPagePushURL: "/households/new",
		CreatorPageURL:     "/dashboard_pages/households/new",
		Columns: []string{
			"ID",
			"Name",
			"External ID",
			"Belongs To User",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Name",
			"ExternalID",
			"BelongsToUser",
		},
		RowDataFieldName:     "Households",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/users_table.gotpl": {
		Title: "Users",
		Columns: []string{
			"ID",
			"Username",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Username",
		},
		EnableSearch:         true,
		RowDataFieldName:     "Users",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     false,
		ExcludeLink:          true,
	},
	"internal/services/frontend/templates/partials/generated/tables/webhooks_table.gotpl": {
		Title:              "Webhooks",
		CreatorPagePushURL: "/households/webhooks/new",
		CreatorPageURL:     "/dashboard_pages/households/webhooks/new",
		Columns: []string{
			"ID",
			"Name",
			"Method",
			"URL",
			"Content Type",
			"Belongs To Household",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Name",
			"Method",
			"URL",
			"ContentType",
			"BelongsToHousehold",
		},
		RowDataFieldName:     "Webhooks",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/valid_instruments_table.gotpl": {
		Title:              "Valid Instruments",
		CreatorPagePushURL: "/valid_instruments/new",
		CreatorPageURL:     "/dashboard_pages/valid_instruments/new",
		Columns: []string{
			"ID",
			"Name",
			"Variant",
			"Description",
			"IconPath",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Name",
			"Variant",
			"Description",
			"IconPath",
		},
		RowDataFieldName:     "ValidInstruments",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/valid_preparations_table.gotpl": {
		Title:              "Valid Preparations",
		CreatorPagePushURL: "/valid_preparations/new",
		CreatorPageURL:     "/dashboard_pages/valid_preparations/new",
		Columns: []string{
			"ID",
			"Name",
			"Description",
			"IconPath",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Name",
			"Description",
			"IconPath",
		},
		RowDataFieldName:     "ValidPreparations",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/valid_ingredients_table.gotpl": {
		Title:              "Valid Ingredients",
		CreatorPagePushURL: "/valid_ingredients/new",
		CreatorPageURL:     "/dashboard_pages/valid_ingredients/new",
		Columns: []string{
			"ID",
			"Name",
			"Variant",
			"Description",
			"Warning",
			"ContainsEgg",
			"ContainsDairy",
			"ContainsPeanut",
			"ContainsTreeNut",
			"ContainsSoy",
			"ContainsWheat",
			"ContainsShellfish",
			"ContainsSesame",
			"ContainsFish",
			"ContainsGluten",
			"AnimalFlesh",
			"AnimalDerived",
			"Volumetric",
			"IconPath",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Name",
			"Variant",
			"Description",
			"Warning",
			"ContainsEgg",
			"ContainsDairy",
			"ContainsPeanut",
			"ContainsTreeNut",
			"ContainsSoy",
			"ContainsWheat",
			"ContainsShellfish",
			"ContainsSesame",
			"ContainsFish",
			"ContainsGluten",
			"AnimalFlesh",
			"AnimalDerived",
			"Volumetric",
			"IconPath",
		},
		RowDataFieldName:     "ValidIngredients",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/valid_ingredient_preparations_table.gotpl": {
		Title:              "Valid Ingredient Preparations",
		CreatorPagePushURL: "/valid_ingredient_preparations/new",
		CreatorPageURL:     "/dashboard_pages/valid_ingredient_preparations/new",
		Columns: []string{
			"ID",
			"Notes",
			"ValidIngredientID",
			"ValidPreparationID",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Notes",
			"ValidIngredientID",
			"ValidPreparationID",
		},
		RowDataFieldName:     "ValidIngredientPreparations",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/valid_preparation_instruments_table.gotpl": {
		Title:              "Valid Preparation Instruments",
		CreatorPagePushURL: "/valid_preparation_instruments/new",
		CreatorPageURL:     "/dashboard_pages/valid_preparation_instruments/new",
		Columns: []string{
			"ID",
			"InstrumentID",
			"PreparationID",
			"Notes",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"InstrumentID",
			"PreparationID",
			"Notes",
		},
		RowDataFieldName:     "ValidPreparationInstruments",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/recipes_table.gotpl": {
		Title:              "Recipes",
		CreatorPagePushURL: "/recipes/new",
		CreatorPageURL:     "/dashboard_pages/recipes/new",
		Columns: []string{
			"ID",
			"Name",
			"Source",
			"Description",
			"InspiredByRecipeID",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Name",
			"Source",
			"Description",
			"InspiredByRecipeID",
		},
		RowDataFieldName:     "Recipes",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/recipe_steps_table.gotpl": {
		Title:              "Recipe Steps",
		CreatorPagePushURL: "/recipe_steps/new",
		CreatorPageURL:     "/dashboard_pages/recipe_steps/new",
		Columns: []string{
			"ID",
			"Index",
			"PreparationID",
			"PrerequisiteStep",
			"MinEstimatedTimeInSeconds",
			"MaxEstimatedTimeInSeconds",
			"TemperatureInCelsius",
			"Notes",
			"Why",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Index",
			"PreparationID",
			"PrerequisiteStep",
			"MinEstimatedTimeInSeconds",
			"MaxEstimatedTimeInSeconds",
			"TemperatureInCelsius",
			"Notes",
			"Why",
		},
		RowDataFieldName:     "RecipeSteps",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/recipe_step_ingredients_table.gotpl": {
		Title:              "Recipe Step Ingredients",
		CreatorPagePushURL: "/recipe_step_ingredients/new",
		CreatorPageURL:     "/dashboard_pages/recipe_step_ingredients/new",
		Columns: []string{
			"ID",
			"IngredientID",
			"Name",
			"QuantityType",
			"QuantityValue",
			"QuantityNotes",
			"ProductOfRecipeStep",
			"IngredientNotes",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"IngredientID",
			"Name",
			"QuantityType",
			"QuantityValue",
			"QuantityNotes",
			"ProductOfRecipeStep",
			"IngredientNotes",
		},
		RowDataFieldName:     "RecipeStepIngredients",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/recipe_step_products_table.gotpl": {
		Title:              "Recipe Step Products",
		CreatorPagePushURL: "/recipe_step_products/new",
		CreatorPageURL:     "/dashboard_pages/recipe_step_products/new",
		Columns: []string{
			"ID",
			"Name",
			"QuantityType",
			"QuantityValue",
			"QuantityNotes",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Name",
			"QuantityType",
			"QuantityValue",
			"QuantityNotes",
		},
		RowDataFieldName:     "RecipeStepProducts",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/invitations_table.gotpl": {
		Title:              "Invitations",
		CreatorPagePushURL: "/invitations/new",
		CreatorPageURL:     "/dashboard_pages/invitations/new",
		Columns: []string{
			"ID",
			"Code",
			"Consumed",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"Code",
			"Consumed",
		},
		RowDataFieldName:     "Invitations",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
	"internal/services/frontend/templates/partials/generated/tables/reports_table.gotpl": {
		Title:              "Reports",
		CreatorPagePushURL: "/reports/new",
		CreatorPageURL:     "/dashboard_pages/reports/new",
		Columns: []string{
			"ID",
			"ReportType",
			"Concern",
			"Last Updated On",
			"Created On",
		},
		CellFields: []string{
			"ReportType",
			"Concern",
		},
		RowDataFieldName:     "Reports",
		IncludeLastUpdatedOn: true,
		IncludeCreatedOn:     true,
		IncludeDeleteRow:     true,
	},
}
