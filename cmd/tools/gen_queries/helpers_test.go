package main

import (
	"testing"

	"github.com/cristalhq/builq"
	"github.com/stretchr/testify/assert"
)

func Test_applyToEach(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleInput := []string{
			"things",
			"and",
			"stuff",
		}

		callCount := 0
		exampleFunc := func(x string) string {
			callCount += 1
			return x
		}

		expected := []string{
			"things",
			"and",
			"stuff",
		}
		actual := applyToEach(exampleInput, exampleFunc)

		assert.Equal(t, callCount, len(exampleInput))
		assert.Equal(t, expected, actual)
	})
}

func Test_buildArchiveQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := `UPDATE things SET archived_at = NOW() WHERE things.archived_at IS NULL AND things.id = sqlc.arg(id) AND things.whatever = sqlc.arg(whatever_id)` + "\n"

		actual := buildArchiveQuery("things", "whatever")

		assert.Equal(t, expected, actual)
	})
}

func Test_buildCreateQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := `INSERT INTO things (bingo,
	bongo) VALUES (sqlc.arg(bingo),
	sqlc.arg(bongo))
`

		actual := buildCreateQuery("things", []string{"bingo", "bongo"})

		assert.Equal(t, expected, actual)
	})
}

func Test_buildExistenceCheckQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := `SELECT EXISTS ( SELECT things.id FROM things WHERE things.archived_at IS NULL AND things.id = sqlc.arg(id) addendum )
`

		actual := buildExistenceCheckQuery("things", "addendum")

		assert.Equal(t, expected, actual)
	})
}

func Test_buildFilteredColumnCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := `(
	    SELECT
	        COUNT(things.id)
	    FROM
	        things
	    WHERE
            things.archived_at IS NULL
            AND things.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
            AND things.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
           
AND (
                things.last_updated_at IS NULL
                OR things.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
            )
            AND (
                things.last_updated_at IS NULL
                OR things.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
            )

            AND addendum

	) as filtered_count
`

		actual := buildFilteredColumnCountQuery("things", true, "addendum")

		assert.Equal(t, expected, actual)
	})
}

func Test_buildRawQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var whatever builq.Builder

		builder := whatever.Addf("SELECT * FROM things")

		expected := "SELECT * FROM things\n"
		actual := buildRawQuery(builder)

		assert.Equal(t, expected, actual)
	})
}

func Test_buildSelectQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := `SELECT column1 FROM table JOIN another table WHERE addendum condition 1 AND addendum condition 2  AND  table.archived_at IS NULL AND table.id = sqlc.arg(id)` + "\n"

		actual := buildSelectQuery(
			"table",
			[]string{
				"column1",
			},
			[]string{
				"another table",
			},
			"addendum condition 1",
			"addendum condition 2",
		)

		assert.Equal(t, expected, actual)
	})
}

func Test_buildTotalColumnCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := `(
	    SELECT
	        COUNT(things.id)
	    FROM
	        things
	    WHERE
            things.archived_at IS NULL
           

            AND whatever

	) as total_count
`
		actual := buildTotalColumnCountQuery("things", "whatever")

		assert.Equal(t, expected, actual)
	})
}

func Test_filterForInsert(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exception := "whatever"
		exampleColumns := []string{
			"things",
			"and",
			"stuff",
			createdAtColumn,
			lastUpdatedAtColumn,
			archivedAtColumn,
			exception,
		}

		expected := []string{
			"things",
			"and",
			"stuff",
		}
		actual := filterForInsert(exampleColumns, exception)

		assert.Equal(t, expected, actual)
	})
}

func Test_formatQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		example := `SELECT stuff 
FROM things 
				WHERE id = 1
`

		expected := "SELECT stuff FROM things WHERE id = 1;"
		actual := formatQuery(example)

		assert.Equal(t, expected, actual)
	})
}

func Test_fullColumnName(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := "things.stuff"
		actual := fullColumnName("things", "stuff")

		assert.Equal(t, expected, actual)
	})
}

func Test_mergeColumns(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := []string{
			"webhooks.id",
			"webhooks.name",
			"webhooks.content_type",
			"webhooks.url",
			"webhooks.method",
			"webhook_trigger_events.id",
			"webhook_trigger_events.trigger_event",
			"webhook_trigger_events.belongs_to_webhook",
			"webhook_trigger_events.created_at",
			"webhook_trigger_events.archived_at",
			"webhooks.created_at",
			"webhooks.last_updated_at",
			"webhooks.archived_at",
			"webhooks.belongs_to_household",
		}

		actual := mergeColumns(
			applyToEach(webhooksColumns, func(s string) string {
				return fullColumnName(webhooksTableName, s)
			}),
			applyToEach(webhookTriggerEventsColumns, func(s string) string {
				return fullColumnName(webhookTriggerEventsTableName, s)
			}),
			5,
		)

		assert.Equal(t, expected, actual)
	})
}

func Test_buildUpdateQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := `UPDATE things SET last_updated_at = NOW(), column1 = sqlc.arg(column1),
	column2 = sqlc.arg(column2) WHERE archived_at IS NULL  AND things.belongs_to_user = sqlc.arg(belongs_to_user_id) AND id = sqlc.arg(id)` + "\n"

		actual := buildUpdateQuery("things", []string{"column1", "column2"}, belongsToUserColumn)

		assert.Equal(t, expected, actual)
	})
}
