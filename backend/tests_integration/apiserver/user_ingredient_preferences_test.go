package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	settingsconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createUserIngredientPreferenceForTest(t *testing.T, clientToUse client.Client) *mealplanning.UserIngredientPreference {
	t.Helper()
	ctx := t.Context()

	validIngredient := createValidIngredientForTest(t)

	exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
	exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
	exampleUserIngredientPreferenceInput.ValidIngredientID = validIngredient.ID
	createdUserIngredientPreference, err := clientToUse.CreateUserIngredientPreference(ctx, &settingssvc.CreateUserIngredientPreferenceRequest{
		Input: settingsconverters.ConvertUserIngredientPreferenceCreationRequestInputToGRPCUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreferenceInput),
	})
	require.NoError(t, err)
	converted := settingsconverters.ConvertGRPCUserIngredientPreferenceToUserIngredientPreference(createdUserIngredientPreference.Created[0])
	assertRoughEquality(t, exampleUserIngredientPreference, converted, defaultIgnoredFields("ID", "BelongsToUser", "Ingredient")...)

	res, err := clientToUse.GetUserIngredientPreference(ctx, &settingssvc.GetUserIngredientPreferenceRequest{UserIngredientPreferenceID: createdUserIngredientPreference.Created[0].ID})
	require.NoError(t, err)
	require.NotNil(t, res)

	serviceSetting := settingsconverters.ConvertGRPCUserIngredientPreferenceToUserIngredientPreference(res.Result)
	assertRoughEquality(t, converted, serviceSetting, defaultIgnoredFields("ID", "BelongsToUser", "Ingredient")...)

	return serviceSetting
}

func TestUserIngredientPreferences_Creating(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		_, testClient := createUserAndClientForTest(t)
		createUserIngredientPreferenceForTest(t, testClient)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		creationRequestInput := fakes.BuildFakeUserIngredientPreferenceCreationRequestInput()
		convertedInput := settingsconverters.ConvertUserIngredientPreferenceCreationRequestInputToGRPCUserIngredientPreferenceCreationRequestInput(creationRequestInput)

		c := buildUnauthenticatedGRPCClientForTest(t)
		created, err := c.CreateUserIngredientPreference(ctx, &settingssvc.CreateUserIngredientPreferenceRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("invalid input", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		creationRequestInput := fakes.BuildFakeUserIngredientPreferenceCreationRequestInput()
		convertedInput := settingsconverters.ConvertUserIngredientPreferenceCreationRequestInputToGRPCUserIngredientPreferenceCreationRequestInput(creationRequestInput)
		// this is not allowed
		convertedInput.ValidIngredientID = ""
		convertedInput.ValidIngredientGroupID = ""

		created, err := testClient.CreateUserIngredientPreference(ctx, &settingssvc.CreateUserIngredientPreferenceRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})

	T.Run("non-admin users are forbidden from creating", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(T)

		creationRequestInput := fakes.BuildFakeUserIngredientPreferenceCreationRequestInput()
		convertedInput := settingsconverters.ConvertUserIngredientPreferenceCreationRequestInputToGRPCUserIngredientPreferenceCreationRequestInput(creationRequestInput)

		created, err := testClient.CreateUserIngredientPreference(ctx, &settingssvc.CreateUserIngredientPreferenceRequest{
			Input: convertedInput,
		})
		assert.Error(t, err)
		assert.Nil(t, created)
	})
}

func TestUserIngredientPreferences_Reading(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createUserIngredientPreferenceForTest(t, testClient)

		retrieved, err := testClient.GetUserIngredientPreference(ctx, &settingssvc.GetUserIngredientPreferenceRequest{UserIngredientPreferenceID: created.ID})
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		converted := settingsconverters.ConvertGRPCUserIngredientPreferenceToUserIngredientPreference(retrieved.Result)

		assertRoughEquality(t, created, converted, defaultIgnoredFields("ID", "BelongsToUser", "Ingredient")...)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		created := createUserIngredientPreferenceForTest(t, testClient)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetUserIngredientPreference(ctx, &settingssvc.GetUserIngredientPreferenceRequest{UserIngredientPreferenceID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, err := testClient.GetUserIngredientPreference(ctx, &settingssvc.GetUserIngredientPreferenceRequest{UserIngredientPreferenceID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestUserIngredientPreferences_Archiving(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		created := createUserIngredientPreferenceForTest(t, testClient)

		_, err := testClient.ArchiveUserIngredientPreference(ctx, &settingssvc.ArchiveUserIngredientPreferenceRequest{UserIngredientPreferenceID: created.ID})
		assert.NoError(t, err)

		x, err := testClient.GetUserIngredientPreference(ctx, &settingssvc.GetUserIngredientPreferenceRequest{UserIngredientPreferenceID: created.ID})
		assert.Nil(t, x)
		assert.Error(t, err)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)
		created := createUserIngredientPreferenceForTest(t, testClient)

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.ArchiveUserIngredientPreference(ctx, &settingssvc.ArchiveUserIngredientPreferenceRequest{UserIngredientPreferenceID: created.ID})
		assert.Error(t, err)
	})

	T.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		_, err := testClient.ArchiveUserIngredientPreference(ctx, &settingssvc.ArchiveUserIngredientPreferenceRequest{UserIngredientPreferenceID: nonexistentID})
		assert.Error(t, err)
	})
}

func TestUserIngredientPreferences_Listing(T *testing.T) {
	T.Parallel()

	_, testClient := createUserAndClientForTest(T)
	createdUserIngredientPreferences := []*mealplanning.UserIngredientPreference{}
	for range exampleQuantity {
		created := createUserIngredientPreferenceForTest(T, testClient)
		createdUserIngredientPreferences = append(createdUserIngredientPreferences, created)
	}

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		retrieved, err := testClient.GetUserIngredientPreferences(ctx, &settingssvc.GetUserIngredientPreferencesRequest{})
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.True(t, len(retrieved.Results) >= len(createdUserIngredientPreferences))
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		c := buildUnauthenticatedGRPCClientForTest(t)

		_, err := c.GetUserIngredientPreferences(ctx, &settingssvc.GetUserIngredientPreferencesRequest{})
		assert.Error(t, err)
	})
}
