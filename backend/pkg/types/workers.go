package types

import "net/http"

type (
	FinalizeMealPlansRequest struct {
		ReturnCount bool `json:"returnCount"`
	}

	FinalizeMealPlansResponse struct {
		Count int `json:"count"`
	}

	InitializeMealPlanGroceryListRequest struct {
		HouseholdID string `json:"householdID"`
	}

	InitializeMealPlanGroceryListResponse struct {
		Success bool `json:"success"`
	}

	CreateMealPlanTasksRequest struct {
		HouseholdID string `json:"householdID"`
	}

	CreateMealPlanTasksResponse struct {
		Success bool `json:"success"`
	}

	// WorkerService describes a structure capable of serving worker-oriented requests.
	WorkerService interface {
		MealPlanFinalizationHandler(http.ResponseWriter, *http.Request)
		MealPlanGroceryListInitializationHandler(res http.ResponseWriter, req *http.Request)
		MealPlanTaskCreationHandler(res http.ResponseWriter, req *http.Request)
	}
)
