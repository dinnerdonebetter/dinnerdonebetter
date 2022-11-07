package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func buildMockRowsFromWebhooks(includeCounts bool, filteredCount uint64, webhooks ...*types.Webhook) *sqlmock.Rows {
	columns := webhooksTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range webhooks {
		for _, e := range w.Events {
			rowValues := []driver.Value{
				w.ID,
				w.Name,
				w.ContentType,
				w.URL,
				w.Method,
				e.ID,
				e.TriggerEvent,
				e.BelongsToWebhook,
				e.CreatedAt,
				e.ArchivedAt,
				w.CreatedAt,
				w.LastUpdatedAt,
				w.ArchivedAt,
				w.BelongsToHousehold,
			}

			if includeCounts {
				rowValues = append(rowValues, filteredCount, len(webhooks))
			}

			exampleRows.AddRow(rowValues...)
		}
	}

	return exampleRows
}

func buildErroneousMockRowsFromWebhooks(includeCounts bool, filteredCount uint64, webhooks ...*types.Webhook) *sqlmock.Rows {
	columns := webhooksTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, w := range webhooks {
		for _, e := range w.Events {
			rowValues := []driver.Value{
				w.ArchivedAt,
				e.ID,
				e.TriggerEvent,
				e.BelongsToWebhook,
				e.CreatedAt,
				e.ArchivedAt,
				w.ID,
				w.Name,
				w.ContentType,
				w.URL,
				w.Method,
				w.CreatedAt,
				w.LastUpdatedAt,
				w.BelongsToHousehold,
			}

			if includeCounts {
				rowValues = append(rowValues, filteredCount, len(webhooks))
			}

			exampleRows.AddRow(rowValues...)
		}
	}

	return exampleRows
}

func TestQuerier_WebhookExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhookID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdID,
			exampleWebhookID,
		}

		db.ExpectQuery(formatQueryForSQLMock(webhookExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.WebhookExists(ctx, exampleWebhookID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
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

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdID,
			exampleWebhookID,
		}

		db.ExpectQuery(formatQueryForSQLMock(webhookExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.WebhookExists(ctx, exampleWebhookID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhookID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleHouseholdID,
			exampleWebhookID,
		}

		db.ExpectQuery(formatQueryForSQLMock(webhookExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.WebhookExists(ctx, exampleWebhookID, exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetWebhook(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleWebhook := fakes.BuildFakeWebhook()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{exampleHouseholdID, exampleWebhook.ID}

		db.ExpectQuery(formatQueryForSQLMock(getWebhookQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromWebhooks(false, 0, exampleWebhook))

		actual, err := c.GetWebhook(ctx, exampleWebhook.ID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.Equal(t, exampleWebhook, actual)

		mock.AssertExpectationsForObjects(t, db)
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
		c, db := buildTestClient(t)

		args := []interface{}{exampleHouseholdID, exampleWebhook.ID}

		db.ExpectQuery(formatQueryForSQLMock(getWebhookQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRowsFromWebhooks(false, 0, exampleWebhook))

		actual, err := c.GetWebhook(ctx, exampleWebhook.ID, exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
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
