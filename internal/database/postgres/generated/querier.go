// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package generated

import (
	"context"
	"database/sql"
)

type Querier interface {
	ArchiveAPIClient(ctx context.Context, db DBTX, arg *ArchiveAPIClientParams) error
	ArchiveHousehold(ctx context.Context, db DBTX, arg *ArchiveHouseholdParams) error
	ArchiveHouseholdUserMembershipForUser(ctx context.Context, db DBTX, belongsToUser string) error
	ArchiveMeal(ctx context.Context, db DBTX, arg *ArchiveMealParams) error
	ArchiveMealPlan(ctx context.Context, db DBTX, arg *ArchiveMealPlanParams) error
	ArchiveMealPlanEvent(ctx context.Context, db DBTX, arg *ArchiveMealPlanEventParams) error
	ArchiveMealPlanGroceryListItem(ctx context.Context, db DBTX, id string) error
	ArchiveMealPlanOption(ctx context.Context, db DBTX, arg *ArchiveMealPlanOptionParams) error
	ArchiveMealPlanOptionVote(ctx context.Context, db DBTX, arg *ArchiveMealPlanOptionVoteParams) error
	ArchiveRecipe(ctx context.Context, db DBTX, arg *ArchiveRecipeParams) error
	ArchiveRecipeMedia(ctx context.Context, db DBTX, id string) error
	ArchiveRecipePrepTask(ctx context.Context, db DBTX, id string) error
	ArchiveRecipeStep(ctx context.Context, db DBTX, arg *ArchiveRecipeStepParams) error
	ArchiveRecipeStepIngredient(ctx context.Context, db DBTX, arg *ArchiveRecipeStepIngredientParams) error
	ArchiveRecipeStepInstrument(ctx context.Context, db DBTX, arg *ArchiveRecipeStepInstrumentParams) error
	ArchiveRecipeStepProduct(ctx context.Context, db DBTX, arg *ArchiveRecipeStepProductParams) error
	ArchiveUser(ctx context.Context, db DBTX, id string) error
	ArchiveValidIngredient(ctx context.Context, db DBTX, id string) error
	ArchiveValidIngredientMeasurementUnit(ctx context.Context, db DBTX, id string) error
	ArchiveValidIngredientPreparation(ctx context.Context, db DBTX, id string) error
	ArchiveValidInstrument(ctx context.Context, db DBTX, id string) error
	ArchiveValidMeasurementConversion(ctx context.Context, db DBTX, id string) error
	ArchiveValidMeasurementUnit(ctx context.Context, db DBTX, id string) error
	ArchiveValidPreparation(ctx context.Context, db DBTX, id string) error
	ArchiveValidPreparationInstrument(ctx context.Context, db DBTX, id string) error
	ArchiveWebhook(ctx context.Context, db DBTX, arg *ArchiveWebhookParams) error
	AttachHouseholdInvitationsToUser(ctx context.Context, db DBTX, arg *AttachHouseholdInvitationsToUserParams) error
	ChangeMealPlanTaskStatus(ctx context.Context, db DBTX, arg *ChangeMealPlanTaskStatusParams) error
	CreateAPIClient(ctx context.Context, db DBTX, arg *CreateAPIClientParams) error
	CreateHousehold(ctx context.Context, db DBTX, arg *CreateHouseholdParams) error
	CreateHouseholdInvitation(ctx context.Context, db DBTX, arg *CreateHouseholdInvitationParams) error
	CreateHouseholdUserMembership(ctx context.Context, db DBTX, arg *CreateHouseholdUserMembershipParams) error
	CreateHouseholdUserMembershipForNewUser(ctx context.Context, db DBTX, arg *CreateHouseholdUserMembershipForNewUserParams) error
	CreateMeal(ctx context.Context, db DBTX, arg *CreateMealParams) error
	CreateMealPlan(ctx context.Context, db DBTX, arg *CreateMealPlanParams) error
	CreateMealPlanEvent(ctx context.Context, db DBTX, arg *CreateMealPlanEventParams) error
	CreateMealPlanGroceryListItem(ctx context.Context, db DBTX, arg *CreateMealPlanGroceryListItemParams) error
	CreateMealPlanOption(ctx context.Context, db DBTX, arg *CreateMealPlanOptionParams) error
	CreateMealPlanOptionVote(ctx context.Context, db DBTX, arg *CreateMealPlanOptionVoteParams) error
	CreateMealPlanTask(ctx context.Context, db DBTX, arg *CreateMealPlanTaskParams) error
	CreateMealRecipe(ctx context.Context, db DBTX, arg *CreateMealRecipeParams) error
	CreatePasswordResetToken(ctx context.Context, db DBTX, arg *CreatePasswordResetTokenParams) error
	CreateRecipe(ctx context.Context, db DBTX, arg *CreateRecipeParams) error
	CreateRecipeMedia(ctx context.Context, db DBTX, arg *CreateRecipeMediaParams) error
	CreateRecipePrepTask(ctx context.Context, db DBTX, arg *CreateRecipePrepTaskParams) error
	CreateRecipePrepTaskStep(ctx context.Context, db DBTX, arg *CreateRecipePrepTaskStepParams) error
	CreateRecipeStep(ctx context.Context, db DBTX, arg *CreateRecipeStepParams) error
	CreateRecipeStepIngredient(ctx context.Context, db DBTX, arg *CreateRecipeStepIngredientParams) error
	CreateRecipeStepInstrument(ctx context.Context, db DBTX, arg *CreateRecipeStepInstrumentParams) error
	CreateRecipeStepProduct(ctx context.Context, db DBTX, arg *CreateRecipeStepProductParams) error
	CreateUser(ctx context.Context, db DBTX, arg *CreateUserParams) error
	CreateValidIngredient(ctx context.Context, db DBTX, arg *CreateValidIngredientParams) error
	CreateValidIngredientMeasurementUnit(ctx context.Context, db DBTX, arg *CreateValidIngredientMeasurementUnitParams) error
	CreateValidIngredientPreparation(ctx context.Context, db DBTX, arg *CreateValidIngredientPreparationParams) error
	CreateValidInstrument(ctx context.Context, db DBTX, arg *CreateValidInstrumentParams) error
	CreateValidMeasurementConversion(ctx context.Context, db DBTX, arg *CreateValidMeasurementConversionParams) error
	CreateValidMeasurementUnit(ctx context.Context, db DBTX, arg *CreateValidMeasurementUnitParams) error
	CreateValidPreparation(ctx context.Context, db DBTX, arg *CreateValidPreparationParams) error
	CreateValidPreparationInstrument(ctx context.Context, db DBTX, arg *CreateValidPreparationInstrumentParams) error
	CreateWebhook(ctx context.Context, db DBTX, arg *CreateWebhookParams) error
	CreateWebhookTriggerEvent(ctx context.Context, db DBTX, arg *CreateWebhookTriggerEventParams) error
	FinalizeMealPlan(ctx context.Context, db DBTX, arg *FinalizeMealPlanParams) error
	FinalizeMealPlanOption(ctx context.Context, db DBTX, arg *FinalizeMealPlanOptionParams) error
	GetAPIClientByClientID(ctx context.Context, db DBTX, clientID string) (*GetAPIClientByClientIDRow, error)
	GetAPIClientByID(ctx context.Context, db DBTX, arg *GetAPIClientByIDParams) (*GetAPIClientByIDRow, error)
	GetAdminUserByUsername(ctx context.Context, db DBTX, username string) (*GetAdminUserByUsernameRow, error)
	GetDefaultHouseholdIDForUser(ctx context.Context, db DBTX, arg *GetDefaultHouseholdIDForUserParams) (string, error)
	GetExpiredAndUnresolvedMealPlans(ctx context.Context, db DBTX) ([]*GetExpiredAndUnresolvedMealPlansRow, error)
	GetFinalizedMealPlansForPlanning(ctx context.Context, db DBTX) error
	GetFinalizedMealPlansWithoutInitializedGroceryLists(ctx context.Context, db DBTX) ([]*GetFinalizedMealPlansWithoutInitializedGroceryListsRow, error)
	GetHouseholdByIDWithMemberships(ctx context.Context, db DBTX, id string) ([]*GetHouseholdByIDWithMembershipsRow, error)
	GetHouseholdInvitationByEmailAndToken(ctx context.Context, db DBTX, arg *GetHouseholdInvitationByEmailAndTokenParams) (*GetHouseholdInvitationByEmailAndTokenRow, error)
	GetHouseholdInvitationByHouseholdAndID(ctx context.Context, db DBTX, arg *GetHouseholdInvitationByHouseholdAndIDParams) (*GetHouseholdInvitationByHouseholdAndIDRow, error)
	GetHouseholdInvitationByTokenAndID(ctx context.Context, db DBTX, arg *GetHouseholdInvitationByTokenAndIDParams) (*GetHouseholdInvitationByTokenAndIDRow, error)
	GetHouseholdUserMembershipsForUser(ctx context.Context, db DBTX, belongsToUser string) ([]*GetHouseholdUserMembershipsForUserRow, error)
	GetMeal(ctx context.Context, db DBTX, id string) error
	GetMealPlan(ctx context.Context, db DBTX, arg *GetMealPlanParams) error
	GetMealPlanEvent(ctx context.Context, db DBTX, arg *GetMealPlanEventParams) error
	GetMealPlanEventsForMealPlan(ctx context.Context, db DBTX, belongsToMealPlan string) ([]*MealPlanEvents, error)
	GetMealPlanGroceryListItem(ctx context.Context, db DBTX, arg *GetMealPlanGroceryListItemParams) error
	GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, db DBTX, belongsToMealPlan string) ([]*GetMealPlanGroceryListItemsForMealPlanRow, error)
	GetMealPlanOption(ctx context.Context, db DBTX, arg *GetMealPlanOptionParams) (*GetMealPlanOptionRow, error)
	GetMealPlanOptionByID(ctx context.Context, db DBTX, id string) (*GetMealPlanOptionByIDRow, error)
	GetMealPlanOptionVote(ctx context.Context, db DBTX, arg *GetMealPlanOptionVoteParams) error
	GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, db DBTX, arg *GetMealPlanOptionVotesForMealPlanOptionParams) ([]*MealPlanOptionVotes, error)
	GetMealPlanOptionsForMealPlanEvent(ctx context.Context, db DBTX, arg *GetMealPlanOptionsForMealPlanEventParams) ([]*GetMealPlanOptionsForMealPlanEventRow, error)
	GetMealPlanPastVotingDeadline(ctx context.Context, db DBTX, arg *GetMealPlanPastVotingDeadlineParams) (*GetMealPlanPastVotingDeadlineRow, error)
	GetMealPlanTask(ctx context.Context, db DBTX, id string) error
	GetPasswordResetToken(ctx context.Context, db DBTX, token string) error
	GetRandomValidIngredient(ctx context.Context, db DBTX) error
	GetRandomValidInstrument(ctx context.Context, db DBTX) error
	GetRandomValidMeasurementUnit(ctx context.Context, db DBTX) (*GetRandomValidMeasurementUnitRow, error)
	GetRandomValidPreparation(ctx context.Context, db DBTX) (*GetRandomValidPreparationRow, error)
	GetRecipeByID(ctx context.Context, db DBTX, id string) ([]*GetRecipeByIDRow, error)
	GetRecipeByIDAndAuthor(ctx context.Context, db DBTX, arg *GetRecipeByIDAndAuthorParams) ([]*GetRecipeByIDAndAuthorRow, error)
	GetRecipeIDsForMeal(ctx context.Context, db DBTX, id string) ([]string, error)
	GetRecipeMedia(ctx context.Context, db DBTX, id string) error
	GetRecipeMediaForRecipe(ctx context.Context, db DBTX, belongsToRecipe sql.NullString) ([]*GetRecipeMediaForRecipeRow, error)
	GetRecipeMediaForRecipeStep(ctx context.Context, db DBTX, arg *GetRecipeMediaForRecipeStepParams) ([]*GetRecipeMediaForRecipeStepRow, error)
	GetRecipePrepTask(ctx context.Context, db DBTX, id string) error
	GetRecipePrepTasksForRecipe(ctx context.Context, db DBTX, id string) ([]*GetRecipePrepTasksForRecipeRow, error)
	GetRecipeStep(ctx context.Context, db DBTX, arg *GetRecipeStepParams) (*GetRecipeStepRow, error)
	GetRecipeStepByID(ctx context.Context, db DBTX, id string) (*GetRecipeStepByIDRow, error)
	GetRecipeStepIngredient(ctx context.Context, db DBTX, arg *GetRecipeStepIngredientParams) error
	GetRecipeStepIngredientForRecipe(ctx context.Context, db DBTX, id string) ([]*GetRecipeStepIngredientForRecipeRow, error)
	GetRecipeStepInstrument(ctx context.Context, db DBTX, arg *GetRecipeStepInstrumentParams) error
	GetRecipeStepInstrumentsForRecipe(ctx context.Context, db DBTX, belongsToRecipe string) ([]*GetRecipeStepInstrumentsForRecipeRow, error)
	GetRecipeStepProduct(ctx context.Context, db DBTX, arg *GetRecipeStepProductParams) (*GetRecipeStepProductRow, error)
	GetRecipeStepProductsForRecipe(ctx context.Context, db DBTX, belongsToRecipe string) ([]*GetRecipeStepProductsForRecipeRow, error)
	GetUserByEmailAddress(ctx context.Context, db DBTX, emailAddress string) (*GetUserByEmailAddressRow, error)
	GetUserByID(ctx context.Context, db DBTX, id string) (*GetUserByIDRow, error)
	GetUserByUsername(ctx context.Context, db DBTX, username string) (*GetUserByUsernameRow, error)
	GetUserWithVerifiedTwoFactorSecret(ctx context.Context, db DBTX, id string) (*GetUserWithVerifiedTwoFactorSecretRow, error)
	GetValidIngredient(ctx context.Context, db DBTX, id string) error
	GetValidIngredientMeasurementUnit(ctx context.Context, db DBTX, id string) error
	GetValidIngredientPreparation(ctx context.Context, db DBTX, id string) error
	GetValidInstrument(ctx context.Context, db DBTX, id string) error
	GetValidInstruments(ctx context.Context, db DBTX, arg *GetValidInstrumentsParams) ([]*GetValidInstrumentsRow, error)
	GetValidMeasurementConversion(ctx context.Context, db DBTX, id string) error
	GetValidMeasurementConversionsFromMeasurementUnit(ctx context.Context, db DBTX, id string) error
	GetValidMeasurementConversionsToMeasurementUnit(ctx context.Context, db DBTX, id string) ([]*GetValidMeasurementConversionsToMeasurementUnitRow, error)
	GetValidMeasurementUnit(ctx context.Context, db DBTX, id string) error
	GetValidPreparation(ctx context.Context, db DBTX, id string) error
	GetValidPreparationInstrument(ctx context.Context, db DBTX, id string) error
	GetWebhook(ctx context.Context, db DBTX, arg *GetWebhookParams) ([]*GetWebhookRow, error)
	GetWebhookTriggerEventsForWebhook(ctx context.Context, db DBTX, id string) ([]*WebhookTriggerEvents, error)
	GetWebhooks(ctx context.Context, db DBTX, arg *GetWebhooksParams) ([]*GetWebhooksRow, error)
	HouseholdInvitationExists(ctx context.Context, db DBTX, id string) error
	ListIncompleteMealPlanTaskByMealPlanOption(ctx context.Context, db DBTX, belongsToMealPlanOption string) ([]*ListIncompleteMealPlanTaskByMealPlanOptionRow, error)
	ListMealPlanTasksForMealPlan(ctx context.Context, db DBTX, id string) ([]*ListMealPlanTasksForMealPlanRow, error)
	MarkHouseholdUserMembershipAsDefaultForUser(ctx context.Context, db DBTX, arg *MarkHouseholdUserMembershipAsDefaultForUserParams) error
	MarkMealPlanAsHavingGroceryListInitialized(ctx context.Context, db DBTX, id string) error
	MarkMealPlanTasksAsCreated(ctx context.Context, db DBTX, id string) error
	MarkUserTwoFactorSecretAsVerified(ctx context.Context, db DBTX, arg *MarkUserTwoFactorSecretAsVerifiedParams) error
	MealExists(ctx context.Context, db DBTX, id string) error
	MealPlanEventExists(ctx context.Context, db DBTX, id string) error
	MealPlanExists(ctx context.Context, db DBTX, id string) error
	MealPlanGroceryListItemExists(ctx context.Context, db DBTX, id string) error
	MealPlanOptionExists(ctx context.Context, db DBTX, arg *MealPlanOptionExistsParams) error
	MealPlanOptionVoteExists(ctx context.Context, db DBTX, arg *MealPlanOptionVoteExistsParams) error
	MealPlanTaskExists(ctx context.Context, db DBTX, arg *MealPlanTaskExistsParams) error
	ModifyHouseholdUserMembershipPermissions(ctx context.Context, db DBTX, arg *ModifyHouseholdUserMembershipPermissionsParams) error
	RecipeExists(ctx context.Context, db DBTX, id string) error
	RecipeMediaExists(ctx context.Context, db DBTX, id string) error
	RecipePrepTaskExists(ctx context.Context, db DBTX, arg *RecipePrepTaskExistsParams) error
	RecipeStepExists(ctx context.Context, db DBTX, arg *RecipeStepExistsParams) error
	RecipeStepIngredientExists(ctx context.Context, db DBTX, arg *RecipeStepIngredientExistsParams) error
	RecipeStepInstrumentExists(ctx context.Context, db DBTX, arg *RecipeStepInstrumentExistsParams) error
	RecipeStepProductExists(ctx context.Context, db DBTX, arg *RecipeStepProductExistsParams) error
	RedeemPasswordResetToken(ctx context.Context, db DBTX, id string) error
	RemoveUserFromHousehold(ctx context.Context, db DBTX, arg *RemoveUserFromHouseholdParams) error
	SearchForUserByUsername(ctx context.Context, db DBTX, username string) ([]*SearchForUserByUsernameRow, error)
	SearchForValidIngredient(ctx context.Context, db DBTX, name string) error
	SearchForValidInstruments(ctx context.Context, db DBTX, name string) error
	SearchForValidMeasurementUnits(ctx context.Context, db DBTX, name string) ([]*SearchForValidMeasurementUnitsRow, error)
	SearchForValidPreparations(ctx context.Context, db DBTX, name string) ([]*SearchForValidPreparationsRow, error)
	SetHouseholdInvitationStatus(ctx context.Context, db DBTX, arg *SetHouseholdInvitationStatusParams) error
	TransferHouseholdOwnership(ctx context.Context, db DBTX, arg *TransferHouseholdOwnershipParams) error
	TransferHouseholdUserMembershipToNewUser(ctx context.Context, db DBTX, arg *TransferHouseholdUserMembershipToNewUserParams) error
	UpdateHousehold(ctx context.Context, db DBTX, arg *UpdateHouseholdParams) error
	UpdateMealPlan(ctx context.Context, db DBTX, arg *UpdateMealPlanParams) error
	UpdateMealPlanEvent(ctx context.Context, db DBTX, arg *UpdateMealPlanEventParams) error
	UpdateMealPlanGroceryListItem(ctx context.Context, db DBTX, arg *UpdateMealPlanGroceryListItemParams) error
	UpdateMealPlanOption(ctx context.Context, db DBTX, arg *UpdateMealPlanOptionParams) error
	UpdateMealPlanOptionVote(ctx context.Context, db DBTX, arg *UpdateMealPlanOptionVoteParams) error
	UpdateRecipe(ctx context.Context, db DBTX, arg *UpdateRecipeParams) error
	UpdateRecipeMedia(ctx context.Context, db DBTX, arg *UpdateRecipeMediaParams) error
	UpdateRecipePrepTask(ctx context.Context, db DBTX, arg *UpdateRecipePrepTaskParams) error
	UpdateRecipeStep(ctx context.Context, db DBTX, arg *UpdateRecipeStepParams) error
	UpdateRecipeStepIngredient(ctx context.Context, db DBTX, arg *UpdateRecipeStepIngredientParams) error
	UpdateRecipeStepInstrument(ctx context.Context, db DBTX, arg *UpdateRecipeStepInstrumentParams) error
	UpdateRecipeStepProduct(ctx context.Context, db DBTX, arg *UpdateRecipeStepProductParams) error
	UpdateUser(ctx context.Context, db DBTX, arg *UpdateUserParams) error
	UpdateUserPassword(ctx context.Context, db DBTX, arg *UpdateUserPasswordParams) error
	UpdateUserTwoFactorSecret(ctx context.Context, db DBTX, arg *UpdateUserTwoFactorSecretParams) error
	UpdateValidIngredient(ctx context.Context, db DBTX, arg *UpdateValidIngredientParams) error
	UpdateValidIngredientMeasurementUnit(ctx context.Context, db DBTX, arg *UpdateValidIngredientMeasurementUnitParams) error
	UpdateValidIngredientPreparation(ctx context.Context, db DBTX, arg *UpdateValidIngredientPreparationParams) error
	UpdateValidInstrument(ctx context.Context, db DBTX, arg *UpdateValidInstrumentParams) error
	UpdateValidMeasurementConversion(ctx context.Context, db DBTX, arg *UpdateValidMeasurementConversionParams) error
	UpdateValidMeasurementUnit(ctx context.Context, db DBTX, arg *UpdateValidMeasurementUnitParams) error
	UpdateValidPreparation(ctx context.Context, db DBTX, arg *UpdateValidPreparationParams) error
	UpdateValidPreparationInstrument(ctx context.Context, db DBTX, arg *UpdateValidPreparationInstrumentParams) error
	UserExistsWithStatus(ctx context.Context, db DBTX, arg *UserExistsWithStatusParams) (bool, error)
	UserIsMemberOfHousehold(ctx context.Context, db DBTX, arg *UserIsMemberOfHouseholdParams) (bool, error)
	ValidIngredientExists(ctx context.Context, db DBTX, id string) error
	ValidIngredientMeasurementUnitExists(ctx context.Context, db DBTX, id string) error
	ValidIngredientPreparationExists(ctx context.Context, db DBTX, id string) error
	ValidInstrumentExists(ctx context.Context, db DBTX, id string) error
	ValidMeasurementConversionExists(ctx context.Context, db DBTX, id string) error
	ValidMeasurementUnitExists(ctx context.Context, db DBTX, id string) error
	ValidPreparationExists(ctx context.Context, db DBTX, id string) error
	ValidPreparationInstrumentExists(ctx context.Context, db DBTX, id string) error
	WebhookExists(ctx context.Context, db DBTX, arg *WebhookExistsParams) (bool, error)
}

var _ Querier = (*Queries)(nil)
