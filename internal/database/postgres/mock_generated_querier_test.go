package postgres

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
)

var _ generated.Querier = (*mockQuerier)(nil)

type mockQuerier struct {
	mock.Mock
}

// AddUserToHouseholdDuringCreation is a mock function that implements our interface.
func (m *mockQuerier) AddUserToHouseholdDuringCreation(ctx context.Context, arg *generated.AddUserToHouseholdDuringCreationParams) error {
	return m.Called(ctx, arg).Error(0)
}

// AddUserToHouseholdQuery is a mock function that implements our interface.
func (m *mockQuerier) AddUserToHouseholdQuery(ctx context.Context, arg *generated.AddUserToHouseholdQueryParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveAPIClient is a mock function that implements our interface.
func (m *mockQuerier) ArchiveAPIClient(ctx context.Context, arg *generated.ArchiveAPIClientParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveHousehold is a mock function that implements our interface.
func (m *mockQuerier) ArchiveHousehold(ctx context.Context, arg *generated.ArchiveHouseholdParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveMeal is a mock function that implements our interface.
func (m *mockQuerier) ArchiveMeal(ctx context.Context, arg *generated.ArchiveMealParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveMealPlan is a mock function that implements our interface.
func (m *mockQuerier) ArchiveMealPlan(ctx context.Context, arg *generated.ArchiveMealPlanParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveMealPlanOption is a mock function that implements our interface.
func (m *mockQuerier) ArchiveMealPlanOption(ctx context.Context, arg *generated.ArchiveMealPlanOptionParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveMealPlanOptionVote is a mock function that implements our interface.
func (m *mockQuerier) ArchiveMealPlanOptionVote(ctx context.Context, arg *generated.ArchiveMealPlanOptionVoteParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveMemberships is a mock function that implements our interface.
func (m *mockQuerier) ArchiveMemberships(ctx context.Context, belongsToUser string) error {
	return m.Called(ctx, belongsToUser).Error(0)
}

// ArchiveRecipe is a mock function that implements our interface.
func (m *mockQuerier) ArchiveRecipe(ctx context.Context, arg *generated.ArchiveRecipeParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveRecipeStep is a mock function that implements our interface.
func (m *mockQuerier) ArchiveRecipeStep(ctx context.Context, arg *generated.ArchiveRecipeStepParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveRecipeStepIngredient is a mock function that implements our interface.
func (m *mockQuerier) ArchiveRecipeStepIngredient(ctx context.Context, arg *generated.ArchiveRecipeStepIngredientParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveRecipeStepInstrument is a mock function that implements our interface.
func (m *mockQuerier) ArchiveRecipeStepInstrument(ctx context.Context, arg *generated.ArchiveRecipeStepInstrumentParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveRecipeStepProduct is a mock function that implements our interface.
func (m *mockQuerier) ArchiveRecipeStepProduct(ctx context.Context, arg *generated.ArchiveRecipeStepProductParams) error {
	return m.Called(ctx, arg).Error(0)
}

// ArchiveUser is a mock function that implements our interface.
func (m *mockQuerier) ArchiveUser(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ArchiveValidIngredient is a mock function that implements our interface.
func (m *mockQuerier) ArchiveValidIngredient(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ArchiveValidIngredientMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) ArchiveValidIngredientMeasurementUnit(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ArchiveValidIngredientPreparation is a mock function that implements our interface.
func (m *mockQuerier) ArchiveValidIngredientPreparation(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ArchiveValidInstrument is a mock function that implements our interface.
func (m *mockQuerier) ArchiveValidInstrument(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ArchiveValidMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) ArchiveValidMeasurementUnit(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ArchiveValidPreparation is a mock function that implements our interface.
func (m *mockQuerier) ArchiveValidPreparation(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ArchiveValidPreparationInstrument is a mock function that implements our interface.
func (m *mockQuerier) ArchiveValidPreparationInstrument(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// ArchiveWebhook is a mock function that implements our interface.
func (m *mockQuerier) ArchiveWebhook(ctx context.Context, arg *generated.ArchiveWebhookParams) error {
	return m.Called(ctx, arg).Error(0)
}

// AttachInvitationsToUserID is a mock function that implements our interface.
func (m *mockQuerier) AttachInvitationsToUserID(ctx context.Context, arg *generated.AttachInvitationsToUserIDParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateAPIClient is a mock function that implements our interface.
func (m *mockQuerier) CreateAPIClient(ctx context.Context, arg *generated.CreateAPIClientParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateHousehold is a mock function that implements our interface.
func (m *mockQuerier) CreateHousehold(ctx context.Context, arg *generated.CreateHouseholdParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateHouseholdInvitation is a mock function that implements our interface.
func (m *mockQuerier) CreateHouseholdInvitation(ctx context.Context, arg *generated.CreateHouseholdInvitationParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateHouseholdMembershipForNewUser is a mock function that implements our interface.
func (m *mockQuerier) CreateHouseholdMembershipForNewUser(ctx context.Context, arg *generated.CreateHouseholdMembershipForNewUserParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateMeal is a mock function that implements our interface.
func (m *mockQuerier) CreateMeal(ctx context.Context, arg *generated.CreateMealParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateMealPlan is a mock function that implements our interface.
func (m *mockQuerier) CreateMealPlan(ctx context.Context, arg *generated.CreateMealPlanParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateMealRecipe is a mock function that implements our interface.
func (m *mockQuerier) CreateMealRecipe(ctx context.Context, arg *generated.CreateMealRecipeParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreatePasswordResetToken is a mock function that implements our interface.
func (m *mockQuerier) CreatePasswordResetToken(ctx context.Context, arg *generated.CreatePasswordResetTokenParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateRecipe is a mock function that implements our interface.
func (m *mockQuerier) CreateRecipe(ctx context.Context, arg *generated.CreateRecipeParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateRecipeStep is a mock function that implements our interface.
func (m *mockQuerier) CreateRecipeStep(ctx context.Context, arg *generated.CreateRecipeStepParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateRecipeStepIngredient is a mock function that implements our interface.
func (m *mockQuerier) CreateRecipeStepIngredient(ctx context.Context, arg *generated.CreateRecipeStepIngredientParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateRecipeStepInstrument is a mock function that implements our interface.
func (m *mockQuerier) CreateRecipeStepInstrument(ctx context.Context, arg *generated.CreateRecipeStepInstrumentParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateRecipeStepProduct is a mock function that implements our interface.
func (m *mockQuerier) CreateRecipeStepProduct(ctx context.Context, arg *generated.CreateRecipeStepProductParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateUser is a mock function that implements our interface.
func (m *mockQuerier) CreateUser(ctx context.Context, arg *generated.CreateUserParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateValidIngredient is a mock function that implements our interface.
func (m *mockQuerier) CreateValidIngredient(ctx context.Context, arg *generated.CreateValidIngredientParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateValidIngredientMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) CreateValidIngredientMeasurementUnit(ctx context.Context, arg *generated.CreateValidIngredientMeasurementUnitParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateValidIngredientPreparation is a mock function that implements our interface.
func (m *mockQuerier) CreateValidIngredientPreparation(ctx context.Context, arg *generated.CreateValidIngredientPreparationParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateValidInstrument is a mock function that implements our interface.
func (m *mockQuerier) CreateValidInstrument(ctx context.Context, arg *generated.CreateValidInstrumentParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateValidMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) CreateValidMeasurementUnit(ctx context.Context, arg *generated.CreateValidMeasurementUnitParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateValidPreparation is a mock function that implements our interface.
func (m *mockQuerier) CreateValidPreparation(ctx context.Context, arg *generated.CreateValidPreparationParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateValidPreparationInstrument is a mock function that implements our interface.
func (m *mockQuerier) CreateValidPreparationInstrument(ctx context.Context, arg *generated.CreateValidPreparationInstrumentParams) error {
	return m.Called(ctx, arg).Error(0)
}

// CreateWebhook is a mock function that implements our interface.
func (m *mockQuerier) CreateWebhook(ctx context.Context, arg *generated.CreateWebhookParams) error {
	return m.Called(ctx, arg).Error(0)
}

// FinalizeMealPlan is a mock function that implements our interface.
func (m *mockQuerier) FinalizeMealPlan(ctx context.Context, arg *generated.FinalizeMealPlanParams) error {
	return m.Called(ctx, arg).Error(0)
}

// FinalizeMealPlanOption is a mock function that implements our interface.
func (m *mockQuerier) FinalizeMealPlanOption(ctx context.Context, arg *generated.FinalizeMealPlanOptionParams) error {
	return m.Called(ctx, arg).Error(0)
}

// GetAPIClientByClientID is a mock function that implements our interface.
func (m *mockQuerier) GetAPIClientByClientID(ctx context.Context, clientID string) (*generated.GetAPIClientByClientIDRow, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).(*generated.GetAPIClientByClientIDRow), args.Error(1)
}

// GetAPIClientByDatabaseID is a mock function that implements our interface.
func (m *mockQuerier) GetAPIClientByDatabaseID(ctx context.Context, arg *generated.GetAPIClientByDatabaseIDParams) (*generated.GetAPIClientByDatabaseIDRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*generated.GetAPIClientByDatabaseIDRow), args.Error(1)
}

// GetAdminUserByUsername is a mock function that implements our interface.
func (m *mockQuerier) GetAdminUserByUsername(ctx context.Context, username string) (*generated.GetAdminUserByUsernameRow, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*generated.GetAdminUserByUsernameRow), args.Error(1)
}

// GetAllHouseholdInvitationsCount is a mock function that implements our interface.
func (m *mockQuerier) GetAllHouseholdInvitationsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetAllHouseholdsCount is a mock function that implements our interface.
func (m *mockQuerier) GetAllHouseholdsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetAllUsersCount is a mock function that implements our interface.
func (m *mockQuerier) GetAllUsersCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetAllWebhooksCount is a mock function that implements our interface.
func (m *mockQuerier) GetAllWebhooksCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetDefaultHouseholdIDForUserQuery is a mock function that implements our interface.
func (m *mockQuerier) GetDefaultHouseholdIDForUserQuery(ctx context.Context, arg *generated.GetDefaultHouseholdIDForUserQueryParams) (string, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(string), args.Error(1)
}

// GetExpiredAndUnresolvedMealPlanIDs is a mock function that implements our interface.
func (m *mockQuerier) GetExpiredAndUnresolvedMealPlanIDs(ctx context.Context) ([]*generated.MealPlans, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*generated.MealPlans), args.Error(1)
}

// GetHousehold is a mock function that implements our interface.
func (m *mockQuerier) GetHousehold(ctx context.Context, id string) ([]*generated.GetHouseholdRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]*generated.GetHouseholdRow), args.Error(1)
}

// GetHouseholdByID is a mock function that implements our interface.
func (m *mockQuerier) GetHouseholdByID(ctx context.Context, id string) ([]*generated.GetHouseholdByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]*generated.GetHouseholdByIDRow), args.Error(1)
}

// GetHouseholdInvitationByEmailAndToken is a mock function that implements our interface.
func (m *mockQuerier) GetHouseholdInvitationByEmailAndToken(ctx context.Context, arg *generated.GetHouseholdInvitationByEmailAndTokenParams) (*generated.GetHouseholdInvitationByEmailAndTokenRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*generated.GetHouseholdInvitationByEmailAndTokenRow), args.Error(1)
}

// GetHouseholdInvitationByHouseholdAndID is a mock function that implements our interface.
func (m *mockQuerier) GetHouseholdInvitationByHouseholdAndID(ctx context.Context, arg *generated.GetHouseholdInvitationByHouseholdAndIDParams) ([]*generated.GetHouseholdInvitationByHouseholdAndIDRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetHouseholdInvitationByHouseholdAndIDRow), args.Error(1)
}

// GetHouseholdInvitationByTokenAndID is a mock function that implements our interface.
func (m *mockQuerier) GetHouseholdInvitationByTokenAndID(ctx context.Context, arg *generated.GetHouseholdInvitationByTokenAndIDParams) (*generated.GetHouseholdInvitationByTokenAndIDRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*generated.GetHouseholdInvitationByTokenAndIDRow), args.Error(1)
}

// GetHouseholdMembershipsForUserQuery is a mock function that implements our interface.
func (m *mockQuerier) GetHouseholdMembershipsForUserQuery(ctx context.Context, belongsToUser string) ([]*generated.GetHouseholdMembershipsForUserQueryRow, error) {
	args := m.Called(ctx, belongsToUser)
	return args.Get(0).([]*generated.GetHouseholdMembershipsForUserQueryRow), args.Error(1)
}

// GetMealByID is a mock function that implements our interface.
func (m *mockQuerier) GetMealByID(ctx context.Context, id string) (*generated.GetMealByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetMealByIDRow), args.Error(1)
}

// GetMealPlan is a mock function that implements our interface.
func (m *mockQuerier) GetMealPlan(ctx context.Context, arg *generated.GetMealPlanParams) ([]*generated.GetMealPlanRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetMealPlanRow), args.Error(1)
}

// GetMealPlanOptionQuery is a mock function that implements our interface.
func (m *mockQuerier) GetMealPlanOptionQuery(ctx context.Context, arg *generated.GetMealPlanOptionQueryParams) (*generated.GetMealPlanOptionQueryRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*generated.GetMealPlanOptionQueryRow), args.Error(1)
}

// GetMealPlanOptionVote is a mock function that implements our interface.
func (m *mockQuerier) GetMealPlanOptionVote(ctx context.Context, arg *generated.GetMealPlanOptionVoteParams) (*generated.MealPlanOptionVotes, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*generated.MealPlanOptionVotes), args.Error(1)
}

// GetMealPlanPastVotingDeadline is a mock function that implements our interface.
func (m *mockQuerier) GetMealPlanPastVotingDeadline(ctx context.Context, arg *generated.GetMealPlanPastVotingDeadlineParams) ([]*generated.GetMealPlanPastVotingDeadlineRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetMealPlanPastVotingDeadlineRow), args.Error(1)
}

// GetPasswordResetToken is a mock function that implements our interface.
func (m *mockQuerier) GetPasswordResetToken(ctx context.Context, token string) (*generated.PasswordResetTokens, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*generated.PasswordResetTokens), args.Error(1)
}

// GetRandomValidIngredient is a mock function that implements our interface.
func (m *mockQuerier) GetRandomValidIngredient(ctx context.Context) (*generated.GetRandomValidIngredientRow, error) {
	args := m.Called(ctx)
	return args.Get(0).(*generated.GetRandomValidIngredientRow), args.Error(1)
}

// GetRandomValidInstrument is a mock function that implements our interface.
func (m *mockQuerier) GetRandomValidInstrument(ctx context.Context) (*generated.GetRandomValidInstrumentRow, error) {
	args := m.Called(ctx)
	return args.Get(0).(*generated.GetRandomValidInstrumentRow), args.Error(1)
}

// GetRandomValidMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) GetRandomValidMeasurementUnit(ctx context.Context) (*generated.GetRandomValidMeasurementUnitRow, error) {
	args := m.Called(ctx)
	return args.Get(0).(*generated.GetRandomValidMeasurementUnitRow), args.Error(1)
}

// GetRandomValidPreparation is a mock function that implements our interface.
func (m *mockQuerier) GetRandomValidPreparation(ctx context.Context) (*generated.ValidPreparations, error) {
	args := m.Called(ctx)
	return args.Get(0).(*generated.ValidPreparations), args.Error(1)
}

// GetRecipeByID is a mock function that implements our interface.
func (m *mockQuerier) GetRecipeByID(ctx context.Context, id string) ([]*generated.GetRecipeByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]*generated.GetRecipeByIDRow), args.Error(1)
}

// GetRecipeByIDAndAuthorID is a mock function that implements our interface.
func (m *mockQuerier) GetRecipeByIDAndAuthorID(ctx context.Context, arg *generated.GetRecipeByIDAndAuthorIDParams) ([]*generated.GetRecipeByIDAndAuthorIDRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetRecipeByIDAndAuthorIDRow), args.Error(1)
}

// GetRecipeStep is a mock function that implements our interface.
func (m *mockQuerier) GetRecipeStep(ctx context.Context, arg *generated.GetRecipeStepParams) (*generated.GetRecipeStepRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*generated.GetRecipeStepRow), args.Error(1)
}

// GetRecipeStepIngredient is a mock function that implements our interface.
func (m *mockQuerier) GetRecipeStepIngredient(ctx context.Context, arg *generated.GetRecipeStepIngredientParams) ([]*generated.GetRecipeStepIngredientRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetRecipeStepIngredientRow), args.Error(1)
}

// GetRecipeStepInstrument is a mock function that implements our interface.
func (m *mockQuerier) GetRecipeStepInstrument(ctx context.Context, arg *generated.GetRecipeStepInstrumentParams) ([]*generated.GetRecipeStepInstrumentRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetRecipeStepInstrumentRow), args.Error(1)
}

// GetRecipeStepInstrumentsForRecipe is a mock function that implements our interface.
func (m *mockQuerier) GetRecipeStepInstrumentsForRecipe(ctx context.Context, arg *generated.GetRecipeStepInstrumentsForRecipeParams) ([]*generated.GetRecipeStepInstrumentsForRecipeRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetRecipeStepInstrumentsForRecipeRow), args.Error(1)
}

// GetRecipeStepProduct is a mock function that implements our interface.
func (m *mockQuerier) GetRecipeStepProduct(ctx context.Context, arg *generated.GetRecipeStepProductParams) ([]*generated.GetRecipeStepProductRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetRecipeStepProductRow), args.Error(1)
}

// GetTotalAPIClientCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalAPIClientCount(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// GetTotalMealPlanOptionVotesCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalMealPlanOptionVotesCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalMealPlanOptionsCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalMealPlanOptionsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalMealPlansCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalMealPlansCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalMealsCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalMealsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalRecipeStepInstrumentCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalRecipeStepInstrumentCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalRecipeStepsCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalRecipeStepsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalRecipesCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalRecipesCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalValidIngredientCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalValidIngredientCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalValidIngredientMeasurementUnitsCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalValidIngredientMeasurementUnitsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalValidIngredientPreparationsCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalValidIngredientPreparationsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalValidInstrumentCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalValidInstrumentCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalValidMeasurementUnitCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalValidMeasurementUnitCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalValidPreparationInstrumentsCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalValidPreparationInstrumentsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetTotalValidPreparationsCount is a mock function that implements our interface.
func (m *mockQuerier) GetTotalValidPreparationsCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// GetUserByID is a mock function that implements our interface.
func (m *mockQuerier) GetUserByID(ctx context.Context, id string) (*generated.GetUserByIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetUserByIDRow), args.Error(1)
}

// GetUserByUsername is a mock function that implements our interface.
func (m *mockQuerier) GetUserByUsername(ctx context.Context, username string) (*generated.GetUserByUsernameRow, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*generated.GetUserByUsernameRow), args.Error(1)
}

// GetUserIDByEmail is a mock function that implements our interface.
func (m *mockQuerier) GetUserIDByEmail(ctx context.Context, emailAddress string) (*generated.GetUserIDByEmailRow, error) {
	args := m.Called(ctx, emailAddress)
	return args.Get(0).(*generated.GetUserIDByEmailRow), args.Error(1)
}

// GetUserWithVerified2FA is a mock function that implements our interface.
func (m *mockQuerier) GetUserWithVerified2FA(ctx context.Context, id string) (*generated.GetUserWithVerified2FARow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetUserWithVerified2FARow), args.Error(1)
}

// GetValidIngredient is a mock function that implements our interface.
func (m *mockQuerier) GetValidIngredient(ctx context.Context, id string) (*generated.GetValidIngredientRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetValidIngredientRow), args.Error(1)
}

// GetValidIngredientMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) GetValidIngredientMeasurementUnit(ctx context.Context, id string) (*generated.GetValidIngredientMeasurementUnitRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetValidIngredientMeasurementUnitRow), args.Error(1)
}

// GetValidIngredientPreparation is a mock function that implements our interface.
func (m *mockQuerier) GetValidIngredientPreparation(ctx context.Context, id string) (*generated.GetValidIngredientPreparationRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetValidIngredientPreparationRow), args.Error(1)
}

// GetValidInstrument is a mock function that implements our interface.
func (m *mockQuerier) GetValidInstrument(ctx context.Context, id string) (*generated.GetValidInstrumentRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetValidInstrumentRow), args.Error(1)
}

// GetValidInstruments is a mock function that implements our interface.
func (m *mockQuerier) GetValidInstruments(ctx context.Context, arg *generated.GetValidInstrumentsParams) ([]*generated.GetValidInstrumentsRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.GetValidInstrumentsRow), args.Error(1)
}

// GetValidMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) GetValidMeasurementUnit(ctx context.Context, id string) (*generated.GetValidMeasurementUnitRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetValidMeasurementUnitRow), args.Error(1)
}

// GetValidPreparation is a mock function that implements our interface.
func (m *mockQuerier) GetValidPreparation(ctx context.Context, id string) (*generated.ValidPreparations, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.ValidPreparations), args.Error(1)
}

// GetValidPreparationInstrument is a mock function that implements our interface.
func (m *mockQuerier) GetValidPreparationInstrument(ctx context.Context, id string) (*generated.GetValidPreparationInstrumentRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*generated.GetValidPreparationInstrumentRow), args.Error(1)
}

// GetWebhook is a mock function that implements our interface.
func (m *mockQuerier) GetWebhook(ctx context.Context, arg *generated.GetWebhookParams) (*generated.Webhooks, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(*generated.Webhooks), args.Error(1)
}

// HouseholdInvitationExists is a mock function that implements our interface.
func (m *mockQuerier) HouseholdInvitationExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// MarkHouseholdAsUserDefaultQuery is a mock function that implements our interface.
func (m *mockQuerier) MarkHouseholdAsUserDefaultQuery(ctx context.Context, arg *generated.MarkHouseholdAsUserDefaultQueryParams) error {
	return m.Called(ctx, arg).Error(0)
}

// MarkUserTwoFactorSecretAsVerified is a mock function that implements our interface.
func (m *mockQuerier) MarkUserTwoFactorSecretAsVerified(ctx context.Context, arg *generated.MarkUserTwoFactorSecretAsVerifiedParams) error {
	return m.Called(ctx, arg).Error(0)
}

// MealExists is a mock function that implements our interface.
func (m *mockQuerier) MealExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// MealPlanExists is a mock function that implements our interface.
func (m *mockQuerier) MealPlanExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// MealPlanOptionCreation is a mock function that implements our interface.
func (m *mockQuerier) MealPlanOptionCreation(ctx context.Context, arg *generated.MealPlanOptionCreationParams) error {
	return m.Called(ctx, arg).Error(0)
}

// MealPlanOptionExists is a mock function that implements our interface.
func (m *mockQuerier) MealPlanOptionExists(ctx context.Context, arg *generated.MealPlanOptionExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// MealPlanOptionVoteCreation is a mock function that implements our interface.
func (m *mockQuerier) MealPlanOptionVoteCreation(ctx context.Context, arg *generated.MealPlanOptionVoteCreationParams) error {
	return m.Called(ctx, arg).Error(0)
}

// MealPlanOptionVoteExists is a mock function that implements our interface.
func (m *mockQuerier) MealPlanOptionVoteExists(ctx context.Context, arg *generated.MealPlanOptionVoteExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// ModifyUserPermissionsQuery is a mock function that implements our interface.
func (m *mockQuerier) ModifyUserPermissionsQuery(ctx context.Context, arg *generated.ModifyUserPermissionsQueryParams) error {
	return m.Called(ctx, arg).Error(0)
}

// RecipeExists is a mock function that implements our interface.
func (m *mockQuerier) RecipeExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// RecipeStepExists is a mock function that implements our interface.
func (m *mockQuerier) RecipeStepExists(ctx context.Context, arg *generated.RecipeStepExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// RecipeStepIngredientExists is a mock function that implements our interface.
func (m *mockQuerier) RecipeStepIngredientExists(ctx context.Context, arg *generated.RecipeStepIngredientExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// RecipeStepInstrumentExists is a mock function that implements our interface.
func (m *mockQuerier) RecipeStepInstrumentExists(ctx context.Context, arg *generated.RecipeStepInstrumentExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// RecipeStepProductExists is a mock function that implements our interface.
func (m *mockQuerier) RecipeStepProductExists(ctx context.Context, arg *generated.RecipeStepProductExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// RecipeStepProductsForRecipeQuery is a mock function that implements our interface.
func (m *mockQuerier) RecipeStepProductsForRecipeQuery(ctx context.Context, arg *generated.RecipeStepProductsForRecipeQueryParams) ([]*generated.RecipeStepProductsForRecipeQueryRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]*generated.RecipeStepProductsForRecipeQueryRow), args.Error(1)
}

// RedeemPasswordResetToken is a mock function that implements our interface.
func (m *mockQuerier) RedeemPasswordResetToken(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// RemoveUserFromHouseholdQuery is a mock function that implements our interface.
func (m *mockQuerier) RemoveUserFromHouseholdQuery(ctx context.Context, arg *generated.RemoveUserFromHouseholdQueryParams) error {
	return m.Called(ctx, arg).Error(0)
}

// SearchForUserByUsername is a mock function that implements our interface.
func (m *mockQuerier) SearchForUserByUsername(ctx context.Context, username string) ([]*generated.SearchForUserByUsernameRow, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]*generated.SearchForUserByUsernameRow), args.Error(1)
}

// SearchForValidIngredients is a mock function that implements our interface.
func (m *mockQuerier) SearchForValidIngredients(ctx context.Context, name string) ([]*generated.SearchForValidIngredientsRow, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]*generated.SearchForValidIngredientsRow), args.Error(1)
}

// SearchForValidInstruments is a mock function that implements our interface.
func (m *mockQuerier) SearchForValidInstruments(ctx context.Context, name string) ([]*generated.SearchForValidInstrumentsRow, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]*generated.SearchForValidInstrumentsRow), args.Error(1)
}

// SearchForValidMeasurementUnits is a mock function that implements our interface.
func (m *mockQuerier) SearchForValidMeasurementUnits(ctx context.Context, name string) ([]*generated.SearchForValidMeasurementUnitsRow, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]*generated.SearchForValidMeasurementUnitsRow), args.Error(1)
}

// SetInvitationStatus is a mock function that implements our interface.
func (m *mockQuerier) SetInvitationStatus(ctx context.Context, arg *generated.SetInvitationStatusParams) error {
	return m.Called(ctx, arg).Error(0)
}

// SetUserAccountStatus is a mock function that implements our interface.
func (m *mockQuerier) SetUserAccountStatus(ctx context.Context, arg *generated.SetUserAccountStatusParams) error {
	return m.Called(ctx, arg).Error(0)
}

// TotalRecipeStepIngredientCount is a mock function that implements our interface.
func (m *mockQuerier) TotalRecipeStepIngredientCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// TotalRecipeStepProductCount is a mock function that implements our interface.
func (m *mockQuerier) TotalRecipeStepProductCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// TransferHouseholdMembershipQuery is a mock function that implements our interface.
func (m *mockQuerier) TransferHouseholdMembershipQuery(ctx context.Context, arg *generated.TransferHouseholdMembershipQueryParams) error {
	return m.Called(ctx, arg).Error(0)
}

// TransferHouseholdOwnershipQuery is a mock function that implements our interface.
func (m *mockQuerier) TransferHouseholdOwnershipQuery(ctx context.Context, arg *generated.TransferHouseholdOwnershipQueryParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateHousehold is a mock function that implements our interface.
func (m *mockQuerier) UpdateHousehold(ctx context.Context, arg *generated.UpdateHouseholdParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateMealPlan is a mock function that implements our interface.
func (m *mockQuerier) UpdateMealPlan(ctx context.Context, arg *generated.UpdateMealPlanParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateMealPlanOption is a mock function that implements our interface.
func (m *mockQuerier) UpdateMealPlanOption(ctx context.Context, arg *generated.UpdateMealPlanOptionParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateMealPlanOptionVote is a mock function that implements our interface.
func (m *mockQuerier) UpdateMealPlanOptionVote(ctx context.Context, arg *generated.UpdateMealPlanOptionVoteParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateRecipe is a mock function that implements our interface.
func (m *mockQuerier) UpdateRecipe(ctx context.Context, arg *generated.UpdateRecipeParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateRecipeStep is a mock function that implements our interface.
func (m *mockQuerier) UpdateRecipeStep(ctx context.Context, arg *generated.UpdateRecipeStepParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateRecipeStepIngredient is a mock function that implements our interface.
func (m *mockQuerier) UpdateRecipeStepIngredient(ctx context.Context, arg *generated.UpdateRecipeStepIngredientParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateRecipeStepInstrument is a mock function that implements our interface.
func (m *mockQuerier) UpdateRecipeStepInstrument(ctx context.Context, arg *generated.UpdateRecipeStepInstrumentParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateRecipeStepProduct is a mock function that implements our interface.
func (m *mockQuerier) UpdateRecipeStepProduct(ctx context.Context, arg *generated.UpdateRecipeStepProductParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateUser is a mock function that implements our interface.
func (m *mockQuerier) UpdateUser(ctx context.Context, arg *generated.UpdateUserParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateUserPassword is a mock function that implements our interface.
func (m *mockQuerier) UpdateUserPassword(ctx context.Context, arg *generated.UpdateUserPasswordParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateUserTwoFactorSecret is a mock function that implements our interface.
func (m *mockQuerier) UpdateUserTwoFactorSecret(ctx context.Context, arg *generated.UpdateUserTwoFactorSecretParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateValidIngredient is a mock function that implements our interface.
func (m *mockQuerier) UpdateValidIngredient(ctx context.Context, arg *generated.UpdateValidIngredientParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateValidIngredientMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) UpdateValidIngredientMeasurementUnit(ctx context.Context, arg *generated.UpdateValidIngredientMeasurementUnitParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateValidIngredientPreparation is a mock function that implements our interface.
func (m *mockQuerier) UpdateValidIngredientPreparation(ctx context.Context, arg *generated.UpdateValidIngredientPreparationParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateValidInstrument is a mock function that implements our interface.
func (m *mockQuerier) UpdateValidInstrument(ctx context.Context, arg *generated.UpdateValidInstrumentParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateValidMeasurementUnit is a mock function that implements our interface.
func (m *mockQuerier) UpdateValidMeasurementUnit(ctx context.Context, arg *generated.UpdateValidMeasurementUnitParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateValidPreparation is a mock function that implements our interface.
func (m *mockQuerier) UpdateValidPreparation(ctx context.Context, arg *generated.UpdateValidPreparationParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UpdateValidPreparationInstrument is a mock function that implements our interface.
func (m *mockQuerier) UpdateValidPreparationInstrument(ctx context.Context, arg *generated.UpdateValidPreparationInstrumentParams) error {
	return m.Called(ctx, arg).Error(0)
}

// UserHasStatus is a mock function that implements our interface.
func (m *mockQuerier) UserHasStatus(ctx context.Context, arg *generated.UserHasStatusParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// UserIsMemberOfHouseholdQuery is a mock function that implements our interface.
func (m *mockQuerier) UserIsMemberOfHouseholdQuery(ctx context.Context, arg *generated.UserIsMemberOfHouseholdQueryParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

// ValidIngredientExists is a mock function that implements our interface.
func (m *mockQuerier) ValidIngredientExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// ValidIngredientMeasurementUnitExists is a mock function that implements our interface.
func (m *mockQuerier) ValidIngredientMeasurementUnitExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// ValidIngredientPreparationExists is a mock function that implements our interface.
func (m *mockQuerier) ValidIngredientPreparationExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// ValidInstrumentExists is a mock function that implements our interface.
func (m *mockQuerier) ValidInstrumentExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// ValidMeasurementUnitExists is a mock function that implements our interface.
func (m *mockQuerier) ValidMeasurementUnitExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// ValidPreparationExists is a mock function that implements our interface.
func (m *mockQuerier) ValidPreparationExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// ValidPreparationInstrumentExists is a mock function that implements our interface.
func (m *mockQuerier) ValidPreparationInstrumentExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// ValidPreparationsSearch is a mock function that implements our interface.
func (m *mockQuerier) ValidPreparationsSearch(ctx context.Context, name string) ([]*generated.ValidPreparations, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]*generated.ValidPreparations), args.Error(1)
}

// WebhookExists is a mock function that implements our interface.
func (m *mockQuerier) WebhookExists(ctx context.Context, arg *generated.WebhookExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}
