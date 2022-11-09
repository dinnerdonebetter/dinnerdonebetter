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

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func buildMockRowsFromValidInstruments(includeCounts bool, filteredCount uint64, validInstruments ...*types.ValidInstrument) *sqlmock.Rows {
	columns := validInstrumentsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validInstruments {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.PluralName,
			x.Description,
			x.IconPath,
			x.UsableForStorage,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validInstruments))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidInstruments(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidInstruments(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidInstrumentExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidInstruments(false, 0, exampleValidInstrument))

		actual, err := c.GetValidInstrument(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidInstrument(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidInstrument(ctx, exampleValidInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRandomValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidInstruments(false, 0, exampleValidInstrument))

		actual, err := c.GetRandomValidInstrument(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRandomValidInstrument(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidInstruments(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstruments := fakes.BuildFakeValidInstrumentList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidInstruments(false, 0, exampleValidInstruments.ValidInstruments...))

		actual, err := c.SearchForValidInstruments(ctx, exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstruments.ValidInstruments, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidInstruments(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidInstruments(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidInstruments(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidInstrumentsForPreparation(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstruments := fakes.BuildFakeValidInstrumentList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidInstruments(false, 0, exampleValidInstruments.ValidInstruments...))

		actual, err := c.SearchForValidInstrumentsForPreparation(ctx, "", exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstruments.ValidInstruments, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidInstrumentsForPreparation(ctx, "", "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidInstrumentsForPreparation(ctx, "", exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validInstrumentSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidInstrumentsForPreparation(ctx, "", exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_instruments", nil, nil, nil, householdOwnershipColumn, validInstrumentsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidInstruments(true, exampleValidInstrumentList.FilteredCount, exampleValidInstrumentList.ValidInstruments...))

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentList()
		exampleValidInstrumentList.Page = 0
		exampleValidInstrumentList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_instruments", nil, nil, nil, householdOwnershipColumn, validInstrumentsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidInstruments(true, exampleValidInstrumentList.FilteredCount, exampleValidInstrumentList.ValidInstruments...))

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_instruments", nil, nil, nil, householdOwnershipColumn, validInstrumentsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_instruments", nil, nil, nil, householdOwnershipColumn, validInstrumentsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		exampleValidInstrument.ID = "1"
		exampleInput := converters.ConvertValidInstrumentToValidInstrumentDatabaseCreationInput(exampleValidInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.PluralName,
			exampleInput.Description,
			exampleInput.IconPath,
			exampleInput.UsableForStorage,
		}

		db.ExpectExec(formatQueryForSQLMock(validInstrumentCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidInstrument.CreatedAt
		}

		actual, err := c.CreateValidInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		exampleInput := converters.ConvertValidInstrumentToValidInstrumentDatabaseCreationInput(exampleValidInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.PluralName,
			exampleInput.Description,
			exampleInput.IconPath,
			exampleInput.UsableForStorage,
		}

		db.ExpectExec(formatQueryForSQLMock(validInstrumentCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidInstrument.CreatedAt
		}

		actual, err := c.CreateValidInstrument(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidInstrument.Name,
			exampleValidInstrument.PluralName,
			exampleValidInstrument.Description,
			exampleValidInstrument.IconPath,
			exampleValidInstrument.UsableForStorage,
			exampleValidInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidInstrument(ctx, exampleValidInstrument))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidInstrument(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidInstrument.Name,
			exampleValidInstrument.PluralName,
			exampleValidInstrument.Description,
			exampleValidInstrument.IconPath,
			exampleValidInstrument.UsableForStorage,
			exampleValidInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidInstrument(ctx, exampleValidInstrument))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidInstrument(ctx, exampleValidInstrument.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidInstrument(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidInstrument(ctx, exampleValidInstrument.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
