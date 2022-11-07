package postgres

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
)

var _ generated.Querier = (*mockGeneratedQuerier)(nil)

type mockGeneratedQuerier struct {
	mock.Mock
}

func (m *mockGeneratedQuerier) ArchiveAPIClient(ctx context.Context, db generated.DBTX, arg *generated.ArchiveAPIClientParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveHousehold(ctx context.Context, db generated.DBTX, arg *generated.ArchiveHouseholdParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveHouseholdUserMembershipForUser(ctx context.Context, db generated.DBTX, belongsToUser string) error {
	return m.Called(ctx, db, belongsToUser).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveMeal(ctx context.Context, db generated.DBTX, arg *generated.ArchiveMealParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveMealPlan(ctx context.Context, db generated.DBTX, arg *generated.ArchiveMealPlanParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveMealPlanEvent(ctx context.Context, db generated.DBTX, arg *generated.ArchiveMealPlanEventParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveMealPlanGroceryListItem(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveMealPlanOption(ctx context.Context, db generated.DBTX, arg *generated.ArchiveMealPlanOptionParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveMealPlanOptionVote(ctx context.Context, db generated.DBTX, arg *generated.ArchiveMealPlanOptionVoteParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveRecipe(ctx context.Context, db generated.DBTX, arg *generated.ArchiveRecipeParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveRecipeMedia(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveRecipePrepTask(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveRecipeStep(ctx context.Context, db generated.DBTX, arg *generated.ArchiveRecipeStepParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveRecipeStepIngredient(ctx context.Context, db generated.DBTX, arg *generated.ArchiveRecipeStepIngredientParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveRecipeStepInstrument(ctx context.Context, db generated.DBTX, arg *generated.ArchiveRecipeStepInstrumentParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveRecipeStepProduct(ctx context.Context, db generated.DBTX, arg *generated.ArchiveRecipeStepProductParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveUser(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveValidIngredient(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveValidIngredientMeasurementUnit(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveValidIngredientPreparation(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveValidInstrument(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveValidMeasurementConversion(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveValidMeasurementUnit(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveValidPreparation(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveValidPreparationInstrument(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ArchiveWebhook(ctx context.Context, db generated.DBTX, arg *generated.ArchiveWebhookParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) AttachHouseholdInvitationsToUser(ctx context.Context, db generated.DBTX, arg *generated.AttachHouseholdInvitationsToUserParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ChangeMealPlanTaskStatus(ctx context.Context, db generated.DBTX, arg *generated.ChangeMealPlanTaskStatusParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateAPIClient(ctx context.Context, db generated.DBTX, arg *generated.CreateAPIClientParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateHousehold(ctx context.Context, db generated.DBTX, arg *generated.CreateHouseholdParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateHouseholdInvitation(ctx context.Context, db generated.DBTX, arg *generated.CreateHouseholdInvitationParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateHouseholdUserMembership(ctx context.Context, db generated.DBTX, arg *generated.CreateHouseholdUserMembershipParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateHouseholdUserMembershipForNewUser(ctx context.Context, db generated.DBTX, arg *generated.CreateHouseholdUserMembershipForNewUserParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateMeal(ctx context.Context, db generated.DBTX, arg *generated.CreateMealParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateMealPlan(ctx context.Context, db generated.DBTX, arg *generated.CreateMealPlanParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateMealPlanEvent(ctx context.Context, db generated.DBTX, arg *generated.CreateMealPlanEventParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateMealPlanGroceryListItem(ctx context.Context, db generated.DBTX, arg *generated.CreateMealPlanGroceryListItemParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateMealPlanOption(ctx context.Context, db generated.DBTX, arg *generated.CreateMealPlanOptionParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateMealPlanOptionVote(ctx context.Context, db generated.DBTX, arg *generated.CreateMealPlanOptionVoteParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateMealPlanTask(ctx context.Context, db generated.DBTX, arg *generated.CreateMealPlanTaskParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateMealRecipe(ctx context.Context, db generated.DBTX, arg *generated.CreateMealRecipeParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreatePasswordResetToken(ctx context.Context, db generated.DBTX, arg *generated.CreatePasswordResetTokenParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateRecipe(ctx context.Context, db generated.DBTX, arg *generated.CreateRecipeParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateRecipeMedia(ctx context.Context, db generated.DBTX, arg *generated.CreateRecipeMediaParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateRecipePrepTask(ctx context.Context, db generated.DBTX, arg *generated.CreateRecipePrepTaskParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateRecipePrepTaskStep(ctx context.Context, db generated.DBTX, arg *generated.CreateRecipePrepTaskStepParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateRecipeStep(ctx context.Context, db generated.DBTX, arg *generated.CreateRecipeStepParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateRecipeStepIngredient(ctx context.Context, db generated.DBTX, arg *generated.CreateRecipeStepIngredientParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateRecipeStepInstrument(ctx context.Context, db generated.DBTX, arg *generated.CreateRecipeStepInstrumentParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateRecipeStepProduct(ctx context.Context, db generated.DBTX, arg *generated.CreateRecipeStepProductParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateUser(ctx context.Context, db generated.DBTX, arg *generated.CreateUserParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateValidIngredient(ctx context.Context, db generated.DBTX, arg *generated.CreateValidIngredientParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateValidIngredientMeasurementUnit(ctx context.Context, db generated.DBTX, arg *generated.CreateValidIngredientMeasurementUnitParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateValidIngredientPreparation(ctx context.Context, db generated.DBTX, arg *generated.CreateValidIngredientPreparationParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateValidInstrument(ctx context.Context, db generated.DBTX, arg *generated.CreateValidInstrumentParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateValidMeasurementConversion(ctx context.Context, db generated.DBTX, arg *generated.CreateValidMeasurementConversionParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateValidMeasurementUnit(ctx context.Context, db generated.DBTX, arg *generated.CreateValidMeasurementUnitParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateValidPreparation(ctx context.Context, db generated.DBTX, arg *generated.CreateValidPreparationParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateValidPreparationInstrument(ctx context.Context, db generated.DBTX, arg *generated.CreateValidPreparationInstrumentParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateWebhook(ctx context.Context, db generated.DBTX, arg *generated.CreateWebhookParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) CreateWebhookTriggerEvent(ctx context.Context, db generated.DBTX, arg *generated.CreateWebhookTriggerEventParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) FinalizeMealPlan(ctx context.Context, db generated.DBTX, arg *generated.FinalizeMealPlanParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) FinalizeMealPlanOption(ctx context.Context, db generated.DBTX, arg *generated.FinalizeMealPlanOptionParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) GetAPIClientByClientID(ctx context.Context, db generated.DBTX, clientID string) (*generated.GetAPIClientByClientIDRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetAPIClientByID(ctx context.Context, db generated.DBTX, arg *generated.GetAPIClientByIDParams) (*generated.GetAPIClientByIDRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetAdminUserByUsername(ctx context.Context, db generated.DBTX, username string) (*generated.GetAdminUserByUsernameRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetDefaultHouseholdIDForUser(ctx context.Context, db generated.DBTX, arg *generated.GetDefaultHouseholdIDForUserParams) (string, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetExpiredAndUnresolvedMealPlans(ctx context.Context, db generated.DBTX) ([]*generated.GetExpiredAndUnresolvedMealPlansRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetFinalizedMealPlansForPlanning(ctx context.Context, db generated.DBTX) error {
	return m.Called(ctx, db).Error(0)
}

func (m *mockGeneratedQuerier) GetFinalizedMealPlansWithoutInitializedGroceryLists(ctx context.Context, db generated.DBTX) ([]*generated.GetFinalizedMealPlansWithoutInitializedGroceryListsRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetHouseholdByIDWithMemberships(ctx context.Context, db generated.DBTX, id string) ([]*generated.GetHouseholdByIDWithMembershipsRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetHouseholdInvitationByEmailAndToken(ctx context.Context, db generated.DBTX, arg *generated.GetHouseholdInvitationByEmailAndTokenParams) (*generated.GetHouseholdInvitationByEmailAndTokenRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetHouseholdInvitationByHouseholdAndID(ctx context.Context, db generated.DBTX, arg *generated.GetHouseholdInvitationByHouseholdAndIDParams) (*generated.GetHouseholdInvitationByHouseholdAndIDRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetHouseholdInvitationByTokenAndID(ctx context.Context, db generated.DBTX, arg *generated.GetHouseholdInvitationByTokenAndIDParams) (*generated.GetHouseholdInvitationByTokenAndIDRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetHouseholdUserMembershipsForUser(ctx context.Context, db generated.DBTX, belongsToUser string) ([]*generated.GetHouseholdUserMembershipsForUserRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetMeal(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetMealPlan(ctx context.Context, db generated.DBTX, arg *generated.GetMealPlanParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) GetMealPlanEvent(ctx context.Context, db generated.DBTX, arg *generated.GetMealPlanEventParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) GetMealPlanEventsForMealPlan(ctx context.Context, db generated.DBTX, belongsToMealPlan string) ([]*generated.MealPlanEvents, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetMealPlanGroceryListItem(ctx context.Context, db generated.DBTX, arg *generated.GetMealPlanGroceryListItemParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, db generated.DBTX, belongsToMealPlan string) ([]*generated.GetMealPlanGroceryListItemsForMealPlanRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetMealPlanOption(ctx context.Context, db generated.DBTX, arg *generated.GetMealPlanOptionParams) (*generated.GetMealPlanOptionRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetMealPlanOptionByID(ctx context.Context, db generated.DBTX, id string) (*generated.GetMealPlanOptionByIDRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetMealPlanOptionVote(ctx context.Context, db generated.DBTX, arg *generated.GetMealPlanOptionVoteParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, db generated.DBTX, arg *generated.GetMealPlanOptionVotesForMealPlanOptionParams) ([]*generated.MealPlanOptionVotes, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetMealPlanOptionsForMealPlanEvent(ctx context.Context, db generated.DBTX, arg *generated.GetMealPlanOptionsForMealPlanEventParams) ([]*generated.GetMealPlanOptionsForMealPlanEventRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetMealPlanPastVotingDeadline(ctx context.Context, db generated.DBTX, arg *generated.GetMealPlanPastVotingDeadlineParams) (*generated.GetMealPlanPastVotingDeadlineRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetMealPlanTask(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetPasswordResetToken(ctx context.Context, db generated.DBTX, token string) error {
	return m.Called(ctx, db, token).Error(0)
}

func (m *mockGeneratedQuerier) GetRandomValidIngredient(ctx context.Context, db generated.DBTX) error {
	return m.Called(ctx, db).Error(0)
}

func (m *mockGeneratedQuerier) GetRandomValidInstrument(ctx context.Context, db generated.DBTX) error {
	return m.Called(ctx, db).Error(0)
}

func (m *mockGeneratedQuerier) GetRandomValidMeasurementUnit(ctx context.Context, db generated.DBTX) (*generated.GetRandomValidMeasurementUnitRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRandomValidPreparation(ctx context.Context, db generated.DBTX) (*generated.GetRandomValidPreparationRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeByID(ctx context.Context, db generated.DBTX, id string) ([]*generated.GetRecipeByIDRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeByIDAndAuthor(ctx context.Context, db generated.DBTX, arg *generated.GetRecipeByIDAndAuthorParams) ([]*generated.GetRecipeByIDAndAuthorRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeIDsForMeal(ctx context.Context, db generated.DBTX, id string) ([]string, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeMedia(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetRecipeMediaForRecipe(ctx context.Context, db generated.DBTX, belongsToRecipe sql.NullString) ([]*generated.GetRecipeMediaForRecipeRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeMediaForRecipeStep(ctx context.Context, db generated.DBTX, arg *generated.GetRecipeMediaForRecipeStepParams) ([]*generated.GetRecipeMediaForRecipeStepRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipePrepTask(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetRecipePrepTasksForRecipe(ctx context.Context, db generated.DBTX, id string) ([]*generated.GetRecipePrepTasksForRecipeRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeStep(ctx context.Context, db generated.DBTX, arg *generated.GetRecipeStepParams) (*generated.GetRecipeStepRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeStepByID(ctx context.Context, db generated.DBTX, id string) (*generated.GetRecipeStepByIDRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeStepIngredient(ctx context.Context, db generated.DBTX, arg *generated.GetRecipeStepIngredientParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) GetRecipeStepIngredientForRecipe(ctx context.Context, db generated.DBTX, id string) ([]*generated.GetRecipeStepIngredientForRecipeRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeStepInstrument(ctx context.Context, db generated.DBTX, arg *generated.GetRecipeStepInstrumentParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) GetRecipeStepInstrumentsForRecipe(ctx context.Context, db generated.DBTX, belongsToRecipe string) ([]*generated.GetRecipeStepInstrumentsForRecipeRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeStepProduct(ctx context.Context, db generated.DBTX, arg *generated.GetRecipeStepProductParams) (*generated.GetRecipeStepProductRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetRecipeStepProductsForRecipe(ctx context.Context, db generated.DBTX, belongsToRecipe string) ([]*generated.GetRecipeStepProductsForRecipeRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetUserByEmailAddress(ctx context.Context, db generated.DBTX, emailAddress string) (*generated.GetUserByEmailAddressRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetUserByID(ctx context.Context, db generated.DBTX, id string) (*generated.GetUserByIDRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetUserByUsername(ctx context.Context, db generated.DBTX, username string) (*generated.GetUserByUsernameRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetUserWithVerifiedTwoFactorSecret(ctx context.Context, db generated.DBTX, id string) (*generated.GetUserWithVerifiedTwoFactorSecretRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetValidIngredient(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetValidIngredientMeasurementUnit(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetValidIngredientPreparation(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetValidInstrument(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetValidInstruments(ctx context.Context, db generated.DBTX, arg *generated.GetValidInstrumentsParams) ([]*generated.GetValidInstrumentsRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetValidMeasurementConversion(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetValidMeasurementConversionsFromMeasurementUnit(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetValidMeasurementConversionsToMeasurementUnit(ctx context.Context, db generated.DBTX, id string) ([]*generated.GetValidMeasurementConversionsToMeasurementUnitRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) GetValidMeasurementUnit(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetValidPreparation(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetValidPreparationInstrument(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) GetWebhook(ctx context.Context, db generated.DBTX, arg *generated.GetWebhookParams) (*generated.GetWebhookRow, error) {
	rv := m.Called(ctx, db, arg)

	return rv.Get(0).(*generated.GetWebhookRow), rv.Error(1)
}

func (m *mockGeneratedQuerier) GetWebhooks(ctx context.Context, db generated.DBTX, arg *generated.GetWebhooksParams) ([]*generated.GetWebhooksRow, error) {
	rv := m.Called(ctx, db, arg)

	return rv.Get(0).([]*generated.GetWebhooksRow), rv.Error(1)
}

func (m *mockGeneratedQuerier) GetWebhookTriggerEventsForWebhook(ctx context.Context, db generated.DBTX, id string) ([]*generated.WebhookTriggerEvents, error) {
	rv := m.Called(ctx, db, id)

	return rv.Get(0).([]*generated.WebhookTriggerEvents), rv.Error(1)
}

func (m *mockGeneratedQuerier) HouseholdInvitationExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ListIncompleteMealPlanTaskByMealPlanOption(ctx context.Context, db generated.DBTX, belongsToMealPlanOption string) ([]*generated.ListIncompleteMealPlanTaskByMealPlanOptionRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) ListMealPlanTasksForMealPlan(ctx context.Context, db generated.DBTX, id string) ([]*generated.ListMealPlanTasksForMealPlanRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) MarkHouseholdUserMembershipAsDefaultForUser(ctx context.Context, db generated.DBTX, arg *generated.MarkHouseholdUserMembershipAsDefaultForUserParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) MarkMealPlanAsHavingGroceryListInitialized(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) MarkMealPlanTasksAsCreated(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) MarkUserTwoFactorSecretAsVerified(ctx context.Context, db generated.DBTX, arg *generated.MarkUserTwoFactorSecretAsVerifiedParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) MealExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) MealPlanEventExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) MealPlanExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) MealPlanGroceryListItemExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) MealPlanOptionExists(ctx context.Context, db generated.DBTX, arg *generated.MealPlanOptionExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) MealPlanOptionVoteExists(ctx context.Context, db generated.DBTX, arg *generated.MealPlanOptionVoteExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) MealPlanTaskExists(ctx context.Context, db generated.DBTX, arg *generated.MealPlanTaskExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) ModifyHouseholdUserMembershipPermissions(ctx context.Context, db generated.DBTX, arg *generated.ModifyHouseholdUserMembershipPermissionsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) RecipeExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) RecipeMediaExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) RecipePrepTaskExists(ctx context.Context, db generated.DBTX, arg *generated.RecipePrepTaskExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) RecipeStepExists(ctx context.Context, db generated.DBTX, arg *generated.RecipeStepExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) RecipeStepIngredientExists(ctx context.Context, db generated.DBTX, arg *generated.RecipeStepIngredientExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) RecipeStepInstrumentExists(ctx context.Context, db generated.DBTX, arg *generated.RecipeStepInstrumentExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) RecipeStepProductExists(ctx context.Context, db generated.DBTX, arg *generated.RecipeStepProductExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) RedeemPasswordResetToken(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) RemoveUserFromHousehold(ctx context.Context, db generated.DBTX, arg *generated.RemoveUserFromHouseholdParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) SearchForUserByUsername(ctx context.Context, db generated.DBTX, username string) ([]*generated.SearchForUserByUsernameRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) SearchForValidIngredient(ctx context.Context, db generated.DBTX, name string) error {
	return m.Called(ctx, db, name).Error(0)
}

func (m *mockGeneratedQuerier) SearchForValidInstruments(ctx context.Context, db generated.DBTX, name string) error {
	return m.Called(ctx, db, name).Error(0)
}

func (m *mockGeneratedQuerier) SearchForValidMeasurementUnits(ctx context.Context, db generated.DBTX, name string) ([]*generated.SearchForValidMeasurementUnitsRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) SearchForValidPreparations(ctx context.Context, db generated.DBTX, name string) ([]*generated.SearchForValidPreparationsRow, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) SetHouseholdInvitationStatus(ctx context.Context, db generated.DBTX, arg *generated.SetHouseholdInvitationStatusParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) TransferHouseholdOwnership(ctx context.Context, db generated.DBTX, arg *generated.TransferHouseholdOwnershipParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) TransferHouseholdUserMembershipToNewUser(ctx context.Context, db generated.DBTX, arg *generated.TransferHouseholdUserMembershipToNewUserParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateHousehold(ctx context.Context, db generated.DBTX, arg *generated.UpdateHouseholdParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateMealPlan(ctx context.Context, db generated.DBTX, arg *generated.UpdateMealPlanParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateMealPlanEvent(ctx context.Context, db generated.DBTX, arg *generated.UpdateMealPlanEventParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateMealPlanGroceryListItem(ctx context.Context, db generated.DBTX, arg *generated.UpdateMealPlanGroceryListItemParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateMealPlanOption(ctx context.Context, db generated.DBTX, arg *generated.UpdateMealPlanOptionParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateMealPlanOptionVote(ctx context.Context, db generated.DBTX, arg *generated.UpdateMealPlanOptionVoteParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateRecipe(ctx context.Context, db generated.DBTX, arg *generated.UpdateRecipeParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateRecipeMedia(ctx context.Context, db generated.DBTX, arg *generated.UpdateRecipeMediaParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateRecipePrepTask(ctx context.Context, db generated.DBTX, arg *generated.UpdateRecipePrepTaskParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateRecipeStep(ctx context.Context, db generated.DBTX, arg *generated.UpdateRecipeStepParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateRecipeStepIngredient(ctx context.Context, db generated.DBTX, arg *generated.UpdateRecipeStepIngredientParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateRecipeStepInstrument(ctx context.Context, db generated.DBTX, arg *generated.UpdateRecipeStepInstrumentParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateRecipeStepProduct(ctx context.Context, db generated.DBTX, arg *generated.UpdateRecipeStepProductParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateUser(ctx context.Context, db generated.DBTX, arg *generated.UpdateUserParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateUserPassword(ctx context.Context, db generated.DBTX, arg *generated.UpdateUserPasswordParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateUserTwoFactorSecret(ctx context.Context, db generated.DBTX, arg *generated.UpdateUserTwoFactorSecretParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateValidIngredient(ctx context.Context, db generated.DBTX, arg *generated.UpdateValidIngredientParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateValidIngredientMeasurementUnit(ctx context.Context, db generated.DBTX, arg *generated.UpdateValidIngredientMeasurementUnitParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateValidIngredientPreparation(ctx context.Context, db generated.DBTX, arg *generated.UpdateValidIngredientPreparationParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateValidInstrument(ctx context.Context, db generated.DBTX, arg *generated.UpdateValidInstrumentParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateValidMeasurementConversion(ctx context.Context, db generated.DBTX, arg *generated.UpdateValidMeasurementConversionParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateValidMeasurementUnit(ctx context.Context, db generated.DBTX, arg *generated.UpdateValidMeasurementUnitParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateValidPreparation(ctx context.Context, db generated.DBTX, arg *generated.UpdateValidPreparationParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UpdateValidPreparationInstrument(ctx context.Context, db generated.DBTX, arg *generated.UpdateValidPreparationInstrumentParams) error {
	return m.Called(ctx, db, arg).Error(0)
}

func (m *mockGeneratedQuerier) UserExistsWithStatus(ctx context.Context, db generated.DBTX, arg *generated.UserExistsWithStatusParams) (bool, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) UserIsMemberOfHousehold(ctx context.Context, db generated.DBTX, arg *generated.UserIsMemberOfHouseholdParams) (bool, error) {
	panic("implement me")
}

func (m *mockGeneratedQuerier) ValidIngredientExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ValidIngredientMeasurementUnitExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ValidIngredientPreparationExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ValidInstrumentExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ValidMeasurementConversionExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ValidMeasurementUnitExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ValidPreparationExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) ValidPreparationInstrumentExists(ctx context.Context, db generated.DBTX, id string) error {
	return m.Called(ctx, db, id).Error(0)
}

func (m *mockGeneratedQuerier) WebhookExists(ctx context.Context, db generated.DBTX, arg *generated.WebhookExistsParams) error {
	return m.Called(ctx, db, arg).Error(0)
}
