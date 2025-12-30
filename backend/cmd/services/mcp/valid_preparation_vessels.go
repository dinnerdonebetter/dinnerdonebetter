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
	GetValidPreparationVesselInvocation struct {
		ValidPreparationVesselID string `jsonschema:"description=The preparation vessel MealPlanTaskID"`
	}
)

var validPreparationVesselsSchema = map[string]any{
	"MealPlanTaskID": stringField("The MealPlanTaskID of the valid preparation vessel"),
	"CreatedAt":      timestampField("When the valid preparation vessel was created"),
	"LastUpdatedAt":  timestampField("When the valid preparation vessel was last updated"),
	"ArchivedAt":     timestampField("When the valid preparation vessel was soft deleted"),
	"Notes":          stringField("Notes about the preparation vessel"),
	"Vessel":         objectType(validVesselsSchema),
	"Preparation":    objectType(validPreparationsSchema),
}

var getValidPreparationVesselTool = &mcp.Tool{
	Name:        "GetValidPreparationVessel",
	Description: "Get a valid preparation vessel by it's MealPlanTaskID",
	InputSchema: schemaObject(map[string]any{
		"ValidPreparationVesselID": stringField("The MealPlanTaskID of the valid preparation vessel to get"),
	}),
	OutputSchema: schemaObject(validPreparationVesselsSchema),
}

func (h *mcpToolManager) GetValidPreparationVessel() mcp.ToolHandlerFor[*GetValidPreparationVesselInvocation, *mealplanning.ValidPreparationVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPreparationVesselInvocation) (*mcp.CallToolResult, *mealplanning.ValidPreparationVessel, error) {
		result, err := h.client.GetValidPreparationVessel(ctx, &mealplanninggrpc.GetValidPreparationVesselRequest{
			ValidPreparationVesselId: x.ValidPreparationVesselID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationVesselToValidPreparationVessel(result.Result), nil
	}
}

type (
	GetValidPreparationVesselsInvocation struct {
		Filter *filtering.QueryFilter
	}

	GetValidPreparationVesselsResult struct {
		Results []*mealplanning.ValidPreparationVessel
	}
)

var getValidPreparationVesselsTool = &mcp.Tool{
	Name:        "GetValidPreparationVessels",
	Description: "Get valid preparation vessels with optional filtering",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validPreparationVesselsSchema)),
	}),
}

func (h *mcpToolManager) GetValidPreparationVessels() mcp.ToolHandlerFor[*GetValidPreparationVesselsInvocation, *GetValidPreparationVesselsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidPreparationVesselsInvocation) (*mcp.CallToolResult, *GetValidPreparationVesselsResult, error) {
		results, err := h.client.GetValidPreparationVessels(ctx, &mealplanninggrpc.GetValidPreparationVesselsRequest{
			Filter: grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
		})
		if err != nil {
			return nil, nil, err
		}

		out := &GetValidPreparationVesselsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidPreparationVesselToValidPreparationVessel(result))
		}

		return nil, out, nil
	}
}

var validPreparationVesselCreationTool = &mcp.Tool{
	Name:        "CreateValidPreparationVessel",
	Description: "Create a valid preparation vessel linking a preparation to a vessel.",
	InputSchema: schemaObject(map[string]any{
		"Notes":              stringField("Notes about the preparation vessel"),
		"ValidPreparationID": stringField("The MealPlanTaskID of the valid preparation"),
		"ValidVesselID":      stringField("The MealPlanTaskID of the valid vessel"),
	}),
	OutputSchema: schemaObject(validPreparationVesselsSchema),
}

func (h *mcpToolManager) CreateValidPreparationVessel() mcp.ToolHandlerFor[*mealplanning.ValidPreparationVesselCreationRequestInput, *mealplanning.ValidPreparationVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidPreparationVesselCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidPreparationVessel, error) {
		result, err := h.client.CreateValidPreparationVessel(ctx, &mealplanninggrpc.CreateValidPreparationVesselRequest{Input: mealplanningconverters.ConvertCreateValidPreparationVesselRequestToGRPCValidPreparationVesselCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationVesselToValidPreparationVessel(result.Result), nil
	}
}

type (
	UpdateValidPreparationVesselInvocation struct {
		*mealplanning.ValidPreparationVesselUpdateRequestInput
		ValidPreparationVesselID string `jsonschema:"required,description=The preparation vessel MealPlanTaskID"`
	}
)

var validPreparationVesselUpdateTool = &mcp.Tool{
	Name:        "UpdateValidPreparationVessel",
	Description: "Update a valid preparation vessel.",
	InputSchema: schemaObject(map[string]any{
		"ValidPreparationVesselID": stringField("The MealPlanTaskID of the valid preparation vessel to update"),
		"Notes":                    stringField("Notes about the preparation vessel"),
		"ValidPreparationID":       stringField("The MealPlanTaskID of the valid preparation"),
		"ValidVesselID":            stringField("The MealPlanTaskID of the valid vessel"),
	}),
	OutputSchema: schemaObject(validPreparationVesselsSchema),
}

func (h *mcpToolManager) UpdateValidPreparationVessel() mcp.ToolHandlerFor[*UpdateValidPreparationVesselInvocation, *mealplanning.ValidPreparationVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidPreparationVesselInvocation) (*mcp.CallToolResult, *mealplanning.ValidPreparationVessel, error) {
		result, err := h.client.UpdateValidPreparationVessel(ctx, &mealplanninggrpc.UpdateValidPreparationVesselRequest{
			ValidPreparationVesselId: x.ValidPreparationVesselID,
			Input:                    mealplanningconverters.ConvertValidPreparationVesselUpdateRequestInputToGRPCValidPreparationVesselUpdateRequestInput(x.ValidPreparationVesselUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidPreparationVesselToValidPreparationVessel(result.Result), nil
	}
}

//
