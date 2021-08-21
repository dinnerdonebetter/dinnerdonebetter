package querybuilding

import (
	"fmt"
)

var (
	//
	// Households Table.
	//

	// HouseholdsUserMembershipTableColumns are the columns for the household user memberships table.
	HouseholdsUserMembershipTableColumns = []string{
		fmt.Sprintf("%s.%s", HouseholdsUserMembershipTableName, IDColumn),
		fmt.Sprintf("%s.%s", HouseholdsUserMembershipTableName, HouseholdsUserMembershipTableUserOwnershipColumn),
		fmt.Sprintf("%s.%s", HouseholdsUserMembershipTableName, HouseholdsUserMembershipTableHouseholdOwnershipColumn),
		fmt.Sprintf("%s.%s", HouseholdsUserMembershipTableName, HouseholdsUserMembershipTableHouseholdRolesColumn),
		fmt.Sprintf("%s.%s", HouseholdsUserMembershipTableName, HouseholdsUserMembershipTableDefaultUserHouseholdColumn),
		fmt.Sprintf("%s.%s", HouseholdsUserMembershipTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", HouseholdsUserMembershipTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", HouseholdsUserMembershipTableName, ArchivedOnColumn),
	}

	//
	// Households Table.
	//

	// HouseholdsTableColumns are the columns for the households table.
	HouseholdsTableColumns = []string{
		fmt.Sprintf("%s.%s", HouseholdsTableName, IDColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, HouseholdsTableNameColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, HouseholdsTableBillingStatusColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, HouseholdsTableContactEmailColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, HouseholdsTableContactPhoneColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, HouseholdsTablePaymentProcessorCustomerIDColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, HouseholdsTableSubscriptionPlanIDColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", HouseholdsTableName, HouseholdsTableUserOwnershipColumn),
	}

	//
	// Users Table.
	//

	// UsersTableColumns are the columns for the users table.
	UsersTableColumns = []string{
		fmt.Sprintf("%s.%s", UsersTableName, IDColumn),
		fmt.Sprintf("%s.%s", UsersTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableUsernameColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableAvatarColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableHashedPasswordColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableRequiresPasswordChangeColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTablePasswordLastChangedOnColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableTwoFactorSekretColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableTwoFactorVerifiedOnColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableServiceRolesColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableReputationColumn),
		fmt.Sprintf("%s.%s", UsersTableName, UsersTableStatusExplanationColumn),
		fmt.Sprintf("%s.%s", UsersTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", UsersTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", UsersTableName, ArchivedOnColumn),
	}

	//
	// Audit Log Entries Table.
	//

	// AuditLogEntriesTableColumns are the columns for the audit log entries table.
	AuditLogEntriesTableColumns = []string{
		fmt.Sprintf("%s.%s", AuditLogEntriesTableName, IDColumn),
		fmt.Sprintf("%s.%s", AuditLogEntriesTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", AuditLogEntriesTableName, AuditLogEntriesTableEventTypeColumn),
		fmt.Sprintf("%s.%s", AuditLogEntriesTableName, AuditLogEntriesTableContextColumn),
		fmt.Sprintf("%s.%s", AuditLogEntriesTableName, CreatedOnColumn),
	}

	//
	// API Clients Table.
	//

	// APIClientsTableColumns are the columns for the API clients table.
	APIClientsTableColumns = []string{
		fmt.Sprintf("%s.%s", APIClientsTableName, IDColumn),
		fmt.Sprintf("%s.%s", APIClientsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", APIClientsTableName, APIClientsTableNameColumn),
		fmt.Sprintf("%s.%s", APIClientsTableName, APIClientsTableClientIDColumn),
		fmt.Sprintf("%s.%s", APIClientsTableName, APIClientsTableSecretKeyColumn),
		fmt.Sprintf("%s.%s", APIClientsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", APIClientsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", APIClientsTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", APIClientsTableName, APIClientsTableOwnershipColumn),
	}

	//
	// Webhooks Table.
	//

	// WebhooksTableColumns are the columns for the webhooks table.
	WebhooksTableColumns = []string{
		fmt.Sprintf("%s.%s", WebhooksTableName, IDColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, WebhooksTableNameColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, WebhooksTableContentTypeColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, WebhooksTableURLColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, WebhooksTableMethodColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, WebhooksTableEventsColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, WebhooksTableDataTypesColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, WebhooksTableTopicsColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", WebhooksTableName, WebhooksTableOwnershipColumn),
	}

	//
	// ValidInstruments Table.
	//

	// ValidInstrumentsTableColumns are the columns for the valid instruments table.
	ValidInstrumentsTableColumns = []string{
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, IDColumn),
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, ValidInstrumentsTableNameColumn),
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, ValidInstrumentsTableVariantColumn),
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, ValidInstrumentsTableDescriptionColumn),
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, ValidInstrumentsTableIconPathColumn),
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", ValidInstrumentsTableName, ArchivedOnColumn),
	}

	//
	// ValidPreparations Table.
	//

	// ValidPreparationsTableColumns are the columns for the valid preparations table.
	ValidPreparationsTableColumns = []string{
		fmt.Sprintf("%s.%s", ValidPreparationsTableName, IDColumn),
		fmt.Sprintf("%s.%s", ValidPreparationsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", ValidPreparationsTableName, ValidPreparationsTableNameColumn),
		fmt.Sprintf("%s.%s", ValidPreparationsTableName, ValidPreparationsTableDescriptionColumn),
		fmt.Sprintf("%s.%s", ValidPreparationsTableName, ValidPreparationsTableIconPathColumn),
		fmt.Sprintf("%s.%s", ValidPreparationsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", ValidPreparationsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", ValidPreparationsTableName, ArchivedOnColumn),
	}

	//
	// ValidIngredients Table.
	//

	// ValidIngredientsTableColumns are the columns for the valid ingredients table.
	ValidIngredientsTableColumns = []string{
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, IDColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableNameColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableVariantColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableDescriptionColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableWarningColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsEggColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsDairyColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsPeanutColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsTreeNutColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsSoyColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsWheatColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsShellfishColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsSesameColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsFishColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableContainsGlutenColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableAnimalFleshColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableAnimalDerivedColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableVolumetricColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ValidIngredientsTableIconPathColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", ValidIngredientsTableName, ArchivedOnColumn),
	}

	//
	// ValidIngredientPreparations Table.
	//

	// ValidIngredientPreparationsTableColumns are the columns for the valid ingredient preparations table.
	ValidIngredientPreparationsTableColumns = []string{
		fmt.Sprintf("%s.%s", ValidIngredientPreparationsTableName, IDColumn),
		fmt.Sprintf("%s.%s", ValidIngredientPreparationsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", ValidIngredientPreparationsTableName, ValidIngredientPreparationsTableNotesColumn),
		fmt.Sprintf("%s.%s", ValidIngredientPreparationsTableName, ValidIngredientPreparationsTableValidIngredientIDColumn),
		fmt.Sprintf("%s.%s", ValidIngredientPreparationsTableName, ValidIngredientPreparationsTableValidPreparationIDColumn),
		fmt.Sprintf("%s.%s", ValidIngredientPreparationsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", ValidIngredientPreparationsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", ValidIngredientPreparationsTableName, ArchivedOnColumn),
	}

	//
	// ValidPreparationInstruments Table.
	//

	// ValidPreparationInstrumentsTableColumns are the columns for the valid preparation instruments table.
	ValidPreparationInstrumentsTableColumns = []string{
		fmt.Sprintf("%s.%s", ValidPreparationInstrumentsTableName, IDColumn),
		fmt.Sprintf("%s.%s", ValidPreparationInstrumentsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", ValidPreparationInstrumentsTableName, ValidPreparationInstrumentsTableInstrumentIDColumn),
		fmt.Sprintf("%s.%s", ValidPreparationInstrumentsTableName, ValidPreparationInstrumentsTablePreparationIDColumn),
		fmt.Sprintf("%s.%s", ValidPreparationInstrumentsTableName, ValidPreparationInstrumentsTableNotesColumn),
		fmt.Sprintf("%s.%s", ValidPreparationInstrumentsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", ValidPreparationInstrumentsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", ValidPreparationInstrumentsTableName, ArchivedOnColumn),
	}

	//
	// Recipes Table.
	//

	// RecipesTableColumns are the columns for the recipes table.
	RecipesTableColumns = []string{
		fmt.Sprintf("%s.%s", RecipesTableName, IDColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, RecipesTableNameColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, RecipesTableSourceColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, RecipesTableDescriptionColumn),
		// fmt.Sprintf("%s.%s", RecipesTableName, RecipesTableDisplayImageURLColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, RecipesTableInspiredByRecipeIDColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", RecipesTableName, RecipesTableHouseholdOwnershipColumn),
	}

	//
	// RecipeSteps Table.
	//

	// RecipeStepsTableColumns are the columns for the recipe steps table.
	RecipeStepsTableColumns = []string{
		fmt.Sprintf("%s.%s", RecipeStepsTableName, IDColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTableIndexColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTablePreparationIDColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTablePrerequisiteStepColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTableMinEstimatedTimeInSecondsColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTableMaxEstimatedTimeInSecondsColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTableTemperatureInCelsiusColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTableNotesColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTableWhyColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepsTableName, RecipeStepsTableBelongsToRecipeColumn),
	}

	//
	// RecipeStepIngredients Table.
	//

	// RecipeStepIngredientsTableColumns are the columns for the recipe step ingredients table.
	RecipeStepIngredientsTableColumns = []string{
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, IDColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, RecipeStepIngredientsTableIngredientIDColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, RecipeStepIngredientsTableNameColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityTypeColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityValueColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityNotesColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, RecipeStepIngredientsTableProductOfRecipeStepColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, RecipeStepIngredientsTableIngredientNotesColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepIngredientsTableName, RecipeStepIngredientsTableBelongsToRecipeStepColumn),
	}

	//
	// RecipeStepProducts Table.
	//

	// RecipeStepProductsTableColumns are the columns for the recipe step products table.
	RecipeStepProductsTableColumns = []string{
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, IDColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, RecipeStepProductsTableNameColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, RecipeStepProductsTableQuantityTypeColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, RecipeStepProductsTableQuantityValueColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, RecipeStepProductsTableQuantityNotesColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", RecipeStepProductsTableName, RecipeStepProductsTableBelongsToRecipeStepColumn),
	}

	//
	// Invitations Table.
	//

	// InvitationsTableColumns are the columns for the invitations table.
	InvitationsTableColumns = []string{
		fmt.Sprintf("%s.%s", InvitationsTableName, IDColumn),
		fmt.Sprintf("%s.%s", InvitationsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", InvitationsTableName, InvitationsTableCodeColumn),
		fmt.Sprintf("%s.%s", InvitationsTableName, InvitationsTableConsumedColumn),
		fmt.Sprintf("%s.%s", InvitationsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", InvitationsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", InvitationsTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", InvitationsTableName, InvitationsTableHouseholdOwnershipColumn),
	}

	//
	// Reports Table.
	//

	// ReportsTableColumns are the columns for the reports table.
	ReportsTableColumns = []string{
		fmt.Sprintf("%s.%s", ReportsTableName, IDColumn),
		fmt.Sprintf("%s.%s", ReportsTableName, ExternalIDColumn),
		fmt.Sprintf("%s.%s", ReportsTableName, ReportsTableReportTypeColumn),
		fmt.Sprintf("%s.%s", ReportsTableName, ReportsTableConcernColumn),
		fmt.Sprintf("%s.%s", ReportsTableName, CreatedOnColumn),
		fmt.Sprintf("%s.%s", ReportsTableName, LastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", ReportsTableName, ArchivedOnColumn),
		fmt.Sprintf("%s.%s", ReportsTableName, ReportsTableHouseholdOwnershipColumn),
	}
)
