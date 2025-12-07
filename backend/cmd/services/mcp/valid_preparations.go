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
	GetValidPreparationInvocation struct {
		ValidPreparationID string `jsonschema:"description=The preparation ID"`
	}
)

var validPreparationsSchema = map[string]any{
	"ID":                          stringField("The ID of the valid preparation"),
	"CreatedAt":                   timestampField("When the valid preparation was created"),
	"LastUpdatedAt":               timestampField("When the valid preparation was last updated"),
	"ArchivedAt":                  timestampField("When the valid preparation was soft deleted"),
	"Name":                        stringField("Name of the preparation"),
	"Description":                 stringField("Description of the preparation"),
	"IconPath":                    stringField("The URL for the icon for the item"),
	"Slug":                        stringField("An easy-to-use URL slug for the preparation"),
	"PastTense":                   stringField("The past tense form of the preparation name (e.g., 'chopped' for 'chop')"),
	"InstrumentCount":             uint16RangeWithOptionalMaxSchema(),
	"IngredientCount":             uint16RangeWithOptionalMaxSchema(),
	"VesselCount":                 uint16RangeWithOptionalMaxSchema(),
	"RestrictToIngredients":       boolField("Whether or not the valid preparation is restricted to ingredients"),
	"TemperatureRequired":         boolField("Whether or not the valid preparation requires a temperature"),
	"TimeEstimateRequired":        boolField("Whether or not the valid preparation requires a time estimate"),
	"ConditionExpressionRequired": boolField("Whether or not the valid preparation requires a condition expression"),
	"ConsumesVessel":              boolField("Whether or not the valid preparation consumes a vessel"),
	"OnlyForVessels":              boolField("Whether or not the valid preparation is only for vessels"),
	"YieldsNothing":               boolField("Whether or not the valid preparation yields nothing"),
}

var getValidPreparationTool = &mcp.Tool{
	Name:        "GetValidPreparation",
	Description: "Get a valid preparation by it's ID",
	InputSchema: schemaObject(map[string]any{
		"ValidPreparationID": stringField("The ID of the valid preparation to get"),
	}),
	OutputSchema: schemaObject(validPreparationsSchema),
}

func (h *mcpToolManager) GetValidPreparation() mcp.ToolHandlerFor[*GetValidPreparationInvocation, *mealplanning.ValidPreparation] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPreparationInvocation) (*mcp.CallToolResult, *mealplanning.ValidPreparation, error) {
		result, err := h.client.GetValidPreparation(ctx, &mealplanninggrpc.GetValidPreparationRequest{
			ValidPreparationID: x.ValidPreparationID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationToValidPreparation(result.Result), nil
	}
}

type (
	SearchValidPreparationsInvocation struct {
		Filter           *filtering.QueryFilter
		Query            string `jsonschema_description:"The preparation name query"`
		UseSearchService bool   `jsonschema_description:"Whether or not to use a search index or just a database search"`
	}

	SearchValidPreparationsResult struct {
		Results []*mealplanning.ValidPreparation
	}
)

var searchForValidPreparationsTool = &mcp.Tool{
	Name:        "SearchForValidPreparations",
	Description: "Search for valid preparations with optional filtering and query string",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
		"Query": map[string]any{
			"type":        strType,
			"description": "The preparation name query",
		},
		"UseSearchService": map[string]any{
			"type":        boolType,
			"description": "Whether or not to use a search index or just a database search",
		},
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(validPreparationsSchema),
	}),
}

func (h *mcpToolManager) SearchForValidPreparations() mcp.ToolHandlerFor[*SearchValidPreparationsInvocation, *SearchValidPreparationsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidPreparationsInvocation) (*mcp.CallToolResult, *SearchValidPreparationsResult, error) {
		results, err := h.client.SearchForValidPreparations(ctx, &mealplanninggrpc.SearchForValidPreparationsRequest{
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchValidPreparationsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidPreparationToValidPreparation(result))
		}

		return nil, out, nil
	}
}

var validPreparationCreationTool = &mcp.Tool{
	Name:        "CreateValidPreparation",
	Description: "Create a valid preparation for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"Name":                        stringField("Name of the preparation"),
		"Description":                 stringField("Description of the preparation"),
		"IconPath":                    stringField("The URL for the icon for the item"),
		"Slug":                        stringField("An easy-to-use URL slug for the preparation"),
		"PastTense":                   stringField("The past tense form of the preparation name (e.g., 'chopped' for 'chop')"),
		"InstrumentCount":             uint16RangeWithOptionalMaxSchema(),
		"IngredientCount":             uint16RangeWithOptionalMaxSchema(),
		"VesselCount":                 uint16RangeWithOptionalMaxSchema(),
		"RestrictToIngredients":       boolField("Whether or not the valid preparation is restricted to ingredients"),
		"TemperatureRequired":         boolField("Whether or not the valid preparation requires a temperature"),
		"TimeEstimateRequired":        boolField("Whether or not the valid preparation requires a time estimate"),
		"ConditionExpressionRequired": boolField("Whether or not the valid preparation requires a condition expression"),
		"ConsumesVessel":              boolField("Whether or not the valid preparation consumes a vessel"),
		"OnlyForVessels":              boolField("Whether or not the valid preparation is only for vessels"),
		"YieldsNothing":               boolField("Whether or not the valid preparation yields nothing"),
	}),
	OutputSchema: schemaObject(validPreparationsSchema),
}

func (h *mcpToolManager) CreateValidPreparation() mcp.ToolHandlerFor[*mealplanning.ValidPreparationCreationRequestInput, *mealplanning.ValidPreparation] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidPreparationCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidPreparation, error) {
		result, err := h.client.CreateValidPreparation(ctx, &mealplanninggrpc.CreateValidPreparationRequest{Input: mealplanningconverters.ConvertValidPreparationCreationRequestInputToGRPCValidPreparationCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationToValidPreparation(result.Result), nil
	}
}

type (
	UpdateValidPreparationInvocation struct {
		*mealplanning.ValidPreparationUpdateRequestInput
		ValidPreparationID string `jsonschema:"required,description=The preparation ID"`
	}
)

var validPreparationUpdateTool = &mcp.Tool{
	Name:        "UpdateValidPreparation",
	Description: "Update a valid preparation for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"ValidPreparationID":          stringField("The ID of the valid preparation to update"),
		"Name":                        stringField("Name of the preparation"),
		"Description":                 stringField("Description of the preparation"),
		"IconPath":                    stringField("The URL for the icon for the item"),
		"Slug":                        stringField("An easy-to-use URL slug for the preparation"),
		"PastTense":                   stringField("The past tense form of the preparation name (e.g., 'chopped' for 'chop')"),
		"InstrumentCount":             uint16RangeWithOptionalMaxSchema(),
		"IngredientCount":             uint16RangeWithOptionalMaxSchema(),
		"VesselCount":                 uint16RangeWithOptionalMaxSchema(),
		"RestrictToIngredients":       boolField("Whether or not the valid preparation is restricted to ingredients"),
		"TemperatureRequired":         boolField("Whether or not the valid preparation requires a temperature"),
		"TimeEstimateRequired":        boolField("Whether or not the valid preparation requires a time estimate"),
		"ConditionExpressionRequired": boolField("Whether or not the valid preparation requires a condition expression"),
		"ConsumesVessel":              boolField("Whether or not the valid preparation consumes a vessel"),
		"OnlyForVessels":              boolField("Whether or not the valid preparation is only for vessels"),
		"YieldsNothing":               boolField("Whether or not the valid preparation yields nothing"),
	}),
	OutputSchema: schemaObject(validPreparationsSchema),
}

func (h *mcpToolManager) UpdateValidPreparation() mcp.ToolHandlerFor[*UpdateValidPreparationInvocation, *mealplanning.ValidPreparation] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidPreparationInvocation) (*mcp.CallToolResult, *mealplanning.ValidPreparation, error) {
		result, err := h.client.UpdateValidPreparation(ctx, &mealplanninggrpc.UpdateValidPreparationRequest{
			ValidPreparationID: x.ValidPreparationID,
			Input:              mealplanningconverters.ConvertValidPreparationUpdateRequestInputToGRPCValidPreparationUpdateRequestInput(x.ValidPreparationUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationToValidPreparation(result.Result), nil
	}
}

//
