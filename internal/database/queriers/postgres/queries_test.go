package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/prixfixeco/api_server/pkg/types"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

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
		"column_one",
		"column_two",
		"column_three",
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_on IS NULL AND example_table.belongs_to_household = $1 AND key = $2 AND example_table.created_on > $3 AND example_table.created_on < $4 AND example_table.last_updated_on > $5 AND example_table.last_updated_on < $6) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_on IS NULL AND example_table.belongs_to_household = $7 AND key = $8) as total_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE example_table.archived_on IS NULL AND example_table.belongs_to_household = $9 AND key = $10 AND example_table.created_on > $11 AND example_table.created_on < $12 AND example_table.last_updated_on > $13 AND example_table.last_updated_on < $14 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
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

		actualQuery, actualArgs := q.buildListQuery(ctx, exampleTableName, exampleJoins, nil, exampleWhere, exampleOwnershipColumn, exampleColumns, exampleUser.ID, false, filter, true)

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

		expectedQuery := "SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.archived_on IS NULL AND example_table.created_on > $1 AND example_table.created_on < $2 AND example_table.last_updated_on > $3 AND example_table.last_updated_on < $4) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.archived_on IS NULL) as total_count FROM example_table WHERE example_table.created_on > $5 AND example_table.created_on < $6 AND example_table.last_updated_on > $7 AND example_table.last_updated_on < $8 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.buildListQuery(ctx, exampleTableName, nil, nil, nil, exampleOwnershipColumn, exampleColumns, exampleUser.ID, true, filter, true)

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
		filter.IncludeArchived = true

		expectedQuery := "SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table WHERE example_table.created_on > $1 AND example_table.created_on < $2 AND example_table.last_updated_on > $3 AND example_table.last_updated_on < $4) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table) as total_count FROM example_table WHERE example_table.created_on > $5 AND example_table.created_on < $6 AND example_table.last_updated_on > $7 AND example_table.last_updated_on < $8 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := q.buildListQuery(ctx, exampleTableName, nil, nil, nil, exampleOwnershipColumn, exampleColumns, exampleUser.ID, true, filter, true)

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
		"column_one",
		"column_two",
		"column_three",
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		exampleUser := fakes.BuildFakeUser()
		filter := fakes.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $1 AND example_table.archived_on IS NULL AND example_table.belongs_to_household = $2 AND example_table.created_on > $3 AND example_table.created_on < $4 AND example_table.last_updated_on > $5 AND example_table.last_updated_on < $6) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $7 AND example_table.archived_on IS NULL AND example_table.belongs_to_household = $8) as total_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $9 AND example_table.archived_on IS NULL AND example_table.belongs_to_household = $10 AND example_table.created_on > $11 AND example_table.created_on < $12 AND example_table.last_updated_on > $13 AND example_table.last_updated_on < $14 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
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

		expectedQuery := "SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $1 AND example_table.archived_on IS NULL AND example_table.belongs_to_household = $2 AND example_table.created_on > $3 AND example_table.created_on < $4 AND example_table.last_updated_on > $5 AND example_table.last_updated_on < $6) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $7 AND example_table.archived_on IS NULL AND example_table.belongs_to_household = $8) as total_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $9 AND example_table.archived_on IS NULL AND example_table.belongs_to_household = $10 AND example_table.created_on > $11 AND example_table.created_on < $12 AND example_table.last_updated_on > $13 AND example_table.last_updated_on < $14 GROUP BY example_table.id, things ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
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
		filter.IncludeArchived = true

		expectedQuery := "SELECT column_one, column_two, column_three, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $1 AND example_table.created_on > $2 AND example_table.created_on < $3 AND example_table.last_updated_on > $4 AND example_table.last_updated_on < $5) as filtered_count, (SELECT COUNT(example_table.id) FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $6) as total_count FROM example_table JOIN things on stuff.thing_id=things.id WHERE key ILIKE $7 AND (1=1) AND example_table.created_on > $8 AND example_table.created_on < $9 AND example_table.last_updated_on > $10 AND example_table.last_updated_on < $11 GROUP BY example_table.id ORDER BY example_table.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
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

		expectedQuery := "SELECT COUNT(table.id) FROM table JOIN thing on another thing WHERE things ILIKE ? AND table.archived_on IS NULL"
		expectedArgs := []interface{}{
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

		expectedQuery := "SELECT COUNT(table.id) FROM table JOIN thing on another thing WHERE table.archived_on IS NULL"
		expectedArgs := []interface{}(nil)

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

		expectedQuery := "SELECT COUNT(table.id) FROM table JOIN thing on another thing WHERE things ILIKE ? AND table.archived_on IS NULL"
		expectedArgs := []interface{}{
			expectedValue,
		}

		actualQuery, actualArgs := q.buildFilteredCountQueryWithILike(ctx, "table", []string{"thing on another thing"}, squirrel.ILike{"things": expectedValue}, "belongs_to_user", "user_id", true, false, types.DefaultQueryFilter())
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
