package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
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

func buildMockRowsFromServiceSettingConfigurations(includeCounts bool, filteredCount uint64, serviceSettingConfigurations ...*types.ServiceSettingConfiguration) *sqlmock.Rows {
	columns := serviceSettingConfigurationsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range serviceSettingConfigurations {
		rowValues := []driver.Value{
			x.ID,
			x.Value,
			x.Notes,
			x.ServiceSetting.ID,
			x.ServiceSetting.Name,
			x.ServiceSetting.Type,
			x.ServiceSetting.Description,
			x.ServiceSetting.DefaultValue,
			strings.Join(x.ServiceSetting.Enumeration, serviceSettingsEnumDelimiter),
			x.ServiceSetting.AdminsOnly,
			x.ServiceSetting.CreatedAt,
			x.ServiceSetting.LastUpdatedAt,
			x.ServiceSetting.ArchivedAt,
			x.BelongsToUser,
			x.BelongsToHousehold,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(serviceSettingConfigurations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanServiceSettingConfigurations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanServiceSettingConfigurations(ctx, mockRows, false)
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

		_, tc, fc, err := q.scanServiceSettingConfigurations(ctx, mockRows, true)
		assert.Error(t, err)
		assert.Zero(t, tc)
		assert.Zero(t, fc)
	})
}

func TestQuerier_ServiceSettingConfigurationExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)
		args := []any{
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingConfigurationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ServiceSettingConfigurationExists(ctx, exampleServiceSettingConfiguration.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ServiceSettingConfigurationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)
		args := []any{
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingConfigurationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ServiceSettingConfigurationExists(ctx, exampleServiceSettingConfiguration.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)
		args := []any{
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(serviceSettingConfigurationExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ServiceSettingConfigurationExists(ctx, exampleServiceSettingConfiguration.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetServiceSettingConfiguration(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettingConfigurations(false, 0, exampleServiceSettingConfiguration))

		actual, err := c.GetServiceSettingConfiguration(ctx, exampleServiceSettingConfiguration.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfiguration, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()
		c, _ := buildTestClient(t)

		actual, err := c.GetServiceSettingConfiguration(ctx, exampleServiceSettingConfiguration.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetServiceSettingConfiguration(ctx, exampleServiceSettingConfiguration.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetServiceSettingConfigurationForUserByName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.ServiceSetting.Name,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForUserByNameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettingConfigurations(false, 0, exampleServiceSettingConfiguration))

		actual, err := c.GetServiceSettingConfigurationForUserByName(ctx, exampleUserID, exampleServiceSettingConfiguration.ServiceSetting.Name)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfiguration, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.GetServiceSettingConfigurationForUserByName(ctx, exampleUserID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.ServiceSetting.Name,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForUserByNameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetServiceSettingConfigurationForUserByName(ctx, exampleUserID, exampleServiceSettingConfiguration.ServiceSetting.Name)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetServiceSettingConfigurationForHouseholdByName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.ServiceSetting.Name,
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForHouseholdByNameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettingConfigurations(false, 0, exampleServiceSettingConfiguration))

		actual, err := c.GetServiceSettingConfigurationForHouseholdByName(ctx, exampleHouseholdID, exampleServiceSettingConfiguration.ServiceSetting.Name)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfiguration, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting name", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.GetServiceSettingConfigurationForHouseholdByName(ctx, exampleHouseholdID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.ServiceSetting.Name,
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForHouseholdByNameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetServiceSettingConfigurationForHouseholdByName(ctx, exampleHouseholdID, exampleServiceSettingConfiguration.ServiceSetting.Name)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetServiceSettingConfigurationsForUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleServiceSettingConfigurationList := fakes.BuildFakeServiceSettingConfigurationList()
		exampleServiceSettingConfigurationList.TotalCount = 0
		exampleServiceSettingConfigurationList.FilteredCount = 0

		c, db := buildTestClient(t)

		args := []any{
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForUserQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettingConfigurations(false, 0, exampleServiceSettingConfigurationList.Data...))

		actual, err := c.GetServiceSettingConfigurationsForUser(ctx, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfigurationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.GetServiceSettingConfigurationsForUser(ctx, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		args := []any{
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForUserQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetServiceSettingConfigurationsForUser(ctx, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetServiceSettingConfigurationsForHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleServiceSettingConfigurationList := fakes.BuildFakeServiceSettingConfigurationList()
		exampleServiceSettingConfigurationList.FilteredCount = 0
		exampleServiceSettingConfigurationList.TotalCount = 0

		c, db := buildTestClient(t)

		args := []any{
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromServiceSettingConfigurations(false, exampleServiceSettingConfigurationList.FilteredCount, exampleServiceSettingConfigurationList.Data...))

		actual, err := c.GetServiceSettingConfigurationsForHousehold(ctx, exampleHouseholdID)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfigurationList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		args := []any{
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetServiceSettingConfigurationsForHousehold(ctx, exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		args := []any{
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getServiceSettingConfigurationForHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetServiceSettingConfigurationsForHousehold(ctx, exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateServiceSettingConfiguration(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()
		exampleServiceSettingConfiguration.ID = "1"
		exampleServiceSettingConfiguration.ServiceSetting = types.ServiceSetting{ID: exampleServiceSettingConfiguration.ServiceSetting.ID}
		exampleInput := converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationDatabaseCreationInput(exampleServiceSettingConfiguration)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Value,
			exampleInput.Notes,
			exampleInput.ServiceSettingID,
			exampleInput.BelongsToUser,
			exampleInput.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(serviceSettingConfigurationCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleServiceSettingConfiguration.CreatedAt
		}

		actual, err := c.CreateServiceSettingConfiguration(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleServiceSettingConfiguration, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateServiceSettingConfiguration(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New(t.Name())
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()
		exampleInput := converters.ConvertServiceSettingConfigurationToServiceSettingConfigurationDatabaseCreationInput(exampleServiceSettingConfiguration)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Value,
			exampleInput.Notes,
			exampleInput.ServiceSettingID,
			exampleInput.BelongsToUser,
			exampleInput.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(serviceSettingConfigurationCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleServiceSettingConfiguration.CreatedAt
		}

		actual, err := c.CreateServiceSettingConfiguration(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateServiceSettingConfiguration(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.Value,
			exampleServiceSettingConfiguration.Notes,
			exampleServiceSettingConfiguration.ServiceSetting.ID,
			exampleServiceSettingConfiguration.BelongsToUser,
			exampleServiceSettingConfiguration.BelongsToHousehold,
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateServiceSettingConfigurationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateServiceSettingConfiguration(ctx, exampleServiceSettingConfiguration))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateServiceSettingConfiguration(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.Value,
			exampleServiceSettingConfiguration.Notes,
			exampleServiceSettingConfiguration.ServiceSetting.ID,
			exampleServiceSettingConfiguration.BelongsToUser,
			exampleServiceSettingConfiguration.BelongsToHousehold,
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateServiceSettingConfigurationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateServiceSettingConfiguration(ctx, exampleServiceSettingConfiguration))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveServiceSettingConfiguration(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveServiceSettingConfigurationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveServiceSettingConfiguration(ctx, exampleServiceSettingConfiguration.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid service setting ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveServiceSettingConfiguration(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleServiceSettingConfiguration := fakes.BuildFakeServiceSettingConfiguration()

		c, db := buildTestClient(t)

		args := []any{
			exampleServiceSettingConfiguration.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveServiceSettingConfigurationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveServiceSettingConfiguration(ctx, exampleServiceSettingConfiguration.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
