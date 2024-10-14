# syntax=docker/dockerfile:1
FROM golang:1.22-bullseye

WORKDIR /go/src/github.com/dinnerdonebetter/backend
COPY . .

# TestIntegration/TestMealPlanGroceryListItems_CompleteLifecycle
# TestIntegration/TestRecipePrepTasks_Listing
# TestIntegration/TestRecipes_GetMealPlanTasksForRecipe
# TestIntegration/TestTOTPTokenValidation
# TestIntegration/TestTOTPTokenValidation
# TestIntegration/TestUsers_Searching
# TestIntegration/TestValidIngredientMeasurementUnits_CompleteLifecycle
# TestIntegration/TestValidIngredientPreparations_Listing_ByValues
# TestIntegration/TestValidIngredientPreparations_Listing_ByValues
# TestIntegration/TestValidMeasurementUnitConversions_GetFromUnits
# TestIntegration/TestValidMeasurementUnitConversions_GetToUnits

# to debug a specific test:
# ENTRYPOINT go test -parallel 1 -v -failfast github.com/dinnerdonebetter/backend/tests/integration -run TestIntegration/TestValidVessels_Searching

ENTRYPOINT go test -v github.com/dinnerdonebetter/backend/tests/integration
