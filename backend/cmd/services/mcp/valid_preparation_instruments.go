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
	GetValidPreparationInstrumentInvocation struct {
		ValidPreparationInstrumentID string `jsonschema:"description=The preparation instrument ID"`
	}
)

var validPreparationInstrumentsSchema = map[string]any{
	"ID":            stringField("The ID of the valid preparation instrument"),
	"CreatedAt":     timestampField("When the valid preparation instrument was created"),
	"LastUpdatedAt": timestampField("When the valid preparation instrument was last updated"),
	"ArchivedAt":    timestampField("When the valid preparation instrument was soft deleted"),
	"Notes":         stringField("Notes about the preparation instrument"),
	"Instrument":    objectType(validInstrumentsSchema),
	"Preparation":   objectType(validPreparationsSchema),
}

var getValidPreparationInstrumentTool = &mcp.Tool{
	Name:        "GetValidPreparationInstrument",
	Description: "Get a valid preparation instrument by it's ID",
	InputSchema: schemaObject(map[string]any{
		"ValidPreparationInstrumentID": stringField("The ID of the valid preparation instrument to get"),
	}),
	OutputSchema: schemaObject(validPreparationInstrumentsSchema),
}

func (h *mcpToolManager) GetValidPreparationInstrument() mcp.ToolHandlerFor[*GetValidPreparationInstrumentInvocation, *mealplanning.ValidPreparationInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPreparationInstrumentInvocation) (*mcp.CallToolResult, *mealplanning.ValidPreparationInstrument, error) {
		result, err := h.client.GetValidPreparationInstrument(ctx, &mealplanninggrpc.GetValidPreparationInstrumentRequest{
			ValidPreparationInstrumentId: x.ValidPreparationInstrumentID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationInstrumentToValidPreparationInstrument(result.Result), nil
	}
}

type (
	GetValidPreparationInstrumentsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetValidPreparationInstrumentsResult struct {
		Results []*mealplanning.ValidPreparationInstrument
	}
)

var getValidPreparationInstrumentsTool = &mcp.Tool{
	Name:        "GetValidPreparationInstruments",
	Description: "Get valid preparation instruments with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validPreparationInstrumentsSchema)),
	}),
}

func (h *mcpToolManager) GetValidPreparationInstruments() mcp.ToolHandlerFor[*GetValidPreparationInstrumentsInvocation, *GetValidPreparationInstrumentsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPreparationInstrumentsInvocation) (*mcp.CallToolResult, *GetValidPreparationInstrumentsResult, error) {
		results, err := h.client.GetValidPreparationInstruments(ctx, &mealplanninggrpc.GetValidPreparationInstrumentsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidPreparationInstrumentsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidPreparationInstrumentToValidPreparationInstrument(result))
		}

		return nil, out, nil
	}
}

var validPreparationInstrumentCreationTool = &mcp.Tool{
	Name:        "CreateValidPreparationInstrument",
	Description: "Create a valid preparation instrument linking a preparation to an instrument.",
	InputSchema: schemaObject(map[string]any{
		"Notes":              stringField("Notes about the preparation instrument"),
		"ValidPreparationID": stringField("The ID of the valid preparation"),
		"ValidInstrumentID":  stringField("The ID of the valid instrument"),
	}),
	OutputSchema: schemaObject(validPreparationInstrumentsSchema),
}

func (h *mcpToolManager) CreateValidPreparationInstrument() mcp.ToolHandlerFor[*mealplanning.ValidPreparationInstrumentCreationRequestInput, *mealplanning.ValidPreparationInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidPreparationInstrumentCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidPreparationInstrument, error) {
		result, err := h.client.CreateValidPreparationInstrument(ctx, &mealplanninggrpc.CreateValidPreparationInstrumentRequest{Input: mealplanningconverters.ConvertCreateValidPreparationInstrumentRequestToGRPCValidPreparationInstrumentCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationInstrumentToValidPreparationInstrument(result.Result), nil
	}
}

type (
	UpdateValidPreparationInstrumentInvocation struct {
		*mealplanning.ValidPreparationInstrumentUpdateRequestInput
		ValidPreparationInstrumentID string `jsonschema:"required,description=The preparation instrument ID"`
	}
)

var validPreparationInstrumentUpdateTool = &mcp.Tool{
	Name:        "UpdateValidPreparationInstrument",
	Description: "Update a valid preparation instrument.",
	InputSchema: schemaObject(map[string]any{
		"ValidPreparationInstrumentID": stringField("The ID of the valid preparation instrument to update"),
		"Notes":                        stringField("Notes about the preparation instrument"),
		"ValidPreparationID":           stringField("The ID of the valid preparation"),
		"ValidInstrumentID":            stringField("The ID of the valid instrument"),
	}),
	OutputSchema: schemaObject(validPreparationInstrumentsSchema),
}

func (h *mcpToolManager) UpdateValidPreparationInstrument() mcp.ToolHandlerFor[*UpdateValidPreparationInstrumentInvocation, *mealplanning.ValidPreparationInstrument] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidPreparationInstrumentInvocation) (*mcp.CallToolResult, *mealplanning.ValidPreparationInstrument, error) {
		result, err := h.client.UpdateValidPreparationInstrument(ctx, &mealplanninggrpc.UpdateValidPreparationInstrumentRequest{
			ValidPreparationInstrumentId: x.ValidPreparationInstrumentID,
			Input:                        mealplanningconverters.ConvertValidPreparationInstrumentUpdateRequestInputToGRPCValidPreparationInstrumentUpdateRequestInput(x.ValidPreparationInstrumentUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationInstrumentToValidPreparationInstrument(result.Result), nil
	}
}

//
