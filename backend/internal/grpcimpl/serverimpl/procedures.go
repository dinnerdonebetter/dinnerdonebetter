package serverimpl

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/converters"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	errUnimplemented = errors.New("unimplemented procedure")
	emptyResponse    = &emptypb.Empty{}
)

func (s *Server) AcceptHouseholdInvitation(ctx context.Context, request *messages.AcceptHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) AdminLoginForToken(ctx context.Context, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
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

func (s *Server) AdminUpdateUserStatus(ctx context.Context, input *messages.UserAccountStatusUpdateInput) (*messages.UserStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) AggregateUserDataReport(ctx context.Context, _ *emptypb.Empty) (*messages.UserDataCollectionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveHousehold(ctx context.Context, request *messages.ArchiveHouseholdRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveHouseholdInstrumentOwnership(ctx context.Context, request *messages.ArchiveHouseholdInstrumentOwnershipRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveMeal(ctx context.Context, request *messages.ArchiveMealRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveMealPlan(ctx context.Context, request *messages.ArchiveMealPlanRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveMealPlanEvent(ctx context.Context, request *messages.ArchiveMealPlanEventRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveMealPlanGroceryListItem(ctx context.Context, request *messages.ArchiveMealPlanGroceryListItemRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveMealPlanOption(ctx context.Context, request *messages.ArchiveMealPlanOptionRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveMealPlanOptionVote(ctx context.Context, request *messages.ArchiveMealPlanOptionVoteRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveOAuth2Client(ctx context.Context, request *messages.ArchiveOAuth2ClientRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipe(ctx context.Context, request *messages.ArchiveRecipeRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipePrepTask(ctx context.Context, request *messages.ArchiveRecipePrepTaskRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipeRating(ctx context.Context, request *messages.ArchiveRecipeRatingRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipeStep(ctx context.Context, request *messages.ArchiveRecipeStepRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *messages.ArchiveRecipeStepCompletionConditionRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipeStepIngredient(ctx context.Context, request *messages.ArchiveRecipeStepIngredientRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipeStepInstrument(ctx context.Context, request *messages.ArchiveRecipeStepInstrumentRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipeStepProduct(ctx context.Context, request *messages.ArchiveRecipeStepProductRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveRecipeStepVessel(ctx context.Context, request *messages.ArchiveRecipeStepVesselRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveServiceSetting(ctx context.Context, request *messages.ArchiveServiceSettingRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveServiceSettingConfiguration(ctx context.Context, request *messages.ArchiveServiceSettingConfigurationRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveUser(ctx context.Context, request *messages.ArchiveUserRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveUserIngredientPreference(ctx context.Context, request *messages.ArchiveUserIngredientPreferenceRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveUserMembership(ctx context.Context, request *messages.ArchiveUserMembershipRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidIngredient(ctx context.Context, request *messages.ArchiveValidIngredientRequest) (*emptypb.Empty, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	if err := s.dataManager.ArchiveValidIngredient(ctx, request.ValidIngredientID); err != nil {
		return nil, err
	}

	return emptyResponse, nil
}

func (s *Server) ArchiveValidIngredientGroup(ctx context.Context, request *messages.ArchiveValidIngredientGroupRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *messages.ArchiveValidIngredientMeasurementUnitRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidIngredientPreparation(ctx context.Context, request *messages.ArchiveValidIngredientPreparationRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidIngredientState(ctx context.Context, request *messages.ArchiveValidIngredientStateRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidIngredientStateIngredient(ctx context.Context, request *messages.ArchiveValidIngredientStateIngredientRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidInstrument(ctx context.Context, request *messages.ArchiveValidInstrumentRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidMeasurementUnit(ctx context.Context, request *messages.ArchiveValidMeasurementUnitRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *messages.ArchiveValidMeasurementUnitConversionRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidPreparation(ctx context.Context, request *messages.ArchiveValidPreparationRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidPreparationInstrument(ctx context.Context, request *messages.ArchiveValidPreparationInstrumentRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidPreparationVessel(ctx context.Context, request *messages.ArchiveValidPreparationVesselRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveValidVessel(ctx context.Context, request *messages.ArchiveValidVesselRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveWebhook(ctx context.Context, request *messages.ArchiveWebhookRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) ArchiveWebhookTriggerEvent(ctx context.Context, request *messages.ArchiveWebhookTriggerEventRequest) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CancelHouseholdInvitation(ctx context.Context, request *messages.CancelHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CheckForReadiness(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CheckPermissions(ctx context.Context, input *messages.UserPermissionsRequestInput) (*messages.UserPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CloneRecipe(ctx context.Context, request *messages.CloneRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateHousehold(ctx context.Context, input *messages.HouseholdCreationRequestInput) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateHouseholdInstrumentOwnership(ctx context.Context, input *messages.HouseholdInstrumentOwnershipCreationRequestInput) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateHouseholdInvitation(ctx context.Context, request *messages.CreateHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateMeal(ctx context.Context, input *messages.MealCreationRequestInput) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateMealPlan(ctx context.Context, input *messages.MealPlanCreationRequestInput) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateMealPlanEvent(ctx context.Context, request *messages.CreateMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateMealPlanGroceryListItem(ctx context.Context, request *messages.CreateMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateMealPlanOption(ctx context.Context, request *messages.CreateMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateMealPlanOptionVote(ctx context.Context, request *messages.CreateMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateMealPlanTask(ctx context.Context, request *messages.CreateMealPlanTaskRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateOAuth2Client(ctx context.Context, input *messages.OAuth2ClientCreationRequestInput) (*messages.OAuth2ClientCreationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipe(ctx context.Context, input *messages.RecipeCreationRequestInput) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipePrepTask(ctx context.Context, request *messages.CreateRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipeRating(ctx context.Context, request *messages.CreateRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipeStep(ctx context.Context, request *messages.CreateRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipeStepCompletionCondition(ctx context.Context, request *messages.CreateRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipeStepIngredient(ctx context.Context, request *messages.CreateRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipeStepInstrument(ctx context.Context, request *messages.CreateRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipeStepProduct(ctx context.Context, request *messages.CreateRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateRecipeStepVessel(ctx context.Context, request *messages.CreateRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateServiceSetting(ctx context.Context, input *messages.ServiceSettingCreationRequestInput) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateServiceSettingConfiguration(ctx context.Context, input *messages.ServiceSettingConfigurationCreationRequestInput) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateUser(ctx context.Context, input *messages.UserRegistrationInput) (*messages.UserCreationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateUserIngredientPreference(ctx context.Context, input *messages.UserIngredientPreferenceCreationRequestInput) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateUserNotification(ctx context.Context, input *messages.UserNotificationCreationRequestInput) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidIngredient(ctx context.Context, input *messages.ValidIngredientCreationRequestInput) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	created, err := s.dataManager.CreateValidIngredient(ctx, &types.ValidIngredientDatabaseCreationInput{
		StorageTemperatureInCelsius: types.OptionalFloat32Range{
			Min: input.StorageTemperatureInCelsius.Min,
			Max: input.StorageTemperatureInCelsius.Max,
		},
		ID:                     identifiers.New(),
		Warning:                input.Warning,
		IconPath:               input.IconPath,
		PluralName:             input.PluralName,
		StorageInstructions:    input.StorageInstructions,
		Name:                   input.Name,
		Description:            input.Description,
		Slug:                   input.Slug,
		ShoppingSuggestions:    input.ShoppingSuggestions,
		ContainsFish:           input.ContainsFish,
		ContainsShellfish:      input.ContainsShellfish,
		AnimalFlesh:            input.AnimalFlesh,
		ContainsEgg:            input.ContainsEgg,
		IsLiquid:               input.IsLiquid,
		ContainsSoy:            input.ContainsSoy,
		ContainsPeanut:         input.ContainsPeanut,
		AnimalDerived:          input.AnimalDerived,
		RestrictToPreparations: input.RestrictToPreparations,
		ContainsDairy:          input.ContainsDairy,
		ContainsSesame:         input.ContainsSesame,
		ContainsTreeNut:        input.ContainsTreeNut,
		ContainsWheat:          input.ContainsWheat,
		ContainsAlcohol:        input.ContainsAlcohol,
		ContainsGluten:         input.ContainsGluten,
		IsStarch:               input.IsStarch,
		IsProtein:              input.IsProtein,
		IsGrain:                input.IsGrain,
		IsFruit:                input.IsFruit,
		IsSalt:                 input.IsSalt,
		IsFat:                  input.IsFat,
		IsAcid:                 input.IsAcid,
		IsHeat:                 input.IsHeat,
	})
	if err != nil {
		return nil, err
	}

	output := &messages.ValidIngredient{
		CreatedAt:     converters.ConvertTimeToPBTimestamp(created.CreatedAt),
		LastUpdatedAt: converters.ConvertTimePointerToPBTimestamp(created.LastUpdatedAt),
		ArchivedAt:    converters.ConvertTimePointerToPBTimestamp(created.ArchivedAt),
		StorageTemperatureInCelsius: &messages.OptionalFloat32Range{
			Max: created.StorageTemperatureInCelsius.Max,
			Min: created.StorageTemperatureInCelsius.Min,
		},
		IconPath:               created.IconPath,
		Warning:                created.Warning,
		PluralName:             created.PluralName,
		StorageInstructions:    created.StorageInstructions,
		Name:                   created.Name,
		ID:                     created.ID,
		Description:            created.Description,
		Slug:                   created.Slug,
		ShoppingSuggestions:    created.ShoppingSuggestions,
		ContainsShellfish:      created.ContainsShellfish,
		IsLiquid:               created.IsLiquid,
		ContainsPeanut:         created.ContainsPeanut,
		ContainsTreeNut:        created.ContainsTreeNut,
		ContainsEgg:            created.ContainsEgg,
		ContainsWheat:          created.ContainsWheat,
		ContainsSoy:            created.ContainsSoy,
		AnimalDerived:          created.AnimalDerived,
		RestrictToPreparations: created.RestrictToPreparations,
		ContainsSesame:         created.ContainsSesame,
		ContainsFish:           created.ContainsFish,
		ContainsGluten:         created.ContainsGluten,
		ContainsDairy:          created.ContainsDairy,
		ContainsAlcohol:        created.ContainsAlcohol,
		AnimalFlesh:            created.AnimalFlesh,
		IsStarch:               created.IsStarch,
		IsProtein:              created.IsProtein,
		IsGrain:                created.IsGrain,
		IsFruit:                created.IsFruit,
		IsSalt:                 created.IsSalt,
		IsFat:                  created.IsFat,
		IsAcid:                 created.IsAcid,
		IsHeat:                 created.IsHeat,
	}

	return output, nil
}

func (s *Server) CreateValidIngredientGroup(ctx context.Context, input *messages.ValidIngredientGroupCreationRequestInput) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidIngredientMeasurementUnit(ctx context.Context, input *messages.ValidIngredientMeasurementUnitCreationRequestInput) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidIngredientPreparation(ctx context.Context, input *messages.ValidIngredientPreparationCreationRequestInput) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidIngredientState(ctx context.Context, input *messages.ValidIngredientStateCreationRequestInput) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidIngredientStateIngredient(ctx context.Context, input *messages.ValidIngredientStateIngredientCreationRequestInput) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidInstrument(ctx context.Context, input *messages.ValidInstrumentCreationRequestInput) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidMeasurementUnit(ctx context.Context, input *messages.ValidMeasurementUnitCreationRequestInput) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidMeasurementUnitConversion(ctx context.Context, input *messages.ValidMeasurementUnitConversionCreationRequestInput) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidPreparation(ctx context.Context, input *messages.ValidPreparationCreationRequestInput) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidPreparationInstrument(ctx context.Context, input *messages.ValidPreparationInstrumentCreationRequestInput) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidPreparationVessel(ctx context.Context, input *messages.ValidPreparationVesselCreationRequestInput) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateValidVessel(ctx context.Context, input *messages.ValidVesselCreationRequestInput) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateWebhook(ctx context.Context, input *messages.WebhookCreationRequestInput) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) CreateWebhookTriggerEvent(ctx context.Context, request *messages.CreateWebhookTriggerEventRequest) (*messages.WebhookTriggerEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) DestroyAllUserData(ctx context.Context, _ *emptypb.Empty) (*messages.DataDeletionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) FetchUserDataReport(ctx context.Context, request *messages.FetchUserDataReportRequest) (*messages.UserDataCollection, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) FinalizeMealPlan(ctx context.Context, request *messages.FinalizeMealPlanRequest) (*messages.FinalizeMealPlansResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetActiveHousehold(ctx context.Context, _ *emptypb.Empty) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetAuditLogEntriesForHousehold(ctx context.Context, request *messages.GetAuditLogEntriesForHouseholdRequest) (*messages.AuditLogEntry, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetAuditLogEntriesForUser(ctx context.Context, request *messages.GetAuditLogEntriesForUserRequest) (*messages.AuditLogEntry, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetAuditLogEntryByID(ctx context.Context, request *messages.GetAuditLogEntryByIDRequest) (*messages.AuditLogEntry, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetAuthStatus(ctx context.Context, _ *emptypb.Empty) (*messages.UserStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetHousehold(ctx context.Context, request *messages.GetHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetHouseholdInstrumentOwnership(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetHouseholdInstrumentOwnerships(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipsRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetHouseholdInvitation(ctx context.Context, request *messages.GetHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetHouseholdInvitationByID(ctx context.Context, request *messages.GetHouseholdInvitationByIDRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetHouseholds(ctx context.Context, request *messages.GetHouseholdsRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMeal(ctx context.Context, request *messages.GetMealRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlan(ctx context.Context, request *messages.GetMealPlanRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanEvent(ctx context.Context, request *messages.GetMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanEvents(ctx context.Context, request *messages.GetMealPlanEventsRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanGroceryListItem(ctx context.Context, request *messages.GetMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *messages.GetMealPlanGroceryListItemsForMealPlanRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanOption(ctx context.Context, request *messages.GetMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanOptionVote(ctx context.Context, request *messages.GetMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanOptionVotes(ctx context.Context, request *messages.GetMealPlanOptionVotesRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanOptions(ctx context.Context, request *messages.GetMealPlanOptionsRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanTask(ctx context.Context, request *messages.GetMealPlanTaskRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlanTasks(ctx context.Context, request *messages.GetMealPlanTasksRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMealPlansForHousehold(ctx context.Context, request *messages.GetMealPlansForHouseholdRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMeals(ctx context.Context, request *messages.GetMealsRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetMermaidDiagramForRecipe(ctx context.Context, request *messages.GetMermaidDiagramForRecipeRequest) (*messages.GetMermaidDiagramForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetOAuth2Client(ctx context.Context, request *messages.GetOAuth2ClientRequest) (*messages.OAuth2Client, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetOAuth2Clients(ctx context.Context, request *messages.GetOAuth2ClientsRequest) (*messages.OAuth2Client, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRandomValidIngredient(ctx context.Context, _ *emptypb.Empty) (*messages.ValidIngredient, error) {
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

func (s *Server) GetRandomValidInstrument(ctx context.Context, _ *emptypb.Empty) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRandomValidPreparation(ctx context.Context, _ *emptypb.Empty) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRandomValidVessel(ctx context.Context, _ *emptypb.Empty) (*messages.ValidVessel, error) {
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

func (s *Server) GetReceivedHouseholdInvitations(ctx context.Context, request *messages.GetReceivedHouseholdInvitationsRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipe(ctx context.Context, request *messages.GetRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeMealPlanTasks(ctx context.Context, request *messages.GetRecipeMealPlanTasksRequest) (*messages.RecipePrepTaskStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipePrepTask(ctx context.Context, request *messages.GetRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipePrepTasks(ctx context.Context, request *messages.GetRecipePrepTasksRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeRating(ctx context.Context, request *messages.GetRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeRatingsForRecipe(ctx context.Context, request *messages.GetRecipeRatingsForRecipeRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStep(ctx context.Context, request *messages.GetRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepCompletionCondition(ctx context.Context, request *messages.GetRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepCompletionConditions(ctx context.Context, request *messages.GetRecipeStepCompletionConditionsRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepIngredient(ctx context.Context, request *messages.GetRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepIngredients(ctx context.Context, request *messages.GetRecipeStepIngredientsRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepInstrument(ctx context.Context, request *messages.GetRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepInstruments(ctx context.Context, request *messages.GetRecipeStepInstrumentsRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepProduct(ctx context.Context, request *messages.GetRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepProducts(ctx context.Context, request *messages.GetRecipeStepProductsRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepVessel(ctx context.Context, request *messages.GetRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeStepVessels(ctx context.Context, request *messages.GetRecipeStepVesselsRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipeSteps(ctx context.Context, request *messages.GetRecipeStepsRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetRecipes(ctx context.Context, request *messages.GetRecipesRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetSelf(ctx context.Context, _ *emptypb.Empty) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetSentHouseholdInvitations(ctx context.Context, request *messages.GetSentHouseholdInvitationsRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetServiceSetting(ctx context.Context, request *messages.GetServiceSettingRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetServiceSettingConfigurationByName(ctx context.Context, request *messages.GetServiceSettingConfigurationByNameRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetServiceSettingConfigurationsForHousehold(ctx context.Context, request *messages.GetServiceSettingConfigurationsForHouseholdRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetServiceSettingConfigurationsForUser(ctx context.Context, request *messages.GetServiceSettingConfigurationsForUserRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetServiceSettings(ctx context.Context, request *messages.GetServiceSettingsRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetUser(ctx context.Context, request *messages.GetUserRequest) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetUserIngredientPreferences(ctx context.Context, request *messages.GetUserIngredientPreferencesRequest) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetUserNotification(ctx context.Context, request *messages.GetUserNotificationRequest) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetUserNotifications(ctx context.Context, request *messages.GetUserNotificationsRequest) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetUsers(ctx context.Context, request *messages.GetUsersRequest) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredient(ctx context.Context, request *messages.GetValidIngredientRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientGroup(ctx context.Context, request *messages.GetValidIngredientGroupRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientGroups(ctx context.Context, request *messages.GetValidIngredientGroupsRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientMeasurementUnits(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByIngredientRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientPreparations(ctx context.Context, request *messages.GetValidIngredientPreparationsRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *messages.GetValidIngredientPreparationsByIngredientRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationsByPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientState(ctx context.Context, request *messages.GetValidIngredientStateRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientStateIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientStateIngredients(ctx context.Context, request *messages.GetValidIngredientStateIngredientsRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientStateRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredientStates(ctx context.Context, request *messages.GetValidIngredientStatesRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidIngredients(ctx context.Context, request *messages.GetValidIngredientsRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidInstrument(ctx context.Context, request *messages.GetValidInstrumentRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidInstruments(ctx context.Context, request *messages.GetValidInstrumentsRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidMeasurementUnit(ctx context.Context, request *messages.GetValidMeasurementUnitRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidMeasurementUnitConversion(ctx context.Context, request *messages.GetValidMeasurementUnitConversionRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsFromUnitRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsToUnitRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidMeasurementUnits(ctx context.Context, request *messages.GetValidMeasurementUnitsRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparation(ctx context.Context, request *messages.GetValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparationInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparationInstruments(ctx context.Context, request *messages.GetValidPreparationInstrumentsRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentsByInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *messages.GetValidPreparationInstrumentsByPreparationRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparationVessel(ctx context.Context, request *messages.GetValidPreparationVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparationVessels(ctx context.Context, request *messages.GetValidPreparationVesselsRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparationVesselsByPreparation(ctx context.Context, request *messages.GetValidPreparationVesselsByPreparationRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparationVesselsByVessel(ctx context.Context, request *messages.GetValidPreparationVesselsByVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidPreparations(ctx context.Context, request *messages.GetValidPreparationsRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidVessel(ctx context.Context, request *messages.GetValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetValidVessels(ctx context.Context, request *messages.GetValidVesselsRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetWebhook(ctx context.Context, request *messages.GetWebhookRequest) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) GetWebhooks(ctx context.Context, request *messages.GetWebhooksRequest) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) LoginForToken(ctx context.Context, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) PublishArbitraryQueueMessage(ctx context.Context, input *messages.ArbitraryQueueMessageRequestInput) (*messages.ArbitraryQueueMessageResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RedeemPasswordResetToken(ctx context.Context, input *messages.PasswordResetTokenRedemptionRequestInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RefreshTOTPSecret(ctx context.Context, input *messages.TOTPSecretRefreshInput) (*messages.TOTPSecretRefreshResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RejectHouseholdInvitation(ctx context.Context, request *messages.RejectHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RequestEmailVerificationEmail(ctx context.Context, _ *emptypb.Empty) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RequestPasswordResetToken(ctx context.Context, input *messages.PasswordResetTokenCreationRequestInput) (*messages.PasswordResetToken, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RequestUsernameReminder(ctx context.Context, input *messages.UsernameReminderRequestInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RunFinalizeMealPlanWorker(ctx context.Context, request *messages.FinalizeMealPlansRequest) (*messages.FinalizeMealPlansResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RunMealPlanGroceryListInitializerWorker(ctx context.Context, request *messages.InitializeMealPlanGroceryListRequest) (*messages.InitializeMealPlanGroceryListResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) RunMealPlanTaskCreatorWorker(ctx context.Context, request *messages.CreateMealPlanTasksRequest) (*messages.CreateMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForMeals(ctx context.Context, request *messages.SearchForMealsRequest) (*messages.Meal, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForRecipes(ctx context.Context, request *messages.SearchForRecipesRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForServiceSettings(ctx context.Context, request *messages.SearchForServiceSettingsRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForUsers(ctx context.Context, request *messages.SearchForUsersRequest) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForValidIngredientGroups(ctx context.Context, request *messages.SearchForValidIngredientGroupsRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForValidIngredientStates(ctx context.Context, request *messages.SearchForValidIngredientStatesRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForValidIngredients(ctx context.Context, request *messages.SearchForValidIngredientsRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForValidInstruments(ctx context.Context, request *messages.SearchForValidInstrumentsRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForValidMeasurementUnits(ctx context.Context, request *messages.SearchForValidMeasurementUnitsRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForValidPreparations(ctx context.Context, request *messages.SearchForValidPreparationsRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchForValidVessels(ctx context.Context, request *messages.SearchForValidVesselsRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchValidIngredientsByPreparation(ctx context.Context, request *messages.SearchValidIngredientsByPreparationRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *messages.SearchValidMeasurementUnitsByIngredientRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) SetDefaultHousehold(ctx context.Context, request *messages.SetDefaultHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) TransferHouseholdOwnership(ctx context.Context, request *messages.TransferHouseholdOwnershipRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateHousehold(ctx context.Context, request *messages.UpdateHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateHouseholdInstrumentOwnership(ctx context.Context, request *messages.UpdateHouseholdInstrumentOwnershipRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateHouseholdMemberPermissions(ctx context.Context, request *messages.UpdateHouseholdMemberPermissionsRequest) (*messages.UserPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateMealPlan(ctx context.Context, request *messages.UpdateMealPlanRequest) (*messages.MealPlan, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateMealPlanEvent(ctx context.Context, request *messages.UpdateMealPlanEventRequest) (*messages.MealPlanEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateMealPlanGroceryListItem(ctx context.Context, request *messages.UpdateMealPlanGroceryListItemRequest) (*messages.MealPlanGroceryListItem, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateMealPlanOption(ctx context.Context, request *messages.UpdateMealPlanOptionRequest) (*messages.MealPlanOption, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateMealPlanOptionVote(ctx context.Context, request *messages.UpdateMealPlanOptionVoteRequest) (*messages.MealPlanOptionVote, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateMealPlanTaskStatus(ctx context.Context, request *messages.UpdateMealPlanTaskStatusRequest) (*messages.MealPlanTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdatePassword(ctx context.Context, input *messages.PasswordUpdateInput) (*messages.PasswordResetResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipe(ctx context.Context, request *messages.UpdateRecipeRequest) (*messages.Recipe, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipePrepTask(ctx context.Context, request *messages.UpdateRecipePrepTaskRequest) (*messages.RecipePrepTask, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipeRating(ctx context.Context, request *messages.UpdateRecipeRatingRequest) (*messages.RecipeRating, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipeStep(ctx context.Context, request *messages.UpdateRecipeStepRequest) (*messages.RecipeStep, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipeStepCompletionCondition(ctx context.Context, request *messages.UpdateRecipeStepCompletionConditionRequest) (*messages.RecipeStepCompletionCondition, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipeStepIngredient(ctx context.Context, request *messages.UpdateRecipeStepIngredientRequest) (*messages.RecipeStepIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipeStepInstrument(ctx context.Context, request *messages.UpdateRecipeStepInstrumentRequest) (*messages.RecipeStepInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipeStepProduct(ctx context.Context, request *messages.UpdateRecipeStepProductRequest) (*messages.RecipeStepProduct, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateRecipeStepVessel(ctx context.Context, request *messages.UpdateRecipeStepVesselRequest) (*messages.RecipeStepVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateServiceSettingConfiguration(ctx context.Context, request *messages.UpdateServiceSettingConfigurationRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateUserDetails(ctx context.Context, input *messages.UserDetailsUpdateRequestInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateUserEmailAddress(ctx context.Context, input *messages.UserEmailAddressUpdateInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateUserIngredientPreference(ctx context.Context, request *messages.UpdateUserIngredientPreferenceRequest) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateUserNotification(ctx context.Context, request *messages.UpdateUserNotificationRequest) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateUserUsername(ctx context.Context, input *messages.UsernameUpdateInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidIngredient(ctx context.Context, request *messages.UpdateValidIngredientRequest) (*messages.ValidIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidIngredientGroup(ctx context.Context, request *messages.UpdateValidIngredientGroupRequest) (*messages.ValidIngredientGroup, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *messages.UpdateValidIngredientMeasurementUnitRequest) (*messages.ValidIngredientMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidIngredientPreparation(ctx context.Context, request *messages.UpdateValidIngredientPreparationRequest) (*messages.ValidIngredientPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidIngredientState(ctx context.Context, request *messages.UpdateValidIngredientStateRequest) (*messages.ValidIngredientState, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidIngredientStateIngredient(ctx context.Context, request *messages.UpdateValidIngredientStateIngredientRequest) (*messages.ValidIngredientStateIngredient, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidInstrument(ctx context.Context, request *messages.UpdateValidInstrumentRequest) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidMeasurementUnit(ctx context.Context, request *messages.UpdateValidMeasurementUnitRequest) (*messages.ValidMeasurementUnit, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidMeasurementUnitConversion(ctx context.Context, request *messages.UpdateValidMeasurementUnitConversionRequest) (*messages.ValidMeasurementUnitConversion, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidPreparation(ctx context.Context, request *messages.UpdateValidPreparationRequest) (*messages.ValidPreparation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidPreparationInstrument(ctx context.Context, request *messages.UpdateValidPreparationInstrumentRequest) (*messages.ValidPreparationInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidPreparationVessel(ctx context.Context, request *messages.UpdateValidPreparationVesselRequest) (*messages.ValidPreparationVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UpdateValidVessel(ctx context.Context, request *messages.UpdateValidVesselRequest) (*messages.ValidVessel, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) UploadUserAvatar(ctx context.Context, input *messages.AvatarUpdateInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) VerifyEmailAddress(ctx context.Context, input *messages.EmailAddressVerificationRequestInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}

func (s *Server) VerifyTOTPSecret(ctx context.Context, input *messages.TOTPSecretVerificationInput) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, errUnimplemented
}
