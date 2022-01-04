package integration

import (
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkValidIngredientEquality(t *testing.T, expected, actual *types.ValidIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Warning, actual.Warning, "expected Warning for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Warning, actual.Warning)
	assert.Equal(t, expected.ContainsEgg, actual.ContainsEgg, "expected ContainsEgg for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsEgg, actual.ContainsEgg)
	assert.Equal(t, expected.ContainsDairy, actual.ContainsDairy, "expected ContainsDairy for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsDairy, actual.ContainsDairy)
	assert.Equal(t, expected.ContainsPeanut, actual.ContainsPeanut, "expected ContainsPeanut for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsPeanut, actual.ContainsPeanut)
	assert.Equal(t, expected.ContainsTreeNut, actual.ContainsTreeNut, "expected ContainsTreeNut for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsTreeNut, actual.ContainsTreeNut)
	assert.Equal(t, expected.ContainsSoy, actual.ContainsSoy, "expected ContainsSoy for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsSoy, actual.ContainsSoy)
	assert.Equal(t, expected.ContainsWheat, actual.ContainsWheat, "expected ContainsWheat for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsWheat, actual.ContainsWheat)
	assert.Equal(t, expected.ContainsShellfish, actual.ContainsShellfish, "expected ContainsShellfish for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsShellfish, actual.ContainsShellfish)
	assert.Equal(t, expected.ContainsSesame, actual.ContainsSesame, "expected ContainsSesame for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsSesame, actual.ContainsSesame)
	assert.Equal(t, expected.ContainsFish, actual.ContainsFish, "expected ContainsFish for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsFish, actual.ContainsFish)
	assert.Equal(t, expected.ContainsGluten, actual.ContainsGluten, "expected ContainsGluten for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsGluten, actual.ContainsGluten)
	assert.Equal(t, expected.AnimalFlesh, actual.AnimalFlesh, "expected AnimalFlesh for valid ingredient %s to be %v, but it was %v", expected.ID, expected.AnimalFlesh, actual.AnimalFlesh)
	assert.Equal(t, expected.AnimalDerived, actual.AnimalDerived, "expected AnimalDerived for valid ingredient %s to be %v, but it was %v", expected.ID, expected.AnimalDerived, actual.AnimalDerived)
	assert.Equal(t, expected.Volumetric, actual.Volumetric, "expected Volumetric for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Volumetric, actual.Volumetric)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.NotZero(t, actual.CreatedOn)
}

// convertValidIngredientToValidIngredientUpdateInput creates an ValidIngredientUpdateRequestInput struct from a valid ingredient.
func convertValidIngredientToValidIngredientUpdateInput(x *types.ValidIngredient) *types.ValidIngredientUpdateRequestInput {
	return &types.ValidIngredientUpdateRequestInput{
		Name:              x.Name,
		Variant:           x.Variant,
		Description:       x.Description,
		Warning:           x.Warning,
		ContainsEgg:       x.ContainsEgg,
		ContainsDairy:     x.ContainsDairy,
		ContainsPeanut:    x.ContainsPeanut,
		ContainsTreeNut:   x.ContainsTreeNut,
		ContainsSoy:       x.ContainsSoy,
		ContainsWheat:     x.ContainsWheat,
		ContainsShellfish: x.ContainsShellfish,
		ContainsSesame:    x.ContainsSesame,
		ContainsFish:      x.ContainsFish,
		ContainsGluten:    x.ContainsGluten,
		AnimalFlesh:       x.AnimalFlesh,
		AnimalDerived:     x.AnimalDerived,
		Volumetric:        x.Volumetric,
		IconPath:          x.IconPath,
	}
}

/*

func (s *TestSuite) TestValidIngredients_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.admin.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			n = <-notificationsChan
			assert.Equal(t, types.ValidIngredientDataType, n.DataType)
			require.NotNil(t, n.ValidIngredient)
			checkValidIngredientEquality(t, exampleValidIngredient, n.ValidIngredient)

			createdValidIngredient, err = testClients.admin.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("changing valid ingredient")
			newValidIngredient := fakes.BuildFakeValidIngredient()
			createdValidIngredient.Update(convertValidIngredientToValidIngredientUpdateInput(newValidIngredient))
			assert.NoError(t, testClients.admin.UpdateValidIngredient(ctx, createdValidIngredient))

			n = <-notificationsChan
			assert.Equal(t, types.ValidIngredientDataType, n.DataType)

			t.Log("fetching changed valid ingredient")
			actual, err := testClients.admin.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient equality
			checkValidIngredientEquality(t, newValidIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid ingredient")
			assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)

			// change valid ingredient
			newValidIngredient := fakes.BuildFakeValidIngredient()
			createdValidIngredient.Update(convertValidIngredientToValidIngredientUpdateInput(newValidIngredient))
			assert.NoError(t, testClients.admin.UpdateValidIngredient(ctx, createdValidIngredient))

			time.Sleep(2 * time.Second)

			// retrieve changed valid ingredient
			var actual *types.ValidIngredient
			checkFunc = func() bool {
				actual, err = testClients.admin.GetValidIngredient(ctx, createdValidIngredient.ID)
				return assert.NotNil(t, createdValidIngredient) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient equality
			checkValidIngredientEquality(t, newValidIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up valid ingredient")
			assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredients_Listing() {
	s.runForCookieClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.admin.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating valid ingredients")
			var expected []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, createdValidIngredientErr)
				t.Logf("valid ingredient %q created", createdValidIngredient.ID)

				n = <-notificationsChan
				assert.Equal(t, types.ValidIngredientDataType, n.DataType)
				require.NotNil(t, n.ValidIngredient)
				checkValidIngredientEquality(t, exampleValidIngredient, n.ValidIngredient)

				expected = append(expected, createdValidIngredient)
			}

			// assert valid ingredient list equality
			actual, err := testClients.admin.GetValidIngredients(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredients),
			)

			t.Log("cleaning up")
			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredients")
			var expected []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, err)
				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				expected = append(expected, createdValidIngredient)
			}

			// assert valid ingredient list equality
			actual, err := testClients.admin.GetValidIngredients(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredients),
			)

			t.Log("cleaning up")
			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredients_Searching() {
	s.runForCookieClient("should be able to be search for valid ingredients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.admin.SubscribeToNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating valid ingredients")
			var expected []*types.ValidIngredient
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredient.Name = "example"
			searchQuery := exampleValidIngredient.Name
			for i := 0; i < 5; i++ {
				exampleValidIngredient.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, createdValidIngredientErr)
				t.Logf("valid ingredient %q created", createdValidIngredient.ID)

				n = <-notificationsChan
				assert.Equal(t, types.ValidIngredientDataType, n.DataType)
				require.NotNil(t, n.ValidIngredient)
				checkValidIngredientEquality(t, exampleValidIngredient, n.ValidIngredient)

				expected = append(expected, createdValidIngredient)
			}

			exampleLimit := uint8(20)

			// give the index a moment
			time.Sleep(3 * time.Second)

			// assert valid ingredient list equality
			actual, err := testClients.admin.SearchValidIngredients(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})

	s.runForPASETOClient("should be able to be search for valid ingredients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating valid ingredients")
			var expected []*types.ValidIngredient
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredient.Name = "example"
			searchQuery := exampleValidIngredient.Name
			for i := 0; i < 5; i++ {
				exampleValidIngredient.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, err)
				requireNotNilAndNoProblems(t, createdValidIngredient, err)
				t.Logf("valid ingredient %q created", createdValidIngredient.ID)

				expected = append(expected, createdValidIngredient)
			}

			exampleLimit := uint8(20)
			time.Sleep(2 * time.Second) // give the index a moment

			// assert valid ingredient list equality
			actual, err := testClients.admin.SearchValidIngredients(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})
}

*/
