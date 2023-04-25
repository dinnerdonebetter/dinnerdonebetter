package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/apiclient"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkServiceSettingEquality(t *testing.T, expected, actual *types.ServiceSetting) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.Description, actual.Description, "expected Description for service setting to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for service setting to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type, "expected Type for service setting to be %v, but it was %v", expected.ID, expected.Type, actual.Type)
	assert.Equal(t, expected.Enumeration, actual.Enumeration, "expected Enumeration for service setting to be %v, but it was %v", expected.ID, expected.Enumeration, actual.Enumeration)
	assert.Equal(t, expected.DefaultValue, actual.DefaultValue, "expected DefaultValue for service setting to be %v, but it was %v", expected.ID, expected.DefaultValue, actual.DefaultValue)
	assert.Equal(t, expected.AdminsOnly, actual.AdminsOnly, "expected AdminsOnly for service setting to be %v, but it was %v", expected.ID, expected.AdminsOnly, actual.AdminsOnly)

	assert.NotZero(t, actual.CreatedAt)
}

func buildServiceSettingForTest(t *testing.T, adminClient *apiclient.Client, settingType string, ctx context.Context) *types.ServiceSetting {
	t.Helper()

	t.Log("creating service setting")
	exampleServiceSetting := fakes.BuildFakeServiceSetting()
	exampleServiceSetting.Type = settingType
	exampleServiceSettingInput := converters.ConvertServiceSettingToServiceSettingCreationRequestInput(exampleServiceSetting)
	createdServiceSetting, err := adminClient.CreateServiceSetting(ctx, exampleServiceSettingInput)
	require.NoError(t, err)
	t.Logf("service setting %q created", createdServiceSetting.ID)
	checkServiceSettingEquality(t, exampleServiceSetting, createdServiceSetting)

	createdServiceSetting, err = adminClient.GetServiceSetting(ctx, createdServiceSetting.ID)
	requireNotNilAndNoProblems(t, createdServiceSetting, err)
	checkServiceSettingEquality(t, exampleServiceSetting, createdServiceSetting)

	return createdServiceSetting
}

func (s *TestSuite) TestServiceSettings_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdServiceSetting := buildServiceSettingForTest(t, testClients.admin, "user", s.ctx)

			t.Log("changing service setting")
			newServiceSetting := fakes.BuildFakeServiceSetting()
			createdServiceSetting.Update(converters.ConvertServiceSettingToServiceSettingUpdateRequestInput(newServiceSetting))
			assert.NoError(t, testClients.admin.UpdateServiceSetting(ctx, createdServiceSetting))

			t.Log("fetching changed service setting")
			actual, err := testClients.admin.GetServiceSetting(ctx, createdServiceSetting.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert service setting equality
			checkServiceSettingEquality(t, newServiceSetting, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up service setting")
			assert.NoError(t, testClients.admin.ArchiveServiceSetting(ctx, createdServiceSetting.ID))
		}
	})
}

func (s *TestSuite) TestServiceSettings_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating service settings")
			var expected []*types.ServiceSetting
			for i := 0; i < 5; i++ {
				exampleServiceSetting := fakes.BuildFakeServiceSetting()
				exampleServiceSettingInput := converters.ConvertServiceSettingToServiceSettingCreationRequestInput(exampleServiceSetting)
				createdServiceSetting, createdServiceSettingErr := testClients.admin.CreateServiceSetting(ctx, exampleServiceSettingInput)
				require.NoError(t, createdServiceSettingErr)
				t.Logf("service setting %q created", createdServiceSetting.ID)

				checkServiceSettingEquality(t, exampleServiceSetting, createdServiceSetting)

				expected = append(expected, createdServiceSetting)
			}

			// assert service setting list equality
			actual, err := testClients.admin.GetServiceSettings(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			t.Log("cleaning up")
			for _, createdServiceSetting := range expected {
				assert.NoError(t, testClients.admin.ArchiveServiceSetting(ctx, createdServiceSetting.ID))
			}
		}
	})
}

func (s *TestSuite) TestServiceSettings_Searching() {
	s.runForEachClient("should be able to be search for service settings", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating service settings")
			var expected []*types.ServiceSetting
			exampleServiceSetting := fakes.BuildFakeServiceSetting()
			exampleServiceSetting.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleServiceSetting.Name
			for i := 0; i < 5; i++ {
				exampleServiceSetting.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleServiceSettingInput := converters.ConvertServiceSettingToServiceSettingCreationRequestInput(exampleServiceSetting)
				createdServiceSetting, createdServiceSettingErr := testClients.admin.CreateServiceSetting(ctx, exampleServiceSettingInput)
				require.NoError(t, createdServiceSettingErr)
				t.Logf("service setting %q created", createdServiceSetting.ID)
				checkServiceSettingEquality(t, exampleServiceSetting, createdServiceSetting)

				expected = append(expected, createdServiceSetting)
			}

			exampleLimit := uint8(20)

			// assert service setting list equality
			actual, err := testClients.admin.SearchServiceSettings(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdServiceSetting := range expected {
				assert.NoError(t, testClients.admin.ArchiveServiceSetting(ctx, createdServiceSetting.ID))
			}
		}
	})
}
