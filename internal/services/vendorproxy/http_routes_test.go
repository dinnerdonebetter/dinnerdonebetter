package vendorproxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	analyticsmock "github.com/prixfixeco/backend/internal/analytics/mock"
	ffmock "github.com/prixfixeco/backend/internal/featureflags/mock"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
	testutils "github.com/prixfixeco/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidInstrumentsService_FeatureFlagHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleFlagName := fakes.BuildFakeID()
		helper.service.featureFlagURLFetcher = func(req *http.Request) string { return exampleFlagName }

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodGet, "https://local.prixfixe.dev", bytes.NewBuffer([]byte("")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		ffm := &ffmock.FeatureFlagManager{}
		ffm.On("CanUseFeature", testutils.ContextMatcher, helper.exampleUser.ID, exampleFlagName).Return(true, nil)
		helper.service.featureFlagManager = ffm

		helper.service.FeatureFlagStatusHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, ffm)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.FeatureFlagStatusHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error checking status", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleFlagName := fakes.BuildFakeID()
		helper.service.featureFlagURLFetcher = func(req *http.Request) string { return exampleFlagName }

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodGet, "https://local.prixfixe.dev", bytes.NewBuffer([]byte("")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		ffm := &ffmock.FeatureFlagManager{}
		ffm.On("CanUseFeature", testutils.ContextMatcher, helper.exampleUser.ID, exampleFlagName).Return(false, errors.New("blah"))
		helper.service.featureFlagManager = ffm

		helper.service.FeatureFlagStatusHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)

		mock.AssertExpectationsForObjects(t, ffm)
	})
}

func TestValidInstrumentsService_AnalyticsTrackHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		content := &analyticsTrackRequest{
			Event:      string(types.HouseholdCreatedCustomerEventType),
			Properties: map[string]any{"things": "stuff"},
		}

		contentBytes, err := json.Marshal(content)
		require.NoError(t, err)
		require.NotNil(t, contentBytes)

		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewBuffer(contentBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mr := &analyticsmock.EventReporter{}
		mr.On("EventOccurred", testutils.ContextMatcher, types.CustomerEventType(content.Event), helper.exampleUser.ID, content.Properties).Return(nil)
		helper.service.eventReporter = mr

		helper.service.AnalyticsTrackHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mr)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.AnalyticsTrackHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error reporting event", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		content := &analyticsTrackRequest{
			Event:      string(types.HouseholdCreatedCustomerEventType),
			Properties: map[string]any{"things": "stuff"},
		}

		contentBytes, err := json.Marshal(content)
		require.NoError(t, err)
		require.NotNil(t, contentBytes)

		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewBuffer(contentBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mr := &analyticsmock.EventReporter{}
		mr.On("EventOccurred", testutils.ContextMatcher, types.CustomerEventType(content.Event), helper.exampleUser.ID, content.Properties).Return(errors.New("blah"))
		helper.service.eventReporter = mr

		helper.service.AnalyticsTrackHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)
	})
}

func TestValidInstrumentsService_AnalyticsIdentifyHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		content := &analyticsIdentifyRequest{
			Properties: map[string]any{"things": "stuff"},
		}

		contentBytes, err := json.Marshal(content)
		require.NoError(t, err)
		require.NotNil(t, contentBytes)

		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewBuffer(contentBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mr := &analyticsmock.EventReporter{}
		mr.On("AddUser", testutils.ContextMatcher, helper.exampleUser.ID, content.Properties).Return(nil)
		helper.service.eventReporter = mr

		helper.service.AnalyticsIdentifyHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mr)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = func(req *http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.AnalyticsIdentifyHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
	})

	T.Run("with error making proxied request", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		content := &analyticsIdentifyRequest{
			Properties: map[string]any{"things": "stuff"},
		}

		contentBytes, err := json.Marshal(content)
		require.NoError(t, err)
		require.NotNil(t, contentBytes)

		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://local.prixfixe.dev", bytes.NewBuffer(contentBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		mr := &analyticsmock.EventReporter{}
		mr.On("AddUser", testutils.ContextMatcher, helper.exampleUser.ID, content.Properties).Return(errors.New("blah"))
		helper.service.eventReporter = mr

		helper.service.AnalyticsIdentifyHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code, "expected %d in status response, got %d", http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, mr)
	})
}
