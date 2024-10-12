package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestMealPlanEvents(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(mealPlanEventsTestSuite))
}

type mealPlanEventsBaseSuite struct {
	suite.Suite
	ctx                              context.Context
	exampleMealPlanEvent             *types.MealPlanEvent
	exampleMealPlanEventResponse     *types.APIResponse[*types.MealPlanEvent]
	exampleMealPlanEventListResponse *types.APIResponse[[]*types.MealPlanEvent]
	exampleMealPlanID                string
	exampleMealPlanEventList         []*types.MealPlanEvent
}

var _ suite.SetupTestSuite = (*mealPlanEventsBaseSuite)(nil)

func (s *mealPlanEventsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleMealPlanID = fakes.BuildFakeID()
	s.exampleMealPlanEvent = fakes.BuildFakeMealPlanEvent()
	s.exampleMealPlanEvent.BelongsToMealPlan = s.exampleMealPlanID
	s.exampleMealPlanEventResponse = &types.APIResponse[*types.MealPlanEvent]{
		Data: s.exampleMealPlanEvent,
	}
	exampleMealPlanEventList := fakes.BuildFakeMealPlanEventsList()
	s.exampleMealPlanEventList = exampleMealPlanEventList.Data
	s.exampleMealPlanEventListResponse = &types.APIResponse[[]*types.MealPlanEvent]{
		Data:       s.exampleMealPlanEventList,
		Pagination: &exampleMealPlanEventList.Pagination,
	}
}

type mealPlanEventsTestSuite struct {
	suite.Suite
	mealPlanEventsBaseSuite
}

func (s *mealPlanEventsTestSuite) TestClient_GetMealPlanEvent() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanEventResponse)
		actual, err := c.GetMealPlanEvent(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanEvent, actual)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanEvent(s.ctx, "", s.exampleMealPlanEvent.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan ClientOption ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanEvent(s.ctx, s.exampleMealPlanID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanEvent(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanEvent(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanEventsTestSuite) TestClient_GetMealPlanEvents() {
	const expectedPath = "/api/v1/meal_plans/%s/events"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleMealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanEventListResponse)
		actual, err := c.GetMealPlanEvents(s.ctx, s.exampleMealPlanID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanEventList, actual.Data)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetMealPlanEvents(s.ctx, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetMealPlanEvents(s.ctx, s.exampleMealPlanID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleMealPlanID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetMealPlanEvents(s.ctx, s.exampleMealPlanID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanEventsTestSuite) TestClient_CreateMealPlanEvent() {
	const expectedPath = "/api/v1/meal_plans/%s/events"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeMealPlanEventCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleMealPlanID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanEventResponse)

		actual, err := c.CreateMealPlanEvent(s.ctx, s.exampleMealPlanID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleMealPlanEvent, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateMealPlanEvent(s.ctx, s.exampleMealPlanID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.MealPlanEventCreationRequestInput{}

		actual, err := c.CreateMealPlanEvent(s.ctx, s.exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanEventToMealPlanEventCreationRequestInput(s.exampleMealPlanEvent)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateMealPlanEvent(s.ctx, s.exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertMealPlanEventToMealPlanEventCreationRequestInput(s.exampleMealPlanEvent)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateMealPlanEvent(s.ctx, s.exampleMealPlanID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *mealPlanEventsTestSuite) TestClient_UpdateMealPlanEvent() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanEventResponse)

		err := c.UpdateMealPlanEvent(s.ctx, s.exampleMealPlanEvent)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateMealPlanEvent(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateMealPlanEvent(s.ctx, s.exampleMealPlanEvent)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateMealPlanEvent(s.ctx, s.exampleMealPlanEvent)
		assert.Error(t, err)
	})
}

func (s *mealPlanEventsTestSuite) TestClient_ArchiveMealPlanEvent() {
	const expectedPathFormat = "/api/v1/meal_plans/%s/events/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleMealPlanEventResponse)

		err := c.ArchiveMealPlanEvent(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid meal plan ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanEvent(s.ctx, "", s.exampleMealPlanEvent.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid meal plan ClientOption ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveMealPlanEvent(s.ctx, s.exampleMealPlanID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveMealPlanEvent(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveMealPlanEvent(s.ctx, s.exampleMealPlanID, s.exampleMealPlanEvent.ID)
		assert.Error(t, err)
	})
}
