package householdinvitations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/prixfixeco/api_server/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/encoding"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	mockrandom "github.com/prixfixeco/api_server/internal/random/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func Test_service_InviteMemberHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserIDByEmail",
			testutils.ContextMatcher,
			exampleInput.ToEmail,
		).Return(helper.exampleUser.ID, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dataManager := database.NewMockDatabase()
		dataManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_created",
			helper.exampleUser.ID,
			testutils.MapOfStringToInterfaceMatcher,
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, dataManager, sg, dataChangesPublisher, cdc)
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
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

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
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

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
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

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
			"GetUserIDByEmail",
			testutils.ContextMatcher,
			exampleInput.ToEmail,
		).Return("", errors.New("blah"))
		helper.service.userDataManager = udm

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, sg)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserIDByEmail",
			testutils.ContextMatcher,
			exampleInput.ToEmail,
		).Return(helper.exampleUser.ID, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dataManager := database.NewMockDatabase()
		dataManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, dataManager, sg)
	})

	T.Run("with error publishing to data change feed", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserIDByEmail",
			testutils.ContextMatcher,
			exampleInput.ToEmail,
		).Return(helper.exampleUser.ID, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dataManager := database.NewMockDatabase()
		dataManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_created",
			helper.exampleUser.ID,
			testutils.MapOfStringToInterfaceMatcher,
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, dataManager, sg, dataChangesPublisher, cdc)
	})

	T.Run("with error collecting data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManager{}
		udm.On(
			"GetUserIDByEmail",
			testutils.ContextMatcher,
			exampleInput.ToEmail,
		).Return(helper.exampleUser.ID, nil)
		helper.service.userDataManager = udm

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dataManager := database.NewMockDatabase()
		dataManager.HouseholdInvitationDataManager.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(testutils.DataChangeMessageMatcher),
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_created",
			helper.exampleUser.ID,
			testutils.MapOfStringToInterfaceMatcher,
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.InviteMemberHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, udm, dataManager, sg, dataChangesPublisher, cdc)
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
			mock.IsType(&types.HouseholdInvitationList{}),
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
		).Return((*types.HouseholdInvitationList)(nil), errors.New("blah"))
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
			mock.IsType(&types.HouseholdInvitationList{}),
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
		).Return((*types.HouseholdInvitationList)(nil), errors.New("blah"))
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
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = hidm

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_accepted",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:           helper.exampleHousehold.ID,
				keys.HouseholdInvitationIDKey: helper.exampleHouseholdInvitation.ID,
			},
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm, cdc)
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

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader([]byte("{ bad json lol")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = hidm

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm)
	})

	T.Run("with error notifying customer data platform", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = hidm

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_accepted",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:           helper.exampleHousehold.ID,
				keys.HouseholdInvitationIDKey: helper.exampleHouseholdInvitation.ID,
			},
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.AcceptInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm, cdc)
	})
}

func Test_service_CancelInviteHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = hidm

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_cancelled",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:           helper.exampleHousehold.ID,
				keys.HouseholdInvitationIDKey: helper.exampleHouseholdInvitation.ID,
			},
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm, cdc)
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

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader([]byte("{ bad json lol")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = hidm

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm)
	})

	T.Run("with error notifying customer data platform", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = hidm

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_cancelled",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:           helper.exampleHousehold.ID,
				keys.HouseholdInvitationIDKey: helper.exampleHouseholdInvitation.ID,
			},
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.CancelInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm, cdc)
	})
}

func Test_service_RejectInviteHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = hidm

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_rejected",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:           helper.exampleHousehold.ID,
				keys.HouseholdInvitationIDKey: helper.exampleHouseholdInvitation.ID,
			},
		).Return(nil)
		helper.service.customerDataCollector = cdc

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm, cdc)
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
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

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
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = hidm

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm)
	})

	T.Run("with error notifying customer data platform", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationUpdateRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		hidm := &mocktypes.HouseholdInvitationDataManager{}
		hidm.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
			exampleInput.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = hidm

		cdc := &customerdata.MockCollector{}
		cdc.On(
			"EventOccurred",
			testutils.ContextMatcher,
			"household_invitation_rejected",
			helper.exampleUser.ID,
			map[string]interface{}{
				keys.HouseholdIDKey:           helper.exampleHousehold.ID,
				keys.HouseholdInvitationIDKey: helper.exampleHouseholdInvitation.ID,
			},
		).Return(errors.New("blah"))
		helper.service.customerDataCollector = cdc

		helper.service.RejectInviteHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusAccepted, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, hidm, cdc)
	})
}
