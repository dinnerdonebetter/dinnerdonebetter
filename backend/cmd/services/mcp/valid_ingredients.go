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

type (
	CreateValidIngredientsInvocation struct {
		Input *mealplanning.ValidIngredientCreationRequestInput
	}

	CreateValidIngredientsResult struct {
		Output *mealplanning.ValidIngredient
	}
)

var validIngredientCreationTool = &mcp.Tool{
	Name:         "CreateValidIngredient",
	Description:  "Create a valid ingredient for use in recipes.",
	InputSchema:  schemaForType(CreateValidIngredientsInvocation{}),
	OutputSchema: schemaForType(CreateValidIngredientsResult{}),
}

func (h *mcpToolManager) CreateValidIngredient() mcp.ToolHandlerFor[*CreateValidIngredientsInvocation, *CreateValidIngredientsResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, x *CreateValidIngredientsInvocation) (*mcp.CallToolResult, *CreateValidIngredientsResult, error) {
		result, err := h.client.CreateValidIngredient(ctx, &mealplanninggrpc.CreateValidIngredientRequest{Input: mealplanningconverters.ConvertValidIngredientCreationRequestInputToGRPCValidIngredientCreationRequestInput(x.Input)})
		if err != nil {
			return nil, nil, err
		}

		out := &CreateValidIngredientsResult{
			Output: mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result.Result),
		}

		return nil, out, nil
	}
}
