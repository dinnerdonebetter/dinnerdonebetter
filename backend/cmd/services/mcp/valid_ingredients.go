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
	GetValidIngredientsInvocation struct {
		ValidIngredientID string `jsonschema:"description=The ingredient ID"`
	}
)

var getValidIngredientTool = &mcp.Tool{
	Name:         "GetValidIngredient",
	Description:  "Get a valid ingredient by it's ID",
	InputSchema:  schemaForType(&GetValidIngredientsInvocation{}),
	OutputSchema: schemaForType(&mealplanning.ValidIngredient{}),
}

func (h *mcpToolManager) GetValidIngredient() mcp.ToolHandlerFor[*GetValidIngredientsInvocation, *mealplanning.ValidIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *GetValidIngredientsInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredient, error) {
		result, err := h.client.GetValidIngredient(ctx, &mealplanninggrpc.GetValidIngredientRequest{
			ValidIngredientID: x.ValidIngredientID,
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result.Result), nil
	}
}

type (
	SearchValidIngredientsInvocation struct {
		Query            string `jsonschema_description:"The ingredient name query"`
		Filter           *filtering.QueryFilter
		UseSearchService bool `jsonschema_description:"Whether or not to use a search index or just a database search"`
	}

	SearchValidIngredientsResult struct {
		Results []*mealplanning.ValidIngredient
	}
)

var searchForValidIngredientsTool = &mcp.Tool{
	Name:         "SearchForValidIngredients",
	Description:  "Search for valid ingredients with optional filtering and query string",
	InputSchema:  schemaForType(SearchValidIngredientsInvocation{}),
	OutputSchema: schemaForType(SearchValidIngredientsResult{}),
}

func (h *mcpToolManager) SearchForValidIngredients() mcp.ToolHandlerFor[*SearchValidIngredientsInvocation, *SearchValidIngredientsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidIngredientsInvocation) (*mcp.CallToolResult, *SearchValidIngredientsResult, error) {
		results, err := h.client.SearchForValidIngredients(ctx, &mealplanninggrpc.SearchForValidIngredientsRequest{
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchValidIngredientsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result))
		}

		return nil, out, nil
	}
}

var validIngredientCreationTool = &mcp.Tool{
	Name:         "CreateValidIngredient",
	Description:  "Create a valid ingredient for use in recipes.",
	InputSchema:  schemaForType(&mealplanning.ValidIngredientCreationRequestInput{}),
	OutputSchema: schemaForType(&mealplanning.ValidIngredient{}),
}

func (h *mcpToolManager) CreateValidIngredient() mcp.ToolHandlerFor[*mealplanning.ValidIngredientCreationRequestInput, *mealplanning.ValidIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *mealplanning.ValidIngredientCreationRequestInput) (*mcp.CallToolResult, *mealplanning.ValidIngredient, error) {
		result, err := h.client.CreateValidIngredient(ctx, &mealplanninggrpc.CreateValidIngredientRequest{Input: mealplanningconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(x)})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result.Result), nil
	}
}

type (
	UpdateValidIngredientInvocation struct {
		ValidIngredientID string `jsonschema:"description=The ingredient ID"`
		Input             *mealplanning.ValidIngredientUpdateRequestInput
	}
)

var validIngredientUpdateTool = &mcp.Tool{
	Name:         "UpdateValidIngredient",
	Description:  "Update a valid ingredient for use in recipes.",
	InputSchema:  schemaForType(&UpdateValidIngredientInvocation{}),
	OutputSchema: schemaForType(&mealplanning.ValidIngredient{}),
}

func (h *mcpToolManager) UpdateValidIngredient() mcp.ToolHandlerFor[*UpdateValidIngredientInvocation, *mealplanning.ValidIngredient] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *UpdateValidIngredientInvocation) (*mcp.CallToolResult, *mealplanning.ValidIngredient, error) {
		result, err := h.client.UpdateValidIngredient(ctx, &mealplanninggrpc.UpdateValidIngredientRequest{
			ValidIngredientID: x.ValidIngredientID,
			Input:             mealplanningconverters.ConvertValidIngredientUpdateRequestInputToGRPCValidIngredientUpdateRequestInput(x.Input),
		})
		if err != nil {
			return nil, nil, err
		}

		return nil, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result.Result), nil
	}
}

//
