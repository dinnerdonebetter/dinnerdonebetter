package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeFinalizeMealPlansRequest builds a faked Webhook.
func BuildFakeFinalizeMealPlansRequest() *types.FinalizeMealPlansRequest {
	return &types.FinalizeMealPlansRequest{
		ReturnCount: fake.Bool(),
	}
}

// BuildFakeFinalizeMealPlansResponse builds a faked FinalizeMealPlansResponse.
func BuildFakeFinalizeMealPlansResponse() *types.FinalizeMealPlansResponse {
	return &types.FinalizeMealPlansResponse{
		Count: int64(buildFakeNumber()),
	}
}

// BuildFakeCreateMealPlanTasksRequest builds a faked Webhook.
func BuildFakeCreateMealPlanTasksRequest() *types.CreateMealPlanTasksRequest {
	return &types.CreateMealPlanTasksRequest{
		HouseholdID: BuildFakeID(),
	}
}

// BuildFakeCreateMealPlanTasksResponse builds a faked Webhook.
func BuildFakeCreateMealPlanTasksResponse() *types.CreateMealPlanTasksResponse {
	return &types.CreateMealPlanTasksResponse{
		Success: true,
	}
}

// BuildFakeInitializeMealPlanGroceryListRequest builds a faked FinalizeMealPlansResponse.
func BuildFakeInitializeMealPlanGroceryListRequest() *types.InitializeMealPlanGroceryListRequest {
	return &types.InitializeMealPlanGroceryListRequest{
		HouseholdID: BuildFakeID(),
	}
}

// BuildFakeInitializeMealPlanGroceryListResponse builds a faked FinalizeMealPlansResponse.
func BuildFakeInitializeMealPlanGroceryListResponse() *types.InitializeMealPlanGroceryListResponse {
	return &types.InitializeMealPlanGroceryListResponse{
		Success: true,
	}
}
