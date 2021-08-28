package querybuilding

import (
	"fmt"
)

func buildColumnReference(tableName, columnName string) string {
	return fmt.Sprintf("%s.%s", tableName, columnName)
}

var (
	//
	// Households Table.
	//

	// HouseholdsUserMembershipTableColumns are the columns for the household user memberships table.
	HouseholdsUserMembershipTableColumns = []string{
		buildColumnReference(HouseholdsUserMembershipTableName, IDColumn),
		buildColumnReference(HouseholdsUserMembershipTableName, HouseholdsUserMembershipTableUserOwnershipColumn),
		buildColumnReference(HouseholdsUserMembershipTableName, HouseholdsUserMembershipTableHouseholdOwnershipColumn),
		buildColumnReference(HouseholdsUserMembershipTableName, HouseholdsUserMembershipTableHouseholdRolesColumn),
		buildColumnReference(HouseholdsUserMembershipTableName, HouseholdsUserMembershipTableDefaultUserHouseholdColumn),
		buildColumnReference(HouseholdsUserMembershipTableName, CreatedOnColumn),
		buildColumnReference(HouseholdsUserMembershipTableName, LastUpdatedOnColumn),
		buildColumnReference(HouseholdsUserMembershipTableName, ArchivedOnColumn),
	}

	//
	// Households Table.
	//

	// HouseholdsTableColumns are the columns for the households table.
	HouseholdsTableColumns = []string{
		buildColumnReference(HouseholdsTableName, IDColumn),
		buildColumnReference(HouseholdsTableName, ExternalIDColumn),
		buildColumnReference(HouseholdsTableName, HouseholdsTableNameColumn),
		buildColumnReference(HouseholdsTableName, HouseholdsTableBillingStatusColumn),
		buildColumnReference(HouseholdsTableName, HouseholdsTableContactEmailColumn),
		buildColumnReference(HouseholdsTableName, HouseholdsTableContactPhoneColumn),
		buildColumnReference(HouseholdsTableName, HouseholdsTablePaymentProcessorCustomerIDColumn),
		buildColumnReference(HouseholdsTableName, HouseholdsTableSubscriptionPlanIDColumn),
		buildColumnReference(HouseholdsTableName, CreatedOnColumn),
		buildColumnReference(HouseholdsTableName, LastUpdatedOnColumn),
		buildColumnReference(HouseholdsTableName, ArchivedOnColumn),
		buildColumnReference(HouseholdsTableName, HouseholdsTableUserOwnershipColumn),
	}

	//
	// Users Table.
	//

	// UsersTableColumns are the columns for the users table.
	UsersTableColumns = []string{
		buildColumnReference(UsersTableName, IDColumn),
		buildColumnReference(UsersTableName, ExternalIDColumn),
		buildColumnReference(UsersTableName, UsersTableUsernameColumn),
		buildColumnReference(UsersTableName, UsersTableAvatarColumn),
		buildColumnReference(UsersTableName, UsersTableHashedPasswordColumn),
		buildColumnReference(UsersTableName, UsersTableRequiresPasswordChangeColumn),
		buildColumnReference(UsersTableName, UsersTablePasswordLastChangedOnColumn),
		buildColumnReference(UsersTableName, UsersTableTwoFactorSekretColumn),
		buildColumnReference(UsersTableName, UsersTableTwoFactorVerifiedOnColumn),
		buildColumnReference(UsersTableName, UsersTableServiceRolesColumn),
		buildColumnReference(UsersTableName, UsersTableReputationColumn),
		buildColumnReference(UsersTableName, UsersTableStatusExplanationColumn),
		buildColumnReference(UsersTableName, CreatedOnColumn),
		buildColumnReference(UsersTableName, LastUpdatedOnColumn),
		buildColumnReference(UsersTableName, ArchivedOnColumn),
	}

	//
	// Audit Log Entries Table.
	//

	// AuditLogEntriesTableColumns are the columns for the audit log entries table.
	AuditLogEntriesTableColumns = []string{
		buildColumnReference(AuditLogEntriesTableName, IDColumn),
		buildColumnReference(AuditLogEntriesTableName, ExternalIDColumn),
		buildColumnReference(AuditLogEntriesTableName, AuditLogEntriesTableEventTypeColumn),
		buildColumnReference(AuditLogEntriesTableName, AuditLogEntriesTableContextColumn),
		buildColumnReference(AuditLogEntriesTableName, CreatedOnColumn),
	}

	//
	// API Clients Table.
	//

	// APIClientsTableColumns are the columns for the API clients table.
	APIClientsTableColumns = []string{
		buildColumnReference(APIClientsTableName, IDColumn),
		buildColumnReference(APIClientsTableName, ExternalIDColumn),
		buildColumnReference(APIClientsTableName, APIClientsTableNameColumn),
		buildColumnReference(APIClientsTableName, APIClientsTableClientIDColumn),
		buildColumnReference(APIClientsTableName, APIClientsTableSecretKeyColumn),
		buildColumnReference(APIClientsTableName, CreatedOnColumn),
		buildColumnReference(APIClientsTableName, LastUpdatedOnColumn),
		buildColumnReference(APIClientsTableName, ArchivedOnColumn),
		buildColumnReference(APIClientsTableName, APIClientsTableOwnershipColumn),
	}

	//
	// Webhooks Table.
	//

	// WebhooksTableColumns are the columns for the webhooks table.
	WebhooksTableColumns = []string{
		buildColumnReference(WebhooksTableName, IDColumn),
		buildColumnReference(WebhooksTableName, ExternalIDColumn),
		buildColumnReference(WebhooksTableName, WebhooksTableNameColumn),
		buildColumnReference(WebhooksTableName, WebhooksTableContentTypeColumn),
		buildColumnReference(WebhooksTableName, WebhooksTableURLColumn),
		buildColumnReference(WebhooksTableName, WebhooksTableMethodColumn),
		buildColumnReference(WebhooksTableName, WebhooksTableEventsColumn),
		buildColumnReference(WebhooksTableName, WebhooksTableDataTypesColumn),
		buildColumnReference(WebhooksTableName, WebhooksTableTopicsColumn),
		buildColumnReference(WebhooksTableName, CreatedOnColumn),
		buildColumnReference(WebhooksTableName, LastUpdatedOnColumn),
		buildColumnReference(WebhooksTableName, ArchivedOnColumn),
		buildColumnReference(WebhooksTableName, WebhooksTableOwnershipColumn),
	}

	//
	// ValidInstruments Table.
	//

	// ValidInstrumentsTableColumns are the columns for the valid instruments table.
	ValidInstrumentsTableColumns = []string{
		buildColumnReference(ValidInstrumentsTableName, IDColumn),
		buildColumnReference(ValidInstrumentsTableName, ExternalIDColumn),
		buildColumnReference(ValidInstrumentsTableName, ValidInstrumentsTableNameColumn),
		buildColumnReference(ValidInstrumentsTableName, ValidInstrumentsTableVariantColumn),
		buildColumnReference(ValidInstrumentsTableName, ValidInstrumentsTableDescriptionColumn),
		buildColumnReference(ValidInstrumentsTableName, ValidInstrumentsTableIconPathColumn),
		buildColumnReference(ValidInstrumentsTableName, CreatedOnColumn),
		buildColumnReference(ValidInstrumentsTableName, LastUpdatedOnColumn),
		buildColumnReference(ValidInstrumentsTableName, ArchivedOnColumn),
	}

	//
	// ValidPreparations Table.
	//

	// ValidPreparationsTableColumns are the columns for the valid preparations table.
	ValidPreparationsTableColumns = []string{
		buildColumnReference(ValidPreparationsTableName, IDColumn),
		buildColumnReference(ValidPreparationsTableName, ExternalIDColumn),
		buildColumnReference(ValidPreparationsTableName, ValidPreparationsTableNameColumn),
		buildColumnReference(ValidPreparationsTableName, ValidPreparationsTableDescriptionColumn),
		buildColumnReference(ValidPreparationsTableName, ValidPreparationsTableIconPathColumn),
		buildColumnReference(ValidPreparationsTableName, CreatedOnColumn),
		buildColumnReference(ValidPreparationsTableName, LastUpdatedOnColumn),
		buildColumnReference(ValidPreparationsTableName, ArchivedOnColumn),
	}

	//
	// ValidIngredients Table.
	//

	// ValidIngredientsTableColumns are the columns for the valid ingredients table.
	ValidIngredientsTableColumns = []string{
		buildColumnReference(ValidIngredientsTableName, IDColumn),
		buildColumnReference(ValidIngredientsTableName, ExternalIDColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableNameColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableVariantColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableDescriptionColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableWarningColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsEggColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsDairyColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsPeanutColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsTreeNutColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsSoyColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsWheatColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsShellfishColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsSesameColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsFishColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsGlutenColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableAnimalFleshColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableAnimalDerivedColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableVolumetricColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableIconPathColumn),
		buildColumnReference(ValidIngredientsTableName, CreatedOnColumn),
		buildColumnReference(ValidIngredientsTableName, LastUpdatedOnColumn),
		buildColumnReference(ValidIngredientsTableName, ArchivedOnColumn),
	}

	//
	// ValidIngredientPreparations Table.
	//

	// ValidIngredientPreparationsTableColumns are the columns for the valid ingredient preparations table.
	ValidIngredientPreparationsTableColumns = []string{
		buildColumnReference(ValidIngredientPreparationsTableName, IDColumn),
		buildColumnReference(ValidIngredientPreparationsTableName, ExternalIDColumn),
		buildColumnReference(ValidIngredientPreparationsTableName, ValidIngredientPreparationsTableNotesColumn),
		buildColumnReference(ValidIngredientPreparationsTableName, ValidIngredientPreparationsTableValidIngredientIDColumn),
		buildColumnReference(ValidIngredientPreparationsTableName, ValidIngredientPreparationsTableValidPreparationIDColumn),
		buildColumnReference(ValidIngredientPreparationsTableName, CreatedOnColumn),
		buildColumnReference(ValidIngredientPreparationsTableName, LastUpdatedOnColumn),
		buildColumnReference(ValidIngredientPreparationsTableName, ArchivedOnColumn),
	}

	//
	// ValidPreparationInstruments Table.
	//

	// ValidPreparationInstrumentsTableColumns are the columns for the valid preparation instruments table.
	ValidPreparationInstrumentsTableColumns = []string{
		buildColumnReference(ValidPreparationInstrumentsTableName, IDColumn),
		buildColumnReference(ValidPreparationInstrumentsTableName, ExternalIDColumn),
		buildColumnReference(ValidPreparationInstrumentsTableName, ValidPreparationInstrumentsTableInstrumentIDColumn),
		buildColumnReference(ValidPreparationInstrumentsTableName, ValidPreparationInstrumentsTablePreparationIDColumn),
		buildColumnReference(ValidPreparationInstrumentsTableName, ValidPreparationInstrumentsTableNotesColumn),
		buildColumnReference(ValidPreparationInstrumentsTableName, CreatedOnColumn),
		buildColumnReference(ValidPreparationInstrumentsTableName, LastUpdatedOnColumn),
		buildColumnReference(ValidPreparationInstrumentsTableName, ArchivedOnColumn),
	}

	//
	// Recipes Table.
	//

	// RecipesTableColumns are the columns for the recipes table.
	RecipesTableColumns = []string{
		buildColumnReference(RecipesTableName, IDColumn),
		buildColumnReference(RecipesTableName, ExternalIDColumn),
		buildColumnReference(RecipesTableName, RecipesTableNameColumn),
		buildColumnReference(RecipesTableName, RecipesTableSourceColumn),
		buildColumnReference(RecipesTableName, RecipesTableDescriptionColumn),
		buildColumnReference(RecipesTableName, RecipesTableInspiredByRecipeIDColumn),
		buildColumnReference(RecipesTableName, CreatedOnColumn),
		buildColumnReference(RecipesTableName, LastUpdatedOnColumn),
		buildColumnReference(RecipesTableName, ArchivedOnColumn),
		buildColumnReference(RecipesTableName, RecipesTableHouseholdOwnershipColumn),
	}

	// FullRecipeColumns are the columns for the recipes table.
	FullRecipeColumns = []string{
		// recipe parts
		buildColumnReference(RecipesTableName, IDColumn),
		buildColumnReference(RecipesTableName, ExternalIDColumn),
		buildColumnReference(RecipesTableName, RecipesTableNameColumn),
		buildColumnReference(RecipesTableName, RecipesTableSourceColumn),
		buildColumnReference(RecipesTableName, RecipesTableDescriptionColumn),
		buildColumnReference(RecipesTableName, RecipesTableInspiredByRecipeIDColumn),
		buildColumnReference(RecipesTableName, CreatedOnColumn),
		buildColumnReference(RecipesTableName, LastUpdatedOnColumn),
		buildColumnReference(RecipesTableName, ArchivedOnColumn),
		buildColumnReference(RecipesTableName, RecipesTableHouseholdOwnershipColumn),
		// recipe step parts
		buildColumnReference(RecipeStepsTableName, IDColumn),
		buildColumnReference(RecipeStepsTableName, ExternalIDColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableIndexColumn),
		// recipe step preparation parts
		buildColumnReference(ValidPreparationsTableName, IDColumn),
		buildColumnReference(ValidPreparationsTableName, ExternalIDColumn),
		buildColumnReference(ValidPreparationsTableName, ValidPreparationsTableNameColumn),
		buildColumnReference(ValidPreparationsTableName, ValidPreparationsTableDescriptionColumn),
		buildColumnReference(ValidPreparationsTableName, ValidPreparationsTableIconPathColumn),
		buildColumnReference(ValidPreparationsTableName, CreatedOnColumn),
		buildColumnReference(ValidPreparationsTableName, LastUpdatedOnColumn),
		buildColumnReference(ValidPreparationsTableName, ArchivedOnColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTablePrerequisiteStepColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableMinEstimatedTimeInSecondsColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableMaxEstimatedTimeInSecondsColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableTemperatureInCelsiusColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableNotesColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableWhyColumn),
		buildColumnReference(RecipeStepsTableName, CreatedOnColumn),
		buildColumnReference(RecipeStepsTableName, LastUpdatedOnColumn),
		buildColumnReference(RecipeStepsTableName, ArchivedOnColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableBelongsToRecipeColumn),
		// recipe step ingredient parts
		buildColumnReference(RecipeStepIngredientsTableName, IDColumn),
		buildColumnReference(RecipeStepIngredientsTableName, ExternalIDColumn),
		// recipe step ingredient ingredient parts
		buildColumnReference(ValidIngredientsTableName, IDColumn),
		buildColumnReference(ValidIngredientsTableName, ExternalIDColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableNameColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableVariantColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableDescriptionColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableWarningColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsEggColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsDairyColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsPeanutColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsTreeNutColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsSoyColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsWheatColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsShellfishColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsSesameColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsFishColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableContainsGlutenColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableAnimalFleshColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableAnimalDerivedColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableVolumetricColumn),
		buildColumnReference(ValidIngredientsTableName, ValidIngredientsTableIconPathColumn),
		buildColumnReference(ValidIngredientsTableName, CreatedOnColumn),
		buildColumnReference(ValidIngredientsTableName, LastUpdatedOnColumn),
		buildColumnReference(ValidIngredientsTableName, ArchivedOnColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableNameColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityTypeColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityValueColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityNotesColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableProductOfRecipeStepColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableIngredientNotesColumn),
		buildColumnReference(RecipeStepIngredientsTableName, CreatedOnColumn),
		buildColumnReference(RecipeStepIngredientsTableName, LastUpdatedOnColumn),
		buildColumnReference(RecipeStepIngredientsTableName, ArchivedOnColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableBelongsToRecipeStepColumn),
	}

	//
	// RecipeSteps Table.
	//

	// RecipeStepsTableColumns are the columns for the recipe steps table.
	RecipeStepsTableColumns = []string{
		buildColumnReference(RecipeStepsTableName, IDColumn),
		buildColumnReference(RecipeStepsTableName, ExternalIDColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableIndexColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTablePreparationIDColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTablePrerequisiteStepColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableMinEstimatedTimeInSecondsColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableMaxEstimatedTimeInSecondsColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableTemperatureInCelsiusColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableNotesColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableWhyColumn),
		buildColumnReference(RecipeStepsTableName, CreatedOnColumn),
		buildColumnReference(RecipeStepsTableName, LastUpdatedOnColumn),
		buildColumnReference(RecipeStepsTableName, ArchivedOnColumn),
		buildColumnReference(RecipeStepsTableName, RecipeStepsTableBelongsToRecipeColumn),
	}

	//
	// RecipeStepIngredients Table.
	//

	// RecipeStepIngredientsTableColumns are the columns for the recipe step ingredients table.
	RecipeStepIngredientsTableColumns = []string{
		buildColumnReference(RecipeStepIngredientsTableName, IDColumn),
		buildColumnReference(RecipeStepIngredientsTableName, ExternalIDColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableIngredientIDColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableNameColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityTypeColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityValueColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableQuantityNotesColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableProductOfRecipeStepColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableIngredientNotesColumn),
		buildColumnReference(RecipeStepIngredientsTableName, CreatedOnColumn),
		buildColumnReference(RecipeStepIngredientsTableName, LastUpdatedOnColumn),
		buildColumnReference(RecipeStepIngredientsTableName, ArchivedOnColumn),
		buildColumnReference(RecipeStepIngredientsTableName, RecipeStepIngredientsTableBelongsToRecipeStepColumn),
	}

	//
	// RecipeStepProducts Table.
	//

	// RecipeStepProductsTableColumns are the columns for the recipe step products table.
	RecipeStepProductsTableColumns = []string{
		buildColumnReference(RecipeStepProductsTableName, IDColumn),
		buildColumnReference(RecipeStepProductsTableName, ExternalIDColumn),
		buildColumnReference(RecipeStepProductsTableName, RecipeStepProductsTableNameColumn),
		buildColumnReference(RecipeStepProductsTableName, RecipeStepProductsTableQuantityTypeColumn),
		buildColumnReference(RecipeStepProductsTableName, RecipeStepProductsTableQuantityValueColumn),
		buildColumnReference(RecipeStepProductsTableName, RecipeStepProductsTableQuantityNotesColumn),
		buildColumnReference(RecipeStepProductsTableName, CreatedOnColumn),
		buildColumnReference(RecipeStepProductsTableName, LastUpdatedOnColumn),
		buildColumnReference(RecipeStepProductsTableName, ArchivedOnColumn),
		buildColumnReference(RecipeStepProductsTableName, RecipeStepProductsTableBelongsToRecipeStepColumn),
	}

	//
	// Invitations Table.
	//

	// InvitationsTableColumns are the columns for the invitations table.
	InvitationsTableColumns = []string{
		buildColumnReference(InvitationsTableName, IDColumn),
		buildColumnReference(InvitationsTableName, ExternalIDColumn),
		buildColumnReference(InvitationsTableName, InvitationsTableCodeColumn),
		buildColumnReference(InvitationsTableName, InvitationsTableConsumedColumn),
		buildColumnReference(InvitationsTableName, CreatedOnColumn),
		buildColumnReference(InvitationsTableName, LastUpdatedOnColumn),
		buildColumnReference(InvitationsTableName, ArchivedOnColumn),
		buildColumnReference(InvitationsTableName, InvitationsTableHouseholdOwnershipColumn),
	}

	//
	// Reports Table.
	//

	// ReportsTableColumns are the columns for the reports table.
	ReportsTableColumns = []string{
		buildColumnReference(ReportsTableName, IDColumn),
		buildColumnReference(ReportsTableName, ExternalIDColumn),
		buildColumnReference(ReportsTableName, ReportsTableReportTypeColumn),
		buildColumnReference(ReportsTableName, ReportsTableConcernColumn),
		buildColumnReference(ReportsTableName, CreatedOnColumn),
		buildColumnReference(ReportsTableName, LastUpdatedOnColumn),
		buildColumnReference(ReportsTableName, ArchivedOnColumn),
		buildColumnReference(ReportsTableName, ReportsTableHouseholdOwnershipColumn),
	}
)
