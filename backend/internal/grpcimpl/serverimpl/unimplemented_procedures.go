package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
)

var _ service.EatingServiceServer = (*Server)(nil)

func (s *Server) AcceptHouseholdInvitation(ctx context.Context, request *messages.AcceptHouseholdInvitationRequest) (*messages.AcceptHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) AdminUpdateUserStatus(ctx context.Context, request *messages.AdminUpdateUserStatusRequest) (*messages.AdminUpdateUserStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) AggregateUserDataReport(ctx context.Context, request *messages.AggregateUserDataReportRequest) (*messages.AggregateUserDataReportResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveHousehold(ctx context.Context, request *messages.ArchiveHouseholdRequest) (*messages.ArchiveHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveHouseholdInstrumentOwnership(ctx context.Context, request *messages.ArchiveHouseholdInstrumentOwnershipRequest) (*messages.ArchiveHouseholdInstrumentOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMeal(ctx context.Context, request *messages.ArchiveMealRequest) (*messages.ArchiveMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlan(ctx context.Context, request *messages.ArchiveMealPlanRequest) (*messages.ArchiveMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanEvent(ctx context.Context, request *messages.ArchiveMealPlanEventRequest) (*messages.ArchiveMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanGroceryListItem(ctx context.Context, request *messages.ArchiveMealPlanGroceryListItemRequest) (*messages.ArchiveMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanOption(ctx context.Context, request *messages.ArchiveMealPlanOptionRequest) (*messages.ArchiveMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveMealPlanOptionVote(ctx context.Context, request *messages.ArchiveMealPlanOptionVoteRequest) (*messages.ArchiveMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveOAuth2Client(ctx context.Context, request *messages.ArchiveOAuth2ClientRequest) (*messages.ArchiveOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipe(ctx context.Context, request *messages.ArchiveRecipeRequest) (*messages.ArchiveRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipePrepTask(ctx context.Context, request *messages.ArchiveRecipePrepTaskRequest) (*messages.ArchiveRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeRating(ctx context.Context, request *messages.ArchiveRecipeRatingRequest) (*messages.ArchiveRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStep(ctx context.Context, request *messages.ArchiveRecipeStepRequest) (*messages.ArchiveRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepCompletionCondition(ctx context.Context, request *messages.ArchiveRecipeStepCompletionConditionRequest) (*messages.ArchiveRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepIngredient(ctx context.Context, request *messages.ArchiveRecipeStepIngredientRequest) (*messages.ArchiveRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepInstrument(ctx context.Context, request *messages.ArchiveRecipeStepInstrumentRequest) (*messages.ArchiveRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepProduct(ctx context.Context, request *messages.ArchiveRecipeStepProductRequest) (*messages.ArchiveRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveRecipeStepVessel(ctx context.Context, request *messages.ArchiveRecipeStepVesselRequest) (*messages.ArchiveRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveServiceSetting(ctx context.Context, request *messages.ArchiveServiceSettingRequest) (*messages.ArchiveServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveServiceSettingConfiguration(ctx context.Context, request *messages.ArchiveServiceSettingConfigurationRequest) (*messages.ArchiveServiceSettingConfigurationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveUser(ctx context.Context, request *messages.ArchiveUserRequest) (*messages.ArchiveUserResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveUserIngredientPreference(ctx context.Context, request *messages.ArchiveUserIngredientPreferenceRequest) (*messages.ArchiveUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveUserMembership(ctx context.Context, request *messages.ArchiveUserMembershipRequest) (*messages.ArchiveUserMembershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientGroup(ctx context.Context, request *messages.ArchiveValidIngredientGroupRequest) (*messages.ArchiveValidIngredientGroupResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientMeasurementUnit(ctx context.Context, request *messages.ArchiveValidIngredientMeasurementUnitRequest) (*messages.ArchiveValidIngredientMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientPreparation(ctx context.Context, request *messages.ArchiveValidIngredientPreparationRequest) (*messages.ArchiveValidIngredientPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientState(ctx context.Context, request *messages.ArchiveValidIngredientStateRequest) (*messages.ArchiveValidIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidIngredientStateIngredient(ctx context.Context, request *messages.ArchiveValidIngredientStateIngredientRequest) (*messages.ArchiveValidIngredientStateIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidInstrument(ctx context.Context, request *messages.ArchiveValidInstrumentRequest) (*messages.ArchiveValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidMeasurementUnit(ctx context.Context, request *messages.ArchiveValidMeasurementUnitRequest) (*messages.ArchiveValidMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidMeasurementUnitConversion(ctx context.Context, request *messages.ArchiveValidMeasurementUnitConversionRequest) (*messages.ArchiveValidMeasurementUnitConversionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparation(ctx context.Context, request *messages.ArchiveValidPreparationRequest) (*messages.ArchiveValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparationInstrument(ctx context.Context, request *messages.ArchiveValidPreparationInstrumentRequest) (*messages.ArchiveValidPreparationInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidPreparationVessel(ctx context.Context, request *messages.ArchiveValidPreparationVesselRequest) (*messages.ArchiveValidPreparationVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveValidVessel(ctx context.Context, request *messages.ArchiveValidVesselRequest) (*messages.ArchiveValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveWebhook(ctx context.Context, request *messages.ArchiveWebhookRequest) (*messages.ArchiveWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveWebhookTriggerEvent(ctx context.Context, request *messages.ArchiveWebhookTriggerEventRequest) (*messages.ArchiveWebhookTriggerEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CancelHouseholdInvitation(ctx context.Context, request *messages.CancelHouseholdInvitationRequest) (*messages.CancelHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CheckForReadiness(ctx context.Context, request *messages.CheckForReadinessRequest) (*messages.CheckForReadinessResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CloneRecipe(ctx context.Context, request *messages.CloneRecipeRequest) (*messages.CloneRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateHousehold(ctx context.Context, request *messages.CreateHouseholdRequest) (*messages.CreateHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateHouseholdInstrumentOwnership(ctx context.Context, request *messages.CreateHouseholdInstrumentOwnershipRequest) (*messages.CreateHouseholdInstrumentOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateHouseholdInvitation(ctx context.Context, request *messages.CreateHouseholdInvitationRequest) (*messages.CreateHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMeal(ctx context.Context, request *messages.CreateMealRequest) (*messages.CreateMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlan(ctx context.Context, request *messages.CreateMealPlanRequest) (*messages.CreateMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanEvent(ctx context.Context, request *messages.CreateMealPlanEventRequest) (*messages.CreateMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanGroceryListItem(ctx context.Context, request *messages.CreateMealPlanGroceryListItemRequest) (*messages.CreateMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanOption(ctx context.Context, request *messages.CreateMealPlanOptionRequest) (*messages.CreateMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanOptionVote(ctx context.Context, request *messages.CreateMealPlanOptionVoteRequest) (*messages.CreateMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateMealPlanTask(ctx context.Context, request *messages.CreateMealPlanTaskRequest) (*messages.CreateMealPlanTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateOAuth2Client(ctx context.Context, request *messages.CreateOAuth2ClientRequest) (*messages.CreateOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipe(ctx context.Context, request *messages.CreateRecipeRequest) (*messages.CreateRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipePrepTask(ctx context.Context, request *messages.CreateRecipePrepTaskRequest) (*messages.CreateRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeRating(ctx context.Context, request *messages.CreateRecipeRatingRequest) (*messages.CreateRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStep(ctx context.Context, request *messages.CreateRecipeStepRequest) (*messages.CreateRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepCompletionCondition(ctx context.Context, request *messages.CreateRecipeStepCompletionConditionRequest) (*messages.CreateRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepIngredient(ctx context.Context, request *messages.CreateRecipeStepIngredientRequest) (*messages.CreateRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepInstrument(ctx context.Context, request *messages.CreateRecipeStepInstrumentRequest) (*messages.CreateRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepProduct(ctx context.Context, request *messages.CreateRecipeStepProductRequest) (*messages.CreateRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateRecipeStepVessel(ctx context.Context, request *messages.CreateRecipeStepVesselRequest) (*messages.CreateRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateServiceSetting(ctx context.Context, request *messages.CreateServiceSettingRequest) (*messages.CreateServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateServiceSettingConfiguration(ctx context.Context, request *messages.CreateServiceSettingConfigurationRequest) (*messages.CreateServiceSettingConfigurationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateUserIngredientPreference(ctx context.Context, request *messages.CreateUserIngredientPreferenceRequest) (*messages.CreateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateUserNotification(ctx context.Context, request *messages.CreateUserNotificationRequest) (*messages.CreateUserNotificationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientGroup(ctx context.Context, request *messages.CreateValidIngredientGroupRequest) (*messages.CreateValidIngredientGroupResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientMeasurementUnit(ctx context.Context, request *messages.CreateValidIngredientMeasurementUnitRequest) (*messages.CreateValidIngredientMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientPreparation(ctx context.Context, request *messages.CreateValidIngredientPreparationRequest) (*messages.CreateValidIngredientPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientState(ctx context.Context, request *messages.CreateValidIngredientStateRequest) (*messages.CreateValidIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidIngredientStateIngredient(ctx context.Context, request *messages.CreateValidIngredientStateIngredientRequest) (*messages.CreateValidIngredientStateIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidInstrument(ctx context.Context, request *messages.CreateValidInstrumentRequest) (*messages.CreateValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidMeasurementUnit(ctx context.Context, request *messages.CreateValidMeasurementUnitRequest) (*messages.CreateValidMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidMeasurementUnitConversion(ctx context.Context, request *messages.CreateValidMeasurementUnitConversionRequest) (*messages.CreateValidMeasurementUnitConversionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidPreparation(ctx context.Context, request *messages.CreateValidPreparationRequest) (*messages.CreateValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidPreparationInstrument(ctx context.Context, request *messages.CreateValidPreparationInstrumentRequest) (*messages.CreateValidPreparationInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidPreparationVessel(ctx context.Context, request *messages.CreateValidPreparationVesselRequest) (*messages.CreateValidPreparationVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateValidVessel(ctx context.Context, request *messages.CreateValidVesselRequest) (*messages.CreateValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateWebhook(ctx context.Context, request *messages.CreateWebhookRequest) (*messages.CreateWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateWebhookTriggerEvent(ctx context.Context, request *messages.CreateWebhookTriggerEventRequest) (*messages.CreateWebhookTriggerEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) DestroyAllUserData(ctx context.Context, request *messages.DestroyAllUserDataRequest) (*messages.DestroyAllUserDataResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) FetchUserDataReport(ctx context.Context, request *messages.FetchUserDataReportRequest) (*messages.FetchUserDataReportResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) FinalizeMealPlan(ctx context.Context, request *messages.FinalizeMealPlanRequest) (*messages.FinalizeMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetActiveHousehold(ctx context.Context, request *messages.GetActiveHouseholdRequest) (*messages.GetActiveHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetAuditLogEntriesForHousehold(ctx context.Context, request *messages.GetAuditLogEntriesForHouseholdRequest) (*messages.GetAuditLogEntriesForHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetAuditLogEntriesForUser(ctx context.Context, request *messages.GetAuditLogEntriesForUserRequest) (*messages.GetAuditLogEntriesForUserResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetAuditLogEntryByID(ctx context.Context, request *messages.GetAuditLogEntryByIDRequest) (*messages.GetAuditLogEntryByIDResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHousehold(ctx context.Context, request *messages.GetHouseholdRequest) (*messages.GetHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInstrumentOwnership(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipRequest) (*messages.GetHouseholdInstrumentOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInstrumentOwnerships(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipsRequest) (*messages.GetHouseholdInstrumentOwnershipsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInvitation(ctx context.Context, request *messages.GetHouseholdInvitationRequest) (*messages.GetHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInvitationByID(ctx context.Context, request *messages.GetHouseholdInvitationByIDRequest) (*messages.GetHouseholdInvitationByIDResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholds(ctx context.Context, request *messages.GetHouseholdsRequest) (*messages.GetHouseholdsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMeal(ctx context.Context, request *messages.GetMealRequest) (*messages.GetMealResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlan(ctx context.Context, request *messages.GetMealPlanRequest) (*messages.GetMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanEvent(ctx context.Context, request *messages.GetMealPlanEventRequest) (*messages.GetMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanEvents(ctx context.Context, request *messages.GetMealPlanEventsRequest) (*messages.GetMealPlanEventsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanGroceryListItem(ctx context.Context, request *messages.GetMealPlanGroceryListItemRequest) (*messages.GetMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanGroceryListItemsForMealPlan(ctx context.Context, request *messages.GetMealPlanGroceryListItemsForMealPlanRequest) (*messages.GetMealPlanGroceryListItemsForMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOption(ctx context.Context, request *messages.GetMealPlanOptionRequest) (*messages.GetMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptionVote(ctx context.Context, request *messages.GetMealPlanOptionVoteRequest) (*messages.GetMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptionVotes(ctx context.Context, request *messages.GetMealPlanOptionVotesRequest) (*messages.GetMealPlanOptionVotesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanOptions(ctx context.Context, request *messages.GetMealPlanOptionsRequest) (*messages.GetMealPlanOptionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanTask(ctx context.Context, request *messages.GetMealPlanTaskRequest) (*messages.GetMealPlanTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlanTasks(ctx context.Context, request *messages.GetMealPlanTasksRequest) (*messages.GetMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMealPlansForHousehold(ctx context.Context, request *messages.GetMealPlansForHouseholdRequest) (*messages.GetMealPlansForHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetMeals(ctx context.Context, request *messages.GetMealsRequest) (*messages.GetMealsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetOAuth2Client(ctx context.Context, request *messages.GetOAuth2ClientRequest) (*messages.GetOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetOAuth2Clients(ctx context.Context, request *messages.GetOAuth2ClientsRequest) (*messages.GetOAuth2ClientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidInstrument(ctx context.Context, request *messages.GetRandomValidInstrumentRequest) (*messages.GetRandomValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidPreparation(ctx context.Context, request *messages.GetRandomValidPreparationRequest) (*messages.GetRandomValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidVessel(ctx context.Context, request *messages.GetRandomValidVesselRequest) (*messages.GetRandomValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetReceivedHouseholdInvitations(ctx context.Context, request *messages.GetReceivedHouseholdInvitationsRequest) (*messages.GetReceivedHouseholdInvitationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipe(ctx context.Context, request *messages.GetRecipeRequest) (*messages.GetRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeMealPlanTasks(ctx context.Context, request *messages.GetRecipeMealPlanTasksRequest) (*messages.GetRecipeMealPlanTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipePrepTask(ctx context.Context, request *messages.GetRecipePrepTaskRequest) (*messages.GetRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipePrepTasks(ctx context.Context, request *messages.GetRecipePrepTasksRequest) (*messages.GetRecipePrepTasksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeRating(ctx context.Context, request *messages.GetRecipeRatingRequest) (*messages.GetRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeRatingsForRecipe(ctx context.Context, request *messages.GetRecipeRatingsForRecipeRequest) (*messages.GetRecipeRatingsForRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStep(ctx context.Context, request *messages.GetRecipeStepRequest) (*messages.GetRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepCompletionCondition(ctx context.Context, request *messages.GetRecipeStepCompletionConditionRequest) (*messages.GetRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepCompletionConditions(ctx context.Context, request *messages.GetRecipeStepCompletionConditionsRequest) (*messages.GetRecipeStepCompletionConditionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepIngredient(ctx context.Context, request *messages.GetRecipeStepIngredientRequest) (*messages.GetRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepIngredients(ctx context.Context, request *messages.GetRecipeStepIngredientsRequest) (*messages.GetRecipeStepIngredientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepInstrument(ctx context.Context, request *messages.GetRecipeStepInstrumentRequest) (*messages.GetRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepInstruments(ctx context.Context, request *messages.GetRecipeStepInstrumentsRequest) (*messages.GetRecipeStepInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepProduct(ctx context.Context, request *messages.GetRecipeStepProductRequest) (*messages.GetRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepProducts(ctx context.Context, request *messages.GetRecipeStepProductsRequest) (*messages.GetRecipeStepProductsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepVessel(ctx context.Context, request *messages.GetRecipeStepVesselRequest) (*messages.GetRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeStepVessels(ctx context.Context, request *messages.GetRecipeStepVesselsRequest) (*messages.GetRecipeStepVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipeSteps(ctx context.Context, request *messages.GetRecipeStepsRequest) (*messages.GetRecipeStepsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRecipes(ctx context.Context, request *messages.GetRecipesRequest) (*messages.GetRecipesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetSelf(ctx context.Context, request *messages.GetSelfRequest) (*messages.GetSelfResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetSentHouseholdInvitations(ctx context.Context, request *messages.GetSentHouseholdInvitationsRequest) (*messages.GetSentHouseholdInvitationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSetting(ctx context.Context, request *messages.GetServiceSettingRequest) (*messages.GetServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettingConfigurationByName(ctx context.Context, request *messages.GetServiceSettingConfigurationByNameRequest) (*messages.GetServiceSettingConfigurationByNameResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettingConfigurationsForHousehold(ctx context.Context, request *messages.GetServiceSettingConfigurationsForHouseholdRequest) (*messages.GetServiceSettingConfigurationsForHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettingConfigurationsForUser(ctx context.Context, request *messages.GetServiceSettingConfigurationsForUserRequest) (*messages.GetServiceSettingConfigurationsForUserResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettings(ctx context.Context, request *messages.GetServiceSettingsRequest) (*messages.GetServiceSettingsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUser(ctx context.Context, request *messages.GetUserRequest) (*messages.GetUserResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserIngredientPreferences(ctx context.Context, request *messages.GetUserIngredientPreferencesRequest) (*messages.GetUserIngredientPreferencesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserNotification(ctx context.Context, request *messages.GetUserNotificationRequest) (*messages.GetUserNotificationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserNotifications(ctx context.Context, request *messages.GetUserNotificationsRequest) (*messages.GetUserNotificationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUsers(ctx context.Context, request *messages.GetUsersRequest) (*messages.GetUsersResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientGroup(ctx context.Context, request *messages.GetValidIngredientGroupRequest) (*messages.GetValidIngredientGroupResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientGroups(ctx context.Context, request *messages.GetValidIngredientGroupsRequest) (*messages.GetValidIngredientGroupsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitRequest) (*messages.GetValidIngredientMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnits(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsRequest) (*messages.GetValidIngredientMeasurementUnitsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnitsByIngredient(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByIngredientRequest) (*messages.GetValidIngredientMeasurementUnitsByIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, request *messages.GetValidIngredientMeasurementUnitsByMeasurementUnitRequest) (*messages.GetValidIngredientMeasurementUnitsByMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationRequest) (*messages.GetValidIngredientPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparations(ctx context.Context, request *messages.GetValidIngredientPreparationsRequest) (*messages.GetValidIngredientPreparationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparationsByIngredient(ctx context.Context, request *messages.GetValidIngredientPreparationsByIngredientRequest) (*messages.GetValidIngredientPreparationsByIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientPreparationsByPreparation(ctx context.Context, request *messages.GetValidIngredientPreparationsByPreparationRequest) (*messages.GetValidIngredientPreparationsByPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientState(ctx context.Context, request *messages.GetValidIngredientStateRequest) (*messages.GetValidIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientRequest) (*messages.GetValidIngredientStateIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredients(ctx context.Context, request *messages.GetValidIngredientStateIngredientsRequest) (*messages.GetValidIngredientStateIngredientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredientsByIngredient(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientRequest) (*messages.GetValidIngredientStateIngredientsByIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStateIngredientsByIngredientState(ctx context.Context, request *messages.GetValidIngredientStateIngredientsByIngredientStateRequest) (*messages.GetValidIngredientStateIngredientsByIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredientStates(ctx context.Context, request *messages.GetValidIngredientStatesRequest) (*messages.GetValidIngredientStatesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidIngredients(ctx context.Context, request *messages.GetValidIngredientsRequest) (*messages.GetValidIngredientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidInstrument(ctx context.Context, request *messages.GetValidInstrumentRequest) (*messages.GetValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidInstruments(ctx context.Context, request *messages.GetValidInstrumentsRequest) (*messages.GetValidInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnit(ctx context.Context, request *messages.GetValidMeasurementUnitRequest) (*messages.GetValidMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversion(ctx context.Context, request *messages.GetValidMeasurementUnitConversionRequest) (*messages.GetValidMeasurementUnitConversionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsFromUnitRequest) (*messages.GetValidMeasurementUnitConversionsFromUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, request *messages.GetValidMeasurementUnitConversionsToUnitRequest) (*messages.GetValidMeasurementUnitConversionsToUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidMeasurementUnits(ctx context.Context, request *messages.GetValidMeasurementUnitsRequest) (*messages.GetValidMeasurementUnitsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparation(ctx context.Context, request *messages.GetValidPreparationRequest) (*messages.GetValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentRequest) (*messages.GetValidPreparationInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstruments(ctx context.Context, request *messages.GetValidPreparationInstrumentsRequest) (*messages.GetValidPreparationInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrumentsByInstrument(ctx context.Context, request *messages.GetValidPreparationInstrumentsByInstrumentRequest) (*messages.GetValidPreparationInstrumentsByInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationInstrumentsByPreparation(ctx context.Context, request *messages.GetValidPreparationInstrumentsByPreparationRequest) (*messages.GetValidPreparationInstrumentsByPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVessel(ctx context.Context, request *messages.GetValidPreparationVesselRequest) (*messages.GetValidPreparationVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVessels(ctx context.Context, request *messages.GetValidPreparationVesselsRequest) (*messages.GetValidPreparationVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVesselsByPreparation(ctx context.Context, request *messages.GetValidPreparationVesselsByPreparationRequest) (*messages.GetValidPreparationVesselsByPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparationVesselsByVessel(ctx context.Context, request *messages.GetValidPreparationVesselsByVesselRequest) (*messages.GetValidPreparationVesselsByVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidPreparations(ctx context.Context, request *messages.GetValidPreparationsRequest) (*messages.GetValidPreparationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidVessel(ctx context.Context, request *messages.GetValidVesselRequest) (*messages.GetValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetValidVessels(ctx context.Context, request *messages.GetValidVesselsRequest) (*messages.GetValidVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetWebhook(ctx context.Context, request *messages.GetWebhookRequest) (*messages.GetWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetWebhooks(ctx context.Context, request *messages.GetWebhooksRequest) (*messages.GetWebhooksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) PublishArbitraryQueueMessage(ctx context.Context, request *messages.PublishArbitraryQueueMessageRequest) (*messages.PublishArbitraryQueueMessageResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RedeemPasswordResetToken(ctx context.Context, request *messages.RedeemPasswordResetTokenRequest) (*messages.RedeemPasswordResetTokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RefreshTOTPSecret(ctx context.Context, request *messages.RefreshTOTPSecretRequest) (*messages.RefreshTOTPSecretResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RejectHouseholdInvitation(ctx context.Context, request *messages.RejectHouseholdInvitationRequest) (*messages.RejectHouseholdInvitationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RequestEmailVerificationEmail(ctx context.Context, request *messages.RequestEmailVerificationEmailRequest) (*messages.RequestEmailVerificationEmailResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RequestPasswordResetToken(ctx context.Context, request *messages.RequestPasswordResetTokenRequest) (*messages.RequestPasswordResetTokenResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RequestUsernameReminder(ctx context.Context, request *messages.RequestUsernameReminderRequest) (*messages.RequestUsernameReminderResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RunFinalizeMealPlanWorker(ctx context.Context, request *messages.RunFinalizeMealPlanWorkerRequest) (*messages.RunFinalizeMealPlanWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RunMealPlanGroceryListInitializerWorker(ctx context.Context, request *messages.RunMealPlanGroceryListInitializerWorkerRequest) (*messages.RunMealPlanGroceryListInitializerWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RunMealPlanTaskCreatorWorker(ctx context.Context, request *messages.RunMealPlanTaskCreatorWorkerRequest) (*messages.RunMealPlanTaskCreatorWorkerResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForMeals(ctx context.Context, request *messages.SearchForMealsRequest) (*messages.SearchForMealsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForRecipes(ctx context.Context, request *messages.SearchForRecipesRequest) (*messages.SearchForRecipesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForServiceSettings(ctx context.Context, request *messages.SearchForServiceSettingsRequest) (*messages.SearchForServiceSettingsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForUsers(ctx context.Context, request *messages.SearchForUsersRequest) (*messages.SearchForUsersResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidIngredientGroups(ctx context.Context, request *messages.SearchForValidIngredientGroupsRequest) (*messages.SearchForValidIngredientGroupsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidIngredientStates(ctx context.Context, request *messages.SearchForValidIngredientStatesRequest) (*messages.SearchForValidIngredientStatesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidInstruments(ctx context.Context, request *messages.SearchForValidInstrumentsRequest) (*messages.SearchForValidInstrumentsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidMeasurementUnits(ctx context.Context, request *messages.SearchForValidMeasurementUnitsRequest) (*messages.SearchForValidMeasurementUnitsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidPreparations(ctx context.Context, request *messages.SearchForValidPreparationsRequest) (*messages.SearchForValidPreparationsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForValidVessels(ctx context.Context, request *messages.SearchForValidVesselsRequest) (*messages.SearchForValidVesselsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchValidMeasurementUnitsByIngredient(ctx context.Context, request *messages.SearchValidMeasurementUnitsByIngredientRequest) (*messages.SearchValidMeasurementUnitsByIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SetDefaultHousehold(ctx context.Context, request *messages.SetDefaultHouseholdRequest) (*messages.SetDefaultHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) TransferHouseholdOwnership(ctx context.Context, request *messages.TransferHouseholdOwnershipRequest) (*messages.TransferHouseholdOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHousehold(ctx context.Context, request *messages.UpdateHouseholdRequest) (*messages.UpdateHouseholdResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHouseholdInstrumentOwnership(ctx context.Context, request *messages.UpdateHouseholdInstrumentOwnershipRequest) (*messages.UpdateHouseholdInstrumentOwnershipResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHouseholdMemberPermissions(ctx context.Context, request *messages.UpdateHouseholdMemberPermissionsRequest) (*messages.UpdateHouseholdMemberPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlan(ctx context.Context, request *messages.UpdateMealPlanRequest) (*messages.UpdateMealPlanResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanEvent(ctx context.Context, request *messages.UpdateMealPlanEventRequest) (*messages.UpdateMealPlanEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanGroceryListItem(ctx context.Context, request *messages.UpdateMealPlanGroceryListItemRequest) (*messages.UpdateMealPlanGroceryListItemResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanOption(ctx context.Context, request *messages.UpdateMealPlanOptionRequest) (*messages.UpdateMealPlanOptionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanOptionVote(ctx context.Context, request *messages.UpdateMealPlanOptionVoteRequest) (*messages.UpdateMealPlanOptionVoteResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateMealPlanTaskStatus(ctx context.Context, request *messages.UpdateMealPlanTaskStatusRequest) (*messages.UpdateMealPlanTaskStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdatePassword(ctx context.Context, request *messages.UpdatePasswordRequest) (*messages.UpdatePasswordResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipe(ctx context.Context, request *messages.UpdateRecipeRequest) (*messages.UpdateRecipeResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipePrepTask(ctx context.Context, request *messages.UpdateRecipePrepTaskRequest) (*messages.UpdateRecipePrepTaskResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeRating(ctx context.Context, request *messages.UpdateRecipeRatingRequest) (*messages.UpdateRecipeRatingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStep(ctx context.Context, request *messages.UpdateRecipeStepRequest) (*messages.UpdateRecipeStepResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepCompletionCondition(ctx context.Context, request *messages.UpdateRecipeStepCompletionConditionRequest) (*messages.UpdateRecipeStepCompletionConditionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepIngredient(ctx context.Context, request *messages.UpdateRecipeStepIngredientRequest) (*messages.UpdateRecipeStepIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepInstrument(ctx context.Context, request *messages.UpdateRecipeStepInstrumentRequest) (*messages.UpdateRecipeStepInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepProduct(ctx context.Context, request *messages.UpdateRecipeStepProductRequest) (*messages.UpdateRecipeStepProductResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateRecipeStepVessel(ctx context.Context, request *messages.UpdateRecipeStepVesselRequest) (*messages.UpdateRecipeStepVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateServiceSettingConfiguration(ctx context.Context, request *messages.UpdateServiceSettingConfigurationRequest) (*messages.UpdateServiceSettingConfigurationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateUserDetails(ctx context.Context, request *messages.UpdateUserDetailsRequest) (*messages.UpdateUserDetailsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateUserEmailAddress(ctx context.Context, request *messages.UpdateUserEmailAddressRequest) (*messages.UpdateUserEmailAddressResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateUserIngredientPreference(ctx context.Context, request *messages.UpdateUserIngredientPreferenceRequest) (*messages.UpdateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateUserNotification(ctx context.Context, request *messages.UpdateUserNotificationRequest) (*messages.UpdateUserNotificationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateUserUsername(ctx context.Context, request *messages.UpdateUserUsernameRequest) (*messages.UpdateUserUsernameResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientGroup(ctx context.Context, request *messages.UpdateValidIngredientGroupRequest) (*messages.UpdateValidIngredientGroupResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientMeasurementUnit(ctx context.Context, request *messages.UpdateValidIngredientMeasurementUnitRequest) (*messages.UpdateValidIngredientMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientPreparation(ctx context.Context, request *messages.UpdateValidIngredientPreparationRequest) (*messages.UpdateValidIngredientPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientState(ctx context.Context, request *messages.UpdateValidIngredientStateRequest) (*messages.UpdateValidIngredientStateResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidIngredientStateIngredient(ctx context.Context, request *messages.UpdateValidIngredientStateIngredientRequest) (*messages.UpdateValidIngredientStateIngredientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidInstrument(ctx context.Context, request *messages.UpdateValidInstrumentRequest) (*messages.UpdateValidInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidMeasurementUnit(ctx context.Context, request *messages.UpdateValidMeasurementUnitRequest) (*messages.UpdateValidMeasurementUnitResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidMeasurementUnitConversion(ctx context.Context, request *messages.UpdateValidMeasurementUnitConversionRequest) (*messages.UpdateValidMeasurementUnitConversionResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparation(ctx context.Context, request *messages.UpdateValidPreparationRequest) (*messages.UpdateValidPreparationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparationInstrument(ctx context.Context, request *messages.UpdateValidPreparationInstrumentRequest) (*messages.UpdateValidPreparationInstrumentResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidPreparationVessel(ctx context.Context, request *messages.UpdateValidPreparationVesselRequest) (*messages.UpdateValidPreparationVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateValidVessel(ctx context.Context, request *messages.UpdateValidVesselRequest) (*messages.UpdateValidVesselResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UploadUserAvatar(ctx context.Context, request *messages.UploadUserAvatarRequest) (*messages.UploadUserAvatarResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) VerifyEmailAddress(ctx context.Context, request *messages.VerifyEmailAddressRequest) (*messages.VerifyEmailAddressResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}
