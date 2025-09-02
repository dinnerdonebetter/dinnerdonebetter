package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	types "github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/domain/settings/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/settings/fakes"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/services/settings/grpc/converters"
	settingsconverters "github.com/dinnerdonebetter/backend/internal/services/settings/grpc/converters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createServiceSettingForTest(t *testing.T) *types.ServiceSetting {
	t.Helper()
	ctx := t.Context()

	exampleServiceSetting := fakes.BuildFakeServiceSetting()
	exampleServiceSettingInput := converters.ConvertServiceSettingToServiceSettingCreationRequestInput(exampleServiceSetting)
	createdServiceSetting, err := adminClient.CreateServiceSetting(ctx, &settingssvc.CreateServiceSettingRequest{
		Input: settingsconverters.ConvertServiceSettingCreationRequestInputToGRPCServiceSettingCreationRequestInput(exampleServiceSettingInput),
	})
	require.NoError(t, err)
	converted := settingsconverters.ConvertGRPCServiceSettingToServiceSetting(createdServiceSetting.Created)
	assertRoughEquality(t, exampleServiceSetting, converted, defaultIgnoredFields("ID")...)

	res, err := adminClient.GetServiceSetting(ctx, &settingssvc.GetServiceSettingRequest{ServiceSettingID: createdServiceSetting.Created.ID})
	require.NoError(t, err)
	require.NotNil(t, res)

	serviceSetting := settingsconverters.ConvertGRPCServiceSettingToServiceSetting(res.Result)
	assertRoughEquality(t, converted, serviceSetting, defaultIgnoredFields()...)

	return serviceSetting
}

func TestServiceSettings_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		createServiceSettingForTest(t)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeServiceSettingCreationRequestInput()
		convertedInput := grpcconverters.ConvertServiceSettingCreationRequestInputToGRPCServiceSettingCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateServiceSetting(ctx, &settingssvc.CreateServiceSettingRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeServiceSettingCreationRequestInput()
		convertedInput := grpcconverters.ConvertServiceSettingCreationRequestInputToGRPCServiceSettingCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.Name = ""

		created, err := adminClient.CreateServiceSetting(ctx, &settingssvc.CreateServiceSettingRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeServiceSettingCreationRequestInput()
		convertedInput := grpcconverters.ConvertServiceSettingCreationRequestInputToGRPCServiceSettingCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateServiceSetting(ctx, &settingssvc.CreateServiceSettingRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestServiceSettings_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createServiceSettingForTest(t)

		retrieved, err := testClient.GetServiceSetting(ctx, &settingssvc.GetServiceSettingRequest{ServiceSettingID: created.ID})
		assert.NoError(t, err)

		converted := grpcconverters.ConvertGRPCServiceSettingToServiceSetting(retrieved.Result)

		assertRoughEquality(t, created, converted, defaultIgnoredFields()...)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createServiceSettingForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetServiceSetting(ctx, &settingssvc.GetServiceSettingRequest{ServiceSettingID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.GetServiceSetting(ctx, &settingssvc.GetServiceSettingRequest{ServiceSettingID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestServiceSettings_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createServiceSettingForTest(t)

		_, err := adminClient.ArchiveServiceSetting(ctx, &settingssvc.ArchiveServiceSettingRequest{ServiceSettingID: created.ID})
		assert.NoError(t, err)

		x, err := adminClient.GetServiceSetting(ctx, &settingssvc.GetServiceSettingRequest{ServiceSettingID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createServiceSettingForTest(t)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.ArchiveServiceSetting(ctx, &settingssvc.ArchiveServiceSettingRequest{ServiceSettingID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := adminClient.ArchiveServiceSetting(ctx, &settingssvc.ArchiveServiceSettingRequest{ServiceSettingID: nonexistentID})
		assert.Error(t, err)
	})

	T.Run("non-admin users are forbidden from archiving", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createServiceSettingForTest(t)
		_, testClient := createUserAndClientForTest(T)

		_, err := testClient.ArchiveServiceSetting(ctx, &settingssvc.ArchiveServiceSettingRequest{ServiceSettingID: created.ID})
		assert.Error(t, err)
	})
}

func TestServiceSettings_Listing(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)
	createdServiceSettings := []*settings.ServiceSetting{}
	for range exampleQuantity {
		created := createServiceSettingForTest(T)
		createdServiceSettings = append(createdServiceSettings, created)
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		retrieved, err := testClient.GetServiceSettings(ctx, &settingssvc.GetServiceSettingsRequest{})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.True(t, len(retrieved.Results) >= len(createdServiceSettings))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetServiceSettings(ctx, &settingssvc.GetServiceSettingsRequest{})
		assert.Error(t, err)
	})
}

func TestServiceSettings_Searching(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)
		created := createServiceSettingForTest(t)

		retrieved, err := testClient.SearchForServiceSettings(ctx, &settingssvc.SearchForServiceSettingsRequest{
			Query: created.Name[:2],
		})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.SearchForServiceSettings(ctx, &settingssvc.SearchForServiceSettingsRequest{})
		assert.Error(t, err)
	})
}
