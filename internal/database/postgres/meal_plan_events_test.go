package postgres

import (
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/prixfixeco/backend/pkg/types"
)

func buildMockRowsFromMealPlanEvents(includeCounts bool, filteredCount uint64, mealPlans ...*types.MealPlanEvent) *sqlmock.Rows {
	columns := mealPlanEventsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range mealPlans {
		rowValues := []driver.Value{
			x.ID,
			x.Notes,
			x.StartsAt,
			x.EndsAt,
			x.MealName,
			x.BelongsToMealPlan,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(mealPlans))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}
