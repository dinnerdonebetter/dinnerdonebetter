package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func buildMockRowsFromIDs(ids ...string) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows([]string{"id"})

	for _, x := range ids {
		exampleRows.AddRow(x)
	}

	return exampleRows
}

func TestSQLQuerier_logQueryBuildingError(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestClient(t)

		ctx := context.Background()
		_, span := tracing.StartSpan(ctx)
		err := errors.New(t.Name())

		q.logQueryBuildingError(span, err)
	})
}

func TestPostgres_buildListQuery(T *testing.T) {
	T.Parallel()

	const (
		exampleTableName       = "example_table"
		exampleOwnershipColumn = "belongs_to_household"
	)

	exampleColumns := []string{
		"column_ate",
		"column_two",
		"column_three",
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT column_ate, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_at IS NULL AND example_table.belongs_to_household = $1 AND key = $2 AND example_table.created_at > $3 AND example_table.created_at < $4 AND example_table.last_updated_at > $5 AND example_table.last_updated_at < $6) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_at IS NULL AND example_table.belongs_to_household = $7 AND key = $8) as total_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_at IS NULL AND example_table.belongs_to_household = $9 AND key = $10 AND example_table.created_at > $11 AND example_table.created_at < $12 AND example_table.last_updated_at > $13 AND example_table.last_updated_at < $14 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []any{
			exampleUser.ID,
			"value",
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			exampleUser.ID,
			"value",
			exampleUser.ID,
			"value",
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		exampleJoins := []string{
			"things on stuff.thing_id=things.id",
		}
		exampleWhere := squirrel.Eq{
			"key": "value",
		}

		actualQuery, actualArgs := q.buildListQuery(ctx, exampleTableName, exampleJoins, nil, exampleWhere, exampleOwnershipColumn, exampleColumns, exampleUser.ID, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})

	T.Run("for admin without archived", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestClient(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT column_ate, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.archived_at IS NULL AND example_table.created_at > $1 AND example_table.created_at < $2 AND example_table.last_updated_at > $3 AND example_table.last_updated_at < $4) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.archived_at IS NULL) as total_count FROM example_table WHERE example_table.created_at > $5 AND example_table.created_at < $6 AND example_table.last_updated_at > $7 AND example_table.last_updated_at < $8 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []any{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.buildListQuery(ctx, exampleTableName, nil, nil, nil, exampleOwnershipColumn, exampleColumns, exampleUser.ID, true, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})

	T.Run("for admin with archived", func(t *testing.T) {
		t.Parallel()

		q, _ := buildTestClient(t)
		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()
		filter.IncludeArchived = func(x bool) *bool { return &x }(true)

		expectedQuery := "SELECT column_ate, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.created_at > $1 AND example_table.created_at < $2 AND example_table.last_updated_at > $3 AND example_table.last_updated_at < $4) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table) as total_count FROM example_table WHERE example_table.created_at > $5 AND example_table.created_at < $6 AND example_table.last_updated_at > $7 AND example_table.last_updated_at < $8 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []any{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.buildListQuery(ctx, exampleTableName, nil, nil, nil, exampleOwnershipColumn, exampleColumns, exampleUser.ID, true, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_buildListQueryWithILike(T *testing.T) {
	T.Parallel()

	const (
		exampleTableName       = "example_table"
		exampleOwnershipColumn = "belongs_to_household"
	)

	exampleColumns := []string{
		"column_ate",
		"column_two",
		"column_three",
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT column_ate, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $1 AND example_table.archived_at IS NULL AND example_table.belongs_to_household = $2 AND example_table.created_at > $3 AND example_table.created_at < $4 AND example_table.last_updated_at > $5 AND example_table.last_updated_at < $6) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $7 AND example_table.archived_at IS NULL AND example_table.belongs_to_household = $8) as total_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $9 AND example_table.archived_at IS NULL AND example_table.belongs_to_household = $10 AND example_table.created_at > $11 AND example_table.created_at < $12 AND example_table.last_updated_at > $13 AND example_table.last_updated_at < $14 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []any{
			"value",
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			"value",
			exampleUser.ID,
			"value",
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		exampleJoins := []string{
			"things on stuff.thing_id=things.id",
		}
		exampleWhere := squirrel.ILike{
			"key": "value",
		}

		actualQuery, actualArgs := q.buildListQueryWithILike(ctx, exampleTableName, exampleJoins, nil, exampleWhere, exampleOwnershipColumn, exampleColumns, exampleUser.ID, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})

	T.Run("with group bys", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT column_ate, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $1 AND example_table.archived_at IS NULL AND example_table.belongs_to_household = $2 AND example_table.created_at > $3 AND example_table.created_at < $4 AND example_table.last_updated_at > $5 AND example_table.last_updated_at < $6) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $7 AND example_table.archived_at IS NULL AND example_table.belongs_to_household = $8) as total_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $9 AND example_table.archived_at IS NULL AND example_table.belongs_to_household = $10 AND example_table.created_at > $11 AND example_table.created_at < $12 AND example_table.last_updated_at > $13 AND example_table.last_updated_at < $14 GROUP BY example_table.id, things ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []any{
			"value",
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			"value",
			exampleUser.ID,
			"value",
			exampleUser.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		exampleJoins := []string{
			"things on stuff.thing_id=things.id",
		}
		exampleWhere := squirrel.ILike{
			"key": "value",
		}
		groupBys := []string{"things"}

		actualQuery, actualArgs := q.buildListQueryWithILike(ctx, exampleTableName, exampleJoins, groupBys, exampleWhere, exampleOwnershipColumn, exampleColumns, exampleUser.ID, false, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})

	T.Run("for admin with archived", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()
		filter.IncludeArchived = func(x bool) *bool { return &x }(true)

		expectedQuery := "SELECT column_ate, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $1 AND example_table.created_at > $2 AND example_table.created_at < $3 AND example_table.last_updated_at > $4 AND example_table.last_updated_at < $5) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $6) as total_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $7 AND (1=1) AND example_table.created_at > $8 AND example_table.created_at < $9 AND example_table.last_updated_at > $10 AND example_table.last_updated_at < $11 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []any{
			"value",
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			"value",
			"value",
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		exampleJoins := []string{
			"things on stuff.thing_id=things.id",
		}
		exampleWhere := squirrel.ILike{
			"key": "value",
		}

		actualQuery, actualArgs := q.buildListQueryWithILike(ctx, exampleTableName, exampleJoins, nil, exampleWhere, exampleOwnershipColumn, exampleColumns, exampleUser.ID, true, filter)

		assertArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSQLQuerier_buildTotalCountQueryWithILike(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		expectedValue := "stuff"

		expectedQuery := "SELECT COUNT(table.id) FROM table JOIN thing on another thing WHERE things ILIKE ? AND table.archived_at IS NULL"
		expectedArgs := []any{
			expectedValue,
		}

		actualQuery, actualArgs := q.buildTotalCountQueryWithILike(ctx, "table", []string{"thing on another thing"}, squirrel.ILike{"things": expectedValue}, "belongs_to_user", "user_id", true, false)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})

	T.Run("with nil where", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		expectedQuery := "SELECT COUNT(table.id) FROM table JOIN thing on another thing WHERE table.archived_at IS NULL"
		expectedArgs := []any(nil)

		actualQuery, actualArgs := q.buildTotalCountQueryWithILike(ctx, "table", []string{"thing on another thing"}, nil, "belongs_to_user", "user_id", true, false)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestSQLQuerier_buildFilteredCountQueryWithILike(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		expectedValue := "stuff"

		expectedQuery := "SELECT COUNT(table.id) FROM table JOIN thing on another thing WHERE things ILIKE ? AND table.archived_at IS NULL"
		expectedArgs := []any{
			expectedValue,
		}

		actualQuery, actualArgs := q.buildFilteredCountQueryWithILike(ctx, "table", []string{"thing on another thing"}, squirrel.ILike{"things": expectedValue}, "belongs_to_user", "user_id", true, false, types.DefaultQueryFilter())
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
