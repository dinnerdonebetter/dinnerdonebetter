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
	GetRecipeStepInstrumentInvocation struct {
		RecipeID               string `jsonschema:"description=The recipe MealPlanTaskID"`
		RecipeStepID           string `jsonschema:"description=The recipe step MealPlanTaskID"`
		RecipeStepInstrumentID string `jsonschema:"description=The recipe step instrument MealPlanTaskID"`
	}
)

var recipeStepInstrumentsSchema = map[string]any{
	"MealPlanTaskID":      stringField("The MealPlanTaskID of the recipe step instrument"),
	"CreatedAt":           timestampField("When the recipe step instrument was created"),
	"LastUpdatedAt":       timestampField("When the recipe step instrument was last updated"),
	"ArchivedAt":          timestampField("When the recipe step instrument was soft deleted"),
	"BelongsToRecipeStep": stringField("The MealPlanTaskID of the recipe step this instrument belongs to"),
	"Name":                stringField("Name of the instrument"),
	"Notes":               stringField("Notes about the instrument"),
	"Instrument":          objectType(validInstrumentsSchema),
	"RecipeStepProductID": stringField("The MealPlanTaskID of the recipe step product this instrument is associated with, if any"),
	"Quantity":            uint32RangeWithOptionalMaxSchema(),
	"OptionIndex":         uintField("The option index for this instrument"),
	"PreferenceRank":      uintField("The preference rank for this instrument (0-255)"),
	"Optional":            boolField("Whether this instrument is optional"),
}

var getRecipeStepInstrumentTool = &mcp.Tool{
	Name:        "GetRecipeStepInstrument",
	Description: "Get a recipe step instrument by it's MealPlanTaskID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":               stringField("The MealPlanTaskID of the recipe"),
		"RecipeStepID":           stringField("The MealPlanTaskID of the recipe step"),
		"RecipeStepInstrumentID": stringField("The MealPlanTaskID of the recipe step instrument to get"),
	}),
	OutputSchema: schemaObject(recipeStepInstrumentsSchema),
}

func (h *mcpToolManager) GetRecipeStepInstrument() mcp.ToolHandlerFor[*GetRecipeStepInstrumentInvocation, *mealplanning.RecipeStepInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepInstrumentInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepInstrument, error) {
		result, err := h.client.GetRecipeStepInstrument(ctx, &mealplanninggrpc.GetRecipeStepInstrumentRequest{
			RecipeId:               x.RecipeID,
			RecipeStepId:           x.RecipeStepID,
			RecipeStepInstrumentId: x.RecipeStepInstrumentID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(result.Result), nil
	}
}

type (
	GetRecipeStepInstrumentsInvocation struct {
		Filter       *filtering.QueryFilter
		RecipeID     string
		RecipeStepID string
	}

	GetRecipeStepInstrumentsResult struct {
		Results []*mealplanning.RecipeStepInstrument
	}
)

var getRecipeStepInstrumentsTool = &mcp.Tool{
	Name:        "GetRecipeStepInstruments",
	Description: "Get recipe step instruments with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":     stringField("The MealPlanTaskID of the recipe"),
		"RecipeStepID": stringField("The MealPlanTaskID of the recipe step"),
		"Filter":       queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(recipeStepInstrumentsSchema)),
	}),
}

func (h *mcpToolManager) GetRecipeStepInstruments() mcp.ToolHandlerFor[*GetRecipeStepInstrumentsInvocation, *GetRecipeStepInstrumentsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepInstrumentsInvocation) (*mcp.CallToolResult, *GetRecipeStepInstrumentsResult, error) {
		results, err := h.client.GetRecipeStepInstruments(ctx, &mealplanninggrpc.GetRecipeStepInstrumentsRequest{
			RecipeId:     x.RecipeID,
			RecipeStepId: x.RecipeStepID,
			Filter:       grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetRecipeStepInstrumentsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(result))
		}

		return nil, out, nil
	}
}

type (
	CreateRecipeStepInstrumentInvocation struct {
		*mealplanning.RecipeStepInstrumentCreationRequestInput
		RecipeID     string `jsonschema:"required,description=The recipe MealPlanTaskID"`
		RecipeStepID string `jsonschema:"required,description=The recipe step MealPlanTaskID"`
	}
)

var recipeStepInstrumentCreationTool = &mcp.Tool{
	Name:        "CreateRecipeStepInstrument",
	Description: "Create a recipe step instrument",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                        stringField("The MealPlanTaskID of the recipe"),
		"RecipeStepID":                    stringField("The MealPlanTaskID of the recipe step"),
		"InstrumentID":                    stringField("The MealPlanTaskID of the instrument"),
		"RecipeStepProductID":             stringField("The MealPlanTaskID of the recipe step product this instrument is associated with, if any"),
		"ProductOfRecipeStepIndex":        uintField("The index of the recipe step that produces this instrument, if any"),
		"ProductOfRecipeStepProductIndex": uintField("The index of the recipe step product that produces this instrument, if any"),
		"Name":                            stringField("Name of the instrument"),
		"Notes":                           stringField("Notes about the instrument"),
		"Quantity":                        uint32RangeWithOptionalMaxSchema(),
		"OptionIndex":                     uintField("The option index for this instrument"),
		"PreferenceRank":                  uintField("The preference rank for this instrument (0-255)"),
		"Optional":                        boolField("Whether this instrument is optional"),
	}),
	OutputSchema: schemaObject(recipeStepInstrumentsSchema),
}

func (h *mcpToolManager) CreateRecipeStepInstrument() mcp.ToolHandlerFor[*CreateRecipeStepInstrumentInvocation, *mealplanning.RecipeStepInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateRecipeStepInstrumentInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepInstrument, error) {
		result, err := h.client.CreateRecipeStepInstrument(ctx, &mealplanninggrpc.CreateRecipeStepInstrumentRequest{
			RecipeId:     x.RecipeID,
			RecipeStepId: x.RecipeStepID,
			Input:        mealplanningconverters.ConvertRecipeStepInstrumentCreationRequestInputToGRPCRecipeStepInstrumentCreationRequestInput(x.RecipeStepInstrumentCreationRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(result.Created), nil
	}
}

type (
	UpdateRecipeStepInstrumentInvocation struct {
		*mealplanning.RecipeStepInstrumentUpdateRequestInput
		RecipeID               string `jsonschema:"required,description=The recipe MealPlanTaskID"`
		RecipeStepID           string `jsonschema:"required,description=The recipe step MealPlanTaskID"`
		RecipeStepInstrumentID string `jsonschema:"required,description=The recipe step instrument MealPlanTaskID"`
	}
)

var recipeStepInstrumentUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipeStepInstrument",
	Description: "Update a recipe step instrument",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":               stringField("The MealPlanTaskID of the recipe"),
		"RecipeStepID":           stringField("The MealPlanTaskID of the recipe step"),
		"RecipeStepInstrumentID": stringField("The MealPlanTaskID of the recipe step instrument to update"),
		"InstrumentID":           stringField("The MealPlanTaskID of the instrument"),
		"RecipeStepProductID":    stringField("The MealPlanTaskID of the recipe step product this instrument is associated with, if any"),
		"Name":                   stringField("Name of the instrument"),
		"Notes":                  stringField("Notes about the instrument"),
		"Quantity":               uint32RangeWithOptionalMaxSchema(),
		"OptionIndex":            uintField("The option index for this instrument"),
		"PreferenceRank":         uintField("The preference rank for this instrument (0-255)"),
		"Optional":               boolField("Whether this instrument is optional"),
	}),
	OutputSchema: schemaObject(recipeStepInstrumentsSchema),
}

func (h *mcpToolManager) UpdateRecipeStepInstrument() mcp.ToolHandlerFor[*UpdateRecipeStepInstrumentInvocation, *mealplanning.RecipeStepInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateRecipeStepInstrumentInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepInstrument, error) {
		result, err := h.client.UpdateRecipeStepInstrument(ctx, &mealplanninggrpc.UpdateRecipeStepInstrumentRequest{
			RecipeId:               x.RecipeID,
			RecipeStepId:           x.RecipeStepID,
			RecipeStepInstrumentId: x.RecipeStepInstrumentID,
			Input:                  mealplanningconverters.ConvertRecipeStepInstrumentUpdateRequestInputToGRPCRecipeStepInstrumentUpdateRequestInput(x.RecipeStepInstrumentUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(result.Updated), nil
	}
}

//
