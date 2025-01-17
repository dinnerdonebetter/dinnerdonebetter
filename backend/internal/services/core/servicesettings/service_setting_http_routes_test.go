package servicesettings

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	testutils2 "github.com/dinnerdonebetter/backend/internal/testutils"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestServiceSettingsService_ReadServiceSettingHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"GetServiceSetting",
			testutils2.ContextMatcher,
			helper.exampleServiceSetting.ID,
		).Return(helper.exampleServiceSetting, nil)
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.ReadServiceSettingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, helper.exampleServiceSetting)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ReadServiceSettingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such service setting in the database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"GetServiceSetting",
			testutils2.ContextMatcher,
			helper.exampleServiceSetting.ID,
		).Return((*types.ServiceSetting)(nil), sql.ErrNoRows)
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.ReadServiceSettingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})

	T.Run("with error fetching from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"GetServiceSetting",
			testutils2.ContextMatcher,
			helper.exampleServiceSetting.ID,
		).Return((*types.ServiceSetting)(nil), errors.New("blah"))
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.ReadServiceSettingHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})
}

func TestServiceSettingsService_ListServiceSettingsHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		exampleServiceSettingList := fakes.BuildFakeServiceSettingsList()

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"GetServiceSettings",
			testutils2.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleServiceSettingList, nil)
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.ListServiceSettingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleServiceSettingList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.ListServiceSettingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"GetServiceSettings",
			testutils2.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ServiceSetting])(nil), sql.ErrNoRows)
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.ListServiceSettingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})

	T.Run("with error retrieving service settings from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"GetServiceSettings",
			testutils2.ContextMatcher,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.ServiceSetting])(nil), errors.New("blah"))
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.ListServiceSettingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})
}

func TestServiceSettingsService_SearchServiceSettingsHandler(T *testing.T) {
	T.Parallel()

	exampleQuery := "whatever"
	exampleLimit := uint8(123)
	exampleServiceSettingList := fakes.BuildFakeServiceSettingsList()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"SearchForServiceSettings",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return(exampleServiceSettingList.Data, nil)
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.SearchServiceSettingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Equal(t, actual.Data, exampleServiceSettingList.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils2.BrokenSessionContextDataFetcher

		helper.service.SearchServiceSettingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)

		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"SearchForServiceSettings",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return([]*types.ServiceSetting{}, sql.ErrNoRows)
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.SearchServiceSettingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})

	T.Run("with error retrieving from database", func(t *testing.T) {
		t.Parallel()

		helper := buildTestHelper(t)
		helper.req.URL.RawQuery = url.Values{
			types.QueryKeySearch: []string{exampleQuery},
			types.QueryKeyLimit:  []string{strconv.Itoa(int(exampleLimit))},
		}.Encode()

		serviceSettingDataManager := &mocktypes.ServiceSettingDataManagerMock{}
		serviceSettingDataManager.On(
			"SearchForServiceSettings",
			testutils2.ContextMatcher,
			exampleQuery,
		).Return([]*types.ServiceSetting{}, errors.New("blah"))
		helper.service.serviceSettingDataManager = serviceSettingDataManager

		helper.service.SearchServiceSettingsHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.ServiceSetting]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, serviceSettingDataManager)
	})
}
