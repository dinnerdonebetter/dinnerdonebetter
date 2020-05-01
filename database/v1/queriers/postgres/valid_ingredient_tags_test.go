package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromValidIngredientTag(validIngredientTags ...*models.ValidIngredientTag) *sqlmock.Rows {
	includeCount := len(validIngredientTags) > 1
	columns := validIngredientTagsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validIngredientTags {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
		}

		if includeCount {
			rowValues = append(rowValues, len(validIngredientTags))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromValidIngredientTag(x *models.ValidIngredientTag) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(validIngredientTagsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.CreatedOn,
		x.UpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanValidIngredientTags(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanValidIngredientTags(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanValidIngredientTags(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildValidIngredientTagExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		expectedQuery := "SELECT EXISTS ( SELECT valid_ingredient_tags.id FROM valid_ingredient_tags WHERE valid_ingredient_tags.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidIngredientTag.ID,
		}
		actualQuery, actualArgs := p.buildValidIngredientTagExistsQuery(exampleValidIngredientTag.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ValidIngredientTagExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT valid_ingredient_tags.id FROM valid_ingredient_tags WHERE valid_ingredient_tags.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ValidIngredientTagExists(ctx, exampleValidIngredientTag.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ValidIngredientTagExists(ctx, exampleValidIngredientTag.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientTagQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		expectedQuery := "SELECT valid_ingredient_tags.id, valid_ingredient_tags.name, valid_ingredient_tags.created_on, valid_ingredient_tags.updated_on, valid_ingredient_tags.archived_on FROM valid_ingredient_tags WHERE valid_ingredient_tags.id = $1"
		expectedArgs := []interface{}{
			exampleValidIngredientTag.ID,
		}
		actualQuery, actualArgs := p.buildGetValidIngredientTagQuery(exampleValidIngredientTag.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredientTag(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_ingredient_tags.id, valid_ingredient_tags.name, valid_ingredient_tags.created_on, valid_ingredient_tags.updated_on, valid_ingredient_tags.archived_on FROM valid_ingredient_tags WHERE valid_ingredient_tags.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.ID,
			).
			WillReturnRows(buildMockRowsFromValidIngredientTag(exampleValidIngredientTag))

		actual, err := p.GetValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientTag, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllValidIngredientTagsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(valid_ingredient_tags.id) FROM valid_ingredient_tags WHERE valid_ingredient_tags.archived_on IS NULL"
		actualQuery := p.buildGetAllValidIngredientTagsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllValidIngredientTagsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_ingredient_tags.id) FROM valid_ingredient_tags WHERE valid_ingredient_tags.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllValidIngredientTagsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientTagsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_ingredient_tags.id, valid_ingredient_tags.name, valid_ingredient_tags.created_on, valid_ingredient_tags.updated_on, valid_ingredient_tags.archived_on, COUNT(valid_ingredient_tags.id) FROM valid_ingredient_tags WHERE valid_ingredient_tags.archived_on IS NULL AND valid_ingredient_tags.created_on > $1 AND valid_ingredient_tags.created_on < $2 AND valid_ingredient_tags.updated_on > $3 AND valid_ingredient_tags.updated_on < $4 GROUP BY valid_ingredient_tags.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetValidIngredientTagsQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredientTags(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT valid_ingredient_tags.id, valid_ingredient_tags.name, valid_ingredient_tags.created_on, valid_ingredient_tags.updated_on, valid_ingredient_tags.archived_on, COUNT(valid_ingredient_tags.id) FROM valid_ingredient_tags WHERE valid_ingredient_tags.archived_on IS NULL GROUP BY valid_ingredient_tags.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidIngredientTagList := fakemodels.BuildFakeValidIngredientTagList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(
				buildMockRowsFromValidIngredientTag(
					&exampleValidIngredientTagList.ValidIngredientTags[0],
					&exampleValidIngredientTagList.ValidIngredientTags[1],
					&exampleValidIngredientTagList.ValidIngredientTags[2],
				),
			)

		actual, err := p.GetValidIngredientTags(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientTagList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredientTags(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetValidIngredientTags(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid ingredient tag", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildErroneousMockRowFromValidIngredientTag(exampleValidIngredientTag))

		actual, err := p.GetValidIngredientTags(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateValidIngredientTagQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		expectedQuery := "INSERT INTO valid_ingredient_tags (name) VALUES ($1) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleValidIngredientTag.Name,
		}
		actualQuery, actualArgs := p.buildCreateValidIngredientTagQuery(exampleValidIngredientTag)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateValidIngredientTag(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO valid_ingredient_tags (name) VALUES ($1) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleValidIngredientTag.ID, exampleValidIngredientTag.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidIngredientTag.Name,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateValidIngredientTag(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientTag, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
		exampleInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidIngredientTag.Name,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateValidIngredientTag(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateValidIngredientTagQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		expectedQuery := "UPDATE valid_ingredient_tags SET name = $1, updated_on = extract(epoch FROM NOW()) WHERE id = $2 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleValidIngredientTag.Name,
			exampleValidIngredientTag.ID,
		}
		actualQuery, actualArgs := p.buildUpdateValidIngredientTagQuery(exampleValidIngredientTag)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateValidIngredientTag(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_ingredient_tags SET name = $1, updated_on = extract(epoch FROM NOW()) WHERE id = $2 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.Name,
				exampleValidIngredientTag.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateValidIngredientTag(ctx, exampleValidIngredientTag)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.Name,
				exampleValidIngredientTag.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateValidIngredientTag(ctx, exampleValidIngredientTag)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveValidIngredientTagQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		expectedQuery := "UPDATE valid_ingredient_tags SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleValidIngredientTag.ID,
		}
		actualQuery, actualArgs := p.buildArchiveValidIngredientTagQuery(exampleValidIngredientTag.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveValidIngredientTag(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_ingredient_tags SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientTag.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveValidIngredientTag(ctx, exampleValidIngredientTag.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
