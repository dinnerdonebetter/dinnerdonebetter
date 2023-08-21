package postgres

import (
	"context"
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

func buildMockRowsFromHouseholdInstrumentOwnership(includeCounts bool, filteredCount uint64, householdInstrumentOwnerships ...*types.HouseholdInstrumentOwnership) *sqlmock.Rows {
	columns := householdInstrumentOwnershipsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range householdInstrumentOwnerships {
		rowValues := []driver.Value{
			&x.ID,
			&x.Notes,
			&x.Quantity,
			&x.Instrument.ID,
			&x.Instrument.Name,
			&x.Instrument.PluralName,
			&x.Instrument.Description,
			&x.Instrument.IconPath,
			&x.Instrument.UsableForStorage,
			&x.Instrument.DisplayInSummaryLists,
			&x.Instrument.IncludeInGeneratedInstructions,
			&x.Instrument.Slug,
			&x.Instrument.CreatedAt,
			&x.Instrument.LastUpdatedAt,
			&x.Instrument.ArchivedAt,
			&x.BelongsToHousehold,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(householdInstrumentOwnerships))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanHouseholdInstrumentOwnerships(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanHouseholdInstrumentOwnerships(ctx, mockRows, false)
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

		_, _, _, err := q.scanHouseholdInstrumentOwnerships(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_HouseholdInstrumentOwnershipExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.HouseholdInstrumentOwnershipExists(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()
		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleHouseholdInstrumentOwnership.ID,
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInstrumentOwnershipQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInstrumentOwnership(false, 0, exampleHouseholdInstrumentOwnership))

		actual, err := c.GetHouseholdInstrumentOwnership(ctx, exampleHouseholdInstrumentOwnership.ID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInstrumentOwnership, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInstrumentOwnership(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()
		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleHouseholdInstrumentOwnership.ID,
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInstrumentOwnershipQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdInstrumentOwnership(ctx, exampleHouseholdInstrumentOwnership.ID, exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdInstrumentOwnerships(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()
		exampleHouseholdInstrumentOwnershipList := fakes.BuildFakeHouseholdInstrumentOwnershipList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "household_instrument_ownerships", []string{"valid_instruments ON valid_instruments.id = household_instrument_ownerships.valid_instrument_id"}, []string{"household_instrument_ownerships.id", "valid_instruments.id"}, nil, householdOwnershipColumn, householdInstrumentOwnershipsTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInstrumentOwnership(true, exampleHouseholdInstrumentOwnershipList.FilteredCount, exampleHouseholdInstrumentOwnershipList.Data...))

		actual, err := c.GetHouseholdInstrumentOwnerships(ctx, exampleHouseholdID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInstrumentOwnershipList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "household_instrument_ownerships", []string{"valid_instruments ON valid_instruments.id = household_instrument_ownerships.valid_instrument_id"}, []string{"household_instrument_ownerships.id", "valid_instruments.id"}, nil, householdOwnershipColumn, householdInstrumentOwnershipsTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdInstrumentOwnerships(ctx, exampleHouseholdID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "household_instrument_ownerships", []string{"valid_instruments ON valid_instruments.id = household_instrument_ownerships.valid_instrument_id"}, []string{"household_instrument_ownerships.id", "valid_instruments.id"}, nil, householdOwnershipColumn, householdInstrumentOwnershipsTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholdInstrumentOwnerships(ctx, exampleHouseholdID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()
		exampleHouseholdInstrumentOwnership.ID = "1"
		exampleHouseholdInstrumentOwnership.Instrument = types.ValidInstrument{ID: exampleHouseholdInstrumentOwnership.Instrument.ID}
		exampleInput := converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipDatabaseCreationInput(exampleHouseholdInstrumentOwnership)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.Quantity,
			exampleInput.ValidInstrumentID,
			exampleInput.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(householdInstrumentOwnershipCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleHouseholdInstrumentOwnership.CreatedAt
		}

		actual, err := c.CreateHouseholdInstrumentOwnership(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdInstrumentOwnership, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHouseholdInstrumentOwnership(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()
		exampleInput := converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipDatabaseCreationInput(exampleHouseholdInstrumentOwnership)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.Quantity,
			exampleInput.ValidInstrumentID,
			exampleInput.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(householdInstrumentOwnershipCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleHouseholdInstrumentOwnership.CreatedAt
		}

		actual, err := c.CreateHouseholdInstrumentOwnership(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleHouseholdInstrumentOwnership.Notes,
			exampleHouseholdInstrumentOwnership.Quantity,
			exampleHouseholdInstrumentOwnership.Instrument.ID,
			exampleHouseholdInstrumentOwnership.ID,
			exampleHouseholdInstrumentOwnership.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(updateHouseholdInstrumentOwnershipQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateHouseholdInstrumentOwnership(ctx, exampleHouseholdInstrumentOwnership))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateHouseholdInstrumentOwnership(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleHouseholdInstrumentOwnership.Notes,
			exampleHouseholdInstrumentOwnership.Quantity,
			exampleHouseholdInstrumentOwnership.Instrument.ID,
			exampleHouseholdInstrumentOwnership.ID,
			exampleHouseholdInstrumentOwnership.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(updateHouseholdInstrumentOwnershipQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateHouseholdInstrumentOwnership(ctx, exampleHouseholdInstrumentOwnership))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHouseholdInstrumentOwnership(ctx, "", exampleHouseholdID))
	})
}
