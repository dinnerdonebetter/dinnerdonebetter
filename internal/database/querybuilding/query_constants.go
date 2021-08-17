package querybuilding

import "fmt"

const (
	// DefaultTestUserTwoFactorSecret is the default TwoFactorSecret we give to test users when we initialize them.
	// `otpauth://totp/prixfixe:username?secret=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=&issuer=prixfixe`
	DefaultTestUserTwoFactorSecret = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="

	//
	// Common Columns.
	//

	// ExistencePrefix goes before a sql query.
	ExistencePrefix = "SELECT EXISTS ("
	// ExistenceSuffix goes after a sql query.
	ExistenceSuffix = ")"

	// IDColumn is a common column name for the sequential ID column.
	IDColumn = "id"
	// ExternalIDColumn is a common column name for the external ID column.
	ExternalIDColumn = "external_id"
	// CreatedOnColumn is a common column name for the row creation time column.
	CreatedOnColumn = "created_on"
	// LastUpdatedOnColumn is a common column name for the latest row update column.
	LastUpdatedOnColumn = "last_updated_on"
	// ArchivedOnColumn is a common column name for the archive time column.
	ArchivedOnColumn       = "archived_on"
	commaSeparator         = ","
	userOwnershipColumn    = "belongs_to_user"
	accountOwnershipColumn = "belongs_to_account"

	//
	// Accounts Table.
	//

	// AccountsTableName is what the accounts table calls itself.
	AccountsTableName = "accounts"
	// AccountsTableNameColumn is what the accounts table calls the Name column.
	AccountsTableNameColumn = "name"
	// AccountsTableBillingStatusColumn is what the accounts table calls the BillingStatus column.
	AccountsTableBillingStatusColumn = "billing_status"
	// AccountsTableContactEmailColumn is what the accounts table calls the ContactEmail column.
	AccountsTableContactEmailColumn = "contact_email"
	// AccountsTableContactPhoneColumn is what the accounts table calls the ContactPhone column.
	AccountsTableContactPhoneColumn = "contact_phone"
	// AccountsTablePaymentProcessorCustomerIDColumn is what the accounts table calls the PaymentProcessorCustomerID column.
	AccountsTablePaymentProcessorCustomerIDColumn = "payment_processor_customer_id"
	// AccountsTableSubscriptionPlanIDColumn is what the accounts table calls the SubscriptionPlanID column.
	AccountsTableSubscriptionPlanIDColumn = "subscription_plan_id"
	// AccountsTableUserOwnershipColumn is what the accounts table calls the user ownership column.
	AccountsTableUserOwnershipColumn = userOwnershipColumn

	//
	// Accounts Membership Table.
	//

	// AccountsUserMembershipTableName is what the accounts membership table calls itself.
	AccountsUserMembershipTableName = "account_user_memberships"
	// AccountsUserMembershipTableAccountRolesColumn is what the accounts membership table calls the column indicating account role.
	AccountsUserMembershipTableAccountRolesColumn = "account_roles"
	// AccountsUserMembershipTableAccountOwnershipColumn is what the accounts membership table calls the user ownership column.
	AccountsUserMembershipTableAccountOwnershipColumn = accountOwnershipColumn
	// AccountsUserMembershipTableUserOwnershipColumn is what the accounts membership table calls the user ownership column.
	AccountsUserMembershipTableUserOwnershipColumn = userOwnershipColumn
	// AccountsUserMembershipTableDefaultUserAccountColumn is what the accounts membership table calls the .
	AccountsUserMembershipTableDefaultUserAccountColumn = "default_account"

	//
	// Users Table.
	//

	// UsersTableName is what the users table calls the <> column.
	UsersTableName = "users"
	// UsersTableUsernameColumn is what the users table calls the <> column.
	UsersTableUsernameColumn = "username"
	// UsersTableHashedPasswordColumn is what the users table calls the <> column.
	UsersTableHashedPasswordColumn = "hashed_password"
	// UsersTableRequiresPasswordChangeColumn is what the users table calls the <> column.
	UsersTableRequiresPasswordChangeColumn = "requires_password_change"
	// UsersTablePasswordLastChangedOnColumn is what the users table calls the <> column.
	UsersTablePasswordLastChangedOnColumn = "password_last_changed_on"
	// UsersTableTwoFactorSekretColumn is what the users table calls the <> column.
	UsersTableTwoFactorSekretColumn = "two_factor_secret"
	// UsersTableTwoFactorVerifiedOnColumn is what the users table calls the <> column.
	UsersTableTwoFactorVerifiedOnColumn = "two_factor_secret_verified_on"
	// UsersTableServiceRolesColumn is what the users table calls the <> column.
	UsersTableServiceRolesColumn = "service_roles"
	// UsersTableReputationColumn is what the users table calls the <> column.
	UsersTableReputationColumn = "reputation"
	// UsersTableStatusExplanationColumn is what the users table calls the <> column.
	UsersTableStatusExplanationColumn = "reputation_explanation"
	// UsersTableAvatarColumn is what the users table calls the <> column.
	UsersTableAvatarColumn = "avatar_src"

	//
	// Audit Log Entries Table.
	//

	// AuditLogEntriesTableName is what the audit log entries table calls itself.
	AuditLogEntriesTableName = "audit_log"
	// AuditLogEntriesTableEventTypeColumn is what the audit log entries table calls the event type column.
	AuditLogEntriesTableEventTypeColumn = "event_type"
	// AuditLogEntriesTableContextColumn is what the audit log entries table calls the context column.
	AuditLogEntriesTableContextColumn = "context"

	//
	// API Clients.
	//

	// APIClientsTableName is what the API clients table calls the <> column.
	APIClientsTableName = "api_clients"
	// APIClientsTableNameColumn is what the API clients table calls the <> column.
	APIClientsTableNameColumn = "name"
	// APIClientsTableClientIDColumn is what the API clients table calls the <> column.
	APIClientsTableClientIDColumn = "client_id"
	// APIClientsTableSecretKeyColumn is what the API clients table calls the <> column.
	APIClientsTableSecretKeyColumn = "secret_key"
	// APIClientsTableOwnershipColumn is what the API clients table calls the <> column.
	APIClientsTableOwnershipColumn = userOwnershipColumn

	//
	// Webhooks Table.
	//

	// WebhooksTableName is what the webhooks table calls the <> column.
	WebhooksTableName = "webhooks"
	// WebhooksTableNameColumn is what the webhooks table calls the <> column.
	WebhooksTableNameColumn = "name"
	// WebhooksTableContentTypeColumn is what the webhooks table calls the <> column.
	WebhooksTableContentTypeColumn = "content_type"
	// WebhooksTableURLColumn is what the webhooks table calls the <> column.
	WebhooksTableURLColumn = "url"
	// WebhooksTableMethodColumn is what the webhooks table calls the <> column.
	WebhooksTableMethodColumn = "method"
	// WebhooksTableEventsColumn is what the webhooks table calls the <> column.
	WebhooksTableEventsColumn = "events"
	// WebhooksTableEventsSeparator is what the webhooks table calls the <> column.
	WebhooksTableEventsSeparator = commaSeparator
	// WebhooksTableDataTypesColumn is what the webhooks table calls the <> column.
	WebhooksTableDataTypesColumn = "data_types"
	// WebhooksTableDataTypesSeparator is what the webhooks table calls the <> column.
	WebhooksTableDataTypesSeparator = commaSeparator
	// WebhooksTableTopicsColumn is what the webhooks table calls the <> column.
	WebhooksTableTopicsColumn = "topics"
	// WebhooksTableTopicsSeparator is what the webhooks table calls the <> column.
	WebhooksTableTopicsSeparator = commaSeparator
	// WebhooksTableOwnershipColumn is what the webhooks table calls the <> column.
	WebhooksTableOwnershipColumn = accountOwnershipColumn

	//
	// ValidInstruments Table.
	//

	// ValidInstrumentsTableName is what the valid instruments table calls itself.
	ValidInstrumentsTableName = "valid_instruments"
	// ValidInstrumentsTableNameColumn is what the valid instruments table calls the name column.
	ValidInstrumentsTableNameColumn = "name"
	// ValidInstrumentsTableVariantColumn is what the valid instruments table calls the variant column.
	ValidInstrumentsTableVariantColumn = "variant"
	// ValidInstrumentsTableDescriptionColumn is what the valid instruments table calls the description column.
	ValidInstrumentsTableDescriptionColumn = "description"
	// ValidInstrumentsTableIconPathColumn is what the valid instruments table calls the icon_path column.
	ValidInstrumentsTableIconPathColumn = "icon_path"
	// ValidInstrumentsTableAccountOwnershipColumn is what the valid instruments table calls the ownership column.
	ValidInstrumentsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// ValidPreparations Table.
	//

	// ValidPreparationsTableName is what the valid preparations table calls itself.
	ValidPreparationsTableName = "valid_preparations"
	// ValidPreparationsTableNameColumn is what the valid preparations table calls the name column.
	ValidPreparationsTableNameColumn = "name"
	// ValidPreparationsTableDescriptionColumn is what the valid preparations table calls the description column.
	ValidPreparationsTableDescriptionColumn = "description"
	// ValidPreparationsTableIconPathColumn is what the valid preparations table calls the icon_path column.
	ValidPreparationsTableIconPathColumn = "icon_path"
	// ValidPreparationsTableAccountOwnershipColumn is what the valid preparations table calls the ownership column.
	ValidPreparationsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// ValidIngredients Table.
	//

	// ValidIngredientsTableName is what the valid ingredients table calls itself.
	ValidIngredientsTableName = "valid_ingredients"
	// ValidIngredientsTableNameColumn is what the valid ingredients table calls the name column.
	ValidIngredientsTableNameColumn = "name"
	// ValidIngredientsTableVariantColumn is what the valid ingredients table calls the variant column.
	ValidIngredientsTableVariantColumn = "variant"
	// ValidIngredientsTableDescriptionColumn is what the valid ingredients table calls the description column.
	ValidIngredientsTableDescriptionColumn = "description"
	// ValidIngredientsTableWarningColumn is what the valid ingredients table calls the warning column.
	ValidIngredientsTableWarningColumn = "warning"
	// ValidIngredientsTableContainsEggColumn is what the valid ingredients table calls the contains_egg column.
	ValidIngredientsTableContainsEggColumn = "contains_egg"
	// ValidIngredientsTableContainsDairyColumn is what the valid ingredients table calls the contains_dairy column.
	ValidIngredientsTableContainsDairyColumn = "contains_dairy"
	// ValidIngredientsTableContainsPeanutColumn is what the valid ingredients table calls the contains_peanut column.
	ValidIngredientsTableContainsPeanutColumn = "contains_peanut"
	// ValidIngredientsTableContainsTreeNutColumn is what the valid ingredients table calls the contains_tree_nut column.
	ValidIngredientsTableContainsTreeNutColumn = "contains_tree_nut"
	// ValidIngredientsTableContainsSoyColumn is what the valid ingredients table calls the contains_soy column.
	ValidIngredientsTableContainsSoyColumn = "contains_soy"
	// ValidIngredientsTableContainsWheatColumn is what the valid ingredients table calls the contains_wheat column.
	ValidIngredientsTableContainsWheatColumn = "contains_wheat"
	// ValidIngredientsTableContainsShellfishColumn is what the valid ingredients table calls the contains_shellfish column.
	ValidIngredientsTableContainsShellfishColumn = "contains_shellfish"
	// ValidIngredientsTableContainsSesameColumn is what the valid ingredients table calls the contains_sesame column.
	ValidIngredientsTableContainsSesameColumn = "contains_sesame"
	// ValidIngredientsTableContainsFishColumn is what the valid ingredients table calls the contains_fish column.
	ValidIngredientsTableContainsFishColumn = "contains_fish"
	// ValidIngredientsTableContainsGlutenColumn is what the valid ingredients table calls the contains_gluten column.
	ValidIngredientsTableContainsGlutenColumn = "contains_gluten"
	// ValidIngredientsTableAnimalFleshColumn is what the valid ingredients table calls the animal_flesh column.
	ValidIngredientsTableAnimalFleshColumn = "animal_flesh"
	// ValidIngredientsTableAnimalDerivedColumn is what the valid ingredients table calls the animal_derived column.
	ValidIngredientsTableAnimalDerivedColumn = "animal_derived"
	// ValidIngredientsTableVolumetricColumn is what the valid ingredients table calls the volumetric column.
	ValidIngredientsTableVolumetricColumn = "volumetric"
	// ValidIngredientsTableIconPathColumn is what the valid ingredients table calls the icon_path column.
	ValidIngredientsTableIconPathColumn = "icon_path"
	// ValidIngredientsTableAccountOwnershipColumn is what the valid ingredients table calls the ownership column.
	ValidIngredientsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// ValidIngredientPreparations Table.
	//

	// ValidIngredientPreparationsTableName is what the valid ingredient preparations table calls itself.
	ValidIngredientPreparationsTableName = "valid_ingredient_preparations"
	// ValidIngredientPreparationsTableNotesColumn is what the valid ingredient preparations table calls the notes column.
	ValidIngredientPreparationsTableNotesColumn = "notes"
	// ValidIngredientPreparationsTableValidIngredientIDColumn is what the valid ingredient preparations table calls the valid_ingredient_id column.
	ValidIngredientPreparationsTableValidIngredientIDColumn = "valid_ingredient_id"
	// ValidIngredientPreparationsTableValidPreparationIDColumn is what the valid ingredient preparations table calls the valid_preparation_id column.
	ValidIngredientPreparationsTableValidPreparationIDColumn = "valid_preparation_id"
	// ValidIngredientPreparationsTableAccountOwnershipColumn is what the valid ingredient preparations table calls the ownership column.
	ValidIngredientPreparationsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// ValidPreparationInstruments Table.
	//

	// ValidPreparationInstrumentsTableName is what the valid preparation instruments table calls itself.
	ValidPreparationInstrumentsTableName = "valid_preparation_instruments"
	// ValidPreparationInstrumentsTableInstrumentIDColumn is what the valid preparation instruments table calls the instrument_id column.
	ValidPreparationInstrumentsTableInstrumentIDColumn = "instrument_id"
	// ValidPreparationInstrumentsTablePreparationIDColumn is what the valid preparation instruments table calls the preparation_id column.
	ValidPreparationInstrumentsTablePreparationIDColumn = "preparation_id"
	// ValidPreparationInstrumentsTableNotesColumn is what the valid preparation instruments table calls the notes column.
	ValidPreparationInstrumentsTableNotesColumn = "notes"
	// ValidPreparationInstrumentsTableAccountOwnershipColumn is what the valid preparation instruments table calls the ownership column.
	ValidPreparationInstrumentsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// Recipes Table.
	//

	// RecipesTableName is what the recipes table calls itself.
	RecipesTableName = "recipes"
	// RecipesTableNameColumn is what the recipes table calls the name column.
	RecipesTableNameColumn = "name"
	// RecipesTableSourceColumn is what the recipes table calls the source column.
	RecipesTableSourceColumn = "source"
	// RecipesTableDescriptionColumn is what the recipes table calls the description column.
	RecipesTableDescriptionColumn = "description"
	// RecipesTableDisplayImageURLColumn is what the recipes table calls the description column.
	RecipesTableDisplayImageURLColumn = "display_image_url"
	// RecipesTableInspiredByRecipeIDColumn is what the recipes table calls the inspired_by_recipe_id column.
	RecipesTableInspiredByRecipeIDColumn = "inspired_by_recipe_id"
	// RecipesTableAccountOwnershipColumn is what the recipes table calls the ownership column.
	RecipesTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// RecipeSteps Table.
	//

	// RecipeStepsTableName is what the recipe steps table calls itself.
	RecipeStepsTableName = "recipe_steps"
	// RecipeStepsTableIndexColumn is what the recipe steps table calls the index column.
	RecipeStepsTableIndexColumn = "index"
	// RecipeStepsTablePreparationIDColumn is what the recipe steps table calls the preparation_id column.
	RecipeStepsTablePreparationIDColumn = "preparation_id"
	// RecipeStepsTablePrerequisiteStepColumn is what the recipe steps table calls the prerequisite_step column.
	RecipeStepsTablePrerequisiteStepColumn = "prerequisite_step"
	// RecipeStepsTableMinEstimatedTimeInSecondsColumn is what the recipe steps table calls the min_estimated_time_in_seconds column.
	RecipeStepsTableMinEstimatedTimeInSecondsColumn = "min_estimated_time_in_seconds"
	// RecipeStepsTableMaxEstimatedTimeInSecondsColumn is what the recipe steps table calls the max_estimated_time_in_seconds column.
	RecipeStepsTableMaxEstimatedTimeInSecondsColumn = "max_estimated_time_in_seconds"
	// RecipeStepsTableTemperatureInCelsiusColumn is what the recipe steps table calls the temperature_in_celsius column.
	RecipeStepsTableTemperatureInCelsiusColumn = "temperature_in_celsius"
	// RecipeStepsTableNotesColumn is what the recipe steps table calls the notes column.
	RecipeStepsTableNotesColumn = "notes"
	// RecipeStepsTableWhyColumn is what the recipe steps table calls the why column.
	RecipeStepsTableWhyColumn = "why"
	// RecipeStepsTableBelongsToRecipeColumn is what the recipe steps table calls the recipe ownership column.
	RecipeStepsTableBelongsToRecipeColumn = "belongs_to_recipe"
	// RecipeStepsTableAccountOwnershipColumn is what the recipe steps table calls the ownership column.
	RecipeStepsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// RecipeStepIngredients Table.
	//

	// RecipeStepIngredientsTableName is what the recipe step ingredients table calls itself.
	RecipeStepIngredientsTableName = "recipe_step_ingredients"
	// RecipeStepIngredientsTableIngredientIDColumn is what the recipe step ingredients table calls the ingredient_id column.
	RecipeStepIngredientsTableIngredientIDColumn = "ingredient_id"
	// RecipeStepIngredientsTableNameColumn is what the recipe step ingredients table calls the name column.
	RecipeStepIngredientsTableNameColumn = "name"
	// RecipeStepIngredientsTableQuantityTypeColumn is what the recipe step ingredients table calls the quantity_type column.
	RecipeStepIngredientsTableQuantityTypeColumn = "quantity_type"
	// RecipeStepIngredientsTableQuantityValueColumn is what the recipe step ingredients table calls the quantity_value column.
	RecipeStepIngredientsTableQuantityValueColumn = "quantity_value"
	// RecipeStepIngredientsTableQuantityNotesColumn is what the recipe step ingredients table calls the quantity_notes column.
	RecipeStepIngredientsTableQuantityNotesColumn = "quantity_notes"
	// RecipeStepIngredientsTableProductOfRecipeStepColumn is what the recipe step ingredients table calls the product_of_recipe_step column.
	RecipeStepIngredientsTableProductOfRecipeStepColumn = "product_of_recipe_step"
	// RecipeStepIngredientsTableIngredientNotesColumn is what the recipe step ingredients table calls the ingredient_notes column.
	RecipeStepIngredientsTableIngredientNotesColumn = "ingredient_notes"
	// RecipeStepIngredientsTableBelongsToRecipeStepColumn is what the recipe step ingredients table calls the recipe step ownership column.
	RecipeStepIngredientsTableBelongsToRecipeStepColumn = "belongs_to_recipe_step"
	// RecipeStepIngredientsTableAccountOwnershipColumn is what the recipe step ingredients table calls the ownership column.
	RecipeStepIngredientsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// RecipeStepProducts Table.
	//

	// RecipeStepProductsTableName is what the recipe step products table calls itself.
	RecipeStepProductsTableName = "recipe_step_products"
	// RecipeStepProductsTableNameColumn is what the recipe step products table calls the name column.
	RecipeStepProductsTableNameColumn = "name"
	// RecipeStepProductsTableQuantityTypeColumn is what the recipe step products table calls the quantity_type column.
	RecipeStepProductsTableQuantityTypeColumn = "quantity_type"
	// RecipeStepProductsTableQuantityValueColumn is what the recipe step products table calls the quantity_value column.
	RecipeStepProductsTableQuantityValueColumn = "quantity_value"
	// RecipeStepProductsTableQuantityNotesColumn is what the recipe step products table calls the quantity_notes column.
	RecipeStepProductsTableQuantityNotesColumn = "quantity_notes"
	// RecipeStepProductsTableBelongsToRecipeStepColumn is what the recipe step products table calls the recipe step ownership column.
	RecipeStepProductsTableBelongsToRecipeStepColumn = "belongs_to_recipe_step"
	// RecipeStepProductsTableAccountOwnershipColumn is what the recipe step products table calls the ownership column.
	RecipeStepProductsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// Invitations Table.
	//

	// InvitationsTableName is what the invitations table calls itself.
	InvitationsTableName = "invitations"
	// InvitationsTableCodeColumn is what the invitations table calls the code column.
	InvitationsTableCodeColumn = "code"
	// InvitationsTableConsumedColumn is what the invitations table calls the consumed column.
	InvitationsTableConsumedColumn = "consumed"
	// InvitationsTableAccountOwnershipColumn is what the invitations table calls the ownership column.
	InvitationsTableAccountOwnershipColumn = accountOwnershipColumn

	//
	// Reports Table.
	//

	// ReportsTableName is what the reports table calls itself.
	ReportsTableName = "reports"
	// ReportsTableReportTypeColumn is what the reports table calls the report_type column.
	ReportsTableReportTypeColumn = "report_type"
	// ReportsTableConcernColumn is what the reports table calls the concern column.
	ReportsTableConcernColumn = "concern"
	// ReportsTableAccountOwnershipColumn is what the reports table calls the ownership column.
	ReportsTableAccountOwnershipColumn = accountOwnershipColumn
)

var (
	// RecipesOnRecipeStepsJoinClause is a join clause that allows the recipe steps table to be joined on the recipes table.
	RecipesOnRecipeStepsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.id", RecipesTableName, RecipeStepsTableName, RecipeStepsTableBelongsToRecipeColumn, RecipesTableName)
	// RecipeStepsOnRecipeStepIngredientsJoinClause is a join clause that allows the recipe step ingredients table to be joined on the recipe steps table.
	RecipeStepsOnRecipeStepIngredientsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.id", RecipeStepsTableName, RecipeStepIngredientsTableName, RecipeStepIngredientsTableBelongsToRecipeStepColumn, RecipeStepsTableName)
	// RecipeStepsOnRecipeStepProductsJoinClause is a join clause that allows the recipe step products table to be joined on the recipe steps table.
	RecipeStepsOnRecipeStepProductsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.id", RecipeStepsTableName, RecipeStepProductsTableName, RecipeStepProductsTableBelongsToRecipeStepColumn, RecipeStepsTableName)
)
