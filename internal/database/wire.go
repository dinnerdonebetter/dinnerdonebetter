package database

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/google/wire"
)

var (
	// Providers represents what we provide to dependency injectors.
	Providers = wire.NewSet(
		ProvideAdminAuditManager,
		ProvideAuthAuditManager,
		ProvideAuditLogEntryDataManager,
		ProvideValidInstrumentDataManager,
		ProvideValidPreparationDataManager,
		ProvideValidIngredientDataManager,
		ProvideValidIngredientPreparationDataManager,
		ProvideValidPreparationInstrumentDataManager,
		ProvideRecipeDataManager,
		ProvideRecipeStepDataManager,
		ProvideRecipeStepIngredientDataManager,
		ProvideRecipeStepProductDataManager,
		ProvideInvitationDataManager,
		ProvideReportDataManager,
		ProvideUserDataManager,
		ProvideAdminUserDataManager,
		ProvideAccountDataManager,
		ProvideAccountUserMembershipDataManager,
		ProvideAPIClientDataManager,
		ProvideWebhookDataManager,
	)
)

// ProvideAdminAuditManager is an arbitrary function for dependency injection's sake.
func ProvideAdminAuditManager(db DataManager) types.AdminAuditManager {
	return db
}

// ProvideAuthAuditManager is an arbitrary function for dependency injection's sake.
func ProvideAuthAuditManager(db DataManager) types.AuthAuditManager {
	return db
}

// ProvideAuditLogEntryDataManager is an arbitrary function for dependency injection's sake.
func ProvideAuditLogEntryDataManager(db DataManager) types.AuditLogEntryDataManager {
	return db
}

// ProvideAccountDataManager is an arbitrary function for dependency injection's sake.
func ProvideAccountDataManager(db DataManager) types.AccountDataManager {
	return db
}

// ProvideAccountUserMembershipDataManager is an arbitrary function for dependency injection's sake.
func ProvideAccountUserMembershipDataManager(db DataManager) types.AccountUserMembershipDataManager {
	return db
}

// ProvideValidInstrumentDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidInstrumentDataManager(db DataManager) types.ValidInstrumentDataManager {
	return db
}

// ProvideValidPreparationDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidPreparationDataManager(db DataManager) types.ValidPreparationDataManager {
	return db
}

// ProvideValidIngredientDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientDataManager(db DataManager) types.ValidIngredientDataManager {
	return db
}

// ProvideValidIngredientPreparationDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientPreparationDataManager(db DataManager) types.ValidIngredientPreparationDataManager {
	return db
}

// ProvideValidPreparationInstrumentDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidPreparationInstrumentDataManager(db DataManager) types.ValidPreparationInstrumentDataManager {
	return db
}

// ProvideRecipeDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeDataManager(db DataManager) types.RecipeDataManager {
	return db
}

// ProvideRecipeStepDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepDataManager(db DataManager) types.RecipeStepDataManager {
	return db
}

// ProvideRecipeStepIngredientDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepIngredientDataManager(db DataManager) types.RecipeStepIngredientDataManager {
	return db
}

// ProvideRecipeStepProductDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepProductDataManager(db DataManager) types.RecipeStepProductDataManager {
	return db
}

// ProvideInvitationDataManager is an arbitrary function for dependency injection's sake.
func ProvideInvitationDataManager(db DataManager) types.InvitationDataManager {
	return db
}

// ProvideReportDataManager is an arbitrary function for dependency injection's sake.
func ProvideReportDataManager(db DataManager) types.ReportDataManager {
	return db
}

// ProvideUserDataManager is an arbitrary function for dependency injection's sake.
func ProvideUserDataManager(db DataManager) types.UserDataManager {
	return db
}

// ProvideAdminUserDataManager is an arbitrary function for dependency injection's sake.
func ProvideAdminUserDataManager(db DataManager) types.AdminUserDataManager {
	return db
}

// ProvideAPIClientDataManager is an arbitrary function for dependency injection's sake.
func ProvideAPIClientDataManager(db DataManager) types.APIClientDataManager {
	return db
}

// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake.
func ProvideWebhookDataManager(db DataManager) types.WebhookDataManager {
	return db
}
