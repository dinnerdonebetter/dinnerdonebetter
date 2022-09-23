package advancedprepsteps

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestAdvancedPrepStepsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		advancedPrepStepDataManager := &mocktypes.AdvancedPrepStepDataManager{}
		advancedPrepStepDataManager.On(
			"GetAdvancedPrepStep",
			testutils.ContextMatcher,
			helper.exampleAdvancedPrepStep.ID,
		).Return(helper.exampleAdvancedPrepStep, nil)
		helper.service.advancedPrepStepDataManager = advancedPrepStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.AdvancedPrepStep{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, advancedPrepStepDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such meal plan in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		advancedPrepStepDataManager := &mocktypes.AdvancedPrepStepDataManager{}
		advancedPrepStepDataManager.On(
			"GetAdvancedPrepStep",
			testutils.ContextMatcher,
			helper.exampleAdvancedPrepStep.ID,
		).Return((*types.AdvancedPrepStep)(nil), sql.ErrNoRows)
		helper.service.advancedPrepStepDataManager = advancedPrepStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, advancedPrepStepDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		advancedPrepStepDataManager := &mocktypes.AdvancedPrepStepDataManager{}
		advancedPrepStepDataManager.On(
			"GetAdvancedPrepStep",
			testutils.ContextMatcher,
			helper.exampleAdvancedPrepStep.ID,
		).Return((*types.AdvancedPrepStep)(nil), errors.New("blah"))
		helper.service.advancedPrepStepDataManager = advancedPrepStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, advancedPrepStepDataManager, encoderDecoder)
	})
}

func TestAdvancedPrepStepsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleAdvancedPrepStepList := fakes.BuildFakeAdvancedPrepStepList().AdvancedPrepSteps

		advancedPrepStepDataManager := &mocktypes.AdvancedPrepStepDataManager{}
		advancedPrepStepDataManager.On(
			"GetAdvancedPrepStepsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return(exampleAdvancedPrepStepList, nil)
		helper.service.advancedPrepStepDataManager = advancedPrepStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.AdvancedPrepStep{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, advancedPrepStepDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		advancedPrepStepDataManager := &mocktypes.AdvancedPrepStepDataManager{}
		advancedPrepStepDataManager.On(
			"GetAdvancedPrepStepsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.AdvancedPrepStep(nil), sql.ErrNoRows)
		helper.service.advancedPrepStepDataManager = advancedPrepStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType([]*types.AdvancedPrepStep{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, advancedPrepStepDataManager, encoderDecoder)
	})

	T.Run("with error retrieving meal plans from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		advancedPrepStepDataManager := &mocktypes.AdvancedPrepStepDataManager{}
		advancedPrepStepDataManager.On(
			"GetAdvancedPrepStepsForMealPlan",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
		).Return([]*types.AdvancedPrepStep(nil), errors.New("blah"))
		helper.service.advancedPrepStepDataManager = advancedPrepStepDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListByMealPlanHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, advancedPrepStepDataManager, encoderDecoder)
	})
}

func TestAdvancedPrepStepsService_CompletionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.AdvancedPrepStepDataManager.On(
			"MarkAdvancedPrepStepAsComplete",
			testutils.ContextMatcher,
			helper.exampleAdvancedPrepStep.ID,
		).Return(nil)
		helper.service.advancedPrepStepDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CompletionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CompletionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.AdvancedPrepStepDataManager.On(
			"MarkAdvancedPrepStepAsComplete",
			testutils.ContextMatcher,
			helper.exampleAdvancedPrepStep.ID,
		).Return(errors.New("blah"))
		helper.service.advancedPrepStepDataManager = dbManager

		helper.service.CompletionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		dbManager := database.NewMockDatabase()
		dbManager.AdvancedPrepStepDataManager.On(
			"MarkAdvancedPrepStepAsComplete",
			testutils.ContextMatcher,
			helper.exampleAdvancedPrepStep.ID,
		).Return(nil)
		helper.service.advancedPrepStepDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CompletionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dbManager, dataChangesPublisher)
	})
}
