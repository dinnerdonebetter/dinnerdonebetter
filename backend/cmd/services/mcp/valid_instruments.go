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
	GetValidInstrumentInvocation struct {
		ValidInstrumentID string `jsonschema:"description=The instrument MealPlanTaskID"`
	}
)

var validInstrumentsSchema = map[string]any{
	"MealPlanTaskID":                 stringField("The MealPlanTaskID of the valid instrument"),
	"CreatedAt":                      timestampField("When the valid instrument was created"),
	"LastUpdatedAt":                  timestampField("When the valid instrument was last updated"),
	"ArchivedAt":                     timestampField("When the valid instrument was soft deleted"),
	"Name":                           stringField("Name of the instrument"),
	"Description":                    stringField("Description of the instrument"),
	"IconPath":                       stringField("The URL for the icon for the item"),
	"PluralName":                     stringField("The plural name for the instrument. So for an instrument named 'knife', this would be 'knives'"),
	"Slug":                           stringField("An easy-to-use URL slug for the instrument"),
	"IncludeInGeneratedInstructions": boolField("Whether or not the valid instrument should be included in generated instructions"),
	"DisplayInSummaryLists":          boolField("Whether or not the valid instrument should be displayed in summary lists"),
	"UsableForStorage":               boolField("Whether or not the valid instrument is usable for storage"),
}

var getValidInstrumentTool = &mcp.Tool{
	Name:        "GetValidInstrument",
	Description: "Get a valid instrument by it's MealPlanTaskID",
	InputSchema: schemaObject(map[string]any{
		"ValidInstrumentID": stringField("The MealPlanTaskID of the valid instrument to get"),
	}),
	OutputSchema: schemaObject(validInstrumentsSchema),
}

func (h *mcpToolManager) GetValidInstrument() mcp.ToolHandlerFor[*GetValidInstrumentInvocation, *mealplanning.ValidInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidInstrumentInvocation) (*mcp.CallToolResult, *mealplanning.ValidInstrument, error) {
		result, err := h.client.GetValidInstrument(ctx, &mealplanninggrpc.GetValidInstrumentRequest{
			ValidInstrumentId: x.ValidInstrumentID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidInstrumentToValidInstrument(result.Result), nil
	}
}

type (
	SearchValidInstrumentsInvocation struct {
		Filter           *filtering.QueryFilter
		Query            string `jsonschema_description:"The instrument name query"`
		UseSearchService bool   `jsonschema_description:"Whether or not to use a search index or just a database search"`
	}

	SearchValidInstrumentsResult struct {
		Results []*mealplanning.ValidInstrument
	}
)

var searchForValidInstrumentsTool = &mcp.Tool{
	Name:        "SearchForValidInstruments",
	Description: "Search for valid instruments with optional filtering and query string",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
		"Query": map[string]any{
			"type":        strType,
			"description": "The instrument name query",
		},
		"UseSearchService": map[string]any{
			"type":        boolType,
			"description": "Whether or not to use a search index or just a database search",
		},
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validInstrumentsSchema)),
	}),
}

func (h *mcpToolManager) SearchForValidInstruments() mcp.ToolHandlerFor[*SearchValidInstrumentsInvocation, *SearchValidInstrumentsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidInstrumentsInvocation) (*mcp.CallToolResult, *SearchValidInstrumentsResult, error) {
		results, err := h.client.SearchForValidInstruments(ctx, &mealplanninggrpc.SearchForValidInstrumentsRequest{
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchValidInstrumentsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidInstrumentToValidInstrument(result))
		}

		return nil, out, nil
	}
}

var validInstrumentCreationTool = &mcp.Tool{
	Name:        "CreateValidInstrument",
	Description: "Create a valid instrument for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"Name":                           stringField("Name of the instrument"),
		"Description":                    stringField("Description of the instrument"),
		"IconPath":                       stringField("The URL for the icon for the item"),
		"PluralName":                     stringField("The plural name for the instrument. So for an instrument named 'knife', this would be 'knives'"),
		"Slug":                           stringField("An easy-to-use URL slug for the instrument"),
		"IncludeInGeneratedInstructions": boolField("Whether or not the valid instrument should be included in generated instructions"),
		"DisplayInSummaryLists":          boolField("Whether or not the valid instrument should be displayed in summary lists"),
		"UsableForStorage":               boolField("Whether or not the valid instrument is usable for storage"),
	}),
	OutputSchema: schemaObject(validInstrumentsSchema),
}

func (h *mcpToolManager) CreateValidInstrument() mcp.ToolHandlerFor[*mealplanning.ValidInstrumentCreationRequestInput, *mealplanning.ValidInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidInstrumentCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidInstrument, error) {
		result, err := h.client.CreateValidInstrument(ctx, &mealplanninggrpc.CreateValidInstrumentRequest{Input: mealplanningconverters.ConvertValidInstrumentCreationRequestInputToGRPCValidInstrumentCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidInstrumentToValidInstrument(result.Result), nil
	}
}

type (
	UpdateValidInstrumentInvocation struct {
		*mealplanning.ValidInstrumentUpdateRequestInput
		ValidInstrumentID string `jsonschema:"required,description=The instrument MealPlanTaskID"`
	}
)

var validInstrumentUpdateTool = &mcp.Tool{
	Name:        "UpdateValidInstrument",
	Description: "Update a valid instrument for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"ValidInstrumentID":              stringField("The MealPlanTaskID of the valid instrument to update"),
		"Name":                           stringField("Name of the instrument"),
		"Description":                    stringField("Description of the instrument"),
		"IconPath":                       stringField("The URL for the icon for the item"),
		"PluralName":                     stringField("The plural name for the instrument. So for an instrument named 'knife', this would be 'knives'"),
		"Slug":                           stringField("An easy-to-use URL slug for the instrument"),
		"IncludeInGeneratedInstructions": boolField("Whether or not the valid instrument should be included in generated instructions"),
		"DisplayInSummaryLists":          boolField("Whether or not the valid instrument should be displayed in summary lists"),
		"UsableForStorage":               boolField("Whether or not the valid instrument is usable for storage"),
	}),
	OutputSchema: schemaObject(validInstrumentsSchema),
}

func (h *mcpToolManager) UpdateValidInstrument() mcp.ToolHandlerFor[*UpdateValidInstrumentInvocation, *mealplanning.ValidInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidInstrumentInvocation) (*mcp.CallToolResult, *mealplanning.ValidInstrument, error) {
		result, err := h.client.UpdateValidInstrument(ctx, &mealplanninggrpc.UpdateValidInstrumentRequest{
			ValidInstrumentId: x.ValidInstrumentID,
			Input:             mealplanningconverters.ConvertValidInstrumentUpdateRequestInputToGRPCValidInstrumentUpdateRequestInput(x.ValidInstrumentUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidInstrumentToValidInstrument(result.Result), nil
	}
}

//
