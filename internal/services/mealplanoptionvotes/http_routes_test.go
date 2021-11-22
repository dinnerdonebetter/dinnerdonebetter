package mealplanoptionvotes

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

func TestMealPlanOptionVotesService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput()
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
			"meal_plan_option_vote_created",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:      helper.exampleHousehold.ID,
				keys.MealPlanIDKey:       helper.exampleMealPlan.ID,
				keys.MealPlanOptionIDKey: helper.exampleMealPlanOption.ID,
			},
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer, cdc)
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

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput()
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

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput()
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

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteDatabaseCreationInput()
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
			"meal_plan_option_vote_created",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:      helper.exampleHousehold.ID,
				keys.MealPlanIDKey:       helper.exampleMealPlan.ID,
				keys.MealPlanOptionIDKey: helper.exampleMealPlanOption.ID,
			},
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer, cdc)
	})
}

func TestMealPlanOptionVotesService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(helper.exampleMealPlanOptionVote, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.MealPlanOptionVote{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, encoderDecoder)
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

	T.Run("with no such meal plan option vote in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return((*types.MealPlanOptionVote)(nil), sql.ErrNoRows)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return((*types.MealPlanOptionVote)(nil), errors.New("blah"))
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, encoderDecoder)
	})
}

func TestMealPlanOptionVotesService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleMealPlanOptionVoteList := fakes.BuildFakeMealPlanOptionVoteList()

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVotes",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleMealPlanOptionVoteList, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.MealPlanOptionVoteList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, encoderDecoder)
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

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVotes",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.MealPlanOptionVoteList)(nil), sql.ErrNoRows)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.MealPlanOptionVoteList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, encoderDecoder)
	})

	T.Run("with error retrieving meal plan option votes from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVotes",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.MealPlanOptionVoteList)(nil), errors.New("blah"))
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, encoderDecoder)
	})
}

func TestMealPlanOptionVotesService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(helper.exampleMealPlanOptionVote, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreUpdateMessageMatcher),
		).Return(nil)
		helper.service.preUpdatesPublisher = mockEventProducer

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"meal_plan_option_vote_updated",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:      helper.exampleHousehold.ID,
				keys.MealPlanIDKey:       helper.exampleMealPlan.ID,
				keys.MealPlanOptionIDKey: helper.exampleMealPlanOption.ID,
			},
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, mockEventProducer, cdc)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.MealPlanOptionVoteUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("without input attached to context", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with no such meal plan option vote", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return((*types.MealPlanOptionVote)(nil), sql.ErrNoRows)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager)
	})

	T.Run("with error retrieving meal plan option vote from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return((*types.MealPlanOptionVote)(nil), errors.New("blah"))
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(helper.exampleMealPlanOptionVote, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreUpdateMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.preUpdatesPublisher = mockEventProducer

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, mockEventProducer)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeMealPlanOptionVoteUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"GetMealPlanOptionVote",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(helper.exampleMealPlanOptionVote, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreUpdateMessageMatcher),
		).Return(nil)
		helper.service.preUpdatesPublisher = mockEventProducer

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"meal_plan_option_vote_updated",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:      helper.exampleHousehold.ID,
				keys.MealPlanIDKey:       helper.exampleMealPlan.ID,
				keys.MealPlanOptionIDKey: helper.exampleMealPlanOption.ID,
			},
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, mockEventProducer, cdc)
	})
}

func TestMealPlanOptionVotesService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(true, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

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
			"meal_plan_option_vote_archived",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:      helper.exampleHousehold.ID,
				keys.MealPlanIDKey:       helper.exampleMealPlan.ID,
				keys.MealPlanOptionIDKey: helper.exampleMealPlanOption.ID,
			},
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, mockEventProducer, cdc)
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

	T.Run("with no such meal plan option vote in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(false, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, encoderDecoder)
	})

	T.Run("with error checking for item in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(false, errors.New("blah"))
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(true, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreArchiveMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.preArchivesPublisher = mockEventProducer

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, mockEventProducer)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mealPlanOptionVoteDataManager := &mocktypes.MealPlanOptionVoteDataManager{}
		mealPlanOptionVoteDataManager.On(
			"MealPlanOptionVoteExists",
			testutils.ContextMatcher,
			helper.exampleMealPlan.ID,
			helper.exampleMealPlanOption.ID,
			helper.exampleMealPlanOptionVote.ID,
		).Return(true, nil)
		helper.service.mealPlanOptionVoteDataManager = mealPlanOptionVoteDataManager

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
			"meal_plan_option_vote_archived",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:      helper.exampleHousehold.ID,
				keys.MealPlanIDKey:       helper.exampleMealPlan.ID,
				keys.MealPlanOptionIDKey: helper.exampleMealPlanOption.ID,
			},
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mealPlanOptionVoteDataManager, mockEventProducer, cdc)
	})
}
