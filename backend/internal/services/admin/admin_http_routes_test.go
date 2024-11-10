package admin

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/identifiers"
	"github.com/dinnerdonebetter/backend/internal/pkg/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAdminService_UserAccountStatusChangeHandler(T *testing.T) {
	T.Parallel()

	T.Run("banning users", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.BannedUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		userDataManager := &mocktypes.AdminUserDataManagerMock{}
		userDataManager.On(
			"UpdateUserAccountStatus",
			testutils.ContextMatcher,
			helper.exampleAccountStatusUpdateInput.TargetUserID,
			helper.exampleAccountStatusUpdateInput,
		).Return(nil)
		helper.service.userDB = userDataManager

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, userDataManager)
	})

	T.Run("without adequate permission", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		sessionCtxData := &types.SessionContextData{
			Requester: types.RequesterInfo{
				ServicePermissions: authorization.NewServiceRolePermissionChecker(),
			},
			HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{},
		}
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return sessionCtxData, nil
		}

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.BannedUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusForbidden, helper.res.Code)
	})

	T.Run("back in good standing", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.GoodStandingUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		userDataManager := &mocktypes.AdminUserDataManagerMock{}
		userDataManager.On(
			"UpdateUserAccountStatus",
			testutils.ContextMatcher,
			helper.exampleAccountStatusUpdateInput.TargetUserID,
			helper.exampleAccountStatusUpdateInput,
		).Return(nil)
		helper.service.userDB = userDataManager

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, userDataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.BannedUserAccountStatus)

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.Webhook]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no input attached to request", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(nil))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[any]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input attached to request", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput = &types.UserAccountStatusUpdateInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.Webhook]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with inadequate admin user attempting to ban", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.BannedUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			scd := &types.SessionContextData{
				Requester: types.RequesterInfo{
					ServicePermissions: authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
				},
			}

			return scd, nil
		}

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusForbidden, helper.res.Code)
	})

	T.Run("with inadequate admin user attempting to terminate households", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.TerminatedUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			scd := &types.SessionContextData{
				Requester: types.RequesterInfo{
					ServicePermissions: authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
				},
			}

			return scd, nil
		}

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusForbidden, helper.res.Code)
	})

	T.Run("with admin that has inadequate permissions", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.neuterAdminUser()

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.BannedUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusForbidden, helper.res.Code)
	})

	T.Run("with no such user in database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.BannedUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		userDataManager := &mocktypes.AdminUserDataManagerMock{}
		userDataManager.On(
			"UpdateUserAccountStatus",
			testutils.ContextMatcher,
			helper.exampleAccountStatusUpdateInput.TargetUserID,
			helper.exampleAccountStatusUpdateInput,
		).Return(sql.ErrNoRows)
		helper.service.userDB = userDataManager

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)

		mock.AssertExpectationsForObjects(t, userDataManager)
	})

	T.Run("with error writing new account status to database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.BannedUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		userDataManager := &mocktypes.AdminUserDataManagerMock{}
		userDataManager.On(
			"UpdateUserAccountStatus",
			testutils.ContextMatcher,
			helper.exampleAccountStatusUpdateInput.TargetUserID,
			helper.exampleAccountStatusUpdateInput,
		).Return(errors.New("blah"))
		helper.service.userDB = userDataManager

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.Webhook]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, userDataManager)
	})

	T.Run("with error destroying session", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleAccountStatusUpdateInput.NewStatus = string(types.BannedUserAccountStatus)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleAccountStatusUpdateInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		userDataManager := &mocktypes.AdminUserDataManagerMock{}
		userDataManager.On(
			"UpdateUserAccountStatus",
			testutils.ContextMatcher,
			helper.exampleAccountStatusUpdateInput.TargetUserID,
			helper.exampleAccountStatusUpdateInput,
		).Return(nil)
		helper.service.userDB = userDataManager

		helper.service.UserAccountStatusChangeHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)

		mock.AssertExpectationsForObjects(t, userDataManager)
	})
}

func TestAdminService_WriteArbitraryQueueMessageHandler(T *testing.T) {
	T.Parallel()

	T.Run("data change message", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		helper.exampleArbitraryQueueMessageRequestInput.QueueName = "data_changes"
		dataChangesRequest := &types.DataChangeMessage{
			RequestID: identifiers.New(),
			Context: map[string]any{
				"testing": true,
			},
		}
		helper.exampleArbitraryQueueMessageRequestInput.Body = string(helper.service.encoderDecoder.MustEncodeJSON(helper.ctx, dataChangesRequest))

		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, helper.exampleArbitraryQueueMessageRequestInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mockPublisher := &mockpublishers.Publisher{}
		mockPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(x []byte) bool {
				var y types.DataChangeMessage
				err = json.Unmarshal(x, &y)
				return err == nil
			}),
		).Return(nil)

		mockPublisherProvider := &mockpublishers.ProducerProvider{}
		mockPublisherProvider.On(
			"ProvidePublisher",
			helper.exampleArbitraryQueueMessageRequestInput.QueueName,
		).Return(mockPublisher, nil)
		helper.service.publisherProvider = mockPublisherProvider

		helper.service.WriteArbitraryQueueMessageHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mockPublisherProvider, mockPublisher)
	})
}
