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
		RecipeID               string `jsonschema:"description=The recipe ID"`
		RecipeStepID           string `jsonschema:"description=The recipe step ID"`
		RecipeStepInstrumentID string `jsonschema:"description=The recipe step instrument ID"`
	}
)

var recipeStepInstrumentsSchema = map[string]any{
	"ID":                  stringField("The ID of the recipe step instrument"),
	"CreatedAt":           timestampField("When the recipe step instrument was created"),
	"LastUpdatedAt":       timestampField("When the recipe step instrument was last updated"),
	"ArchivedAt":          timestampField("When the recipe step instrument was soft deleted"),
	"BelongsToRecipeStep": stringField("The ID of the recipe step this instrument belongs to"),
	"Name":                stringField("Name of the instrument"),
	"Notes":               stringField("Notes about the instrument"),
	"Instrument":          objectType(validInstrumentsSchema),
	"RecipeStepProductID": stringField("The ID of the recipe step product this instrument is associated with, if any"),
	"Quantity":            uint32RangeWithOptionalMaxSchema(),
	"OptionIndex":         uintField("The option index for this instrument"),
	"PreferenceRank":      uintField("The preference rank for this instrument (0-255)"),
	"Optional":            boolField("Whether this instrument is optional"),
}

var getRecipeStepInstrumentTool = &mcp.Tool{
	Name:        "GetRecipeStepInstrument",
	Description: "Get a recipe step instrument by it's ID",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":               stringField("The ID of the recipe"),
		"RecipeStepID":           stringField("The ID of the recipe step"),
		"RecipeStepInstrumentID": stringField("The ID of the recipe step instrument to get"),
	}),
	OutputSchema: schemaObject(recipeStepInstrumentsSchema),
}

func (h *mcpToolManager) GetRecipeStepInstrument() mcp.ToolHandlerFor[*GetRecipeStepInstrumentInvocation, *mealplanning.RecipeStepInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepInstrumentInvocation) (*mcp.CallToolResult, *mealplanning.RecipeStepInstrument, error) {
		result, err := h.client.GetRecipeStepInstrument(ctx, &mealplanninggrpc.GetRecipeStepInstrumentRequest{
			RecipeID:               x.RecipeID,
			RecipeStepID:           x.RecipeStepID,
			RecipeStepInstrumentID: x.RecipeStepInstrumentID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(result.Result), nil
	}
}

type (
	GetRecipeStepInstrumentsInvocation struct {
		RecipeID     string
		RecipeStepID string
		Filter       *filtering.QueryFilter
	}

	GetRecipeStepInstrumentsResult struct {
		Results []*mealplanning.RecipeStepInstrument
	}
)

var getRecipeStepInstrumentsTool = &mcp.Tool{
	Name:        "GetRecipeStepInstruments",
	Description: "Get recipe step instruments with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":     stringField("The ID of the recipe"),
		"RecipeStepID": stringField("The ID of the recipe step"),
		"Filter":       queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(recipeStepInstrumentsSchema),
	}),
}

func (h *mcpToolManager) GetRecipeStepInstruments() mcp.ToolHandlerFor[*GetRecipeStepInstrumentsInvocation, *GetRecipeStepInstrumentsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetRecipeStepInstrumentsInvocation) (*mcp.CallToolResult, *GetRecipeStepInstrumentsResult, error) {
		results, err := h.client.GetRecipeStepInstruments(ctx, &mealplanninggrpc.GetRecipeStepInstrumentsRequest{
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
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
		RecipeID     string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID string `jsonschema:"required,description=The recipe step ID"`
	}
)

var recipeStepInstrumentCreationTool = &mcp.Tool{
	Name:        "CreateRecipeStepInstrument",
	Description: "Create a recipe step instrument",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":                        stringField("The ID of the recipe"),
		"RecipeStepID":                    stringField("The ID of the recipe step"),
		"InstrumentID":                    stringField("The ID of the instrument"),
		"RecipeStepProductID":             stringField("The ID of the recipe step product this instrument is associated with, if any"),
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
			RecipeID:     x.RecipeID,
			RecipeStepID: x.RecipeStepID,
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
		RecipeID               string `jsonschema:"required,description=The recipe ID"`
		RecipeStepID           string `jsonschema:"required,description=The recipe step ID"`
		RecipeStepInstrumentID string `jsonschema:"required,description=The recipe step instrument ID"`
	}
)

var recipeStepInstrumentUpdateTool = &mcp.Tool{
	Name:        "UpdateRecipeStepInstrument",
	Description: "Update a recipe step instrument",
	InputSchema: schemaObject(map[string]any{
		"RecipeID":               stringField("The ID of the recipe"),
		"RecipeStepID":           stringField("The ID of the recipe step"),
		"RecipeStepInstrumentID": stringField("The ID of the recipe step instrument to update"),
		"InstrumentID":           stringField("The ID of the instrument"),
		"RecipeStepProductID":    stringField("The ID of the recipe step product this instrument is associated with, if any"),
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
			RecipeID:               x.RecipeID,
			RecipeStepID:           x.RecipeStepID,
			RecipeStepInstrumentID: x.RecipeStepInstrumentID,
			Input:                  mealplanningconverters.ConvertRecipeStepInstrumentUpdateRequestInputToGRPCRecipeStepInstrumentUpdateRequestInput(x.RecipeStepInstrumentUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCRecipeStepInstrumentToRecipeStepInstrument(result.Updated), nil
	}
}

//
