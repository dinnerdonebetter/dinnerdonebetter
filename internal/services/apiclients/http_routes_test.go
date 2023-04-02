package apiclients

import (
	"bytes"
	"database/sql"
	"errors"
	"math"
	"net/http"
	"testing"

	mockauthn "github.com/prixfixeco/backend/internal/authentication/mock"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/encoding"
	mockencoding "github.com/prixfixeco/backend/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	mockrandom "github.com/prixfixeco/backend/internal/random/mock"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	mocktypes "github.com/prixfixeco/backend/pkg/types/mock"
	testutils "github.com/prixfixeco/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAPIClientsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleAPIClientList := fakes.BuildFakeAPIClientList()

		mockDB := database.NewMockDatabase()
		mockDB.APIClientDataManager.On(
			"GetAPIClients",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleAPIClientList, nil)
		helper.service.apiClientDataManager = mockDB

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.APIClient]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, encoderDecoder)
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

	T.Run("with no results returned from datastore", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mockDB := database.NewMockDatabase()
		mockDB.APIClientDataManager.On(
			"GetAPIClients",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.APIClient])(nil), sql.ErrNoRows)
		helper.service.apiClientDataManager = mockDB
		helper.service.userDataManager = mockDB

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.QueryFilteredResult[types.APIClient]{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, encoderDecoder)
	})

	T.Run("with error retrieving clients from the datastore", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		mockDB := database.NewMockDatabase()
		mockDB.APIClientDataManager.On(
			"GetAPIClients",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.APIClient])(nil), errors.New("blah"))
		helper.service.apiClientDataManager = mockDB
		helper.service.userDataManager = mockDB

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ListHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, encoderDecoder)
	})
}

var apiClientCreationInputMatcher any = mock.MatchedBy(func(input *types.APIClientDatabaseCreationInput) bool {
	return true
})

func TestAPIClientsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockDB := database.NewMockDatabase()
		mockDB.UserDataManager.On(
			"GetUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = mockDB

		a := &mockauthn.Authenticator{}
		a.On(
			"ValidateLogin",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleInput.TOTPToken,
		).Return(true, nil)
		helper.service.authenticator = a

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			clientIDSize,
		).Return(helper.exampleAPIClient.ClientID, nil)
		sg.On(
			"GenerateRawBytes",
			testutils.ContextMatcher,
			clientSecretSize,
		).Return(helper.exampleAPIClient.ClientSecret, nil)
		helper.service.secretGenerator = sg

		mockDB.APIClientDataManager.On(
			"CreateAPIClient",
			testutils.ContextMatcher,
			apiClientCreationInputMatcher,
		).Return(helper.exampleAPIClient, nil)
		helper.service.apiClientDataManager = mockDB

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, sg, dataChangesPublisher)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})

	T.Run("without input attached to request", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with invalid input attached to request", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.service.cfg.minimumPasswordLength = math.MaxUint8
		helper.exampleInput.Password = ""
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
	})

	T.Run("with error retrieving user", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockDB := database.NewMockDatabase()
		mockDB.UserDataManager.On(
			"GetUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return((*types.User)(nil), errors.New("blah"))
		helper.service.apiClientDataManager = mockDB
		helper.service.userDataManager = mockDB

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with invalid credentials", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockDB := database.NewMockDatabase()

		mockDB.UserDataManager.On(
			"GetUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(helper.exampleUser, nil)

		mockDB.APIClientDataManager.On(
			"CreateAPIClient",
			testutils.ContextMatcher,
			apiClientCreationInputMatcher,
		).Return(helper.exampleAPIClient, nil)
		helper.service.apiClientDataManager = mockDB
		helper.service.userDataManager = mockDB

		a := &mockauthn.Authenticator{}
		a.On(
			"ValidateLogin",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleInput.TOTPToken,
		).Return(false, nil)
		helper.service.authenticator = a

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a)
	})

	T.Run("with invalid password", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockDB := database.NewMockDatabase()
		mockDB.UserDataManager.On(
			"GetUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(helper.exampleUser, nil)

		mockDB.APIClientDataManager.On(
			"CreateAPIClient",
			testutils.ContextMatcher,
			apiClientCreationInputMatcher,
		).Return(helper.exampleAPIClient, nil)
		helper.service.apiClientDataManager = mockDB
		helper.service.userDataManager = mockDB

		a := &mockauthn.Authenticator{}
		a.On(
			"ValidateLogin",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleInput.TOTPToken,
		).Return(true, errors.New("blah"))
		helper.service.authenticator = a

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a)
	})

	T.Run("with error generating client ID", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockDB := database.NewMockDatabase()
		mockDB.UserDataManager.On(
			"GetUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(helper.exampleUser, nil)

		a := &mockauthn.Authenticator{}
		a.On(
			"ValidateLogin",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleInput.TOTPToken,
		).Return(true, nil)
		helper.service.authenticator = a

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			clientIDSize,
		).Return("", errors.New("blah"))
		helper.service.secretGenerator = sg

		mockDB.APIClientDataManager.On(
			"CreateAPIClient",
			testutils.ContextMatcher,
			apiClientCreationInputMatcher,
		).Return(helper.exampleAPIClient, nil)

		helper.service.apiClientDataManager = mockDB
		helper.service.userDataManager = mockDB

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, sg)
	})

	T.Run("with error generating client secret", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockDB := database.NewMockDatabase()
		mockDB.UserDataManager.On(
			"GetUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = mockDB

		a := &mockauthn.Authenticator{}
		a.On(
			"ValidateLogin",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleInput.TOTPToken,
		).Return(true, nil)
		helper.service.authenticator = a

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			clientIDSize,
		).Return(helper.exampleAPIClient.ClientID, nil)
		sg.On(
			"GenerateRawBytes",
			testutils.ContextMatcher,
			clientSecretSize,
		).Return([]byte(nil), errors.New("blah"))
		helper.service.secretGenerator = sg

		mockDB.APIClientDataManager.On(
			"CreateAPIClient",
			testutils.ContextMatcher,
			apiClientCreationInputMatcher,
		).Return(helper.exampleAPIClient, nil)
		helper.service.apiClientDataManager = mockDB

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, sg)
	})

	T.Run("with error creating API Client in data store", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockDB := database.NewMockDatabase()
		mockDB.UserDataManager.On(
			"GetUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(helper.exampleUser, nil)

		a := &mockauthn.Authenticator{}
		a.On(
			"ValidateLogin",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleInput.TOTPToken,
		).Return(true, nil)
		helper.service.authenticator = a

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			clientIDSize,
		).Return(helper.exampleAPIClient.ClientID, nil)
		sg.On(
			"GenerateRawBytes",
			testutils.ContextMatcher,
			clientSecretSize,
		).Return(helper.exampleAPIClient.ClientSecret, nil)
		helper.service.secretGenerator = sg

		mockDB.APIClientDataManager.On(
			"CreateAPIClient",
			testutils.ContextMatcher,
			apiClientCreationInputMatcher,
		).Return((*types.APIClient)(nil), errors.New("blah"))

		helper.service.apiClientDataManager = mockDB
		helper.service.userDataManager = mockDB

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, sg)
	})

	T.Run("with error publishing service event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockDB := database.NewMockDatabase()
		mockDB.UserDataManager.On(
			"GetUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = mockDB

		a := &mockauthn.Authenticator{}
		a.On(
			"ValidateLogin",
			testutils.ContextMatcher,
			helper.exampleUser.HashedPassword,
			helper.exampleInput.Password,
			helper.exampleUser.TwoFactorSecret,
			helper.exampleInput.TOTPToken,
		).Return(true, nil)
		helper.service.authenticator = a

		sg := &mockrandom.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			clientIDSize,
		).Return(helper.exampleAPIClient.ClientID, nil)
		sg.On(
			"GenerateRawBytes",
			testutils.ContextMatcher,
			clientSecretSize,
		).Return(helper.exampleAPIClient.ClientSecret, nil)
		helper.service.secretGenerator = sg

		mockDB.APIClientDataManager.On(
			"CreateAPIClient",
			testutils.ContextMatcher,
			apiClientCreationInputMatcher,
		).Return(helper.exampleAPIClient, nil)
		helper.service.apiClientDataManager = mockDB

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CreateHandler(helper.res, helper.req)
		assert.Equal(t, http.StatusCreated, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockDB, a, sg, dataChangesPublisher)
	})
}

func TestAPIClientsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		apiClientDataManager := &mocktypes.APIClientDataManager{}
		apiClientDataManager.On(
			"GetAPIClientByDatabaseID",
			testutils.ContextMatcher,
			helper.exampleAPIClient.ID,
			helper.exampleUser.ID,
		).Return(helper.exampleAPIClient, nil)
		helper.service.apiClientDataManager = apiClientDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"RespondWithData",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
			mock.IsType(&types.APIClient{}),
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, apiClientDataManager, encoderDecoder)
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

	T.Run("with no such API client in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		apiClientDataManager := &mocktypes.APIClientDataManager{}
		apiClientDataManager.On(
			"GetAPIClientByDatabaseID",
			testutils.ContextMatcher,
			helper.exampleAPIClient.ID,
			helper.exampleUser.ID,
		).Return((*types.APIClient)(nil), sql.ErrNoRows)
		helper.service.apiClientDataManager = apiClientDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, apiClientDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		apiClientDataManager := &mocktypes.APIClientDataManager{}
		apiClientDataManager.On(
			"GetAPIClientByDatabaseID",
			testutils.ContextMatcher,
			helper.exampleAPIClient.ID,
			helper.exampleUser.ID,
		).Return((*types.APIClient)(nil), errors.New("blah"))
		helper.service.apiClientDataManager = apiClientDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, apiClientDataManager, encoderDecoder)
	})
}

func TestAPIClientsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		apiClientDataManager := &mocktypes.APIClientDataManager{}
		apiClientDataManager.On(
			"ArchiveAPIClient",
			testutils.ContextMatcher,
			helper.exampleAPIClient.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.apiClientDataManager = apiClientDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, apiClientDataManager, dataChangesPublisher)
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

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, encoderDecoder)
	})

	T.Run("with no such API client in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		apiClientDataManager := &mocktypes.APIClientDataManager{}
		apiClientDataManager.On(
			"ArchiveAPIClient",
			testutils.ContextMatcher,
			helper.exampleAPIClient.ID,
			helper.exampleUser.ID,
		).Return(sql.ErrNoRows)
		helper.service.apiClientDataManager = apiClientDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeNotFoundResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		).Return()
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, apiClientDataManager, encoderDecoder)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		apiClientDataManager := &mocktypes.APIClientDataManager{}
		apiClientDataManager.On(
			"ArchiveAPIClient",
			testutils.ContextMatcher,
			helper.exampleAPIClient.ID,
			helper.exampleUser.ID,
		).Return(errors.New("blah"))
		helper.service.apiClientDataManager = apiClientDataManager

		encoderDecoder := mockencoding.NewMockEncoderDecoder()
		encoderDecoder.On(
			"EncodeUnspecifiedInternalServerErrorResponse",
			testutils.ContextMatcher,
			testutils.HTTPResponseWriterMatcher,
		)
		helper.service.encoderDecoder = encoderDecoder

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, apiClientDataManager, encoderDecoder)
	})

	T.Run("with error publishing service event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		apiClientDataManager := &mocktypes.APIClientDataManager{}
		apiClientDataManager.On(
			"ArchiveAPIClient",
			testutils.ContextMatcher,
			helper.exampleAPIClient.ID,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.apiClientDataManager = apiClientDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.ArchiveHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNoContent, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, apiClientDataManager, dataChangesPublisher)
	})
}
