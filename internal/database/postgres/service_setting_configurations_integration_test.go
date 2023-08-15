package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createServiceSettingConfigurationForTest(t *testing.T, ctx context.Context, exampleServiceSettingConfiguration *types.ServiceSettingConfiguration, dbc *Querier) *types.ServiceSettingConfiguration {
	t.Helper()

	// create
	if exampleServiceSettingConfiguration == nil {
		user := createUserForTest(t, ctx, nil, dbc)
		householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
		require.NoError(t, err)
		serviceSetting := createServiceSettingForTest(t, ctx, nil, dbc)
		exampleServiceSettingConfiguration = fakes.BuildFakeServiceSettingConfiguration()
		exampleServiceSettingConfiguration.ServiceSetting = *serviceSetting
		exampleServiceSettingConfiguration.BelongsToUser = user.ID
		exampleServiceSettingConfiguration.BelongsToHousehold = householdID
	}
	dbInput := converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationDatabaseCreationInput(exampleServiceSettingConfiguration)

	created, err := dbc.CreateServiceSettingConfiguration(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	exampleServiceSettingConfiguration.CreatedAt = created.CreatedAt
	exampleServiceSettingConfiguration.ServiceSetting = created.ServiceSetting
	assert.Equal(t, exampleServiceSettingConfiguration, created)

	serviceSettingConfiguration, err := dbc.GetServiceSettingConfiguration(ctx, created.ID)
	exampleServiceSettingConfiguration.CreatedAt = serviceSettingConfiguration.CreatedAt
	exampleServiceSettingConfiguration.ServiceSetting = serviceSettingConfiguration.ServiceSetting

	assert.NoError(t, err)
	assert.Equal(t, serviceSettingConfiguration, exampleServiceSettingConfiguration)

	return created
}

func TestQuerier_Integration_ServiceSettingConfigurations(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	serviceSetting := createServiceSettingForTest(t, ctx, nil, dbc)
	exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()
	exampleServiceSettingConfiguration.ServiceSetting = *serviceSetting
	exampleServiceSettingConfiguration.BelongsToUser = user.ID
	exampleServiceSettingConfiguration.BelongsToHousehold = householdID
	createdServiceSettingConfigurations := []*types.ServiceSettingConfiguration{}

	// create
	createdServiceSettingConfigurations = append(createdServiceSettingConfigurations, createServiceSettingConfigurationForTest(t, ctx, exampleServiceSettingConfiguration, dbc))

	// delete
	for _, serviceSettingConfiguration := range createdServiceSettingConfigurations {
		assert.NoError(t, dbc.ArchiveServiceSettingConfiguration(ctx, serviceSettingConfiguration.ID))

		var exists bool
		exists, err = dbc.ServiceSettingConfigurationExists(ctx, serviceSettingConfiguration.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ServiceSettingConfiguration
		y, err = dbc.GetServiceSettingConfiguration(ctx, serviceSettingConfiguration.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
