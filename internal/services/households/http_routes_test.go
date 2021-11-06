package households

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/encoding"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	mockmetrics "github.com/prixfixeco/api_server/internal/observability/metrics/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestHouseholdsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHouseholds",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleHouseholdList, nil)
		helper.service.householdDataManager = householdDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.HouseholdList{}),
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with now rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHouseholds",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.HouseholdList)(nil), sql.ErrNoRows)
		helper.service.householdDataManager = householdDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.HouseholdList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, encoderDecoder)
	})

	T.Run("with error fetching households from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHouseholds",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.HouseholdList)(nil), errors.New("blah"))
		helper.service.householdDataManager = householdDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, encoderDecoder)
	})
}

func TestHouseholdsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"CreateHousehold",
			testutils.ContextMatcher,
			mock.IsType(&types.HouseholdCreationRequestInput{}),
		).Return(helper.exampleHousehold, nil)
		helper.service.householdDataManager = householdDataManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Increment", testutils.ContextMatcher).Return()
		helper.service.householdCounter = unitCounter

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, unitCounter)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("without input attached to request", func(t *testing.T) {
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

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.HouseholdCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error creating household in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"CreateHousehold",
			testutils.ContextMatcher,
			mock.IsType(&types.HouseholdCreationRequestInput{}),
		).Return((*types.Household)(nil), errors.New("blah"))
		helper.service.householdDataManager = householdDataManager

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager)
	})
}

func TestHouseholdsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleUser.ID,
		).Return(helper.exampleHousehold, nil)
		helper.service.householdDataManager = householdDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.Household{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such household in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleUser.ID,
		).Return((*types.Household)(nil), sql.ErrNoRows)
		helper.service.householdDataManager = householdDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, encoderDecoder)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID, helper.exampleUser.ID,
		).Return((*types.Household)(nil), errors.New("blah"))
		helper.service.householdDataManager = householdDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, encoderDecoder)
	})
}

func TestHouseholdsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleUser.ID,
		).Return(helper.exampleHousehold, nil)
		householdDataManager.On(
			"UpdateHousehold",
			testutils.ContextMatcher,
			mock.IsType(&types.Household{}),
		).Return(nil)
		helper.service.householdDataManager = householdDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("without update input attached to request", func(t *testing.T) {
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

	T.Run("with invalid input attached to request", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.HouseholdUpdateInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with no rows", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID, helper.exampleUser.ID,
		).Return((*types.Household)(nil), sql.ErrNoRows)
		helper.service.householdDataManager = householdDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager)
	})

	T.Run("with error querying for household", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID, helper.exampleUser.ID,
		).Return((*types.Household)(nil), errors.New("blah"))
		helper.service.householdDataManager = householdDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager)
	})

	T.Run("with error updating household", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdUpdateInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.exampleHousehold = fakes.BuildFakeHousehold()
		helper.exampleHousehold.BelongsToUser = helper.exampleUser.ID

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"GetHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleUser.ID,
		).Return(helper.exampleHousehold, nil)
		householdDataManager.On(
			"UpdateHousehold",
			testutils.ContextMatcher,
			mock.IsType(&types.Household{}),
		).Return(errors.New("blah"))
		helper.service.householdDataManager = householdDataManager

		helper.service.UpdateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager)
	})
}

func TestHouseholdsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"ArchiveHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.householdDataManager = householdDataManager

		unitCounter := &mockmetrics.UnitCounter{}
		unitCounter.On("Decrement", testutils.ContextMatcher).Return()
		helper.service.householdCounter = unitCounter

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, unitCounter)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such household in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"ArchiveHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleUser.ID,
		).Return(sql.ErrNoRows)
		helper.service.householdDataManager = householdDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, encoderDecoder)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdDataManager := &mocktypes.HouseholdDataManager{}
		householdDataManager.On(
			"ArchiveHousehold",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleUser.ID,
		).Return(errors.New("blah"))
		helper.service.householdDataManager = householdDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdDataManager, encoderDecoder)
	})
}

func TestHouseholdsService_AddMemberHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeAddUserToHouseholdInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

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

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeAddUserToHouseholdInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with missing input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := &types.AddUserToHouseholdInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error publishing to event queue", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeAddUserToHouseholdInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

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

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer)
	})
}

func TestHouseholdsService_ModifyMemberPermissionsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeUserPermissionModificationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdMembershipDataManager := &mocktypes.HouseholdUserMembershipDataManager{}
		householdMembershipDataManager.On(
			"ModifyUserPermissions",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			helper.exampleHousehold.ID,
			exampleInput,
		).Return(nil)
		helper.service.householdMembershipDataManager = householdMembershipDataManager

		helper.service.ModifyMemberPermissionsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdMembershipDataManager)
	})

	T.Run("with missing input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.ModifyMemberPermissionsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with invalid input attached", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := &types.ModifyUserPermissionsInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.ModifyMemberPermissionsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeUserPermissionModificationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.ModifyMemberPermissionsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeUserPermissionModificationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdMembershipDataManager := &mocktypes.HouseholdUserMembershipDataManager{}
		householdMembershipDataManager.On(
			"ModifyUserPermissions",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			helper.exampleHousehold.ID,
			exampleInput,
		).Return(errors.New("blah"))
		helper.service.householdMembershipDataManager = householdMembershipDataManager

		helper.service.ModifyMemberPermissionsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdMembershipDataManager)
	})
}

func TestHouseholdsService_TransferHouseholdOwnershipHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdMembershipDataManager := &mocktypes.HouseholdUserMembershipDataManager{}
		householdMembershipDataManager.On(
			"TransferHouseholdOwnership",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			exampleInput,
		).Return(nil)
		helper.service.householdMembershipDataManager = householdMembershipDataManager

		helper.service.TransferHouseholdOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdMembershipDataManager)
	})

	T.Run("without input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.TransferHouseholdOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdOwnershipTransferInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.TransferHouseholdOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.TransferHouseholdOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		householdMembershipDataManager := &mocktypes.HouseholdUserMembershipDataManager{}
		householdMembershipDataManager.On(
			"TransferHouseholdOwnership",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			exampleInput,
		).Return(errors.New("blah"))
		helper.service.householdMembershipDataManager = householdMembershipDataManager

		helper.service.TransferHouseholdOwnershipHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdMembershipDataManager)
	})
}

func TestHouseholdsService_RemoveMemberHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleReason := t.Name()

		householdMembershipDataManager := &mocktypes.HouseholdUserMembershipDataManager{}
		householdMembershipDataManager.On(
			"RemoveUserFromHousehold",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			helper.exampleHousehold.ID,
		).Return(nil)
		helper.service.householdMembershipDataManager = householdMembershipDataManager

		helper.req.URL.RawQuery = fmt.Sprintf("reason=%s", url.QueryEscape(exampleReason))

		helper.service.RemoveMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdMembershipDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.RemoveMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleReason := t.Name()
		helper.req.URL.RawQuery = fmt.Sprintf("reason=%s", url.QueryEscape(exampleReason))

		householdMembershipDataManager := &mocktypes.HouseholdUserMembershipDataManager{}
		householdMembershipDataManager.On(
			"RemoveUserFromHousehold",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			helper.exampleHousehold.ID,
		).Return(errors.New("blah"))
		helper.service.householdMembershipDataManager = householdMembershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.RemoveMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdMembershipDataManager, encoderDecoder)
	})
}

func TestHouseholdsService_MarkAsDefaultHouseholdHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdMembershipDataManager := &mocktypes.HouseholdUserMembershipDataManager{}
		householdMembershipDataManager.On(
			"MarkHouseholdAsUserDefault",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			helper.exampleHousehold.ID,
		).Return(nil)
		helper.service.householdMembershipDataManager = householdMembershipDataManager

		helper.service.MarkAsDefaultHouseholdHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdMembershipDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"unauthenticated",
			http.StatusUnauthorized,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.MarkAsDefaultHouseholdHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		householdMembershipDataManager := &mocktypes.HouseholdUserMembershipDataManager{}
		householdMembershipDataManager.On(
			"MarkHouseholdAsUserDefault",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			helper.exampleHousehold.ID,
		).Return(errors.New("blah"))
		helper.service.householdMembershipDataManager = householdMembershipDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.MarkAsDefaultHouseholdHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, householdMembershipDataManager, encoderDecoder)
	})
}
