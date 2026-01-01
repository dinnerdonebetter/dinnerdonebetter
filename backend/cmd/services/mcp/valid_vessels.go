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
	GetValidVesselInvocation struct {
		ValidVesselID string `jsonschema:"description=The vessel ID"`
	}
)

var validVesselsSchema = map[string]any{
	"ID":                             stringField("The ID of the valid vessel"),
	"CreatedAt":                      timestampField("When the valid vessel was created"),
	"LastUpdatedAt":                  timestampField("When the valid vessel was last updated"),
	"ArchivedAt":                     timestampField("When the valid vessel was soft deleted"),
	"Name":                           stringField("Name of the vessel"),
	"Description":                    stringField("Description of the vessel"),
	"IconPath":                       stringField("The URL for the icon for the item"),
	"PluralName":                     stringField("The plural name for the vessel. So for a vessel named 'pan', this would be 'pans'"),
	"Slug":                           stringField("An easy-to-use URL slug for the vessel"),
	"Shape":                          stringField("The shape of the vessel (hemisphere, rectangle, cone, pyramid, cylinder, sphere, cube, or other)"),
	"WidthInMillimeters":             floatField("Width of the vessel in millimeters"),
	"LengthInMillimeters":            floatField("Length of the vessel in millimeters"),
	"HeightInMillimeters":            floatField("Height of the vessel in millimeters"),
	"Capacity":                       floatField("Capacity of the vessel"),
	"IncludeInGeneratedInstructions": boolField("Whether or not the valid vessel should be included in generated instructions"),
	"DisplayInSummaryLists":          boolField("Whether or not the valid vessel should be displayed in summary lists"),
	"UsableForStorage":               boolField("Whether or not the valid vessel is usable for storage"),
	"CapacityUnit":                   objectType(validMeasurementUnitsSchema),
}

var getValidVesselTool = &mcp.Tool{
	Name:        "GetValidVessel",
	Description: "Get a valid vessel by it's ID",
	InputSchema: schemaObject(map[string]any{
		"ValidVesselID": stringField("The ID of the valid vessel to get"),
	}),
	OutputSchema: schemaObject(validVesselsSchema),
}

func (h *mcpToolManager) GetValidVessel() mcp.ToolHandlerFor[*GetValidVesselInvocation, *mealplanning.ValidVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidVesselInvocation) (*mcp.CallToolResult, *mealplanning.ValidVessel, error) {
		result, err := h.client.GetValidVessel(ctx, &mealplanninggrpc.GetValidVesselRequest{
			ValidVesselId: x.ValidVesselID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidVesselToValidVessel(result.Result), nil
	}
}

type (
	SearchValidVesselsInvocation struct {
		Filter           *filtering.QueryFilter
		Query            string `jsonschema_description:"The vessel name query"`
		UseSearchService bool   `jsonschema_description:"Whether or not to use a search index or just a database search"`
	}

	SearchValidVesselsResult struct {
		Results []*mealplanning.ValidVessel
	}
)

var searchForValidVesselsTool = &mcp.Tool{
	Name:        "SearchForValidVessels",
	Description: "Search for valid vessels with optional filtering and query string",
	InputSchema: schemaObject(map[string]any{
		"Filter": queryFilterSchema(),
		"Query": map[string]any{
			"type":        strType,
			"description": "The vessel name query",
		},
		"UseSearchService": map[string]any{
			"type":        boolType,
			"description": "Whether or not to use a search index or just a database search",
		},
	}),
	OutputSchema: schemaObject(map[string]any{
		"Results": arrayType(schemaObject(validVesselsSchema)),
	}),
}

func (h *mcpToolManager) SearchForValidVessels() mcp.ToolHandlerFor[*SearchValidVesselsInvocation, *SearchValidVesselsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidVesselsInvocation) (*mcp.CallToolResult, *SearchValidVesselsResult, error) {
		results, err := h.client.SearchForValidVessels(ctx, &mealplanninggrpc.SearchForValidVesselsRequest{
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchValidVesselsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidVesselToValidVessel(result))
		}

		return nil, out, nil
	}
}

var validVesselCreationTool = &mcp.Tool{
	Name:        "CreateValidVessel",
	Description: "Create a valid vessel for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"Name":                           stringField("Name of the vessel"),
		"Description":                    stringField("Description of the vessel"),
		"IconPath":                       stringField("The URL for the icon for the item"),
		"PluralName":                     stringField("The plural name for the vessel. So for a vessel named 'pan', this would be 'pans'"),
		"Slug":                           stringField("An easy-to-use URL slug for the vessel"),
		"Shape":                          stringField("The shape of the vessel (hemisphere, rectangle, cone, pyramid, cylinder, sphere, cube, or other)"),
		"WidthInMillimeters":             floatField("Width of the vessel in millimeters"),
		"LengthInMillimeters":            floatField("Length of the vessel in millimeters"),
		"HeightInMillimeters":            floatField("Height of the vessel in millimeters"),
		"Capacity":                       floatField("Capacity of the vessel"),
		"CapacityUnitID":                 stringField("The ID of the valid measurement unit for capacity (optional)"),
		"IncludeInGeneratedInstructions": boolField("Whether or not the valid vessel should be included in generated instructions"),
		"DisplayInSummaryLists":          boolField("Whether or not the valid vessel should be displayed in summary lists"),
		"UsableForStorage":               boolField("Whether or not the valid vessel is usable for storage"),
	}),
	OutputSchema: schemaObject(validVesselsSchema),
}

func (h *mcpToolManager) CreateValidVessel() mcp.ToolHandlerFor[*mealplanning.ValidVesselCreationRequestInput, *mealplanning.ValidVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidVesselCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidVessel, error) {
		result, err := h.client.CreateValidVessel(ctx, &mealplanninggrpc.CreateValidVesselRequest{Input: mealplanningconverters.ConvertValidVesselCreationRequestInputToGRPCValidVesselCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidVesselToValidVessel(result.Result), nil
	}
}

type (
	UpdateValidVesselInvocation struct {
		*mealplanning.ValidVesselUpdateRequestInput
		ValidVesselID string `jsonschema:"required,description=The vessel ID"`
	}
)

var validVesselUpdateTool = &mcp.Tool{
	Name:        "UpdateValidVessel",
	Description: "Update a valid vessel for use in recipes.",
	InputSchema: schemaObject(map[string]any{
		"ValidVesselID":                  stringField("The ID of the valid vessel to update"),
		"Name":                           stringField("Name of the vessel"),
		"Description":                    stringField("Description of the vessel"),
		"IconPath":                       stringField("The URL for the icon for the item"),
		"PluralName":                     stringField("The plural name for the vessel. So for a vessel named 'pan', this would be 'pans'"),
		"Slug":                           stringField("An easy-to-use URL slug for the vessel"),
		"Shape":                          stringField("The shape of the vessel (hemisphere, rectangle, cone, pyramid, cylinder, sphere, cube, or other)"),
		"WidthInMillimeters":             floatField("Width of the vessel in millimeters"),
		"LengthInMillimeters":            floatField("Length of the vessel in millimeters"),
		"HeightInMillimeters":            floatField("Height of the vessel in millimeters"),
		"Capacity":                       floatField("Capacity of the vessel"),
		"CapacityUnitID":                 stringField("The ID of the valid measurement unit for capacity (optional)"),
		"IncludeInGeneratedInstructions": boolField("Whether or not the valid vessel should be included in generated instructions"),
		"DisplayInSummaryLists":          boolField("Whether or not the valid vessel should be displayed in summary lists"),
		"UsableForStorage":               boolField("Whether or not the valid vessel is usable for storage"),
	}),
	OutputSchema: schemaObject(validVesselsSchema),
}

func (h *mcpToolManager) UpdateValidVessel() mcp.ToolHandlerFor[*UpdateValidVesselInvocation, *mealplanning.ValidVessel] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidVesselInvocation) (*mcp.CallToolResult, *mealplanning.ValidVessel, error) {
		result, err := h.client.UpdateValidVessel(ctx, &mealplanninggrpc.UpdateValidVesselRequest{
			ValidVesselId: x.ValidVesselID,
			Input:         mealplanningconverters.ConvertValidVesselUpdateRequestInputToGRPCValidVesselUpdateRequestInput(x.ValidVesselUpdateRequestInput),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidVesselToValidVessel(result.Result), nil
	}
}

//
