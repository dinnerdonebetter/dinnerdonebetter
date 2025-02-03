package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/converters"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) AcceptHouseholdInvitation(ctx context.Context, request *messages.AcceptHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) AdminLoginForToken(ctx context.Context, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	// TODO: validation

	user, err := s.dataManager.GetAdminUserByUsername(ctx, input.Username)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching user by username")
	}

	loginValid, err := s.authenticator.CredentialsAreValid(
		ctx,
		user.HashedPassword,
		input.Password,
		user.TwoFactorSecret,
		input.TOTPToken,
	)
	if err != nil {
		return nil, observability.PrepareError(err, span, "validating login")
	}

	if !loginValid {
		return nil, observability.PrepareError(err, span, "invalid login")
	}

	if loginValid && user.TwoFactorSecretVerifiedAt != nil && input.TOTPToken == "" {
		return nil, observability.PrepareError(err, span, "user with two factor verification active attempted to log in without providing TOTP")
	}

	defaultHouseholdID, err := s.dataManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "fetching user memberships")
	}

	var token string
	token, err = s.tokenIssuer.IssueToken(ctx, user, s.config.Services.Auth.TokenLifetime)
	if err != nil {
		return nil, observability.PrepareError(err, span, "signing token")
	}

	output := &messages.TokenResponse{
		UserID:      user.ID,
		HouseholdID: defaultHouseholdID,
		Token:       token,
	}

	return output, nil
}

func (s *server) AdminUpdateUserStatus(ctx context.Context, input *messages.UserAccountStatusUpdateInput) (*messages.UserStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) AggregateUserDataReport(ctx context.Context, _ *emptypb.Empty) (*messages.UserDataCollectionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveHousehold(ctx context.Context, request *messages.ArchiveHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveHouseholdInstrumentOwnership(ctx context.Context, request *messages.ArchiveHouseholdInstrumentOwnershipRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveMeal(ctx context.Context, request *messages.ArchiveMealRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveMealPlan(ctx context.Context, request *messages.ArchiveMealPlanRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveMealPlanEvent(ctx context.Context, request *messages.ArchiveMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveMealPlanGroceryListItem(ctx context.Context, request *messages.ArchiveMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveMealPlanOption(ctx context.Context, request *messages.ArchiveMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveMealPlanOptionVote(ctx context.Context, request *messages.ArchiveMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveOAuth2Client(ctx context.Context, request *messages.ArchiveOAuth2ClientRequest) (*messages.OAuth2Client, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipe(ctx context.Context, request *messages.ArchiveRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipePrepTask(ctx context.Context, request *messages.ArchiveRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipeRating(ctx context.Context, request *messages.ArchiveRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipeStep(ctx context.Context, request *messages.ArchiveRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *messages.ArchiveRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipeStepIngredient(ctx context.Context, request *messages.ArchiveRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipeStepInstrument(ctx context.Context, request *messages.ArchiveRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipeStepProduct(ctx context.Context, request *messages.ArchiveRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveRecipeStepVessel(ctx context.Context, request *messages.ArchiveRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveServiceSetting(ctx context.Context, request *messages.ArchiveServiceSettingRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveServiceSettingConfiguration(ctx context.Context, request *messages.ArchiveServiceSettingConfigurationRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveUser(ctx context.Context, request *messages.ArchiveUserRequest) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveUserIngredientPreference(ctx context.Context, request *messages.ArchiveUserIngredientPreferenceRequest) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveUserMembership(ctx context.Context, request *messages.ArchiveUserMembershipRequest) (*messages.HouseholdUserMembership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidIngredient(ctx context.Context, request *messages.ArchiveValidIngredientRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidIngredientGroup(ctx context.Context, request *messages.ArchiveValidIngredientGroupRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *messages.ArchiveValidIngredientMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidIngredientPreparation(ctx context.Context, request *messages.ArchiveValidIngredientPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidIngredientState(ctx context.Context, request *messages.ArchiveValidIngredientStateRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidIngredientStateIngredient(ctx context.Context, request *messages.ArchiveValidIngredientStateIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidInstrument(ctx context.Context, request *messages.ArchiveValidInstrumentRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidMeasurementUnit(ctx context.Context, request *messages.ArchiveValidMeasurementUnitRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *messages.ArchiveValidMeasurementUnitConversionRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidPreparation(ctx context.Context, request *messages.ArchiveValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidPreparationInstrument(ctx context.Context, request *messages.ArchiveValidPreparationInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidPreparationVessel(ctx context.Context, request *messages.ArchiveValidPreparationVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveValidVessel(ctx context.Context, request *messages.ArchiveValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveWebhook(ctx context.Context, request *messages.ArchiveWebhookRequest) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) ArchiveWebhookTriggerEvent(ctx context.Context, request *messages.ArchiveWebhookTriggerEventRequest) (*messages.WebhookTriggerEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CancelHouseholdInvitation(ctx context.Context, request *messages.CancelHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CheckForReadiness(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CheckPermissions(ctx context.Context, input *messages.UserPermissionsRequestInput) (*messages.UserPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CloneRecipe(ctx context.Context, request *messages.CloneRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateHousehold(ctx context.Context, input *messages.HouseholdCreationRequestInput) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateHouseholdInstrumentOwnership(ctx context.Context, input *messages.HouseholdInstrumentOwnershipCreationRequestInput) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateHouseholdInvitation(ctx context.Context, request *messages.CreateHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateMeal(ctx context.Context, input *messages.MealCreationRequestInput) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateMealPlan(ctx context.Context, input *messages.MealPlanCreationRequestInput) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateMealPlanEvent(ctx context.Context, request *messages.CreateMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateMealPlanGroceryListItem(ctx context.Context, request *messages.CreateMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateMealPlanOption(ctx context.Context, request *messages.CreateMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateMealPlanOptionVote(ctx context.Context, request *messages.CreateMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateMealPlanTask(ctx context.Context, request *messages.CreateMealPlanTaskRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateOAuth2Client(ctx context.Context, input *messages.OAuth2ClientCreationRequestInput) (*messages.OAuth2ClientCreationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipe(ctx context.Context, input *messages.RecipeCreationRequestInput) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipePrepTask(ctx context.Context, request *messages.CreateRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipeRating(ctx context.Context, request *messages.CreateRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipeStep(ctx context.Context, request *messages.CreateRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipeStepCompletionCondition(ctx context.Context, request *messages.CreateRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipeStepIngredient(ctx context.Context, request *messages.CreateRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipeStepInstrument(ctx context.Context, request *messages.CreateRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipeStepProduct(ctx context.Context, request *messages.CreateRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateRecipeStepVessel(ctx context.Context, request *messages.CreateRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateServiceSetting(ctx context.Context, input *messages.ServiceSettingCreationRequestInput) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateServiceSettingConfiguration(ctx context.Context, input *messages.ServiceSettingConfigurationCreationRequestInput) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateUser(ctx context.Context, input *messages.UserRegistrationInput) (*messages.UserCreationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateUserIngredientPreference(ctx context.Context, input *messages.UserIngredientPreferenceCreationRequestInput) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateUserNotification(ctx context.Context, input *messages.UserNotificationCreationRequestInput) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidIngredient(ctx context.Context, input *messages.ValidIngredientCreationRequestInput) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidIngredientGroup(ctx context.Context, input *messages.ValidIngredientGroupCreationRequestInput) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidIngredientMeasurementUnit(ctx context.Context, input *messages.ValidIngredientMeasurementUnitCreationRequestInput) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidIngredientPreparation(ctx context.Context, input *messages.ValidIngredientPreparationCreationRequestInput) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidIngredientState(ctx context.Context, input *messages.ValidIngredientStateCreationRequestInput) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidIngredientStateIngredient(ctx context.Context, input *messages.ValidIngredientStateIngredientCreationRequestInput) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidInstrument(ctx context.Context, input *messages.ValidInstrumentCreationRequestInput) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidMeasurementUnit(ctx context.Context, input *messages.ValidMeasurementUnitCreationRequestInput) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidMeasurementUnitConversion(ctx context.Context, input *messages.ValidMeasurementUnitConversionCreationRequestInput) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidPreparation(ctx context.Context, input *messages.ValidPreparationCreationRequestInput) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidPreparationInstrument(ctx context.Context, input *messages.ValidPreparationInstrumentCreationRequestInput) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidPreparationVessel(ctx context.Context, input *messages.ValidPreparationVesselCreationRequestInput) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateValidVessel(ctx context.Context, input *messages.ValidVesselCreationRequestInput) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateWebhook(ctx context.Context, input *messages.WebhookCreationRequestInput) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) CreateWebhookTriggerEvent(ctx context.Context, request *messages.CreateWebhookTriggerEventRequest) (*messages.WebhookTriggerEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) DestroyAllUserData(ctx context.Context, _ *emptypb.Empty) (*messages.DataDeletionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) FetchUserDataReport(ctx context.Context, request *messages.FetchUserDataReportRequest) (*messages.UserDataCollection, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) FinalizeMealPlan(ctx context.Context, request *messages.FinalizeMealPlanRequest) (*messages.FinalizeMealPlansResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetActiveHousehold(ctx context.Context, _ *emptypb.Empty) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetAuditLogEntriesForHousehold(ctx context.Context, request *messages.GetAuditLogEntriesForHouseholdRequest) (*messages.AuditLogEntry, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetAuditLogEntriesForUser(ctx context.Context, request *messages.GetAuditLogEntriesForUserRequest) (*messages.AuditLogEntry, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetAuditLogEntryByID(ctx context.Context, request *messages.GetAuditLogEntryByIDRequest) (*messages.AuditLogEntry, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetAuthStatus(ctx context.Context, _ *emptypb.Empty) (*messages.UserStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetHousehold(ctx context.Context, request *messages.GetHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetHouseholdInstrumentOwnership(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetHouseholdInstrumentOwnerships(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipsRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetHouseholdInvitation(ctx context.Context, request *messages.GetHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetHouseholdInvitationByID(ctx context.Context, request *messages.GetHouseholdInvitationByIDRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetHouseholds(ctx context.Context, request *messages.GetHouseholdsRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMeal(ctx context.Context, request *messages.GetMealRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlan(ctx context.Context, request *messages.GetMealPlanRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanEvent(ctx context.Context, request *messages.GetMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanEvents(ctx context.Context, request *messages.GetMealPlanEventsRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanGroceryListItem(ctx context.Context, request *messages.GetMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *messages.GetMealPlanGroceryListItemsForMealPlanRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanOption(ctx context.Context, request *messages.GetMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanOptionVote(ctx context.Context, request *messages.GetMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanOptionVotes(ctx context.Context, request *messages.GetMealPlanOptionVotesRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanOptions(ctx context.Context, request *messages.GetMealPlanOptionsRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanTask(ctx context.Context, request *messages.GetMealPlanTaskRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlanTasks(ctx context.Context, request *messages.GetMealPlanTasksRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMealPlansForHousehold(ctx context.Context, request *messages.GetMealPlansForHouseholdRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMeals(ctx context.Context, request *messages.GetMealsRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetMermaidDiagramForRecipe(ctx context.Context, request *messages.GetMermaidDiagramForRecipeRequest) (*messages.GetMermaidDiagramForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetOAuth2Client(ctx context.Context, request *messages.GetOAuth2ClientRequest) (*messages.OAuth2Client, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetOAuth2Clients(ctx context.Context, request *messages.GetOAuth2ClientsRequest) (*messages.OAuth2Client, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRandomValidIngredient(ctx context.Context, _ *emptypb.Empty) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	ingredient, err := s.dataManager.GetRandomValidIngredient(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting random valid ingredient")
	}

	output := &messages.ValidIngredient{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(ingredient.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(ingredient.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(ingredient.ArchivedAt),
		StorageTemperatureInCelsius: &messages.OptionalFloat32Range{
			Max: ingredient.StorageTemperatureInCelsius.Max,
			Min: ingredient.StorageTemperatureInCelsius.Min,
		},
		IconPath:               ingredient.IconPath,
		Warning:                ingredient.Warning,
		PluralName:             ingredient.PluralName,
		StorageInstructions:    ingredient.StorageInstructions,
		Name:                   ingredient.Name,
		ID:                     ingredient.ID,
		Description:            ingredient.Description,
		Slug:                   ingredient.Slug,
		ShoppingSuggestions:    ingredient.ShoppingSuggestions,
		ContainsShellfish:      ingredient.ContainsShellfish,
		IsLiquid:               ingredient.IsLiquid,
		ContainsPeanut:         ingredient.ContainsPeanut,
		ContainsTreeNut:        ingredient.ContainsTreeNut,
		ContainsEgg:            ingredient.ContainsEgg,
		ContainsWheat:          ingredient.ContainsWheat,
		ContainsSoy:            ingredient.ContainsSoy,
		AnimalDerived:          ingredient.AnimalDerived,
		RestrictToPreparations: ingredient.RestrictToPreparations,
		ContainsSesame:         ingredient.ContainsSesame,
		ContainsFish:           ingredient.ContainsFish,
		ContainsGluten:         ingredient.ContainsGluten,
		ContainsDairy:          ingredient.ContainsDairy,
		ContainsAlcohol:        ingredient.ContainsAlcohol,
		AnimalFlesh:            ingredient.AnimalFlesh,
		IsStarch:               ingredient.IsStarch,
		IsProtein:              ingredient.IsProtein,
		IsGrain:                ingredient.IsGrain,
		IsFruit:                ingredient.IsFruit,
		IsSalt:                 ingredient.IsSalt,
		IsFat:                  ingredient.IsFat,
		IsAcid:                 ingredient.IsAcid,
		IsHeat:                 ingredient.IsHeat,
	}

	return output, nil
}

func (s *server) GetRandomValidInstrument(ctx context.Context, _ *emptypb.Empty) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRandomValidPreparation(ctx context.Context, _ *emptypb.Empty) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRandomValidVessel(ctx context.Context, _ *emptypb.Empty) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	ingredient, err := s.dataManager.GetRandomValidVessel(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting random valid ingredient")
	}

	output := &messages.ValidVessel{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(ingredient.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(ingredient.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(ingredient.ArchivedAt),
		CapacityUnit: &messages.ValidMeasurementUnit{
			CreatedAt:     converters.ConvertTimeToPBTimestamp(ingredient.CapacityUnit.CreatedAt),
			LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(ingredient.CapacityUnit.LastUpdatedAt),
			ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(ingredient.CapacityUnit.ArchivedAt),
			Name:          ingredient.CapacityUnit.Name,
			IconPath:      ingredient.CapacityUnit.IconPath,
			ID:            ingredient.CapacityUnit.ID,
			Description:   ingredient.CapacityUnit.Description,
			PluralName:    ingredient.CapacityUnit.PluralName,
			Slug:          ingredient.CapacityUnit.Slug,
			Volumetric:    ingredient.CapacityUnit.Volumetric,
			Universal:     ingredient.CapacityUnit.Universal,
			Metric:        ingredient.CapacityUnit.Metric,
			Imperial:      ingredient.CapacityUnit.Imperial,
		},
		IconPath:                       ingredient.IconPath,
		PluralName:                     ingredient.PluralName,
		Description:                    ingredient.Description,
		Name:                           ingredient.Name,
		Slug:                           ingredient.Slug,
		Shape:                          ingredient.Shape,
		ID:                             ingredient.ID,
		WidthInMillimeters:             ingredient.WidthInMillimeters,
		LengthInMillimeters:            ingredient.LengthInMillimeters,
		HeightInMillimeters:            ingredient.HeightInMillimeters,
		Capacity:                       ingredient.Capacity,
		IncludeInGeneratedInstructions: ingredient.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          ingredient.DisplayInSummaryLists,
		UsableForStorage:               ingredient.UsableForStorage,
	}

	return output, nil
}

func (s *server) GetReceivedHouseholdInvitations(ctx context.Context, request *messages.GetReceivedHouseholdInvitationsRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipe(ctx context.Context, request *messages.GetRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeMealPlanTasks(ctx context.Context, request *messages.GetRecipeMealPlanTasksRequest) (*messages.RecipePrepTaskStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipePrepTask(ctx context.Context, request *messages.GetRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipePrepTasks(ctx context.Context, request *messages.GetRecipePrepTasksRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeRating(ctx context.Context, request *messages.GetRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeRatingsForRecipe(ctx context.Context, request *messages.GetRecipeRatingsForRecipeRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStep(ctx context.Context, request *messages.GetRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepCompletionCondition(ctx context.Context, request *messages.GetRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepCompletionConditions(ctx context.Context, request *messages.GetRecipeStepCompletionConditionsRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepIngredient(ctx context.Context, request *messages.GetRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepIngredients(ctx context.Context, request *messages.GetRecipeStepIngredientsRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepInstrument(ctx context.Context, request *messages.GetRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepInstruments(ctx context.Context, request *messages.GetRecipeStepInstrumentsRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepProduct(ctx context.Context, request *messages.GetRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepProducts(ctx context.Context, request *messages.GetRecipeStepProductsRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepVessel(ctx context.Context, request *messages.GetRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeStepVessels(ctx context.Context, request *messages.GetRecipeStepVesselsRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipeSteps(ctx context.Context, request *messages.GetRecipeStepsRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetRecipes(ctx context.Context, request *messages.GetRecipesRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetSelf(ctx context.Context, _ *emptypb.Empty) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetSentHouseholdInvitations(ctx context.Context, request *messages.GetSentHouseholdInvitationsRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetServiceSetting(ctx context.Context, request *messages.GetServiceSettingRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetServiceSettingConfigurationByName(ctx context.Context, request *messages.GetServiceSettingConfigurationByNameRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetServiceSettingConfigurationsForHousehold(ctx context.Context, request *messages.GetServiceSettingConfigurationsForHouseholdRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetServiceSettingConfigurationsForUser(ctx context.Context, request *messages.GetServiceSettingConfigurationsForUserRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetServiceSettings(ctx context.Context, request *messages.GetServiceSettingsRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetUser(ctx context.Context, request *messages.GetUserRequest) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetUserIngredientPreferences(ctx context.Context, request *messages.GetUserIngredientPreferencesRequest) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetUserNotification(ctx context.Context, request *messages.GetUserNotificationRequest) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetUserNotifications(ctx context.Context, request *messages.GetUserNotificationsRequest) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetUsers(ctx context.Context, request *messages.GetUsersRequest) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredient(ctx context.Context, request *messages.GetValidIngredientRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientGroup(ctx context.Context, request *messages.GetValidIngredientGroupRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientGroups(ctx context.Context, request *messages.GetValidIngredientGroupsRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientMeasurementUnits(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByIngredientRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientPreparations(ctx context.Context, request *messages.GetValidIngredientPreparationsRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *messages.GetValidIngredientPreparationsByIngredientRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationsByPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientState(ctx context.Context, request *messages.GetValidIngredientStateRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientStateIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientStateIngredients(ctx context.Context, request *messages.GetValidIngredientStateIngredientsRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientStateRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredientStates(ctx context.Context, request *messages.GetValidIngredientStatesRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidIngredients(ctx context.Context, request *messages.GetValidIngredientsRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidInstrument(ctx context.Context, request *messages.GetValidInstrumentRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidInstruments(ctx context.Context, request *messages.GetValidInstrumentsRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidMeasurementUnit(ctx context.Context, request *messages.GetValidMeasurementUnitRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidMeasurementUnitConversion(ctx context.Context, request *messages.GetValidMeasurementUnitConversionRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsFromUnitRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsToUnitRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidMeasurementUnits(ctx context.Context, request *messages.GetValidMeasurementUnitsRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparation(ctx context.Context, request *messages.GetValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparationInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparationInstruments(ctx context.Context, request *messages.GetValidPreparationInstrumentsRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentsByInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *messages.GetValidPreparationInstrumentsByPreparationRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparationVessel(ctx context.Context, request *messages.GetValidPreparationVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparationVessels(ctx context.Context, request *messages.GetValidPreparationVesselsRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparationVesselsByPreparation(ctx context.Context, request *messages.GetValidPreparationVesselsByPreparationRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparationVesselsByVessel(ctx context.Context, request *messages.GetValidPreparationVesselsByVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidPreparations(ctx context.Context, request *messages.GetValidPreparationsRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidVessel(ctx context.Context, request *messages.GetValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetValidVessels(ctx context.Context, request *messages.GetValidVesselsRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetWebhook(ctx context.Context, request *messages.GetWebhookRequest) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) GetWebhooks(ctx context.Context, request *messages.GetWebhooksRequest) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) LoginForToken(ctx context.Context, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) PublishArbitraryQueueMessage(ctx context.Context, input *messages.ArbitraryQueueMessageRequestInput) (*messages.ArbitraryQueueMessageResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RedeemPasswordResetToken(ctx context.Context, input *messages.PasswordResetTokenRedemptionRequestInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RefreshTOTPSecret(ctx context.Context, input *messages.TOTPSecretRefreshInput) (*messages.TOTPSecretRefreshResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RejectHouseholdInvitation(ctx context.Context, request *messages.RejectHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RequestEmailVerificationEmail(ctx context.Context, _ *emptypb.Empty) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RequestPasswordResetToken(ctx context.Context, input *messages.PasswordResetTokenCreationRequestInput) (*messages.PasswordResetToken, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RequestUsernameReminder(ctx context.Context, input *messages.UsernameReminderRequestInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RunFinalizeMealPlanWorker(ctx context.Context, request *messages.FinalizeMealPlansRequest) (*messages.FinalizeMealPlansResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RunMealPlanGroceryListInitializerWorker(ctx context.Context, request *messages.InitializeMealPlanGroceryListRequest) (*messages.InitializeMealPlanGroceryListResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) RunMealPlanTaskCreatorWorker(ctx context.Context, request *messages.CreateMealPlanTasksRequest) (*messages.CreateMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForMeals(ctx context.Context, request *messages.SearchForMealsRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForRecipes(ctx context.Context, request *messages.SearchForRecipesRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForServiceSettings(ctx context.Context, request *messages.SearchForServiceSettingsRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForUsers(ctx context.Context, request *messages.SearchForUsersRequest) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForValidIngredientGroups(ctx context.Context, request *messages.SearchForValidIngredientGroupsRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForValidIngredientStates(ctx context.Context, request *messages.SearchForValidIngredientStatesRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForValidIngredients(ctx context.Context, request *messages.SearchForValidIngredientsRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForValidInstruments(ctx context.Context, request *messages.SearchForValidInstrumentsRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForValidMeasurementUnits(ctx context.Context, request *messages.SearchForValidMeasurementUnitsRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForValidPreparations(ctx context.Context, request *messages.SearchForValidPreparationsRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchForValidVessels(ctx context.Context, request *messages.SearchForValidVesselsRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchValidIngredientsByPreparation(ctx context.Context, request *messages.SearchValidIngredientsByPreparationRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *messages.SearchValidMeasurementUnitsByIngredientRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) SetDefaultHousehold(ctx context.Context, request *messages.SetDefaultHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) TransferHouseholdOwnership(ctx context.Context, request *messages.TransferHouseholdOwnershipRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateHousehold(ctx context.Context, request *messages.UpdateHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateHouseholdInstrumentOwnership(ctx context.Context, request *messages.UpdateHouseholdInstrumentOwnershipRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateHouseholdMemberPermissions(ctx context.Context, request *messages.UpdateHouseholdMemberPermissionsRequest) (*messages.UserPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateMealPlan(ctx context.Context, request *messages.UpdateMealPlanRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateMealPlanEvent(ctx context.Context, request *messages.UpdateMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateMealPlanGroceryListItem(ctx context.Context, request *messages.UpdateMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateMealPlanOption(ctx context.Context, request *messages.UpdateMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateMealPlanOptionVote(ctx context.Context, request *messages.UpdateMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateMealPlanTaskStatus(ctx context.Context, request *messages.UpdateMealPlanTaskStatusRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdatePassword(ctx context.Context, input *messages.PasswordUpdateInput) (*messages.PasswordResetResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipe(ctx context.Context, request *messages.UpdateRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipePrepTask(ctx context.Context, request *messages.UpdateRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipeRating(ctx context.Context, request *messages.UpdateRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipeStep(ctx context.Context, request *messages.UpdateRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipeStepCompletionCondition(ctx context.Context, request *messages.UpdateRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipeStepIngredient(ctx context.Context, request *messages.UpdateRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipeStepInstrument(ctx context.Context, request *messages.UpdateRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipeStepProduct(ctx context.Context, request *messages.UpdateRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateRecipeStepVessel(ctx context.Context, request *messages.UpdateRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateServiceSettingConfiguration(ctx context.Context, request *messages.UpdateServiceSettingConfigurationRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateUserDetails(ctx context.Context, input *messages.UserDetailsUpdateRequestInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateUserEmailAddress(ctx context.Context, input *messages.UserEmailAddressUpdateInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateUserIngredientPreference(ctx context.Context, request *messages.UpdateUserIngredientPreferenceRequest) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateUserNotification(ctx context.Context, request *messages.UpdateUserNotificationRequest) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateUserUsername(ctx context.Context, input *messages.UsernameUpdateInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidIngredient(ctx context.Context, request *messages.UpdateValidIngredientRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidIngredientGroup(ctx context.Context, request *messages.UpdateValidIngredientGroupRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *messages.UpdateValidIngredientMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidIngredientPreparation(ctx context.Context, request *messages.UpdateValidIngredientPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidIngredientState(ctx context.Context, request *messages.UpdateValidIngredientStateRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidIngredientStateIngredient(ctx context.Context, request *messages.UpdateValidIngredientStateIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidInstrument(ctx context.Context, request *messages.UpdateValidInstrumentRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidMeasurementUnit(ctx context.Context, request *messages.UpdateValidMeasurementUnitRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidMeasurementUnitConversion(ctx context.Context, request *messages.UpdateValidMeasurementUnitConversionRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidPreparation(ctx context.Context, request *messages.UpdateValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidPreparationInstrument(ctx context.Context, request *messages.UpdateValidPreparationInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidPreparationVessel(ctx context.Context, request *messages.UpdateValidPreparationVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UpdateValidVessel(ctx context.Context, request *messages.UpdateValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) UploadUserAvatar(ctx context.Context, input *messages.AvatarUpdateInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) VerifyEmailAddress(ctx context.Context, input *messages.EmailAddressVerificationRequestInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *server) VerifyTOTPSecret(ctx context.Context, input *messages.TOTPSecretVerificationInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}
