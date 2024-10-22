package v2

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

type (
	CoreService interface {
		types.AdminDataService
		types.AuditLogEntryDataService
		types.AuthDataService
		types.CapitalismDataService
		types.HouseholdDataService
		types.HouseholdInvitationDataService
		types.OAuth2ClientDataService
		types.PasswordResetTokenDataService
		types.ServiceSettingDataService
		types.ServiceSettingConfigurationDataService
		types.UserDataService
		types.UserNotificationDataService
		types.WebhookDataService
	}

	MealPlanningService interface {
		types.HouseholdInstrumentOwnershipDataService
		types.MealDataService
		types.MealPlanDataService
		types.MealPlanEventDataService
		types.MealPlanGroceryListItemDataService
		types.MealPlanOptionDataService
		types.MealPlanOptionVoteDataService
		types.MealPlanTaskDataService
		types.RecipeDataService
		types.RecipeMediaDataService
		types.RecipePrepTaskDataService
		types.RecipeRatingDataService
		types.RecipeStepCompletionConditionDataService
		types.RecipeStepDataService
		types.RecipeStepIngredientDataService
		types.RecipeStepInstrumentDataService
		types.RecipeStepProductDataService
		types.RecipeStepVesselDataService
	}

	EnumerationsService interface {
		types.ValidIngredientDataService
		types.ValidIngredientGroupDataService
		types.ValidIngredientMeasurementUnitDataService
		types.ValidIngredientPreparationDataService
		types.ValidIngredientStateDataService
		types.ValidIngredientStateIngredientDataService
		types.ValidInstrumentDataService
		types.ValidMeasurementUnitConversionDataService
		types.ValidMeasurementUnitDataService
		types.ValidPreparationDataService
		types.ValidPreparationInstrumentDataService
		types.ValidPreparationVesselDataService
		types.ValidVesselDataService
	}
)
