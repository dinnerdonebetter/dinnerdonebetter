package mealplanning

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	servertiming "github.com/mitchellh/go-server-timing"
)

var _ types.MealDataService = (*service)(nil)

// CreateMealHandler is our meal creation route.
func (s *service) CreateMealHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// read parsed input struct from request body.
	providedInput := new(types.MealCreationRequestInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	if err = providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Info("provided meal creation input was invalid")
		errRes := types.NewAPIErrorResponse(err.Error(), types.ErrValidatingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	input := converters.ConvertMealCreationRequestInputToMealDatabaseCreationInput(providedInput)
	input.ID = identifiers.New()

	input.CreatedByUser = sessionCtxData.Requester.UserID
	tracing.AttachToSpan(span, keys.MealIDKey, input.ID)

	createTimer := timing.NewMetric("database").WithDesc("create").Start()
	meal, err := s.mealPlanningDataManager.CreateMeal(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating meal")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	createTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.MealCreatedServiceEventType,
		Meal:        meal,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.Meal]{
		Details: responseDetails,
		Data:    meal,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusCreated)
}

// ReadMealHandler returns a GET handler that returns a meal.
func (s *service) ReadMealHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// determine meal ID.
	mealID := s.mealIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)
	logger = logger.WithValue(keys.MealIDKey, mealID)

	// fetch meal from database.
	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	x, err := s.mealPlanningDataManager.GetMeal(ctx, mealID)
	if errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meal")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[*types.Meal]{
		Details: responseDetails,
		Data:    x,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ListMealsHandler is our list route.
func (s *service) ListMealsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	filter := types.ExtractQueryFilterFromRequest(req)
	logger := s.logger.WithRequest(req).WithSpan(span)
	logger = filter.AttachToLogger(logger)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	tracing.AttachRequestToSpan(span, req)
	tracing.AttachQueryFilterToSpan(span, filter)

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	meals, err := s.mealPlanningDataManager.GetMeals(ctx, filter)
	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		meals = &types.QueryFilteredResult[types.Meal]{Data: []*types.Meal{}}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving meals")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[[]*types.Meal]{
		Details:    responseDetails,
		Data:       meals.Data,
		Pagination: &meals.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// SearchMealsHandler is our list route.
func (s *service) SearchMealsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	query := req.URL.Query().Get(types.QueryKeySearch)
	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	filter := types.ExtractQueryFilterFromRequest(req)
	tracing.AttachQueryFilterToSpan(span, filter)
	logger = filter.AttachToLogger(logger)

	useDB := !s.useSearchService || strings.TrimSpace(strings.ToLower(req.URL.Query().Get(types.QueryKeySearchWithDatabase))) == "true"
	logger = logger.WithValue(keys.SearchQueryKey, query).WithValue("using_database", useDB)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	meals := &types.QueryFilteredResult[types.Meal]{
		Pagination: filter.ToPagination(),
	}

	readTimer := timing.NewMetric("database").WithDesc("fetch").Start()
	if useDB {
		meals, err = s.mealPlanningDataManager.SearchForMeals(ctx, query, filter)
	} else {
		var mealSubsets []*types.MealSearchSubset
		mealSubsets, err = s.searchIndex.Search(ctx, query)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "searching for meals")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}

		ids := []string{}
		for _, mealSubset := range mealSubsets {
			ids = append(ids, mealSubset.ID)
		}

		meals.Data, err = s.mealPlanningDataManager.GetMealsWithIDs(ctx, ids)
	}

	if errors.Is(err, sql.ErrNoRows) {
		// in the event no rows exist, return an empty list.
		meals = &types.QueryFilteredResult[types.Meal]{
			Pagination: filter.ToPagination(),
			Data:       []*types.Meal{},
		}
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "searching for meals")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	readTimer.Stop()

	responseValue := &types.APIResponse[[]*types.Meal]{
		Details:    responseDetails,
		Data:       meals.Data,
		Pagination: &meals.Pagination,
	}

	// encode our response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}

// ArchiveMealHandler returns a handler that archives a meal.
func (s *service) ArchiveMealHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	timing := servertiming.FromContext(ctx)
	logger := s.logger.WithRequest(req).WithSpan(span)
	tracing.AttachRequestToSpan(span, req)

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	// determine user ID.
	sessionContextTimer := timing.NewMetric("session").WithDesc("fetch session context").Start()
	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		errRes := types.NewAPIErrorResponse("unauthenticated", types.ErrFetchingSessionContextData, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusUnauthorized)
		return
	}
	sessionContextTimer.Stop()

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	responseDetails.CurrentHouseholdID = sessionCtxData.ActiveHouseholdID

	// determine meal ID.
	mealID := s.mealIDFetcher(req)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)
	logger = logger.WithValue(keys.MealIDKey, mealID)

	existenceTimer := timing.NewMetric("database").WithDesc("existence check").Start()
	exists, err := s.mealPlanningDataManager.MealExists(ctx, mealID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		observability.AcknowledgeError(err, logger, span, "checking meal existence")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	} else if !exists || errors.Is(err, sql.ErrNoRows) {
		errRes := types.NewAPIErrorResponse("not found", types.ErrDataNotFound, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusNotFound)
		return
	}
	existenceTimer.Stop()

	archiveTimer := timing.NewMetric("database").WithDesc("archive").Start()
	if err = s.mealPlanningDataManager.ArchiveMeal(ctx, mealID, sessionCtxData.Requester.UserID); err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving meal")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}
	archiveTimer.Stop()

	dcm := &types.DataChangeMessage{
		EventType:   types.MealArchivedServiceEventType,
		MealID:      mealID,
		HouseholdID: sessionCtxData.ActiveHouseholdID,
		UserID:      sessionCtxData.Requester.UserID,
	}

	go s.dataChangesPublisher.PublishAsync(ctx, dcm)

	responseValue := &types.APIResponse[*types.Meal]{
		Details: responseDetails,
	}

	// let everybody go home.
	s.encoderDecoder.RespondWithData(ctx, res, responseValue)
}
