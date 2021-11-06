package householdinvitations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/encoding"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestHouseholdInvitationsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInvitationDatabaseCreationInput()
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

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockEventProducer)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInvitationDatabaseCreationInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error decoding request", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"DecodeRequest",
			testutils.ContextMatcher,
			testutils.HTTPRequestMatcher,
			mock.IsType(&types.HouseholdInvitationCreationRequestInput{}),
		).Return(errors.New("blah"))

		encoderDecoder.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(""),
			http.StatusBadRequest,
		).Return(nil)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with invalid content attached to request", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := &types.HouseholdInvitationCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleCreationInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error publishing to pre-writes queue", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), encoding.ContentTypeJSON)

		exampleCreationInput := fakes.BuildFakeHouseholdInvitationDatabaseCreationInput()
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
}

func TestHouseholdInvitationsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		exampleHouseholdInvitationList := fakes.BuildFakeHouseholdInvitationList()

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetHouseholdInvitations",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleHouseholdInvitationList, nil)
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.HouseholdInvitationList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ListHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetHouseholdInvitations",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.HouseholdInvitationList)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.HouseholdInvitationList{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})

	T.Run("with error fetching household invitation from database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetHouseholdInvitations",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.HouseholdInvitationList)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})
}

func TestHouseholdInvitationsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHousehold.ID,
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.HouseholdInvitation{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with no such household invitation in database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInvitation)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})

	T.Run("with error fetching household invitation from database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHousehold.ID,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})
}

func TestHouseholdInvitationsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"HouseholdInvitationExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
		).Return(true, nil)
		helper.service.householdInvitationDataManager = wd

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreArchiveMessageMatcher),
		).Return(nil)
		helper.service.preArchivesPublisher = mockEventProducer

		helper.service.ArchiveHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusNoContent, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error checking household invitation existence", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"HouseholdInvitationExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
		).Return(false, errors.New("blah"))
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with no household invitation in database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"HouseholdInvitationExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
		).Return(false, sql.ErrNoRows)
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})

	T.Run("with error publishing to message queue", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"HouseholdInvitationExists",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
		).Return(true, nil)
		helper.service.householdInvitationDataManager = wd

		mockEventProducer := &mockpublishers.Publisher{}
		mockEventProducer.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.PreArchiveMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.preArchivesPublisher = mockEventProducer

		helper.service.ArchiveHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd)
	})
}
