package config

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	customerdataconfig "github.com/prixfixeco/backend/internal/customerdata/config"
	dbconfig "github.com/prixfixeco/backend/internal/database/config"
	emailconfig "github.com/prixfixeco/backend/internal/email/config"
	"github.com/prixfixeco/backend/internal/email/sendgrid"
	"github.com/prixfixeco/backend/internal/encoding"
	msgconfig "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/routing"
	"github.com/prixfixeco/backend/internal/server"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	householdinvitationsservice "github.com/prixfixeco/backend/internal/services/householdinvitations"
	householdsservice "github.com/prixfixeco/backend/internal/services/households"
	mealplanoptionsservice "github.com/prixfixeco/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/prixfixeco/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/prixfixeco/backend/internal/services/mealplans"
	mealsservice "github.com/prixfixeco/backend/internal/services/meals"
	recipepreptasksservice "github.com/prixfixeco/backend/internal/services/recipepreptasks"
	recipesservice "github.com/prixfixeco/backend/internal/services/recipes"
	recipestepingredientsservice "github.com/prixfixeco/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/prixfixeco/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/prixfixeco/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/prixfixeco/backend/internal/services/recipesteps"
	usersservice "github.com/prixfixeco/backend/internal/services/users"
	validingredientmeasurementunitsservice "github.com/prixfixeco/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/prixfixeco/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/prixfixeco/backend/internal/services/validingredients"
	validingredientstatesservice "github.com/prixfixeco/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/prixfixeco/backend/internal/services/validinstruments"
	"github.com/prixfixeco/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/prixfixeco/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/prixfixeco/backend/internal/services/validpreparations"
	webhooksservice "github.com/prixfixeco/backend/internal/services/webhooks"
	websocketsservice "github.com/prixfixeco/backend/internal/services/websockets"
	testutils "github.com/prixfixeco/backend/tests/utils"
)

type mockSecretVersionAccessor struct {
	mock.Mock
}

func (m *mockSecretVersionAccessor) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	args := m.Called(ctx, req, opts)

	return args.Get(0).(*secretmanagerpb.AccessSecretVersionResponse), args.Error(1)
}

func TestGetAPIServerConfigFromGoogleCloudRunEnvironment(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		baseConfig := &InstanceConfig{
			Observability: observability.Config{},
			Email: emailconfig.Config{
				Sendgrid: &sendgrid.Config{},
				Provider: emailconfig.ProviderSendgrid,
			},
			CustomerData: customerdataconfig.Config{
				Provider: customerdataconfig.ProviderSegment,
				APIToken: "",
			},
			Encoding: encoding.Config{ContentType: "application/json"},
			Routing:  routing.Config{},
			Database: dbconfig.Config{},
			Meta:     MetaSettings{},
			Events:   msgconfig.Config{},
			Server: server.Config{
				StartupDeadline: time.Second,
			},
			Services: ServicesConfigurations{
				ValidMeasurementUnits:           validmeasurementunits.Config{},
				ValidInstruments:                validinstrumentsservice.Config{},
				ValidIngredients:                validingredientsservice.Config{},
				ValidPreparations:               validpreparationsservice.Config{},
				MealPlanOptionVotes:             mealplanoptionvotesservice.Config{},
				ValidIngredientPreparations:     validingredientpreparationsservice.Config{},
				ValidPreparationInstruments:     validpreparationinstrumentsservice.Config{},
				ValidInstrumentMeasurementUnits: validingredientmeasurementunitsservice.Config{},
				Meals:                           mealsservice.Config{},
				Recipes: recipesservice.Config{
					PublicMediaURLPrefix: t.Name(),
				},
				RecipeSteps: recipestepsservice.Config{
					PublicMediaURLPrefix: t.Name(),
				},
				RecipeStepProducts:    recipestepproductsservice.Config{},
				RecipeStepInstruments: recipestepinstrumentsservice.Config{},
				RecipeStepIngredients: recipestepingredientsservice.Config{},
				MealPlans:             mealplansservice.Config{},
				MealPlanOptions:       mealplanoptionsservice.Config{},
				Households:            householdsservice.Config{},
				HouseholdInvitations:  householdinvitationsservice.Config{},
				Websockets:            websocketsservice.Config{},
				Webhooks:              webhooksservice.Config{},
				Users:                 usersservice.Config{},
				RecipePrepTasks:       recipepreptasksservice.Config{},
				ValidIngredientStates: validingredientstatesservice.Config{},
				Auth: authservice.Config{
					MinimumPasswordLength: 8,
					MinimumUsernameLength: 8,
				},
			},
		}

		f, err := os.CreateTemp("", "testing.json")
		require.NoError(t, err)

		require.NoError(t, json.NewEncoder(f).Encode(baseConfig))

		require.NoError(t, os.Setenv(gcpConfigFilePathEnvVarKey, f.Name()))
		require.NoError(t, os.Setenv(gcpPortEnvVarKey, "1234"))
		require.NoError(t, os.Setenv(gcpDatabaseSocketDirEnvVarKey, "/example/blah"))
		require.NoError(t, os.Setenv(gcpDatabaseUserEnvVarKey, "user"))
		require.NoError(t, os.Setenv(gcpDatabaseUserPasswordEnvVarKey, "hunter2"))
		require.NoError(t, os.Setenv(gcpDatabaseNameEnvVarKey, "fake_db_name"))
		require.NoError(t, os.Setenv(gcpDatabaseInstanceConnNameEnvVarKey, "fake_conn_name"))
		require.NoError(t, os.Setenv(gcpCookieHashKeyEnvVarKey, "fake_cookie_hash_key"))
		require.NoError(t, os.Setenv(gcpCookieBlockKeyEnvVarKey, "fake_cookie_block_key"))
		require.NoError(t, os.Setenv(gcpPASETOLocalKeyEnvVarKey, "fake_paseto_local_key"))
		require.NoError(t, os.Setenv(gcpSendgridTokenEnvVarKey, "fake_sendgrid_token"))
		require.NoError(t, os.Setenv(gcpSegmentTokenEnvVarKey, "fake_segment_token"))

		ctx := context.Background()
		client := &mockSecretVersionAccessor{}

		client.On(
			"AccessSecretVersion",
			testutils.ContextMatcher,
			&secretmanagerpb.AccessSecretVersionRequest{Name: buildSecretPathForSecretStore(dataChangesTopicAccessName)},
			[]gax.CallOption(nil),
		).Return(
			&secretmanagerpb.AccessSecretVersionResponse{
				Name: dataChangesTopicAccessName,
				Payload: &secretmanagerpb.SecretPayload{
					Data: []byte("this_is_the_big_secret"),
				},
			},
			nil,
		)

		cfg, err := GetAPIServerConfigFromGoogleCloudRunEnvironment(ctx, client)
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		require.NoError(t, os.Unsetenv(gcpConfigFilePathEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpPortEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseSocketDirEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseUserEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseUserPasswordEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseNameEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseInstanceConnNameEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpCookieHashKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpCookieBlockKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpPASETOLocalKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpSendgridTokenEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpSegmentTokenEnvVarKey))
	})
}
