package main

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	mealplanningconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type (
	GetRecipeStepIngredientInvocation struct {
		RecipeID              string `jsonschema:"description=The recipe ID"`
		RecipeStepID          string `jsonschema:"description=The recipe step ID"`
		RecipeStepIngredientID string `jsonschema:"description=The recipe step ingredient ID"`
	}
)

var recipeStepIngredientsSchema = map[string]any{
	"ID":                        stringField("The ID of the recipe step ingredient"),
	"CreatedAt":                 timestampField("When the recipe step ingredient was created"),
	"LastUpdatedAt":             timestampField("When the recipe step ingredient was last updated"),
	"ArchivedAt":                timestampField("When the recipe step ingredient was soft deleted"),
	"BelongsToRecipeStep":       stringField("The ID of the recipe step this ingredient belongs to"),
	"Name":                      stringField("Name of the ingredient"),
	"QuantityNotes":             stringField("Notes about the quantity"),
	"IngredientNotes":           stringField("Notes about the ingredient"),
	"Ingredient":                objectType(validIngredientsSchema),
	"MeasurementUnit":           objectType(validMeasurementUnitsSchema),
	"Quantity":                  float32RangeWithOptionalMaxSchema(),
	"RecipeStepProductID":        stringField("The ID of the recipe step product this ingredient is associated with, if any"),
	"ProductOfRecipeID":          stringField("The ID of the recipe that produces this ingredient, if any"),
	"ProductPercentageToUse":     floatField("The percentage of the product to use, if any"),
	"VesselIndex":                uintField("The index of the vessel this ingredient is in, if any"),
	"OptionIndex":                uintField("The option index for this ingredient"),
	"Optional":                   boolField("Whether this ingredient is optional"),
	"ToTaste":                    boolField("Whether this ingredient is 'to taste'"),
}

var getRecipeStepIngredientTool = &mcp.Tool{
	Name:        "GetRecipeStepIngredient",
	Description: "Get a recipe step ingredient by it's ID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":              stringField("The ID of the recipe"),
		"RecipeStepID":          stringField("The ID of the recipe step"),
		"RecipeStepIngredientID": stringField("The ID of the recipe step ingredient to get"),
	}),
	OutputSchema: schemaObject(recipeStepIngredientsSchema),
}

func (h *mcpToolManager) GetRecipeStepIngredient() mcp.ToolHandlerFor[*GetRecipeStepIngredientInvocation, *mealplanning.RecipeStepIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepIngredientInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepIngredient, error) {
		result, err := h.client.GetRecipeStepIngredient(ctx, &mealplanninggrpc.GetRecipeStepIngredientRequest{
			RecipeID:              x.RecipeID,
			RecipeStepID:          x.RecipeStepID,
			RecipeStepIngredientID: x.RecipeStepIngredientID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepIngredientToRecipeStepIngredient(result.Result), nil
	}
}

type (
	GetRecipeStepIngredientsInvocation struct {
		RecipeID     string
		RecipeStepID string
		Filter       *filtering.QueryFilter
	}

	GetRecipeStepIngredientsResult struct {
		Results []*mealplanning.RecipeStepIngredient
	}
)

var getRecipeStepIngredientsTool = &mcp.Tool{
	Name:        "GetRecipeStepIngredients",
	Description: "Get recipe step ingredients with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":     stringField("The ID of the recipe"),
		"RecipeStepID": stringField("The ID of the recipe step"),
		"Filter":       queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(recipeStepIngredientsSchema),
	}),
}

func (h *mcpToolManager) GetRecipeStepIngredients() mcp.ToolHandlerFor[*GetRecipeStepIngredientsInvocation, *GetRecipeStepIngredientsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepIngredientsInvocation) (*mcp.CallToolResult, *GetRecipeStepIngredientsResult, error) {
		results, err := h.client.GetRecipeStepIngredients(ctx, &mealplanninggrpc.GetRecipeStepIngredientsRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Filter:       grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetRecipeStepIngredientsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipeStepIngredientToRecipeStepIngredient(result))
		}

		return nil, out, nil
	}
}

type (
	CreateRecipeStepIngredientInvocation struct {
		*mealplanning.RecipeStepIngredientCreationRequestInput
		RecipeID     string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID string `jsonschema:"required,description=The recipe step ID"`
	}
)

var recipeStepIngredientCreationTool = &mcp.Tool{
	Name:        "CreateRecipeStepIngredient",
	Description: "Create a recipe step ingredient",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                      stringField("The ID of the recipe"),
		"RecipeStepID":                  stringField("The ID of the recipe step"),
		"IngredientID":                  stringField("The ID of the ingredient"),
		"ProductOfRecipeID":             stringField("The ID of the recipe that produces this ingredient, if any"),
		"ProductOfRecipeStepIndex":      uintField("The index of the recipe step that produces this ingredient, if any"),
		"ProductOfRecipeStepProductIndex": uintField("The index of the recipe step product that produces this ingredient, if any"),
		"RecipeStepProductID":           stringField("The ID of the recipe step product this ingredient is associated with, if any"),
		"MeasurementUnitID":             stringField("The ID of the measurement unit"),
		"Name":                          stringField("Name of the ingredient"),
		"QuantityNotes":                 stringField("Notes about the quantity"),
		"IngredientNotes":               stringField("Notes about the ingredient"),
		"Quantity":                      float32RangeWithOptionalMaxSchema(),
		"VesselIndex":                   uintField("The index of the vessel this ingredient is in, if any"),
		"ProductPercentageToUse":        floatField("The percentage of the product to use, if any"),
		"OptionIndex":                   uintField("The option index for this ingredient"),
		"Optional":                      boolField("Whether this ingredient is optional"),
		"ToTaste":                       boolField("Whether this ingredient is 'to taste'"),
	}),
	OutputSchema: schemaObject(recipeStepIngredientsSchema),
}

func (h *mcpToolManager) CreateRecipeStepIngredient() mcp.ToolHandlerFor[*CreateRecipeStepIngredientInvocation, *mealplanning.RecipeStepIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateRecipeStepIngredientInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepIngredient, error) {
		result, err := h.client.CreateRecipeStepIngredient(ctx, &mealplanninggrpc.CreateRecipeStepIngredientRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
			Input:       mealplanningconverters.ConvertRecipeStepIngredientCreationRequestInputToGRPCRecipeStepIngredientCreationRequestInput(x.RecipeStepIngredientCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepIngredientToRecipeStepIngredient(result.Created), nil
	}
}

type (
	UpdateRecipeStepIngredientInvocation struct {
		*mealplanning.RecipeStepIngredientUpdateRequestInput
		RecipeID              string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID          string `jsonschema:"required,description=The recipe step ID"`
		RecipeStepIngredientID string `jsonschema:"required,description=The recipe step ingredient ID"`
	}
)

var recipeStepIngredientUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipeStepIngredient",
	Description: "Update a recipe step ingredient",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":              stringField("The ID of the recipe"),
		"RecipeStepID":          stringField("The ID of the recipe step"),
		"RecipeStepIngredientID": stringField("The ID of the recipe step ingredient to update"),
		"IngredientID":          stringField("The ID of the ingredient"),
		"RecipeStepProductID":   stringField("The ID of the recipe step product this ingredient is associated with, if any"),
		"ProductOfRecipeID":     stringField("The ID of the recipe that produces this ingredient, if any"),
		"MeasurementUnitID":     stringField("The ID of the measurement unit"),
		"Name":                  stringField("Name of the ingredient"),
		"QuantityNotes":         stringField("Notes about the quantity"),
		"IngredientNotes":       stringField("Notes about the ingredient"),
		"Quantity":              float32RangeWithOptionalMaxSchema(),
		"VesselIndex":           uintField("The index of the vessel this ingredient is in, if any"),
		"ProductPercentageToUse": floatField("The percentage of the product to use, if any"),
		"OptionIndex":           uintField("The option index for this ingredient"),
		"Optional":              boolField("Whether this ingredient is optional"),
		"ToTaste":               boolField("Whether this ingredient is 'to taste'"),
	}),
	OutputSchema: schemaObject(recipeStepIngredientsSchema),
}

func (h *mcpToolManager) UpdateRecipeStepIngredient() mcp.ToolHandlerFor[*UpdateRecipeStepIngredientInvocation, *mealplanning.RecipeStepIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateRecipeStepIngredientInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepIngredient, error) {
		result, err := h.client.UpdateRecipeStepIngredient(ctx, &mealplanninggrpc.UpdateRecipeStepIngredientRequest{
			RecipeID:              x.RecipeID,
			RecipeStepID:          x.RecipeStepID,
			RecipeStepIngredientID: x.RecipeStepIngredientID,
			Input:                 mealplanningconverters.ConvertRecipeStepIngredientUpdateRequestInputToGRPCRecipeStepIngredientUpdateRequestInput(x.RecipeStepIngredientUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepIngredientToRecipeStepIngredient(result.Updated), nil
	}
}

//
