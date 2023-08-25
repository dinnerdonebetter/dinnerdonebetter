package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
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
