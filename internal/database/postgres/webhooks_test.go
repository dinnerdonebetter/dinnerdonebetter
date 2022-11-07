package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestQuerier_WebhookExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhookID := fakes.BuildFakeID()

		c, _, mockGQ := buildNewTestClient(t)

		mockGQ.On(
			"WebhookExists",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.WebhookExistsParams{
				ID:                 exampleWebhookID,
				BelongsToHousehold: exampleHouseholdID,
			}).Return(true, nil)

		actual, err := c.WebhookExists(ctx, exampleWebhookID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.True(t, actual)
	})

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.WebhookExists(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleWebhookID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		actual, err := c.WebhookExists(ctx, exampleWebhookID, "")
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhookID := fakes.BuildFakeID()

		c, _, mockGQ := buildNewTestClient(t)

		mockGQ.On(
			"WebhookExists",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.WebhookExistsParams{
				ID:                 exampleWebhookID,
				BelongsToHousehold: exampleHouseholdID,
			}).Return(false, sql.ErrNoRows)

		actual, err := c.WebhookExists(ctx, exampleWebhookID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.False(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhookID := fakes.BuildFakeID()

		c, _, mockGQ := buildNewTestClient(t)

		mockGQ.On(
			"WebhookExists",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.WebhookExistsParams{
				ID:                 exampleWebhookID,
				BelongsToHousehold: exampleHouseholdID,
			}).Return(false, errors.New("blah"))

		actual, err := c.WebhookExists(ctx, exampleWebhookID, exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhook := fakes.BuildFakeWebhook()

		ctx := context.Background()

		c, _, mockGQ := buildNewTestClient(t)

		generatedResponse := []*generated.GetWebhookRow{}
		for i := range exampleWebhook.Events {
			wr := &generated.GetWebhookRow{
				ID:                 exampleWebhook.ID,
				Name:               exampleWebhook.Name,
				ContentType:        exampleWebhook.ContentType,
				Url:                exampleWebhook.URL,
				Method:             exampleWebhook.Method,
				CreatedAt:          exampleWebhook.CreatedAt,
				BelongsToHousehold: exampleWebhook.BelongsToHousehold,
				ID_2:               exampleWebhook.Events[i].ID,
				CreatedAt_2:        exampleWebhook.Events[i].CreatedAt,
				BelongsToWebhook:   exampleWebhook.Events[i].BelongsToWebhook,
				TriggerEvent:       generated.WebhookEvent(exampleWebhook.Events[i].TriggerEvent),
			}

			if exampleWebhook.LastUpdatedAt != nil {
				wr.LastUpdatedAt = sql.NullTime{Time: *exampleWebhook.LastUpdatedAt}
			}
			if exampleWebhook.ArchivedAt != nil {
				wr.ArchivedAt = sql.NullTime{Time: *exampleWebhook.ArchivedAt}
			}

			if exampleWebhook.Events[i].ArchivedAt != nil {
				wr.ArchivedAt_2 = sql.NullTime{Time: *exampleWebhook.Events[i].ArchivedAt}
			}

			generatedResponse = append(generatedResponse, wr)
		}

		mockGQ.On(
			"GetWebhook",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.GetWebhookParams{
				ID:                 exampleWebhook.ID,
				BelongsToHousehold: exampleHouseholdID,
			}).Return(generatedResponse, nil)

		actual, err := c.GetWebhook(ctx, exampleWebhook.ID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhook, actual)
	})

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetWebhook(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleWebhook := fakes.BuildFakeWebhook()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetWebhook(ctx, exampleWebhook.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhook := fakes.BuildFakeWebhook()

		ctx := context.Background()
		c, _, mockGQ := buildNewTestClient(t)

		mockGQ.On(
			"GetWebhook",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.GetWebhookParams{
				ID:                 exampleWebhook.ID,
				BelongsToHousehold: exampleHouseholdID,
			}).Return([]*generated.GetWebhookRow(nil), errors.New("blah"))

		actual, err := c.GetWebhook(ctx, exampleWebhook.ID, exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetWebhooks(T *testing.T) {
	T.Parallel()

	exampleHouseholdID := fakes.BuildFakeID()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleWebhookList := fakes.BuildFakeWebhookList()
		exampleWebhookList.Pagination = types.Pagination{}
		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _, mockGQ := buildNewTestClient(t)

		filterArgs := filter.ToDatabaseArgs()

		generatedResponse := []*generated.GetWebhooksRow{}
		for i := range exampleWebhookList.Webhooks {
			for j := range exampleWebhookList.Webhooks[i].Events {
				wr := &generated.GetWebhooksRow{
					ID:                 exampleWebhookList.Webhooks[i].ID,
					Name:               exampleWebhookList.Webhooks[i].Name,
					ContentType:        exampleWebhookList.Webhooks[i].ContentType,
					Url:                exampleWebhookList.Webhooks[i].URL,
					Method:             exampleWebhookList.Webhooks[i].Method,
					CreatedAt:          exampleWebhookList.Webhooks[i].CreatedAt,
					BelongsToHousehold: exampleWebhookList.Webhooks[i].BelongsToHousehold,
					ID_2:               exampleWebhookList.Webhooks[i].Events[j].ID,
					CreatedAt_2:        exampleWebhookList.Webhooks[i].Events[j].CreatedAt,
					BelongsToWebhook:   exampleWebhookList.Webhooks[i].Events[j].BelongsToWebhook,
					TriggerEvent:       generated.WebhookEvent(exampleWebhookList.Webhooks[i].Events[j].TriggerEvent),
					FilteredCount:      0,
					TotalCount:         0,
				}

				if exampleWebhookList.Webhooks[i].LastUpdatedAt != nil {
					wr.LastUpdatedAt = sql.NullTime{Time: *exampleWebhookList.Webhooks[i].LastUpdatedAt}
				}
				if exampleWebhookList.Webhooks[i].ArchivedAt != nil {
					wr.ArchivedAt = sql.NullTime{Time: *exampleWebhookList.Webhooks[i].ArchivedAt}
				}
				if exampleWebhookList.Webhooks[i].Events[j].ArchivedAt != nil {
					wr.ArchivedAt_2 = sql.NullTime{Time: *exampleWebhookList.Webhooks[i].Events[j].ArchivedAt}
				}

				generatedResponse = append(generatedResponse, wr)
			}
		}

		mockGQ.On(
			"GetWebhooks",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.GetWebhooksParams{
				BelongsToHousehold: exampleHouseholdID,
				CreatedAfter:       filterArgs.CreatedAfter,
				CreatedBefore:      filterArgs.CreatedBefore,
				UpdatedAfter:       filterArgs.UpdatedAfter,
				UpdatedBefore:      filterArgs.UpdatedBefore,
			}).Return(generatedResponse, nil)

		actual, err := c.GetWebhooks(ctx, exampleHouseholdID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhookList, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		exampleWebhookList := fakes.BuildFakeWebhookList()
		exampleWebhookList.Pagination = types.Pagination{}
		filter := (*types.QueryFilter)(nil)

		ctx := context.Background()
		c, _, mockGQ := buildNewTestClient(t)

		filterArgs := types.DefaultQueryFilter().ToDatabaseArgs()

		generatedResponse := []*generated.GetWebhooksRow{}
		for i := range exampleWebhookList.Webhooks {
			for j := range exampleWebhookList.Webhooks[i].Events {
				wr := &generated.GetWebhooksRow{
					ID:                 exampleWebhookList.Webhooks[i].ID,
					Name:               exampleWebhookList.Webhooks[i].Name,
					ContentType:        exampleWebhookList.Webhooks[i].ContentType,
					Url:                exampleWebhookList.Webhooks[i].URL,
					Method:             exampleWebhookList.Webhooks[i].Method,
					CreatedAt:          exampleWebhookList.Webhooks[i].CreatedAt,
					BelongsToHousehold: exampleWebhookList.Webhooks[i].BelongsToHousehold,
					ID_2:               exampleWebhookList.Webhooks[i].Events[j].ID,
					CreatedAt_2:        exampleWebhookList.Webhooks[i].Events[j].CreatedAt,
					BelongsToWebhook:   exampleWebhookList.Webhooks[i].Events[j].BelongsToWebhook,
					TriggerEvent:       generated.WebhookEvent(exampleWebhookList.Webhooks[i].Events[j].TriggerEvent),
					FilteredCount:      0,
					TotalCount:         0,
				}

				if exampleWebhookList.Webhooks[i].LastUpdatedAt != nil {
					wr.LastUpdatedAt = sql.NullTime{Time: *exampleWebhookList.Webhooks[i].LastUpdatedAt}
				}
				if exampleWebhookList.Webhooks[i].ArchivedAt != nil {
					wr.ArchivedAt = sql.NullTime{Time: *exampleWebhookList.Webhooks[i].ArchivedAt}
				}
				if exampleWebhookList.Webhooks[i].Events[j].ArchivedAt != nil {
					wr.ArchivedAt_2 = sql.NullTime{Time: *exampleWebhookList.Webhooks[i].Events[j].ArchivedAt}
				}

				generatedResponse = append(generatedResponse, wr)
			}
		}

		mockGQ.On(
			"GetWebhooks",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.GetWebhooksParams{
				BelongsToHousehold: exampleHouseholdID,
				CreatedAfter:       filterArgs.CreatedAfter,
				CreatedBefore:      filterArgs.CreatedBefore,
				UpdatedAfter:       filterArgs.UpdatedAfter,
				UpdatedBefore:      filterArgs.UpdatedBefore,
			}).Return(generatedResponse, nil)

		actual, err := c.GetWebhooks(ctx, exampleHouseholdID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhookList, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		c, _ := buildTestClient(t)

		actual, err := c.GetWebhooks(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _, mockGQ := buildNewTestClient(t)

		filter := types.DefaultQueryFilter()
		filterArgs := filter.ToDatabaseArgs()

		mockGQ.On(
			"GetWebhooks",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.GetWebhooksParams{
				BelongsToHousehold: exampleHouseholdID,
				CreatedAfter:       filterArgs.CreatedAfter,
				CreatedBefore:      filterArgs.CreatedBefore,
				UpdatedAfter:       filterArgs.UpdatedAfter,
				UpdatedBefore:      filterArgs.UpdatedBefore,
			}).Return([]*generated.GetWebhooksRow(nil), errors.New("blah"))

		actual, err := c.GetWebhooks(ctx, exampleHouseholdID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateWebhook(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleWebhook := fakes.BuildFakeWebhook()
		for i := range exampleWebhook.Events {
			exampleWebhook.Events[i].CreatedAt = exampleWebhook.CreatedAt
		}
		exampleInput := converters.ConvertWebhookToWebhookDatabaseCreationInput(exampleWebhook)

		ctx := context.Background()
		c, db, _ := buildNewTestClient(t)
		mockGQ := &mockGeneratedQuerier{}

		db.ExpectBegin()

		mockGQ.On(
			"CreateWebhook",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.CreateWebhookParams{
				ID:                 exampleInput.ID,
				Name:               exampleInput.Name,
				ContentType:        exampleInput.ContentType,
				Url:                exampleInput.URL,
				Method:             exampleInput.Method,
				BelongsToHousehold: exampleInput.BelongsToHousehold,
			}).Return(nil)

		for _, evt := range exampleInput.Events {
			mockGQ.On("CreateWebhookTriggerEvent",
				testutils.ContextMatcher,
				database.SQLQueryExecutorMatcher,
				&generated.CreateWebhookTriggerEventParams{
					ID:               evt.ID,
					TriggerEvent:     generated.WebhookEvent(evt.TriggerEvent),
					BelongsToWebhook: evt.BelongsToWebhook,
				}).Return(nil)
		}

		c.generatedQuerier = mockGQ
		db.ExpectCommit()

		c.timeFunc = func() time.Time {
			return exampleWebhook.CreatedAt
		}

		actual, err := c.CreateWebhook(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhook, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateWebhook(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing creation query", func(t *testing.T) {
		t.Parallel()

		exampleWebhook := fakes.BuildFakeWebhook()
		for i := range exampleWebhook.Events {
			exampleWebhook.Events[i].CreatedAt = exampleWebhook.CreatedAt
		}
		exampleInput := converters.ConvertWebhookToWebhookDatabaseCreationInput(exampleWebhook)

		ctx := context.Background()
		c, db, mockGQ := buildNewTestClient(t)

		db.ExpectBegin()

		mockGQ.On(
			"CreateWebhook",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.CreateWebhookParams{
				ID:                 exampleInput.ID,
				Name:               exampleInput.Name,
				ContentType:        exampleInput.ContentType,
				Url:                exampleInput.URL,
				Method:             exampleInput.Method,
				BelongsToHousehold: exampleInput.BelongsToHousehold,
			}).Return(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() time.Time {
			return exampleWebhook.CreatedAt
		}

		actual, err := c.CreateWebhook(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveWebhook(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhookID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _, mockGQ := buildNewTestClient(t)

		mockGQ.On("ArchiveWebhook",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.ArchiveWebhookParams{
				ID:                 exampleWebhookID,
				BelongsToHousehold: exampleHouseholdID,
			}).Return(nil)

		actual := c.ArchiveWebhook(ctx, exampleWebhookID, exampleHouseholdID)
		assert.NoError(t, actual)

		mock.AssertExpectationsForObjects(t, mockGQ)
	})

	T.Run("with invalid webhook ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveWebhook(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleWebhookID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveWebhook(ctx, exampleWebhookID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhookID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _, mockGQ := buildNewTestClient(t)

		mockGQ.On("ArchiveWebhook",
			testutils.ContextMatcher,
			database.SQLQueryExecutorMatcher,
			&generated.ArchiveWebhookParams{
				ID:                 exampleWebhookID,
				BelongsToHousehold: exampleHouseholdID,
			}).Return(errors.New("blah"))

		assert.Error(t, c.ArchiveWebhook(ctx, exampleWebhookID, exampleHouseholdID))

		mock.AssertExpectationsForObjects(t, mockGQ)
	})
}
