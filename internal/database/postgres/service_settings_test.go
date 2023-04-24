package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromServiceSettings(includeCounts bool, filteredCount uint64, serviceSettings ...*types.ServiceSetting) *sqlmock.Rows {
	columns := serviceSettingsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range serviceSettings {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Type,
			x.Description,
			x.DefaultValue,
			x.AdminsOnly,
			strings.Join(x.Enumeration, serviceSettingsEnumDelimiter),
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(serviceSettings))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanServiceSettings(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanServiceSettings(ctx, mockRows, false)
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

		_, _, _, err := q.scanServiceSettings(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ServiceSettingExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)
		args := []any{
			exampleServiceSetting.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ServiceSettingExists(ctx, exampleServiceSetting.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ServiceSettingExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)
		args := []any{
			exampleServiceSetting.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ServiceSettingExists(ctx, exampleServiceSetting.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)
		args := []any{
			exampleServiceSetting.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ServiceSettingExists(ctx, exampleServiceSetting.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetServiceSetting(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSetting.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettings(false, 0, exampleServiceSetting))

		actual, err := c.GetServiceSetting(ctx, exampleServiceSetting.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSetting, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetServiceSetting(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSetting.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetServiceSetting(ctx, exampleServiceSetting.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForServiceSettings(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettings := fakes.BuildFakeServiceSettingList()

		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettings(false, 0, exampleServiceSettings.Data...))

		actual, err := c.SearchForServiceSettings(ctx, exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettings.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForServiceSettings(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForServiceSettings(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForServiceSettings(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetServiceSettings(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleServiceSettingList := fakes.BuildFakeServiceSettingList()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "service_settings", nil, nil, nil, householdOwnershipColumn, serviceSettingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettings(true, exampleServiceSettingList.FilteredCount, exampleServiceSettingList.Data...))

		actual, err := c.GetServiceSettings(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := (*types.QueryFilter)(nil)
		exampleServiceSettingList := fakes.BuildFakeServiceSettingList()
		exampleServiceSettingList.Page = 0
		exampleServiceSettingList.Limit = 0

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "service_settings", nil, nil, nil, householdOwnershipColumn, serviceSettingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettings(true, exampleServiceSettingList.FilteredCount, exampleServiceSettingList.Data...))

		actual, err := c.GetServiceSettings(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "service_settings", nil, nil, nil, householdOwnershipColumn, serviceSettingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetServiceSettings(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()

		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "service_settings", nil, nil, nil, householdOwnershipColumn, serviceSettingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetServiceSettings(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateServiceSetting(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()
		exampleServiceSetting.ID = "1"
		exampleInput := converters.ConvertServiceSettingToServiceSettingDatabaseCreationInput(exampleServiceSetting)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Type,
			exampleInput.Description,
			exampleInput.DefaultValue,
			exampleInput.AdminsOnly,
			strings.Join(exampleInput.Enumeration, serviceSettingsEnumDelimiter),
		}

		db.ExpectExec(formatQueryForSQLMock(serviceSettingCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleServiceSetting.CreatedAt
		}

		actual, err := c.CreateServiceSetting(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSetting, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateServiceSetting(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New(t.Name())
		exampleServiceSetting := fakes.BuildFakeServiceSetting()
		exampleInput := converters.ConvertServiceSettingToServiceSettingDatabaseCreationInput(exampleServiceSetting)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Type,
			exampleInput.Description,
			exampleInput.DefaultValue,
			exampleInput.AdminsOnly,
			strings.Join(exampleInput.Enumeration, serviceSettingsEnumDelimiter),
		}

		db.ExpectExec(formatQueryForSQLMock(serviceSettingCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleServiceSetting.CreatedAt
		}

		actual, err := c.CreateServiceSetting(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateServiceSetting(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSetting.Name,
			exampleServiceSetting.Type,
			exampleServiceSetting.Description,
			exampleServiceSetting.DefaultValue,
			exampleServiceSetting.AdminsOnly,
			strings.Join(exampleServiceSetting.Enumeration, serviceSettingsEnumDelimiter),
			exampleServiceSetting.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateServiceSettingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateServiceSetting(ctx, exampleServiceSetting))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateServiceSetting(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSetting.Name,
			exampleServiceSetting.Type,
			exampleServiceSetting.Description,
			exampleServiceSetting.DefaultValue,
			exampleServiceSetting.AdminsOnly,
			strings.Join(exampleServiceSetting.Enumeration, serviceSettingsEnumDelimiter),
			exampleServiceSetting.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateServiceSettingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateServiceSetting(ctx, exampleServiceSetting))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveServiceSetting(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSetting.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveServiceSettingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveServiceSetting(ctx, exampleServiceSetting.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveServiceSetting(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSetting := fakes.BuildFakeServiceSetting()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSetting.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveServiceSettingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveServiceSetting(ctx, exampleServiceSetting.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
