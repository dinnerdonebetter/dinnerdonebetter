package workers

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// MealPlanFinalizationHandler finalizes a meal plan.
func (s *service) MealPlanFinalizationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	logger.Info("meal plan finalization worker invoked")

	var request *types.FinalizeMealPlansRequest
	if err := s.encoderDecoder.DecodeRequest(ctx, req, &request); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusBadRequest)
		return
	}

	response := &types.FinalizeMealPlansResponse{}
	if request.ReturnCount {
		count, err := s.mealPlanFinalizationWorker.FinalizeExpiredMealPlans(ctx, nil)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "finalizing expired meal plans")
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusInternalServerError)
			return
		}

		response.Count = count
	} else {
		if err := s.mealPlanFinalizationWorker.FinalizeExpiredMealPlansWithoutReturningCount(ctx, nil); err != nil {
			observability.AcknowledgeError(err, logger, span, "finalizing expired meal plans")
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusInternalServerError)
			return
		}
	}

	logger.WithValue("finalized_count", response.Count).Info("meal plan finalization worker completed")

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, response, http.StatusAccepted)
}

// MealPlanGroceryListInitializationHandler initializes a grocery list for a given meal plan.
func (s *service) MealPlanGroceryListInitializationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	logger.Info("meal plan grocery list initialization worker invoked")

	if err := s.mealPlanGroceryListInitializer.InitializeGroceryListsForFinalizedMealPlans(ctx, nil); err != nil {
		observability.AcknowledgeError(err, logger, span, "finalizing expired meal plans")
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusInternalServerError)
		return
	}

	logger.Info("meal plan grocery list initialization worker completed")

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusAccepted)
}

// MealPlanTaskCreationHandler creates tasks for a meal plan.
func (s *service) MealPlanTaskCreationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	logger.Info("meal plan task creation worker invoked")

	if err := s.mealPlanTaskCreatorWorker.CreateMealPlanTasksForFinalizedMealPlans(ctx, nil); err != nil {
		observability.AcknowledgeError(err, logger, span, "finalizing expired meal plans")
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusInternalServerError)
		return
	}

	logger.Info("meal plan task creation worker completed")

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusAccepted)
}
