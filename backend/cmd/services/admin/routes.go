package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/platform/routing"

	ghttp "maragu.dev/gomponents/http"
)

const (
	assetsDir = "./cmd/services/admin/assets"
)

func (s *AdminFrontendServer) setupRoutes(router routing.Router) {
	r := router.WithMiddleware(s.authMiddleware)

	r.Get("/", ghttp.Adapt(s.HomePage))

	r.Get(fmt.Sprintf("/users/{%s}", userIDURLParamKey), ghttp.Adapt(s.UserPage))
	r.Get("/users", ghttp.Adapt(s.UsersList))
	r.Get("/api/users/search", ghttp.Adapt(s.UsersSearch))
	r.Get(fmt.Sprintf("/api/users/{%s}/accounts", userIDURLParamKey), ghttp.Adapt(s.UserAccountsList))
	r.Get(fmt.Sprintf("/api/users/{%s}/audit-log", userIDURLParamKey), ghttp.Adapt(s.UserAuditLogList))
	r.Get(fmt.Sprintf("/api/users/{%s}/subscriptions", userIDURLParamKey), ghttp.Adapt(s.UserSubscriptionsList))

	r.Get(fmt.Sprintf("/accounts/{%s}", accountIDURLParamKey), ghttp.Adapt(s.AccountPage))
	r.Get("/accounts", ghttp.Adapt(s.AccountsList))
	r.Get("/api/accounts/search", ghttp.Adapt(s.AccountsSearch))
	r.Get(fmt.Sprintf("/api/accounts/{%s}/users", accountIDURLParamKey), ghttp.Adapt(s.AccountUsersList))
	r.Get(fmt.Sprintf("/api/accounts/{%s}/audit-log", accountIDURLParamKey), ghttp.Adapt(s.AccountAuditLogList))
	r.Get(fmt.Sprintf("/api/accounts/{%s}/subscriptions", accountIDURLParamKey), ghttp.Adapt(s.AccountSubscriptionsList))

	r.Get(fmt.Sprintf("/oauth2_clients/{%s}", oauth2ClientIDURLParamKey), ghttp.Adapt(s.OAuth2ClientPage))
	r.Get("/oauth2_clients", ghttp.Adapt(s.OAuth2ClientsList))
	r.Get("/api/oauth2_clients/search", ghttp.Adapt(s.OAuth2ClientsSearch))

	// General search endpoints for search boxes
	r.Get("/admin/search/valid_measurement_units", ghttp.Adapt(s.SearchValidMeasurementUnits))
	r.Get("/admin/search/valid_ingredients", ghttp.Adapt(s.SearchValidIngredients))

	// Valid Ingredients - specific routes before dynamic ones
	r.Get("/valid_ingredients/new", ghttp.Adapt(s.ValidIngredientNewPage))
	r.Post("/api/valid_ingredients", ghttp.Adapt(s.ValidIngredientCreate))
	r.Get(fmt.Sprintf("/valid_ingredients/{%s}", validIngredientIDURLParamKey), ghttp.Adapt(s.ValidIngredientPage))
	r.Get("/valid_ingredients", ghttp.Adapt(s.ValidIngredientsList))
	r.Get("/api/valid_ingredients/search", ghttp.Adapt(s.ValidIngredientsSearch))

	// Valid Ingredient Measurement Unit associations (from ingredient side)
	r.Get(fmt.Sprintf("/api/valid_ingredients/{%s}/measurement_units", validIngredientIDURLParamKey), ghttp.Adapt(s.ValidIngredientMeasurementUnitsForIngredient))
	r.Get(fmt.Sprintf("/api/valid_ingredients/{%s}/measurement_units/search", validIngredientIDURLParamKey), ghttp.Adapt(s.SearchMeasurementUnitsForIngredient))
	r.Post(fmt.Sprintf("/api/valid_ingredients/{%s}/measurement_units", validIngredientIDURLParamKey), ghttp.Adapt(s.CreateIngredientMeasurementUnitFromIngredient))

	// Valid Ingredient Preparation associations (from ingredient side)
	r.Get(fmt.Sprintf("/api/valid_ingredients/{%s}/preparations", validIngredientIDURLParamKey), ghttp.Adapt(s.ValidIngredientPreparationsForIngredient))
	r.Get(fmt.Sprintf("/api/valid_ingredients/{%s}/preparations/search", validIngredientIDURLParamKey), ghttp.Adapt(s.SearchPreparationsForIngredient))
	r.Post(fmt.Sprintf("/api/valid_ingredients/{%s}/preparations", validIngredientIDURLParamKey), ghttp.Adapt(s.CreateIngredientPreparationFromIngredient))

	// Valid Instruments - specific routes before dynamic ones
	r.Get("/valid_instruments/new", ghttp.Adapt(s.ValidInstrumentNewPage))
	r.Post("/api/valid_instruments", ghttp.Adapt(s.ValidInstrumentCreate))
	r.Get(fmt.Sprintf("/valid_instruments/{%s}", validInstrumentIDURLParamKey), ghttp.Adapt(s.ValidInstrumentPage))
	r.Get("/valid_instruments", ghttp.Adapt(s.ValidInstrumentsList))
	r.Get("/api/valid_instruments/search", ghttp.Adapt(s.ValidInstrumentsSearch))

	// Valid Preparation Instrument associations (from instrument side)
	r.Get(fmt.Sprintf("/api/valid_instruments/{%s}/preparations", validInstrumentIDURLParamKey), ghttp.Adapt(s.ValidPreparationInstrumentsForInstrument))
	r.Get(fmt.Sprintf("/api/valid_instruments/{%s}/preparations/search", validInstrumentIDURLParamKey), ghttp.Adapt(s.SearchPreparationsForInstrument))
	r.Post(fmt.Sprintf("/api/valid_instruments/{%s}/preparations", validInstrumentIDURLParamKey), ghttp.Adapt(s.CreatePreparationInstrumentFromInstrument))

	// Valid Vessels - specific routes before dynamic ones
	r.Get("/valid_vessels/new", ghttp.Adapt(s.ValidVesselNewPage))
	r.Post("/api/valid_vessels", ghttp.Adapt(s.ValidVesselCreate))
	r.Get(fmt.Sprintf("/valid_vessels/{%s}", validVesselIDURLParamKey), ghttp.Adapt(s.ValidVesselPage))
	r.Get("/valid_vessels", ghttp.Adapt(s.ValidVesselsList))
	r.Get("/api/valid_vessels/search", ghttp.Adapt(s.ValidVesselsSearch))

	// Valid Preparation Vessel associations (from vessel side)
	r.Get(fmt.Sprintf("/api/valid_vessels/{%s}/preparations", validVesselIDURLParamKey), ghttp.Adapt(s.ValidPreparationVesselsForVessel))
	r.Get(fmt.Sprintf("/api/valid_vessels/{%s}/preparations/search", validVesselIDURLParamKey), ghttp.Adapt(s.SearchPreparationsForVessel))
	r.Post(fmt.Sprintf("/api/valid_vessels/{%s}/preparations", validVesselIDURLParamKey), ghttp.Adapt(s.CreatePreparationVesselFromVessel))

	// Valid Measurement Units - specific routes before dynamic ones
	r.Get("/valid_measurement_units/new", ghttp.Adapt(s.ValidMeasurementUnitNewPage))
	r.Post("/api/valid_measurement_units", ghttp.Adapt(s.ValidMeasurementUnitCreate))
	r.Get(fmt.Sprintf("/valid_measurement_units/{%s}", validMeasurementUnitIDURLParamKey), ghttp.Adapt(s.ValidMeasurementUnitPage))
	r.Get("/valid_measurement_units", ghttp.Adapt(s.ValidMeasurementUnitsList))
	r.Get("/api/valid_measurement_units/search", ghttp.Adapt(s.ValidMeasurementUnitsSearch))

	// Valid Ingredient Measurement Unit associations (from measurement unit side)
	r.Get(fmt.Sprintf("/api/valid_measurement_units/{%s}/ingredients", validMeasurementUnitIDURLParamKey), ghttp.Adapt(s.ValidIngredientMeasurementUnitsForMeasurementUnit))
	r.Get(fmt.Sprintf("/api/valid_measurement_units/{%s}/ingredients/search", validMeasurementUnitIDURLParamKey), ghttp.Adapt(s.SearchIngredientsForMeasurementUnit))
	r.Post(fmt.Sprintf("/api/valid_measurement_units/{%s}/ingredients", validMeasurementUnitIDURLParamKey), ghttp.Adapt(s.CreateIngredientMeasurementUnitFromMeasurementUnit))

	// Valid Measurement Unit Conversions
	r.Get(fmt.Sprintf("/api/valid_measurement_units/{%s}/conversions", validMeasurementUnitIDURLParamKey), ghttp.Adapt(s.ValidMeasurementUnitConversionsForUnit))
	r.Get(fmt.Sprintf("/api/valid_measurement_units/{%s}/conversions/search", validMeasurementUnitIDURLParamKey), ghttp.Adapt(s.SearchMeasurementUnitsForConversion))
	r.Post("/api/valid_measurement_unit_conversions", ghttp.Adapt(s.CreateMeasurementUnitConversion))

	// Valid Ingredient States - specific routes before dynamic ones
	r.Get("/valid_ingredient_states/new", ghttp.Adapt(s.ValidIngredientStateNewPage))
	r.Post("/api/valid_ingredient_states", ghttp.Adapt(s.ValidIngredientStateCreate))
	r.Get(fmt.Sprintf("/api/valid_ingredient_states/{%s}/ingredients/search", validIngredientStateIDURLParamKey), ghttp.Adapt(s.ValidIngredientStateIngredientsSearch))
	r.Get(fmt.Sprintf("/valid_ingredient_states/{%s}", validIngredientStateIDURLParamKey), ghttp.Adapt(s.ValidIngredientStatePage))
	r.Get("/valid_ingredient_states", ghttp.Adapt(s.ValidIngredientStatesList))
	r.Get("/api/valid_ingredient_states/search", ghttp.Adapt(s.ValidIngredientStatesSearch))

	// Valid Preparations - specific routes before dynamic ones
	r.Get("/valid_preparations/new", ghttp.Adapt(s.ValidPreparationNewPage))
	r.Post("/api/valid_preparations", ghttp.Adapt(s.ValidPreparationCreate))
	r.Get(fmt.Sprintf("/valid_preparations/{%s}", validPreparationIDURLParamKey), ghttp.Adapt(s.ValidPreparationPage))
	r.Get("/valid_preparations", ghttp.Adapt(s.ValidPreparationsList))
	r.Get("/api/valid_preparations/search", ghttp.Adapt(s.ValidPreparationsSearch))

	// Valid Preparation Instrument associations (from preparation side)
	r.Get(fmt.Sprintf("/api/valid_preparations/{%s}/instruments", validPreparationIDURLParamKey), ghttp.Adapt(s.ValidPreparationInstrumentsForPreparation))
	r.Get(fmt.Sprintf("/api/valid_preparations/{%s}/instruments/search", validPreparationIDURLParamKey), ghttp.Adapt(s.SearchInstrumentsForPreparation))
	r.Post(fmt.Sprintf("/api/valid_preparations/{%s}/instruments", validPreparationIDURLParamKey), ghttp.Adapt(s.CreatePreparationInstrumentFromPreparation))

	// Valid Ingredient Preparation associations (from preparation side)
	r.Get(fmt.Sprintf("/api/valid_preparations/{%s}/ingredients", validPreparationIDURLParamKey), ghttp.Adapt(s.ValidIngredientPreparationsForPreparation))
	r.Get(fmt.Sprintf("/api/valid_preparations/{%s}/ingredients/search", validPreparationIDURLParamKey), ghttp.Adapt(s.SearchIngredientsForPreparation))
	r.Post(fmt.Sprintf("/api/valid_preparations/{%s}/ingredients", validPreparationIDURLParamKey), ghttp.Adapt(s.CreateIngredientPreparationFromPreparation))

	// Valid Preparation Vessel associations (from preparation side)
	r.Get(fmt.Sprintf("/api/valid_preparations/{%s}/vessels", validPreparationIDURLParamKey), ghttp.Adapt(s.ValidPreparationVesselsForPreparation))
	r.Get(fmt.Sprintf("/api/valid_preparations/{%s}/vessels/search", validPreparationIDURLParamKey), ghttp.Adapt(s.SearchVesselsForPreparation))
	r.Post(fmt.Sprintf("/api/valid_preparations/{%s}/vessels", validPreparationIDURLParamKey), ghttp.Adapt(s.CreatePreparationVesselFromPreparation))

	r.Get("/settings/new", ghttp.Adapt(s.SettingNewPage))
	r.Post("/api/settings", ghttp.Adapt(s.SettingCreate))
	r.Get(fmt.Sprintf("/settings/{%s}", settingIDURLParamKey), ghttp.Adapt(s.SettingPage))
	r.Get("/settings", ghttp.Adapt(s.SettingsList))
	r.Get("/api/settings/search", ghttp.Adapt(s.SettingsSearch))

	// Products
	r.Get("/products/new", ghttp.Adapt(s.ProductNewPage))
	r.Post("/api/products", ghttp.Adapt(s.ProductCreate))
	r.Get(fmt.Sprintf("/products/{%s}", productIDURLParamKey), ghttp.Adapt(s.ProductPage))
	r.Post(fmt.Sprintf("/api/products/{%s}", productIDURLParamKey), ghttp.Adapt(s.ProductUpdate))
	r.Get("/products", ghttp.Adapt(s.ProductsList))
	r.Get("/api/products/search", ghttp.Adapt(s.ProductsSearch))

	// Subscriptions
	r.Get("/subscriptions/new", ghttp.Adapt(s.SubscriptionNewPage))
	r.Post("/api/subscriptions", ghttp.Adapt(s.SubscriptionCreate))
	r.Get(fmt.Sprintf("/subscriptions/{%s}", subscriptionIDURLParamKey), ghttp.Adapt(s.SubscriptionPage))
	r.Post(fmt.Sprintf("/api/subscriptions/{%s}", subscriptionIDURLParamKey), ghttp.Adapt(s.SubscriptionUpdate))
	r.Get("/subscriptions", ghttp.Adapt(s.SubscriptionsList))
	r.Get("/api/subscriptions/search", ghttp.Adapt(s.SubscriptionsSearch))

	// Valid Prep Task Configs
	r.Get(fmt.Sprintf("/valid_prep_task_configs/{%s}", validPrepTaskConfigIDURLParamKey), ghttp.Adapt(s.ValidPrepTaskConfigPage))
	r.Get("/valid_prep_task_configs", ghttp.Adapt(s.ValidPrepTaskConfigsList))
	r.Get("/api/valid_prep_task_configs/search", ghttp.Adapt(s.ValidPrepTaskConfigsSearch))

	// Recipes
	r.Get(fmt.Sprintf("/recipes/{%s}", recipeIDURLParamKey), ghttp.Adapt(s.RecipePage))
	r.Get("/recipes", ghttp.Adapt(s.RecipesList))
	r.Get("/api/recipes/search", ghttp.Adapt(s.RecipesSearch))

	// Waitlists
	r.Get(fmt.Sprintf("/waitlists/{%s}", waitlistIDURLParamKey), ghttp.Adapt(s.WaitlistPage))
	r.Get("/waitlists", ghttp.Adapt(s.WaitlistsList))
	r.Get("/api/waitlists/search", ghttp.Adapt(s.WaitlistsSearch))
	r.Get(fmt.Sprintf("/api/waitlists/{%s}/signups", waitlistIDURLParamKey), ghttp.Adapt(s.WaitlistSignupsForWaitlist))

	// Issue Reports
	r.Get(fmt.Sprintf("/issue_reports/{%s}", issueReportIDURLParamKey), ghttp.Adapt(s.IssueReportPage))
	r.Get("/issue_reports", ghttp.Adapt(s.IssueReportsList))
	r.Get("/api/issue_reports/search", ghttp.Adapt(s.IssueReportsSearch))

	// Queue Test
	r.Get("/queue_test", ghttp.Adapt(s.QueueTestPage))
	r.Post("/api/queue_test", ghttp.Adapt(s.QueueTestSubmit))

	// Association delete routes
	r.Delete("/api/valid_preparation_instruments/{associationID}", ghttp.Adapt(s.DeletePreparationInstrument))
	r.Delete("/api/valid_preparation_vessels/{associationID}", ghttp.Adapt(s.DeletePreparationVessel))
	r.Delete("/api/valid_ingredient_measurement_units/{associationID}", ghttp.Adapt(s.DeleteIngredientMeasurementUnit))
	r.Delete("/api/valid_ingredient_preparations/{associationID}", ghttp.Adapt(s.DeleteIngredientPreparation))
	r.Delete("/api/valid_measurement_unit_conversions/{conversionID}", ghttp.Adapt(s.DeleteMeasurementUnitConversion))

	router.Get("/login", ghttp.Adapt(s.LoginPage))
	router.Post("/login/submit", ghttp.Adapt(s.LoginSubmission))

	// static files - NOTE: this must be registered last
	fileServer := http.FileServer(http.Dir(assetsDir))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Only serve root-level files (no subdirectories)
		if strings.Contains(r.URL.Path[1:], "/") {
			http.NotFound(w, r)
			return
		}

		// Check if the file exists in the assets directory (guard against path traversal)
		filePath := filepath.Clean(filepath.Join(assetsDir, r.URL.Path))
		absAssets, err := filepath.Abs(assetsDir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		rel, err := filepath.Rel(absAssets, absFilePath)
		if err != nil || strings.HasPrefix(rel, "..") {
			http.NotFound(w, r)
			return
		}
		info, statErr := os.Stat(filePath) //nolint:gosec // G703: path validated above to be within assetsDir
		if statErr == nil && !info.IsDir() {
			// File exists, serve it
			fileServer.ServeHTTP(w, r)
			return
		}

		// File doesn't exist
		http.NotFound(w, r)
	})
}
