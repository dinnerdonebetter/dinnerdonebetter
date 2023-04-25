package integration

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/apiclient"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkServiceSettingConfigurationEquality(t *testing.T, expected, actual *types.ServiceSettingConfiguration) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.Value, actual.Value, "expected Value for service setting %s to be %v, but it was %v", expected.ID, expected.Value, actual.Value)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for service setting %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.NotEmpty(t, actual.ServiceSetting.ID, "expected ServiceSetting for service setting %s to be set, but it was not", expected.ID)
	assert.NotEmpty(t, actual.BelongsToHousehold, "expected BelongsToHousehold for service setting %s to be set, but it was not", expected.ID)
	assert.NotEmpty(t, actual.BelongsToUser, "expected BelongsToUser for service setting %s to be set, but it was not", expected.ID)

	assert.NotZero(t, actual.CreatedAt)
}

func buildUserServiceSettingConfigurationForTest(t *testing.T, adminClient, userClient *apiclient.Client, settingType string, ctx context.Context) *types.ServiceSettingConfiguration {
	t.Helper()

	createdServiceSetting := buildServiceSettingForTest(t, adminClient, settingType, ctx)

	t.Log("creating service setting")
	exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()
	exampleServiceSettingConfiguration.ServiceSetting = types.ServiceSetting{ID: createdServiceSetting.ID}
	exampleServiceSettingConfigurationInput := converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationCreationRequestInput(exampleServiceSettingConfiguration)
	createdServiceSettingConfiguration, err := userClient.CreateServiceSettingConfiguration(ctx, exampleServiceSettingConfigurationInput)
	require.NoError(t, err)
	t.Logf("service setting %q created", createdServiceSettingConfiguration.ID)
	checkServiceSettingConfigurationEquality(t, exampleServiceSettingConfiguration, createdServiceSettingConfiguration)

	return createdServiceSettingConfiguration
}

func (s *TestSuite) TestServiceSettingConfigurations_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdServiceSettingConfiguration := buildUserServiceSettingConfigurationForTest(t, testClients.admin, testClients.user, "user", s.ctx)

			t.Log("changing service setting configuration")
			newServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()
			newServiceSettingConfiguration.ServiceSetting = createdServiceSettingConfiguration.ServiceSetting
			newServiceSettingConfiguration.BelongsToUser = createdServiceSettingConfiguration.BelongsToUser
			newServiceSettingConfiguration.BelongsToHousehold = createdServiceSettingConfiguration.BelongsToHousehold
			createdServiceSettingConfiguration.Update(converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationUpdateRequestInput(newServiceSettingConfiguration))
			assert.NoError(t, testClients.user.UpdateServiceSettingConfiguration(ctx, createdServiceSettingConfiguration))

			t.Log("fetching changed service setting configuration")
			actual, err := testClients.user.GetServiceSettingConfigurationForUserByName(ctx, createdServiceSettingConfiguration.ServiceSetting.Name, nil)
			requireNotNilAndNoProblems(t, actual, err)

			// assert service setting equality
			checkServiceSettingConfigurationEquality(t, newServiceSettingConfiguration, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up service setting configuration")
			assert.NoError(t, testClients.user.ArchiveServiceSettingConfiguration(ctx, createdServiceSettingConfiguration.ID))

			t.Log("cleaning up service setting")
			assert.NoError(t, testClients.admin.ArchiveServiceSetting(ctx, createdServiceSettingConfiguration.ServiceSetting.ID))
		}
	})
}

func (s *TestSuite) TestServiceSettingConfigurations_ListingForUser() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating service settings")
			var expected []*types.ServiceSettingConfiguration
			for i := 0; i < 5; i++ {
				createdServiceSettingConfiguration := buildUserServiceSettingConfigurationForTest(t, testClients.admin, testClients.user, "user", s.ctx)
				expected = append(expected, createdServiceSettingConfiguration)
			}

			// assert service setting list equality
			actual, err := testClients.user.GetServiceSettingConfigurationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			t.Log("cleaning up")
			for _, createdServiceSettingConfiguration := range expected {
				assert.NoError(t, testClients.user.ArchiveServiceSettingConfiguration(ctx, createdServiceSettingConfiguration.ID))
				assert.NoError(t, testClients.admin.ArchiveServiceSetting(ctx, createdServiceSettingConfiguration.ServiceSetting.ID))
			}
		}
	})
}

func (s *TestSuite) TestServiceSettingConfigurations_ListingForHousehold() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating service settings")
			var expected []*types.ServiceSettingConfiguration
			for i := 0; i < 5; i++ {
				createdServiceSettingConfiguration := buildUserServiceSettingConfigurationForTest(t, testClients.admin, testClients.user, "household", s.ctx)
				expected = append(expected, createdServiceSettingConfiguration)
			}

			// assert service setting list equality
			actual, err := testClients.user.GetServiceSettingConfigurationsForHousehold(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			t.Log("cleaning up")
			for _, createdServiceSettingConfiguration := range expected {
				assert.NoError(t, testClients.user.ArchiveServiceSettingConfiguration(ctx, createdServiceSettingConfiguration.ID))
				assert.NoError(t, testClients.admin.ArchiveServiceSetting(ctx, createdServiceSettingConfiguration.ServiceSetting.ID))
			}
		}
	})
}
