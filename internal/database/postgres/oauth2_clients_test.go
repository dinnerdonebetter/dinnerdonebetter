package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromOAuth2Clients(includeCounts bool, filteredCount uint64, clients ...*types.OAuth2Client) *sqlmock.Rows {
	columns := oauth2ClientsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, c := range clients {
		rowValues := []driver.Value{
			c.ID,
			c.Name,
			c.ClientID,
			c.ClientSecret,
			c.CreatedAt,
			c.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(clients))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanOAuth2Clients(ctx, mockRows, false)
		assert.Error(t, err)
	})

	T.Run("logs row closing errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, _, err := q.scanOAuth2Clients(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_GetOAuth2ClientByClientID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2Client := fakes.BuildFakeOAuth2Client()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleOAuth2Client.ClientID}

		db.ExpectQuery(formatQueryForSQLMock(getOAuth2ClientByClientIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromOAuth2Clients(false, 0, exampleOAuth2Client))

		actual, err := c.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with empty client ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetOAuth2ClientByClientID(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("respects sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2Client := fakes.BuildFakeOAuth2Client()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleOAuth2Client.ClientID}

		db.ExpectQuery(formatQueryForSQLMock(getOAuth2ClientByClientIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2Client := fakes.BuildFakeOAuth2Client()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleOAuth2Client.ClientID}

		db.ExpectQuery(formatQueryForSQLMock(getOAuth2ClientByClientIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetOAuth2ClientByClientID(ctx, exampleOAuth2Client.ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetOAuth2ClientByDatabaseID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2Client := fakes.BuildFakeOAuth2Client()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleOAuth2Client.ID}

		db.ExpectQuery(formatQueryForSQLMock(getOAuth2ClientByDatabaseIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromOAuth2Clients(false, 0, exampleOAuth2Client))

		actual, err := c.GetOAuth2ClientByDatabaseID(ctx, exampleOAuth2Client.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid client ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetOAuth2ClientByDatabaseID(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("respects sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2ClientID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleOAuth2ClientID}

		db.ExpectQuery(formatQueryForSQLMock(getOAuth2ClientByDatabaseIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.GetOAuth2ClientByDatabaseID(ctx, exampleOAuth2ClientID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2ClientID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleOAuth2ClientID}

		db.ExpectQuery(formatQueryForSQLMock(getOAuth2ClientByDatabaseIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetOAuth2ClientByDatabaseID(ctx, exampleOAuth2ClientID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetOAuth2Clients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2ClientList := fakes.BuildFakeOAuth2ClientList()
		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "oauth2_clients", nil, nil, nil, userOwnershipColumn, oauth2ClientsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromOAuth2Clients(true, exampleOAuth2ClientList.FilteredCount, exampleOAuth2ClientList.Data...))

		actual, err := c.GetOAuth2Clients(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2ClientList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetOAuth2Clients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("respects sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "oauth2_clients", nil, nil, nil, userOwnershipColumn, oauth2ClientsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.GetOAuth2Clients(ctx, filter)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "oauth2_clients", nil, nil, nil, userOwnershipColumn, oauth2ClientsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetOAuth2Clients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "oauth2_clients", nil, nil, nil, userOwnershipColumn, oauth2ClientsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetOAuth2Clients(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2Client := fakes.BuildFakeOAuth2Client()
		exampleOAuth2Client.ClientSecret = ""
		exampleInput := converters.ConvertOAuth2ClientToOAuth2ClientDatabaseCreationInput(exampleOAuth2Client)
		exampleOAuth2Client.ID = exampleInput.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.ClientID,
			exampleInput.ClientSecret,
		}

		db.ExpectExec(formatQueryForSQLMock(createOAuth2ClientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleOAuth2Client.CreatedAt
		}

		actual, err := c.CreateOAuth2Client(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleOAuth2Client, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid creator ID", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2Client := fakes.BuildFakeOAuth2Client()
		exampleOAuth2Client.ClientSecret = ""
		exampleInput := converters.ConvertOAuth2ClientToOAuth2ClientDatabaseCreationInput(exampleOAuth2Client)

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateOAuth2Client(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateOAuth2Client(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleOAuth2Client := fakes.BuildFakeOAuth2Client()
		exampleInput := converters.ConvertOAuth2ClientToOAuth2ClientDatabaseCreationInput(exampleOAuth2Client)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.ClientID,
			exampleInput.ClientSecret,
		}

		db.ExpectExec(formatQueryForSQLMock(createOAuth2ClientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleOAuth2Client.CreatedAt
		}

		actual, err := c.CreateOAuth2Client(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveOAuth2Client(T *testing.T) {
	T.Parallel()

	T.Run("with invalid client ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveOAuth2Client(ctx, ""))
	})
}
