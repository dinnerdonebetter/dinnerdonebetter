package householdinvitations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/encoding"
	mockencoding "github.com/prixfixeco/backend/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	mockrandom "github.com/prixfixeco/backend/internal/random/mock"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	mocktypes "github.com/prixfixeco/backend/pkg/types/mock"
	testutils "github.com/prixfixeco/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_service_InviteMemberHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusCreated, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager, dataChangesPublisher)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error decoding request", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		ed := mockencoding.NewMockEncoderDecoder()
		ed.On(
			"DecodeRequest",
			testutils.ContextMatcher,
			testutils.HTTPRequestMatcher,
			mock.IsType(&types.HouseholdInvitationCreationRequestInput{}),
		).Return(errors.New("blah"))

		ed.On(
			"EncodeErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			"invalid request content",
			http.StatusBadRequest,
		).Return(errors.New("blah"))
		helper.service.encoderDecoder = ed

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader([]byte("")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdInvitationCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error generating invitation token", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return("", errors.New("blah"))
		helper.service.secretGenerator = sg

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, sg)
	})

	T.Run("with error fetching user ID", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return((*types.User)(nil), errors.New("blah"))
		helper.service.userDataManager = udm

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, sg)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = dbManager

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager)
	})

	T.Run("with error publishing message", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusCreated, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager, dataChangesPublisher)
	})

	T.Run("with error collecting data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusCreated, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing email request", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusCreated, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager, dataChangesPublisher)
	})
}

func Test_service_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetHouseholdInvitationByHouseholdAndID",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
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
			"GetHouseholdInvitationByHouseholdAndID",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
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
			"GetHouseholdInvitationByHouseholdAndID",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
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

func Test_service_InboundInvitesHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		exampleHouseholdInvitations := fakes.BuildFakeHouseholdInvitationList()

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetPendingHouseholdInvitationsForUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleHouseholdInvitations, nil)
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.HouseholdInvitation]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.InboundInvitesHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.InboundInvitesHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetPendingHouseholdInvitationsForUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.HouseholdInvitation])(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = wd

		helper.service.InboundInvitesHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd)
	})
}

func Test_service_OutboundInvitesHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		exampleHouseholdInvitations := fakes.BuildFakeHouseholdInvitationList()

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetPendingHouseholdInvitationsFromUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleHouseholdInvitations, nil)
		helper.service.householdInvitationDataManager = wd

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.HouseholdInvitation]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.OutboundInvitesHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd, encoderDecoder)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.OutboundInvitesHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManager{}
		wd.On(
			"GetPendingHouseholdInvitationsFromUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.HouseholdInvitation])(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = wd

		helper.service.OutboundInvitesHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, wd)
	})
}

func Test_service_AcceptInviteHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with unparseable input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader([]byte("{ bad json lol")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with no matching invitation", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = dataManager

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusNotFound, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error reading invitation", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error publishing service event", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})
}

func Test_service_CancelInviteHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with unparseable input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader([]byte("{ bad json lol")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = dataManager

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusNotFound, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error fetching invitation", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error publishing service event", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})
}

func Test_service_RejectInviteHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdInvitationUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = dataManager

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusNotFound, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error fetching invitation", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))

		helper.service.householdInvitationDataManager = dataManager

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader([]byte("{ bad json lol")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error publishing service event", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManager{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})
}
