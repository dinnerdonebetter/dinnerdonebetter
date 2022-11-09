package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/pointers"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestQuerier_ChangeMealPlanTaskStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeMealPlanTaskStatusChangeRequestInput()
		// so we trigger setting the time function
		exampleInput.Status = pointers.String(types.MealPlanTaskStatusFinished)

		ctx := context.Background()
		c, db := buildTestClient(t)

		exampleTime := time.Now()
		c.timeFunc = func() time.Time {
			return exampleTime
		}

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Status,
			exampleInput.StatusExplanation,
			exampleTime,
		}

		db.ExpectExec(formatQueryForSQLMock(changeMealPlanTaskStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		actual := c.ChangeMealPlanTaskStatus(ctx, exampleInput)
		assert.NoError(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}
