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
		Query            string
		Filter           *filtering.QueryFilter
		UseSearchService bool
	}

	SearchValidIngredientsResult struct {
		Results []*mealplanning.ValidIngredient
	}
)

func (h *mcpToolManager) SearchForValidIngredients() (*mcp.Tool, mcp.ToolHandlerFor[*SearchValidIngredientsInvocation, *SearchValidIngredientsResult]) {
	tool := &mcp.Tool{
		Name:         "SearchForValidIngredients",
		Description:  "Search for valid ingredients with optional filtering and query string",
		InputSchema:  schemaForType(SearchValidIngredientsInvocation{}),
		OutputSchema: schemaForType(SearchValidIngredientsResult{}),
	}

	return tool, func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidIngredientsInvocation) (*mcp.CallToolResult, *SearchValidIngredientsResult, error) {
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
