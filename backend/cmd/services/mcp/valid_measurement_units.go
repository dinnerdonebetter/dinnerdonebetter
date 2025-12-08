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
	GetValidMeasurementUnitInvocation struct {
		ValidMeasurementUnitID string `jsonschema:"description=The measurement unit ID"`
	}
)

var validMeasurementUnitsSchema = map[string]any{
	"ID":            stringField("The ID of the valid measurement unit"),
	"CreatedAt":     timestampField("When the valid measurement unit was created"),
	"LastUpdatedAt": timestampField("When the valid measurement unit was last updated"),
	"ArchivedAt":    timestampField("When the valid measurement unit was soft deleted"),
	"Name":          stringField("Name of the measurement unit"),
	"Description":   stringField("Description of the measurement unit"),
	"IconPath":      stringField("The URL for the icon for the item"),
	"PluralName":    stringField("The plural name for the measurement unit. So for a unit named 'cup', this would be 'cups'"),
	"Slug":          stringField("An easy-to-use URL slug for the measurement unit"),
	"Volumetric":    boolField("Whether or not the valid measurement unit is volumetric"),
	"Universal":     boolField("Whether or not the valid measurement unit is universal (valid for all ingredients). For instance, 'grams' is a universal measurement unit"),
	"Metric":        boolField("Whether or not the valid measurement unit is metric"),
	"Imperial":      boolField("Whether or not the valid measurement unit is imperial"),
}

var getValidMeasurementUnitTool = &mcp.Tool{
	Name:        "GetValidMeasurementUnit",
	Description: "Get a valid measurement unit by it's ID",
	InputSchema: schemaObject(map[string]any{
		"ValidMeasurementUnitID": stringField("The ID of the valid measurement unit to get"),
	}),
	OutputSchema: schemaObject(validMeasurementUnitsSchema),
}

func (h *mcpToolManager) GetValidMeasurementUnit() mcp.ToolHandlerFor[*GetValidMeasurementUnitInvocation, *mealplanning.ValidMeasurementUnit] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidMeasurementUnitInvocation) (*mcp.CallToolResult, *mealplanning.ValidMeasurementUnit, error) {
		result, err := h.client.GetValidMeasurementUnit(ctx, &mealplanninggrpc.GetValidMeasurementUnitRequest{
			ValidMeasurementUnitID: x.ValidMeasurementUnitID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidMeasurementUnitToValidMeasurementUnit(result.Result), nil
	}
}

type (
	SearchValidMeasurementUnitsInvocation struct {
		Filter           *filtering.QueryFilter
		Query            string `jsonschema_description:"The measurement unit name query"`
		UseSearchService bool   `jsonschema_description:"Whether or not to use a search index or just a database search"`
	}

	SearchValidMeasurementUnitsResult struct {
		Results []*mealplanning.ValidMeasurementUnit
	}
)

var searchForValidMeasurementUnitsTool = &mcp.Tool{
	Name:        "SearchForValidMeasurementUnits",
	Description: "Search for valid measurement units with optional filtering and query string",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
		"Query": map[string]any{
			"type":        strType,
			"description": "The measurement unit name query",
		},
		"UseSearchService": map[string]any{
			"type":        boolType,
			"description": "Whether or not to use a search index or just a database search",
		},
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validMeasurementUnitsSchema)),
	}),
}

func (h *mcpToolManager) SearchForValidMeasurementUnits() mcp.ToolHandlerFor[*SearchValidMeasurementUnitsInvocation, *SearchValidMeasurementUnitsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidMeasurementUnitsInvocation) (*mcp.CallToolResult, *SearchValidMeasurementUnitsResult, error) {
		results, err := h.client.SearchForValidMeasurementUnits(ctx, &mealplanninggrpc.SearchForValidMeasurementUnitsRequest{
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchValidMeasurementUnitsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidMeasurementUnitToValidMeasurementUnit(result))
		}

		return nil, out, nil
	}
}

var validMeasurementUnitCreationTool = &mcp.Tool{
	Name:        "CreateValidMeasurementUnit",
	Description: "Create a valid measurement unit for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"Name":        stringField("Name of the measurement unit"),
		"Description": stringField("Description of the measurement unit"),
		"IconPath":    stringField("The URL for the icon for the item"),
		"PluralName":  stringField("The plural name for the measurement unit. So for a unit named 'cup', this would be 'cups'"),
		"Slug":        stringField("An easy-to-use URL slug for the measurement unit"),
		"Volumetric":  boolField("Whether or not the valid measurement unit is volumetric"),
		"Universal":   boolField("Whether or not the valid measurement unit is universal (valid for all ingredients). For instance, 'grams' is a universal measurement unit"),
		"Metric":      boolField("Whether or not the valid measurement unit is metric"),
		"Imperial":    boolField("Whether or not the valid measurement unit is imperial"),
	}),
	OutputSchema: schemaObject(validMeasurementUnitsSchema),
}

func (h *mcpToolManager) CreateValidMeasurementUnit() mcp.ToolHandlerFor[*mealplanning.ValidMeasurementUnitCreationRequestInput, *mealplanning.ValidMeasurementUnit] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidMeasurementUnitCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidMeasurementUnit, error) {
		result, err := h.client.CreateValidMeasurementUnit(ctx, &mealplanninggrpc.CreateValidMeasurementUnitRequest{Input: mealplanningconverters.ConvertValidMeasurementUnitCreationRequestInputToGRPCValidMeasurementUnitCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidMeasurementUnitToValidMeasurementUnit(result.Result), nil
	}
}

type (
	UpdateValidMeasurementUnitInvocation struct {
		*mealplanning.ValidMeasurementUnitUpdateRequestInput
		ValidMeasurementUnitID string `jsonschema:"required,description=The measurement unit ID"`
	}
)

var validMeasurementUnitUpdateTool = &mcp.Tool{
	Name:        "UpdateValidMeasurementUnit",
	Description: "Update a valid measurement unit for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"ValidMeasurementUnitID": stringField("The ID of the valid measurement unit to update"),
		"Name":                   stringField("Name of the measurement unit"),
		"Description":            stringField("Description of the measurement unit"),
		"IconPath":               stringField("The URL for the icon for the item"),
		"PluralName":             stringField("The plural name for the measurement unit. So for a unit named 'cup', this would be 'cups'"),
		"Slug":                   stringField("An easy-to-use URL slug for the measurement unit"),
		"Volumetric":             boolField("Whether or not the valid measurement unit is volumetric"),
		"Universal":              boolField("Whether or not the valid measurement unit is universal (valid for all ingredients). For instance, 'grams' is a universal measurement unit"),
		"Metric":                 boolField("Whether or not the valid measurement unit is metric"),
		"Imperial":               boolField("Whether or not the valid measurement unit is imperial"),
	}),
	OutputSchema: schemaObject(validMeasurementUnitsSchema),
}

func (h *mcpToolManager) UpdateValidMeasurementUnit() mcp.ToolHandlerFor[*UpdateValidMeasurementUnitInvocation, *mealplanning.ValidMeasurementUnit] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidMeasurementUnitInvocation) (*mcp.CallToolResult, *mealplanning.ValidMeasurementUnit, error) {
		result, err := h.client.UpdateValidMeasurementUnit(ctx, &mealplanninggrpc.UpdateValidMeasurementUnitRequest{
			ValidMeasurementUnitID: x.ValidMeasurementUnitID,
			Input:                  mealplanningconverters.ConvertValidMeasurementUnitUpdateRequestInputToGRPCValidMeasurementUnitUpdateRequestInput(x.ValidMeasurementUnitUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidMeasurementUnitToValidMeasurementUnit(result.Result), nil
	}
}

//
