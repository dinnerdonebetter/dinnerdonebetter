package dataprivacy

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWorkerService_DataDeletionHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodDelete, "https://whatever.whocares.gov", http.NoBody)
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataPrivacyDataManager := &mocktypes.DataPrivacyDataManagerMock{}
		dataPrivacyDataManager.On(
			"DeleteUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(nil)
		helper.service.dataPrivacyDataManager = dataPrivacyDataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.DataDeletionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)
		var actual *types.APIResponse[*types.DataDeletionResponse]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dataPrivacyDataManager, dataChangesPublisher)
	})

	T.Run("with error", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodDelete, "https://whatever.whocares.gov", http.NoBody)
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataPrivacyDataManager := &mocktypes.DataPrivacyDataManagerMock{}
		dataPrivacyDataManager.On(
			"DeleteUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
		).Return(errors.New("blah"))
		helper.service.dataPrivacyDataManager = dataPrivacyDataManager

		helper.service.DataDeletionHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.DataDeletionResponse]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Error(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dataPrivacyDataManager)
	})
}

func TestWorkerService_UserDataAggregationRequestHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodDelete, "https://whatever.whocares.gov", http.NoBody)
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		userDataAggregationPublisher := &mockpublishers.Publisher{}
		userDataAggregationPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(request *types.UserDataAggregationRequest) bool {
				return true
			}),
		).Return(nil)
		helper.service.userDataAggregationPublisher = userDataAggregationPublisher

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.UserDataAggregationRequestHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)
		var actual *types.APIResponse[*types.UserDataCollectionResponse]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, userDataAggregationPublisher, dataChangesPublisher)
	})
}
