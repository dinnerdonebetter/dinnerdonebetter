package meals

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/encoding"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestParseBool(t *testing.T) {
	t.Parallel()

	expectations := map[string]bool{
		"1":      true,
		t.Name(): false,
		"true":   true,
		"troo":   false,
		"t":      true,
		"false":  false,
	}

	for input, expected := range expectations {
		assert.Equal(t, expected, parseBool(input))
	}
}

func TestMealsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreWriteMessageMatcher),
		).Return(nil)
		helper.service.preWritesPublisher = mockEventProducer

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"meal_created",
			helper.exampleUser.ID,
			testutils.MapOfStringToInterfaceMatcher,
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer)
	})

	T.Run("without input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.MealCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error publishing event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreWriteMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.preWritesPublisher = mockEventProducer

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer)
	})

	T.Run("with error writing to customer data platform", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreWriteMessageMatcher),
		).Return(nil)
		helper.service.preWritesPublisher = mockEventProducer

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"meal_created",
			helper.exampleUser.ID,
			testutils.MapOfStringToInterfaceMatcher,
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer)
	})
}

func TestMealsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.mealIDFetcher = func(_ *http.Request) string {
			return helper.exampleMeal.ID
		}

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"GetMeal",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(helper.exampleMeal, nil)
		helper.service.mealDataManager = mealDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.Meal{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, encoderDecoder)
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

	T.Run("with no such meal in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.mealIDFetcher = func(_ *http.Request) string {
			return helper.exampleMeal.ID
		}

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"GetMeal",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return((*types.Meal)(nil), sql.ErrNoRows)
		helper.service.mealDataManager = mealDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.mealIDFetcher = func(_ *http.Request) string {
			return helper.exampleMeal.ID
		}

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"GetMeal",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return((*types.Meal)(nil), errors.New("blah"))
		helper.service.mealDataManager = mealDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, encoderDecoder)
	})
}

func TestMealsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealList := fakes.BuildFakeMealList()

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"GetMeals",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleMealList, nil)
		helper.service.mealDataManager = mealDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.MealList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, encoderDecoder)
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

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"GetMeals",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.MealList)(nil), sql.ErrNoRows)
		helper.service.mealDataManager = mealDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.MealList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, encoderDecoder)
	})

	T.Run("with error retrieving meals from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"GetMeals",
			testutils.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.MealList)(nil), errors.New("blah"))
		helper.service.mealDataManager = mealDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, encoderDecoder)
	})
}

func TestMealsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(true, nil)
		helper.service.mealDataManager = mealDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreArchiveMessageMatcher),
		).Return(nil)
		helper.service.preArchivesPublisher = mockEventProducer

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"meal_archived",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey: helper.exampleHousehold.ID,
				keys.MealIDKey:      helper.exampleMeal.ID,
			},
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, mockEventProducer, cdc)
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

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such meal in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(false, nil)
		helper.service.mealDataManager = mealDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, encoderDecoder)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(false, errors.New("blah"))
		helper.service.mealDataManager = mealDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(true, nil)
		helper.service.mealDataManager = mealDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreArchiveMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.preArchivesPublisher = mockEventProducer

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, mockEventProducer)
	})

	T.Run("with error writing to customer data platform", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealDataManager := &mocktypes.MealDataManager{}
		mealDataManager.On(
			"MealExists",
			testutils.ContextMatcher,
			helper.exampleMeal.ID,
		).Return(true, nil)
		helper.service.mealDataManager = mealDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreArchiveMessageMatcher),
		).Return(nil)
		helper.service.preArchivesPublisher = mockEventProducer

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"meal_archived",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey: helper.exampleHousehold.ID,
				keys.MealIDKey:      helper.exampleMeal.ID,
			},
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealDataManager, mockEventProducer, cdc)
	})
}
