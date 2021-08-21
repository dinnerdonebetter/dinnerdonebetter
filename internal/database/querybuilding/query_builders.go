package querybuilding

import (
	"context"
	"database/sql"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

type (
	// HouseholdSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	HouseholdSQLQueryBuilder interface {
		BuildGetHouseholdQuery(ctx context.Context, householdID, userID uint64) (query string, args []interface{})
		BuildGetAllHouseholdsCountQuery(ctx context.Context) string
		BuildGetBatchOfHouseholdsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetHouseholdsQuery(ctx context.Context, userID uint64, forAdmin bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildHouseholdCreationQuery(ctx context.Context, input *types.HouseholdCreationInput) (query string, args []interface{})
		BuildUpdateHouseholdQuery(ctx context.Context, input *types.Household) (query string, args []interface{})
		BuildArchiveHouseholdQuery(ctx context.Context, householdID, userID uint64) (query string, args []interface{})
		BuildTransferHouseholdOwnershipQuery(ctx context.Context, currentOwnerID, newOwnerID, householdID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForHouseholdQuery(ctx context.Context, householdID uint64) (query string, args []interface{})
	}

	// HouseholdUserMembershipSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	HouseholdUserMembershipSQLQueryBuilder interface {
		BuildGetDefaultHouseholdIDForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{})
		BuildArchiveHouseholdMembershipsForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{})
		BuildGetHouseholdMembershipsForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{})
		BuildMarkHouseholdAsUserDefaultQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{})
		BuildModifyUserPermissionsQuery(ctx context.Context, userID, householdID uint64, newRoles []string) (query string, args []interface{})
		BuildTransferHouseholdMembershipsQuery(ctx context.Context, currentOwnerID, newOwnerID, householdID uint64) (query string, args []interface{})
		BuildUserIsMemberOfHouseholdQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{})
		BuildCreateMembershipForNewUserQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{})
		BuildAddUserToHouseholdQuery(ctx context.Context, input *types.AddUserToHouseholdInput) (query string, args []interface{})
		BuildRemoveUserFromHouseholdQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{})
	}

	// APIClientSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	APIClientSQLQueryBuilder interface {
		BuildGetBatchOfAPIClientsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetAPIClientByClientIDQuery(ctx context.Context, clientID string) (query string, args []interface{})
		BuildGetAPIClientByDatabaseIDQuery(ctx context.Context, clientID, userID uint64) (query string, args []interface{})
		BuildGetAllAPIClientsCountQuery(ctx context.Context) string
		BuildGetAPIClientsQuery(ctx context.Context, userID uint64, filter *types.QueryFilter) (query string, args []interface{})
		BuildCreateAPIClientQuery(ctx context.Context, input *types.APIClientCreationInput) (query string, args []interface{})
		BuildUpdateAPIClientQuery(ctx context.Context, input *types.APIClient) (query string, args []interface{})
		BuildArchiveAPIClientQuery(ctx context.Context, clientID, userID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForAPIClientQuery(ctx context.Context, clientID uint64) (query string, args []interface{})
	}

	// AuditLogEntrySQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	AuditLogEntrySQLQueryBuilder interface {
		BuildGetAuditLogEntryQuery(ctx context.Context, entryID uint64) (query string, args []interface{})
		BuildGetAllAuditLogEntriesCountQuery(ctx context.Context) string
		BuildGetBatchOfAuditLogEntriesQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesQuery(ctx context.Context, filter *types.QueryFilter) (query string, args []interface{})
		BuildCreateAuditLogEntryQuery(ctx context.Context, input *types.AuditLogEntryCreationInput) (query string, args []interface{})
	}

	// UserSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	UserSQLQueryBuilder interface {
		BuildUserHasStatusQuery(ctx context.Context, userID uint64, statuses ...string) (query string, args []interface{})
		BuildGetUserQuery(ctx context.Context, userID uint64) (query string, args []interface{})
		BuildGetUsersQuery(ctx context.Context, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetUserWithUnverifiedTwoFactorSecretQuery(ctx context.Context, userID uint64) (query string, args []interface{})
		BuildGetUserByUsernameQuery(ctx context.Context, username string) (query string, args []interface{})
		BuildSearchForUserByUsernameQuery(ctx context.Context, usernameQuery string) (query string, args []interface{})
		BuildGetAllUsersCountQuery(ctx context.Context) (query string)
		BuildCreateUserQuery(ctx context.Context, input *types.UserDataStoreCreationInput) (query string, args []interface{})
		BuildUpdateUserQuery(ctx context.Context, input *types.User) (query string, args []interface{})
		BuildUpdateUserPasswordQuery(ctx context.Context, userID uint64, newHash string) (query string, args []interface{})
		BuildUpdateUserTwoFactorSecretQuery(ctx context.Context, userID uint64, newSecret string) (query string, args []interface{})
		BuildVerifyUserTwoFactorSecretQuery(ctx context.Context, userID uint64) (query string, args []interface{})
		BuildArchiveUserQuery(ctx context.Context, userID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{})
		BuildSetUserStatusQuery(ctx context.Context, input *types.UserReputationUpdateInput) (query string, args []interface{})
	}

	// WebhookSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	WebhookSQLQueryBuilder interface {
		BuildGetWebhookQuery(ctx context.Context, webhookID, householdID uint64) (query string, args []interface{})
		BuildGetAllWebhooksCountQuery(ctx context.Context) string
		BuildGetBatchOfWebhooksQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetWebhooksQuery(ctx context.Context, householdID uint64, filter *types.QueryFilter) (query string, args []interface{})
		BuildCreateWebhookQuery(ctx context.Context, x *types.WebhookCreationInput) (query string, args []interface{})
		BuildUpdateWebhookQuery(ctx context.Context, input *types.Webhook) (query string, args []interface{})
		BuildArchiveWebhookQuery(ctx context.Context, webhookID, householdID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForWebhookQuery(ctx context.Context, webhookID uint64) (query string, args []interface{})
	}

	// ValidInstrumentSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	ValidInstrumentSQLQueryBuilder interface {
		BuildValidInstrumentExistsQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{})
		BuildGetValidInstrumentIDForNameQuery(ctx context.Context, validInstrumentName string) (query string, args []interface{})
		BuildSearchForValidInstrumentByNameQuery(ctx context.Context, name string) (query string, args []interface{})
		BuildGetValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{})
		BuildGetAllValidInstrumentsCountQuery(ctx context.Context) string
		BuildGetBatchOfValidInstrumentsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetValidInstrumentsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetValidInstrumentsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{})
		BuildCreateValidInstrumentQuery(ctx context.Context, input *types.ValidInstrumentCreationInput) (query string, args []interface{})
		BuildUpdateValidInstrumentQuery(ctx context.Context, input *types.ValidInstrument) (query string, args []interface{})
		BuildArchiveValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForValidInstrumentQuery(ctx context.Context, validInstrumentID uint64) (query string, args []interface{})
	}

	// ValidPreparationSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	ValidPreparationSQLQueryBuilder interface {
		BuildValidPreparationExistsQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{})
		BuildGetValidPreparationIDForNameQuery(ctx context.Context, validPreparationName string) (query string, args []interface{})
		BuildSearchForValidPreparationByNameQuery(ctx context.Context, name string) (query string, args []interface{})
		BuildGetValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{})
		BuildGetAllValidPreparationsCountQuery(ctx context.Context) string
		BuildGetBatchOfValidPreparationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetValidPreparationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetValidPreparationsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{})
		BuildCreateValidPreparationQuery(ctx context.Context, input *types.ValidPreparationCreationInput) (query string, args []interface{})
		BuildUpdateValidPreparationQuery(ctx context.Context, input *types.ValidPreparation) (query string, args []interface{})
		BuildArchiveValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForValidPreparationQuery(ctx context.Context, validPreparationID uint64) (query string, args []interface{})
	}

	// ValidIngredientSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	ValidIngredientSQLQueryBuilder interface {
		BuildValidIngredientExistsQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{})
		BuildGetValidIngredientIDForNameQuery(ctx context.Context, validIngredientName string) (query string, args []interface{})
		BuildSearchForValidIngredientByNameQuery(ctx context.Context, name string) (query string, args []interface{})
		BuildGetValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{})
		BuildGetAllValidIngredientsCountQuery(ctx context.Context) string
		BuildGetBatchOfValidIngredientsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetValidIngredientsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetValidIngredientsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{})
		BuildCreateValidIngredientQuery(ctx context.Context, input *types.ValidIngredientCreationInput) (query string, args []interface{})
		BuildUpdateValidIngredientQuery(ctx context.Context, input *types.ValidIngredient) (query string, args []interface{})
		BuildArchiveValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForValidIngredientQuery(ctx context.Context, validIngredientID uint64) (query string, args []interface{})
	}

	// ValidIngredientPreparationSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	ValidIngredientPreparationSQLQueryBuilder interface {
		BuildValidIngredientPreparationExistsQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{})
		BuildGetValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{})
		BuildGetAllValidIngredientPreparationsCountQuery(ctx context.Context) string
		BuildGetBatchOfValidIngredientPreparationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetValidIngredientPreparationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetValidIngredientPreparationsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{})
		BuildCreateValidIngredientPreparationQuery(ctx context.Context, input *types.ValidIngredientPreparationCreationInput) (query string, args []interface{})
		BuildUpdateValidIngredientPreparationQuery(ctx context.Context, input *types.ValidIngredientPreparation) (query string, args []interface{})
		BuildArchiveValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForValidIngredientPreparationQuery(ctx context.Context, validIngredientPreparationID uint64) (query string, args []interface{})
	}

	// ValidPreparationInstrumentSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	ValidPreparationInstrumentSQLQueryBuilder interface {
		BuildValidPreparationInstrumentExistsQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{})
		BuildGetValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{})
		BuildGetAllValidPreparationInstrumentsCountQuery(ctx context.Context) string
		BuildGetBatchOfValidPreparationInstrumentsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetValidPreparationInstrumentsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetValidPreparationInstrumentsWithIDsQuery(ctx context.Context, limit uint8, ids []uint64) (query string, args []interface{})
		BuildCreateValidPreparationInstrumentQuery(ctx context.Context, input *types.ValidPreparationInstrumentCreationInput) (query string, args []interface{})
		BuildUpdateValidPreparationInstrumentQuery(ctx context.Context, input *types.ValidPreparationInstrument) (query string, args []interface{})
		BuildArchiveValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForValidPreparationInstrumentQuery(ctx context.Context, validPreparationInstrumentID uint64) (query string, args []interface{})
	}

	// RecipeSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	RecipeSQLQueryBuilder interface {
		BuildRecipeExistsQuery(ctx context.Context, recipeID uint64) (query string, args []interface{})
		BuildGetRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{})
		BuildGetAllRecipesCountQuery(ctx context.Context) string
		BuildGetBatchOfRecipesQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetRecipesQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetRecipesWithIDsQuery(ctx context.Context, householdID uint64, limit uint8, ids []uint64, restrictToHousehold bool) (query string, args []interface{})
		BuildCreateRecipeQuery(ctx context.Context, input *types.RecipeCreationInput) (query string, args []interface{})
		BuildUpdateRecipeQuery(ctx context.Context, input *types.Recipe) (query string, args []interface{})
		BuildArchiveRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForRecipeQuery(ctx context.Context, recipeID uint64) (query string, args []interface{})
	}

	// RecipeStepSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	RecipeStepSQLQueryBuilder interface {
		BuildRecipeStepExistsQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{})
		BuildGetRecipeStepQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{})
		BuildGetAllRecipeStepsCountQuery(ctx context.Context) string
		BuildGetBatchOfRecipeStepsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetRecipeStepsQuery(ctx context.Context, recipeID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetRecipeStepsWithIDsQuery(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) (query string, args []interface{})
		BuildCreateRecipeStepQuery(ctx context.Context, input *types.RecipeStepCreationInput) (query string, args []interface{})
		BuildUpdateRecipeStepQuery(ctx context.Context, input *types.RecipeStep) (query string, args []interface{})
		BuildArchiveRecipeStepQuery(ctx context.Context, recipeID, recipeStepID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForRecipeStepQuery(ctx context.Context, recipeStepID uint64) (query string, args []interface{})
	}

	// RecipeStepIngredientSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	RecipeStepIngredientSQLQueryBuilder interface {
		BuildRecipeStepIngredientExistsQuery(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{})
		BuildGetRecipeStepIngredientQuery(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{})
		BuildGetAllRecipeStepIngredientsCountQuery(ctx context.Context) string
		BuildGetBatchOfRecipeStepIngredientsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetRecipeStepIngredientsQuery(ctx context.Context, recipeID, recipeStepID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetRecipeStepIngredientsWithIDsQuery(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{})
		BuildCreateRecipeStepIngredientQuery(ctx context.Context, input *types.RecipeStepIngredientCreationInput) (query string, args []interface{})
		BuildUpdateRecipeStepIngredientQuery(ctx context.Context, input *types.RecipeStepIngredient) (query string, args []interface{})
		BuildArchiveRecipeStepIngredientQuery(ctx context.Context, recipeStepID, recipeStepIngredientID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForRecipeStepIngredientQuery(ctx context.Context, recipeStepIngredientID uint64) (query string, args []interface{})
	}

	// RecipeStepProductSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	RecipeStepProductSQLQueryBuilder interface {
		BuildRecipeStepProductExistsQuery(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (query string, args []interface{})
		BuildGetRecipeStepProductQuery(ctx context.Context, recipeID, recipeStepID, recipeStepProductID uint64) (query string, args []interface{})
		BuildGetAllRecipeStepProductsCountQuery(ctx context.Context) string
		BuildGetBatchOfRecipeStepProductsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetRecipeStepProductsQuery(ctx context.Context, recipeID, recipeStepID uint64, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetRecipeStepProductsWithIDsQuery(ctx context.Context, recipeStepID uint64, limit uint8, ids []uint64) (query string, args []interface{})
		BuildCreateRecipeStepProductQuery(ctx context.Context, input *types.RecipeStepProductCreationInput) (query string, args []interface{})
		BuildUpdateRecipeStepProductQuery(ctx context.Context, input *types.RecipeStepProduct) (query string, args []interface{})
		BuildArchiveRecipeStepProductQuery(ctx context.Context, recipeStepID, recipeStepProductID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForRecipeStepProductQuery(ctx context.Context, recipeStepProductID uint64) (query string, args []interface{})
	}

	// InvitationSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	InvitationSQLQueryBuilder interface {
		BuildInvitationExistsQuery(ctx context.Context, invitationID uint64) (query string, args []interface{})
		BuildGetInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{})
		BuildGetAllInvitationsCountQuery(ctx context.Context) string
		BuildGetBatchOfInvitationsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetInvitationsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetInvitationsWithIDsQuery(ctx context.Context, householdID uint64, limit uint8, ids []uint64, restrictToHousehold bool) (query string, args []interface{})
		BuildCreateInvitationQuery(ctx context.Context, input *types.InvitationCreationInput) (query string, args []interface{})
		BuildUpdateInvitationQuery(ctx context.Context, input *types.Invitation) (query string, args []interface{})
		BuildArchiveInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForInvitationQuery(ctx context.Context, invitationID uint64) (query string, args []interface{})
	}

	// ReportSQLQueryBuilder describes a structure capable of generating query/arg pairs for certain situations.
	ReportSQLQueryBuilder interface {
		BuildReportExistsQuery(ctx context.Context, reportID uint64) (query string, args []interface{})
		BuildGetReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{})
		BuildGetAllReportsCountQuery(ctx context.Context) string
		BuildGetBatchOfReportsQuery(ctx context.Context, beginID, endID uint64) (query string, args []interface{})
		BuildGetReportsQuery(ctx context.Context, includeArchived bool, filter *types.QueryFilter) (query string, args []interface{})
		BuildGetReportsWithIDsQuery(ctx context.Context, householdID uint64, limit uint8, ids []uint64, restrictToHousehold bool) (query string, args []interface{})
		BuildCreateReportQuery(ctx context.Context, input *types.ReportCreationInput) (query string, args []interface{})
		BuildUpdateReportQuery(ctx context.Context, input *types.Report) (query string, args []interface{})
		BuildArchiveReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{})
		BuildGetAuditLogEntriesForReportQuery(ctx context.Context, reportID uint64) (query string, args []interface{})
	}

	// SQLQueryBuilder describes anything that builds SQL queries to manage our data.
	SQLQueryBuilder interface {
		BuildMigrationFunc(db *sql.DB) func()
		BuildTestUserCreationQuery(ctx context.Context, testUserConfig *types.TestUserCreationConfig) (query string, args []interface{})

		HouseholdSQLQueryBuilder
		HouseholdUserMembershipSQLQueryBuilder
		UserSQLQueryBuilder
		AuditLogEntrySQLQueryBuilder
		APIClientSQLQueryBuilder
		WebhookSQLQueryBuilder
		ValidInstrumentSQLQueryBuilder
		ValidPreparationSQLQueryBuilder
		ValidIngredientSQLQueryBuilder
		ValidIngredientPreparationSQLQueryBuilder
		ValidPreparationInstrumentSQLQueryBuilder
		RecipeSQLQueryBuilder
		RecipeStepSQLQueryBuilder
		RecipeStepIngredientSQLQueryBuilder
		RecipeStepProductSQLQueryBuilder
		InvitationSQLQueryBuilder
		ReportSQLQueryBuilder
	}
)
